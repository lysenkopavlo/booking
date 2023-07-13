package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lysenkopavlo/booking/internal/config"
	"github.com/lysenkopavlo/booking/internal/handler"
)

func routes(app *config.AppConfig) http.Handler {
	// Creating a multiplexer
	mux := chi.NewRouter()

	// Testing middleware of this package
	mux.Use(middleware.Recoverer)

	//NoSurf for CSRF Token
	mux.Use(NoSurf)

	mux.Use(SessionLoad)

	mux.Get("/", handler.Repo.Home)
	mux.Get("/about", handler.Repo.About)
	mux.Get("/generals-quarters", handler.Repo.Generals)
	mux.Get("/majors-suite", handler.Repo.Majors)

	mux.Get("/search-availability", handler.Repo.Availability)
	mux.Post("/search-availability", handler.Repo.PostAvailability)
	mux.Post("/search-availability-json", handler.Repo.AvailabilityJSON)

	mux.Get("/choose-room/{id}", handler.Repo.ChooseRoom)

	mux.Get("/contacts", handler.Repo.Contacts)

	mux.Get("/make-reservation", handler.Repo.Reservation)
	mux.Post("/make-reservation", handler.Repo.PostReservation)
	mux.Get("/reservation-summary", handler.Repo.ReservationSummary)

	// Telling to app where are the files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
