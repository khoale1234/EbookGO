package adminhandler

import (
	"Ebook/internal/config"
	"Ebook/internal/driver"
	"Ebook/internal/models"
	"Ebook/internal/render"
	"Ebook/internal/repository"
	dbrepo "Ebook/internal/repository/dprepo"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var RepoAdmin *AdminHandler

type AdminHandler struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

func NewAdminRepository(a *config.AppConfig, db *driver.DB) *AdminHandler {
	return &AdminHandler{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}
func NewAdminHandlers(r *AdminHandler) {
	RepoAdmin = r
}
func (m *AdminHandler) AdminAllBooks(w http.ResponseWriter, r *http.Request) {
	var books []models.BookDtls
	books, err := m.DB.BookRepo().GetAllBooks()
	if err != nil {
		log.Println(err)
		return
	}
	data := make(map[string]interface{})
	data["books"] = books
	render.Template(w, r, "admin-allbooks.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
func (m *AdminHandler) AdminHome(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-home.page.tmpl", &models.TemplateData{})
}
func (m *AdminHandler) AdminOrders(w http.ResponseWriter, r *http.Request) {
	session, _ := m.App.Session.Get(r, "posty")
	_, ok := session.Values["userId"].(int)
	if !ok {
		http.Redirect(w, r, "login", http.StatusSeeOther)
	}
	books, err := m.DB.OrderRepo().GetAllOrder()
	if err != nil {
		log.Println(err)
		return
	}
	data := make(map[string]interface{})
	data["books"] = books
	render.Template(w, r, "admin-order.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
func (m *AdminHandler) AdminBookDelete(w http.ResponseWriter, r *http.Request) {
	bid, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Println(err)
		return
	}
	err = m.DB.BookRepo().DeleteBook(bid)
	if err != nil {
		log.Println(err)
		return
	}
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}
func (m *AdminHandler) AdminEditBook(w http.ResponseWriter, r *http.Request) {
	bid, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Println(err)
		return
	}
	book, err := m.DB.BookRepo().GetBookById(bid)
	if err != nil {
		log.Println(err)
	}
	data := make(map[string]interface{})
	data["book"] = book
	render.Template(w, r, "admin-editbook.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
func (m *AdminHandler) AdminPostEditBook(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	// Get form values
	bookID, _ := strconv.Atoi(r.FormValue("bid"))
	bookName := r.FormValue("bname")
	authorName := r.FormValue("Aname")
	price, _ := strconv.ParseFloat(r.FormValue("bprice"), 32)
	bookStatus := r.FormValue("bstatus")
	log.Println("hello", bookName)
	var errMsg string
	var flash string
	err = m.DB.BookRepo().UpdateEditBook(bookName, authorName, bookStatus, float32(price), bookID)
	if err != nil {
		log.Println(err)
		errMsg = "Something wrong on server"
		return
	} else {
		flash = "Edit book successfully"
	}
	render.Template(w, r, "admin-editbook.page.tmpl", &models.TemplateData{
		Error: errMsg,
		Flash: flash,
	})
}
func (m *AdminHandler) AdminAddBook(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-addbook.page.tmpl", &models.TemplateData{})
}
func (m *AdminHandler) PostAdminAddBook(w http.ResponseWriter, r *http.Request) {
	log.Println("hello")
	err := r.ParseMultipartForm(10 << 20) // Limit to 10 MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get form values
	bookName := r.FormValue("bname")
	authorName := r.FormValue("Aname")
	price := r.FormValue("bprice")
	bookType := r.FormValue("btype")
	bookStatus := r.FormValue("bstatus")
	file, handler, err := r.FormFile("bookimg")
	if err != nil {
		log.Println("file là:", handler.Filename)
		http.Error(w, "Unable to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	book := models.BookDtls{
		BookName:     bookName,
		Author:       authorName,
		Price:        price,
		BookCategory: bookType,
		Status:       bookStatus,
		PhotoName:    handler.Filename,
		Email:        "admin",
	}
	var errorMsg string
	var flash string
	err = m.DB.BookRepo().AddBook(book)
	if err != nil {
		log.Println(err)
	}
	if err == nil {
		path := filepath.Join("static", "book")
		fmt.Println("File Path:", path)

		// Tạo thư mục nếu nó không tồn tại
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			http.Error(w, "Unable to create directory", http.StatusInternalServerError)
			return
		}
		// Ghi file vào đường dẫn đã xây dựng
		filePath := filepath.Join(path, handler.Filename)
		out, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Unable to create file", http.StatusInternalServerError)
			errorMsg = "Something went wrong on server"
			return
		}
		defer out.Close()

		// Ghi nội dung file
		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(w, "Unable to write file", http.StatusInternalServerError)
			errorMsg = "Something went wrong on server"
			return
		}
		fmt.Println("File written to:", filePath)
		flash = "Sell book successfully"
	}
	render.Template(w, r, "admin-addbook.page.tmpl", &models.TemplateData{
		Flash: flash,
		Error: errorMsg,
	})
}
