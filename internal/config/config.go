package config

import (
	"Ebook/internal/models"
	"html/template"
	"log"

	"github.com/gorilla/sessions"
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
