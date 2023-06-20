// Package "handler" is used to handle http requests and responses
// Also here I'm using a "repository pattern"

package handler

import (
	"net/http"

	"github.com/lysenkopavlo/booking/pkg/config"
	"github.com/lysenkopavlo/booking/pkg/models"
	"github.com/lysenkopavlo/booking/pkg/render"
)

// Repo the repository used by handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
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

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

// About is an about page handler
func (rp *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "About Handler"

	remoteIP := rp.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Reservation renders the make reservation page and displays a form
func (rp *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, "make-reservation.page.tmpl", &models.TemplateData{})
}

// Generals renders page for General's quarters room
func (rp *Repository) Generals(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders page for Major's suite room
func (rp *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders page for Major's suite room
func (rp *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "search-availability.page.tmpl", &models.TemplateData{})
}

// Contact renders the contact page
func (rp *Repository) Contacts(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "contacts.page.tmpl", &models.TemplateData{})
}
