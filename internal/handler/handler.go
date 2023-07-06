// Package "handler" is used to handle http requests and responses
// Also here I'm using a "repository pattern"

package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	err := render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
}

// PostReservation handles the posting of a reservation form
func (rp *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	//err = errors.New("something happened!")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)
	//	form.Has("first_name", r)
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

// PostAvailability
func (rp *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	_, err := w.Write([]byte(fmt.Sprintf("Your starting date is: %s and your ending date is: %s", start, end)))
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
	err := render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
}
