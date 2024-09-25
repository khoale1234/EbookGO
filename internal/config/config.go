package config

import (
	"Ebook/internal/models"
	"html/template"
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AppConfig struct {
	TemplateCache map[string]*template.Template
	UseCache      bool
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       sessions.CookieStore
	MailChan      chan models.MailData
}
type Config struct {
	GoogleLoginConfig oauth2.Config
}

var AppConfigOauth Config

func GoogleConfig() oauth2.Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	AppConfigOauth.GoogleLoginConfig = oauth2.Config{
		RedirectURL:  "http://localhost:8080/oauth/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/user.phonenumbers.read"},
		Endpoint: google.Endpoint,
	}

	return AppConfigOauth.GoogleLoginConfig
}
