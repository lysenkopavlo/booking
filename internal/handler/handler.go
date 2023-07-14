// Package "handler" is used to handle http requests and responses
// Also here I'm using a "repository pattern"

package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lysenkopavlo/booking/internal/config"
	"github.com/lysenkopavlo/booking/internal/driver"
	"github.com/lysenkopavlo/booking/internal/forms"
	"github.com/lysenkopavlo/booking/internal/helpers"
	"github.com/lysenkopavlo/booking/internal/models"
	"github.com/lysenkopavlo/booking/internal/render"
	"github.com/lysenkopavlo/booking/internal/repository"
	"github.com/lysenkopavlo/booking/internal/repository/dbrepo"
)

// Repo the repository used by handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DataBaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandler(r *Repository) {
	Repo = r
}

// Home is a home page handler
func (rp *Repository) Home(w http.ResponseWriter, r *http.Request) {

	remoteIP := r.RemoteAddr
	rp.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	err := render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

}

// About is an about page handler
func (rp *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "About Handler"

	remoteIP := rp.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	err := render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
}

// Reservation renders the make reservation page and displays a form
func (rp *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	// Pull out reservation data from session
	res, ok := rp.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation from session"))
		return
	}

	//Getting RoomName by its id
	room, err := rp.DB.GetRoomByID(res.RoomID)
	if err != nil {
		helpers.ServerError(w, errors.New("cannot get a room name by its id!"))
		return
	}

	res.Room.RoomName = room.RoomName

	// Put updated reservation data into session
	rp.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed
	//stringMap["room_name"] = room.RoomName

	data := make(map[string]interface{})
	data["reservation"] = res

	err = render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
}

// PostReservation handles the posting of a reservation form
func (rp *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {

	// Pull out reservation data from the session
	reservation, ok := rp.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation data from session"))
		return
	}

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	form := forms.New(r.PostForm)
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newReservationID, err := rp.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// building a model for room_restriction table
	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionId: 1,
	}

	err = rp.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Putting user's values into context
	rp.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Generals renders page for General's quarters room
func (rp *Repository) Generals(w http.ResponseWriter, r *http.Request) {

	err := render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
}

// Majors renders page for Major's suite room
func (rp *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
}

// Availability renders page for search page
func (rp *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
}

// PostAvailability renders the search-availability page
func (rp *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Looking for available rooms
	rooms, err := rp.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		// it means there is no free rooms
		// I want to show an error message
		rp.App.Session.Put(r.Context(), "error", "There is NO AVAILABLE ROOMS!")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	// if there are available rooms
	// we render new page with this rooms
	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	rp.App.Session.Put(r.Context(), "reservation", res)

	err = render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

}

// jsonResponse
type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON
func (rp *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	res := jsonResponse{
		OK:      true,
		Message: "Available",
	}
	out, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	log.Println(string(out))
	w.Header().Set("Content-type", "application/json")
	_, err = w.Write(out)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
}

// Contact renders the contact page
func (rp *Repository) Contacts(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, r, "contacts.page.tmpl", &models.TemplateData{})
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
}

func (rp *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := rp.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		rp.App.ErrorLog.Println("can't get item from session")
		rp.App.Session.Put(r.Context(), "error", "can't get item from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	rp.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	err := render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
}

func (rp *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	// Getting room_id from url /choose-room/{id}
	// to put it into reservation
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// getting values reservation from session
	// to update its roomID field with new value
	res, ok := rp.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get a reservation values from session"))
		return
	}
	res.RoomID = roomID
	rp.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
