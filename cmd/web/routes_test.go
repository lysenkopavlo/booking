package main

import (
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/lysenkopavlo/booking/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// Do nothing
	default:
		t.Errorf("the wrong type of %T", v)
	}

}

//	to display test coverage in more details
//	use this command:
//	go test -coverprofile=coverage.out && go tool cover -html=coverage.out
