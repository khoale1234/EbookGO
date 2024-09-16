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
	// mux.Use(NoSurf)
	mux.Use(middleware.Recoverer)
	// Tạo file server cho tài nguyên tĩnh

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	mux.Get("/login", handlers.Repo.Login)
	mux.Get("/register", handlers.Repo.Register)
	mux.Post("/register", handlers.Repo.PostRegister)
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/all_new_books", handlers.Repo.AllNewBooks)
	mux.Get("/all_old_books", handlers.Repo.AllOldBooks)
	mux.Get("/all_recent_books", handlers.Repo.AllRecentBooks)
	mux.Get("/book_detail/{id}", handlers.Repo.BookDetail)
	mux.Get("/cart", handlers.Repo.AddCart)
	mux.Get("/checkout", handlers.Repo.Checkout)
	mux.Get("/edit-profile", handlers.Repo.EditProfile)
	mux.Post("/edit-profile", handlers.Repo.PostEditProfile)
	mux.Get("/logout", handlers.Repo.Logout)
	mux.Post("/login", handlers.Repo.PostLogin)
	mux.Get("/old-books", handlers.Repo.OldBooks)
	mux.Get("/order-success", handlers.Repo.OrderSuccess)
	mux.Get("/order", handlers.Repo.GetOrderByUser)
	mux.Get("/search", handlers.Repo.Search)
	mux.Get("/user-address", handlers.Repo.UserAddress)
	mux.Get("/setting", handlers.Repo.Setting)
	mux.Get("/helpline", handlers.Repo.Helpline)
	mux.Get("/sellbook", handlers.Repo.SellBook)
	mux.Post("/sellbook", handlers.Repo.PostSellBook)
	mux.Get("/delete_old_book", handlers.Repo.DeleteOldBook)
	mux.Get("/remove_book", handlers.Repo.RemoveBook)
	mux.Post("/order", handlers.Repo.Order)
	mux.Post("/search", handlers.Repo.SearchBook)
	mux.Route("/admin", func(r chi.Router) {
		r.Handle("/static/*", http.StripPrefix("admin/static", fileServer))
		r.Get("/allbooks", handlers.Repo.AdminAllBooks)
		r.Get("/home", handlers.Repo.AdminHome)
		r.Get("/order", handlers.Repo.AdminOrders)
		r.Get("/editbook", handlers.Repo.EditBook)
		r.Get("/bookdelete", handlers.Repo.BookDelete)
		r.Get("/editbook", handlers.Repo.EditBook)
		r.Post("/bookedit", handlers.Repo.PostEditBook)
		r.Get("/addbook", handlers.Repo.AdminAddBook)
		r.Post("/addbook", handlers.Repo.PostAdminAddBook)

	})
	return Rep.Authenticate(mux)
}
