package main

import (
	"Ebook/internal/config"
	"Ebook/internal/driver"
	adminhandler "Ebook/internal/handlers/admin_handler"
	userhandler "Ebook/internal/handlers/user_handler"
	"Ebook/internal/render"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
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
	gothic.Store = session
	app.Session = *session
	//config the gg oauth
	config.GoogleConfig()
	log.Println("Connecting to database ...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=book_app user=postgres password=07052002")
	log.Println("Database connected")
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Println(err)
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false
	rep := NewRep(&app, db)
	NewMid(rep)
	repo_user := userhandler.NewUserRepository(&app, db)
	userhandler.NewUserHandlers(repo_user)
	repo_admin := adminhandler.NewAdminRepository(&app, db)
	adminhandler.NewAdminHandlers(repo_admin)

	render.NewRenderer(&app)
	return db, nil
}
