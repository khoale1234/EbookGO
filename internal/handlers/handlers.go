package handlers

import (
	"Ebook/internal/config"
	"Ebook/internal/driver"
	"Ebook/internal/forms"
	"Ebook/internal/helpers"
	"Ebook/internal/models"
	"Ebook/internal/render"
	"Ebook/internal/repository"
	dbrepo "Ebook/internal/repository/dprepo"

	"log"
	"net/http"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	var newbooks []models.BookDtls
	newbooks, err := m.DB.GetSomeNewBooks()
	if err != nil {
		helpers.ServerError(w, err)
	}
	data := make(map[string]interface{})
	data["newbooks"] = newbooks
	var recentbooks []models.BookDtls
	recentbooks, err = m.DB.GetSomeRecentBooks()
	if err != nil {
		helpers.ServerError(w, err)
	}
	data["recentbooks"] = recentbooks
	var oldbooks []models.BookDtls
	oldbooks, err = m.DB.GetSomeOldBooks()
	if err != nil {
		helpers.ServerError(w, err)
	}
	data["oldbooks"] = oldbooks
	render.Template(w, r, "index.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
func (m *Repository) AllNewBooks(w http.ResponseWriter, r *http.Request) {
	var newbooks []models.BookDtls
	newbooks, err := m.DB.GetNewBooks()
	if err != nil {
		helpers.ServerError(w, err)
	}
	data := make(map[string]interface{})
	data["newbooks"] = newbooks
	render.Template(w, r, "all_new_books.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
func (m *Repository) AllOldBooks(w http.ResponseWriter, r *http.Request) {
	var oldbooks []models.BookDtls
	oldbooks, err := m.DB.GetOldBooks()
	if err != nil {
		helpers.ServerError(w, err)
	}
	data := make(map[string]interface{})
	data["oldbooks"] = oldbooks
	render.Template(w, r, "all_old_books.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
func (m *Repository) AllRecentBooks(w http.ResponseWriter, r *http.Request) {
	var recentbooks []models.BookDtls
	recentbooks, err := m.DB.GetRecentBooks()
	if err != nil {
		helpers.ServerError(w, err)
	}
	data := make(map[string]interface{})
	data["recentbooks"] = recentbooks
	render.Template(w, r, "all_recent_books.page.tmpl", &models.TemplateData{
		Data: data,
	})
	render.Template(w, r, "all_old_books.page.tmpl", &models.TemplateData{})
}
func (m *Repository) BookDetail(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "book_detail.page.tmpl", &models.TemplateData{})
}
func (m *Repository) AddCart(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "cart.page.tmpl", &models.TemplateData{})
}
func (m *Repository) Checkout(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "checkout.page.tmpl", &models.TemplateData{})
}
func (m *Repository) EditProfile(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "edit-profile.page.tmpl", &models.TemplateData{})
}
func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}
func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	err := m.App.Session.RenewToken(r.Context())
	if err != nil {
		log.Fatalf("Error on post %v", err)
	}
	err = r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	var email string
	var password string
	email = r.Form.Get("email")
	password = r.Form.Get("password")
	form := forms.New(r.PostForm)
	form.Required("email", "password")
	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}
	_, _, err = m.DB.Login(email, password)
	if err != nil {
		log.Println("login failed")
	} else {
		log.Println("login successful")
	}

}
func (m *Repository) OldBooks(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "old_books.page.tmpl", &models.TemplateData{})
}
func (m *Repository) OrderSuccess(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "order_success.page.tmpl", &models.TemplateData{})
}
func (m *Repository) Register(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "register.page.tmpl", &models.TemplateData{})
}
func (m *Repository) PostRegister(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	name := r.Form.Get("name")
	phone := r.Form.Get("phone")
	form := forms.New(r.PostForm)
	form.Required("email", "name", "password", "phone")
	if !form.Valid() {
		render.Template(w, r, "register.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}
	var user = models.User{
		Email:    email,
		Password: password,
		Name:     name,
		Phone_no: phone,
	}
	log.Println(user)
	err = m.DB.Register(user)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Invalid register credentials")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	m.App.Session.Put(r.Context(), "flash", "Register successfully")
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}
func (m *Repository) Search(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search.page.tmpl", &models.TemplateData{})
}
func (m *Repository) UserAddress(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "user_address.page.tmpl", &models.TemplateData{})
}
func (m *Repository) Setting(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "setting.page.tmpl", &models.TemplateData{})
}
