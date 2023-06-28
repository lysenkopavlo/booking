package handler

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"github.com/lysenkopavlo/booking/internal/config"
	"github.com/lysenkopavlo/booking/internal/models"
	"github.com/lysenkopavlo/booking/internal/render"
)

var app config.AppConfig
var session *scs.SessionManager
var pathTemplates = "./../../templates"
var functions = template.FuncMap{}

func getRoutes() http.Handler {

	// what types i am going to put in the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	app.InProduction = false

	infoLogger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLogger

	errorLogger := log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLogger

	// Session is a variable to set user's session parameters
	session = scs.New()

	// setting lifetime for session
	session.Lifetime = 24 * time.Hour

	// is cockie persisting after closing a browser?
	session.Cookie.Persist = true // means don't want to clear cockie after closing browser

	// how restrict do you want to be about coockie?
	session.Cookie.SameSite = http.SameSiteLaxMode

	// insisting about coockie incryption
	session.Cookie.Secure = app.InProduction // for now

	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Println("Can't create a template cache")

	}
	app.TemplateCache = tc
	app.UseCache = true
	app.Session = session

	repo := NewRepo(&app)
	NewHandler(repo)
	render.NewTemplate(&app)
	// Creating a multiplexer
	mux := chi.NewRouter()
	// Testing middleware of this package
	mux.Use(middleware.Recoverer)

	// mux.Use(WriteToConsole)

	// turning off(comment out) CSRF token for tests
	// mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)

	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)

	mux.Get("/contacts", Repo.Contacts)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	// Telling to app where are the files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoads loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
func CreateTestTemplateCache() (map[string]*template.Template, error) {
	myCache := make(map[string]*template.Template)

	// get all of the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathTemplates))
	if err != nil {
		return myCache, err
	}

	// range throuh all the pages
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob(fmt.Sprintf("%s/*layout.tmpl", pathTemplates))
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*layout.tmpl", pathTemplates))
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
