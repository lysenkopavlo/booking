package main

import (
	"fmt"
	"log"
	"time"

	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/lysenkopavlo/booking/internal/config"
	"github.com/lysenkopavlo/booking/internal/handler"
	"github.com/lysenkopavlo/booking/internal/render"
)

const portNumber = ":8080"

// app is a variable for...
var app config.AppConfig

var session *scs.SessionManager

func main() {

	// Session is a variable to set user's session parameters
	session = scs.New()

	// setting lifetime for session
	session.Lifetime = 24 * time.Hour

	// Is cockie persisting after closing a browser?
	session.Cookie.Persist = true // means don't want to clear cockie after closing browser

	//How restrict do you want to be about coockie?
	session.Cookie.SameSite = http.SameSiteLaxMode

	// insisting about coockie incryption
	session.Cookie.Secure = app.InProduction // for now

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Can't create a template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false
	app.Session = session
	render.NewTemplate(&app)

	repo := handler.NewRepo(&app)

	handler.NewHandler(repo)

	fmt.Printf("Starting application on port %s\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
