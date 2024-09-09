package main

import (
	"Ebook/internal/config"
	"Ebook/internal/driver"
	"Ebook/internal/handlers"
	"Ebook/internal/render"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

const portNumber = ":8080"

var InfoLog *log.Logger
var app config.AppConfig
var ErrorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()
	defer close(app.MailChan)
	fmt.Println("Start email Listener ...")
	listenForMail()

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
func run() (*driver.DB, error) {
	// change this to true when in production
	app.InProduction = false
	InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = InfoLog
	ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = ErrorLog
	// set up the session

	session := sessions.NewCookieStore([]byte("HiOYcmnxUuj4KoQGgks63DB1DOGponWX"))
	session.Options.HttpOnly = true
	session.Options.SameSite = http.SameSiteLaxMode
	app.Session = *session
	log.Println("Connecting to database ...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=book_app user=postgres password=07052002")
	log.Println("Database connected")
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false
	rep := NewRep(&app, db)
	NewMid(rep)
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	return db, nil
}
