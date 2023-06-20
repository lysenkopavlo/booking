package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lysenkopavlo/booking/pkg/config"
	"github.com/lysenkopavlo/booking/pkg/handler"
)

func routes(a *config.AppConfig) http.Handler {
	// Creating a multiplexer
	mux := chi.NewRouter()
	// Testing middleware of this package
	mux.Use(middleware.Recoverer)

	mux.Use(WriteToConsole)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/", handler.Repo.Home)
	mux.Get("/about", handler.Repo.About)
	mux.Get("/generals-quarters", handler.Repo.Generals)
	mux.Get("/majors-suite", handler.Repo.Majors)
	mux.Get("/search-availability", handler.Repo.Availability)
	mux.Get("/contacts", handler.Repo.Contacts)

	// Telling to app where are the files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
