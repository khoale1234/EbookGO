package main

import (
	"Ebook/internal/config"
	adminhandler "Ebook/internal/handlers/admin_handler"
	userhandler "Ebook/internal/handlers/user_handler"

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
	mux.Get("/login", userhandler.RepoUser.Login)
	mux.Get("/auth/google_login", userhandler.RepoUser.GoogleLogin)
	mux.Get("/oauth/callback", userhandler.RepoUser.GoogleCallBack)
	mux.Get("/register", userhandler.RepoUser.Register)
	mux.Post("/register", userhandler.RepoUser.PostRegister)
	mux.Get("/", userhandler.RepoUser.Home)
	mux.Get("/all_new_books", userhandler.RepoUser.AllNewBooks)
	mux.Get("/all_old_books", userhandler.RepoUser.AllOldBooks)
	mux.Get("/all_recent_books", userhandler.RepoUser.AllRecentBooks)
	mux.Get("/book_detail/{id}", userhandler.RepoUser.BookDetail)
	mux.Get("/cart", userhandler.RepoUser.AddCart)
	mux.Get("/checkout", userhandler.RepoUser.Checkout)
	mux.Get("/edit-profile", userhandler.RepoUser.EditProfile)
	mux.Post("/edit-profile", userhandler.RepoUser.PostEditProfile)
	mux.Get("/logout", userhandler.RepoUser.Logout)
	mux.Post("/login", userhandler.RepoUser.PostLogin)
	mux.Get("/old-books", userhandler.RepoUser.OldBooks)
	mux.Get("/order-success", userhandler.RepoUser.OrderSuccess)
	mux.Get("/order", userhandler.RepoUser.GetOrderByUser)
	mux.Get("/search", userhandler.RepoUser.Search)
	mux.Get("/user-address", userhandler.RepoUser.UserAddress)
	mux.Get("/setting", userhandler.RepoUser.Setting)
	mux.Get("/helpline", userhandler.RepoUser.Helpline)
	mux.Get("/sellbook", userhandler.RepoUser.SellBook)
	mux.Post("/sellbook", userhandler.RepoUser.PostSellBook)
	mux.Get("/delete_old_book", userhandler.RepoUser.DeleteOldBook)
	mux.Get("/remove_book", userhandler.RepoUser.RemoveBook)
	mux.Post("/order", userhandler.RepoUser.Order)
	mux.Post("/search", userhandler.RepoUser.SearchBook)
	mux.Route("/admin", func(r chi.Router) {
		r.Handle("/static/*", http.StripPrefix("admin/static", fileServer))
		r.Get("/allbooks", adminhandler.RepoAdmin.AdminAllBooks)
		r.Get("/home", adminhandler.RepoAdmin.AdminHome)
		r.Get("/order", adminhandler.RepoAdmin.AdminOrders)
		r.Get("/bookdelete", adminhandler.RepoAdmin.AdminBookDelete)
		r.Get("/editbook", adminhandler.RepoAdmin.AdminEditBook)
		r.Post("/bookedit", adminhandler.RepoAdmin.AdminPostEditBook)
		r.Get("/addbook", adminhandler.RepoAdmin.AdminAddBook)
		r.Post("/addbook", adminhandler.RepoAdmin.PostAdminAddBook)

	})
	return Rep.Authenticate(mux)
}
