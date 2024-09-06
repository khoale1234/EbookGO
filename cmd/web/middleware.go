package main

import (
	"Ebook/internal/repository"
	"context"
	"log"
	"net/http"

	"github.com/justinas/nosurf"
)

var DB repository.DatabaseRepo

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

// loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
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
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("hello")
		id, ok := app.Session.Get(r.Context(), "userId").(int)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}
		user, err := DB.FindUserByID(id)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
