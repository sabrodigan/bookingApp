package main

import (
	"github.com/alexedwards/scs/v2"
	"github.com/sabrodigan/bookings-app/internal/config"
	"github.com/sabrodigan/bookings-app/internal/handlers"
	"github.com/sabrodigan/bookings-app/internal/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the main function
func main() {

	// change this to true when in production
	app.InProduction = false

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache: ", err)
	}
	done := make(chan bool)
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	handlers.TypeWriter("Starting the application server on", 50)
	handlers.TypeWriter(portNumber, 50)

	if app.InProduction != true {
		handlers.TypeWriter("\n\nStarting the development server...running\n", 75)
		go handlers.Spinner(200, done)
	}

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
