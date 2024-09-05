package main

import (
	"Ebook/internal/config"
	"Ebook/internal/handlers"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	//mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/all_new_books", handlers.Repo.AllNewBooks)
	mux.Get("/all_old_books", handlers.Repo.AllOldBooks)
	mux.Get("/all_recent_books", handlers.Repo.AllRecentBooks)
	mux.Get("/book_detail/{id}", handlers.Repo.BookDetail)
	mux.Get("/cart", handlers.Repo.AddCart)
	mux.Get("/checkout", handlers.Repo.AddCart)
	mux.Get("/edit-profile", handlers.Repo.EditProfile)
	mux.Get("/login", handlers.Repo.Login)
	mux.Post("/login", handlers.Repo.PostLogin)
	mux.Get("/old-books", handlers.Repo.OldBooks)
	mux.Get("/order-success", handlers.Repo.OrderSuccess)
	mux.Get("/register", handlers.Repo.Register)
	mux.Post("/register", handlers.Repo.PostRegister)
	mux.Get("/search", handlers.Repo.Search)
	mux.Get("/user-address", handlers.Repo.UserAddress)
	mux.Get("/setting", handlers.Repo.Setting)
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
