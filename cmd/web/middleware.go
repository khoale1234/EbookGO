package main

import (
	"Ebook/internal/config"
	"Ebook/internal/driver"
	"Ebook/internal/repository"
	dbrepo "Ebook/internal/repository/dprepo"
	"context"
	"net/http"

	"github.com/justinas/nosurf"
)

var Rep *Repo

type Repo struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRep(a *config.AppConfig, db *driver.DB) *Repo {
	return &Repo{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewMid(r *Repo) {
	Rep = r
}

// NoSurf adds CSRF protection to all Post requests
func NoSurf(next http.Handler) http.Handler {
	crsfHandler := nosurf.New(next)
	crsfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return crsfHandler
}

//	func Auth(next http.Handler) http.Handler {
//		return http.HandlerFunc((func(w http.ResponseWriter, r *http.Request) {
//			if !helpers.IsAuthenticated(r) {
//				session.Put(r.Context(), "error", "Log in first")
//				http.Redirect(w, r, "/login", http.StatusSeeOther)
//				return
//			}
//			next.ServeHTTP(w, r)
//		}))
//	}
func (m *Repo) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := m.App.Session.Get(r, "posty")
		id, ok := session.Values["userId"].(int)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}
		user, err := m.DB.FindUserByID(id)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
