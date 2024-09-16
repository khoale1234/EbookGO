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
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
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
	fmt.Println("render home")

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
	ID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
	}
	var book models.BookDtls
	book, err = m.DB.GetBookById(ID)
	log.Println(book)
	if err != nil {
		helpers.ServerError(w, err)
	}
	data := make(map[string]interface{})
	data["book"] = book
	render.Template(w, r, "book_detail.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
func (m *Repository) AddCart(w http.ResponseWriter, r *http.Request) {
	bid, err := strconv.Atoi(r.URL.Query().Get("bid"))
	if err != nil {
		log.Println(err)
		return
	}
	uid, err := strconv.Atoi(r.URL.Query().Get("uid"))
	if err != nil {
		log.Println(err)
		return
	}
	book, err := m.DB.GetBookById(bid)
	if err != nil {
		log.Println(err)
		return
	}
	price, err := strconv.ParseFloat(book.Price, 32)
	if err != nil {
		log.Println(err)
		return
	}
	cart := models.Cart{
		Bid:      bid,
		Uid:      uid,
		BookName: book.BookName,
		Author:   book.Author,
		Price:    price,
	}
	err = m.DB.AddCart(cart)
	if err != nil {
		log.Println(err)
		return
	}
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}
func (m *Repository) Checkout(w http.ResponseWriter, r *http.Request) {
	session, _ := m.App.Session.Get(r, "posty")
	id, ok := session.Values["userId"].(int)
	if !ok {
		http.Redirect(w, r, "login", http.StatusSeeOther)
	}
	var cart []models.Cart
	cart, totalPrice, err := m.DB.GetBookByUser(id)
	if err != nil {
		log.Println(err)
		return
	}
	data := make(map[string]interface{})
	data["cart"] = cart
	data["totalPrice"] = totalPrice
	render.Template(w, r, "checkout.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
func (m *Repository) EditProfile(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "edit-profile.page.tmpl", &models.TemplateData{})
}
func (m *Repository) PostEditProfile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Lấy giá trị từ form
	userID, err := strconv.Atoi(r.FormValue("id"))
	log.Println(userID)
	if err != nil {
		log.Println(err)
		return
	}
	fullName := r.FormValue("fname")
	phone := r.FormValue("fphone")
	email := r.FormValue("femail")
	password := r.FormValue("fpassword")
	err = bcrypt.CompareHashAndPassword([]byte(m.DB.CheckPassword(userID)), []byte(password))
	if err != nil {
		warning := "Wrong Password"
		render.Template(w, r, "edit-profile.page.tmpl", &models.TemplateData{
			Form:    nil,
			Warning: warning,
		})
	} else {
		err = m.DB.UpdateProfile(fullName, email, phone, userID)
		if err != nil {
			log.Println("loi o day ne", err)
			errMsg := "Something wrong on server"
			render.Template(w, r, "edit-profile.page.tmpl", &models.TemplateData{
				Form:  nil,
				Error: errMsg,
			})
		} else {
			flash := "Edited profile successfully"
			render.Template(w, r, "edit-profile.page.tmpl", &models.TemplateData{
				Form:  nil,
				Flash: flash,
			})
		}

	}
}
func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	// Call the delete session function
	session, _ := m.App.Session.Get(r, "posty")
	delete(session.Values, "userId")
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
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
	if email == "admin@gmail.com" && password == "adminpassword" {
		session, _ := m.App.Session.Get(r, "posty")
		session.Values["userId"] = 5
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/admin/home", http.StatusSeeOther)
	} else {
		id, _, err := m.DB.Login(email, password)
		var errMsg string
		if err != nil {
			log.Println(err)
			errMsg = "Login failed"
			render.Template(w, r, "login.page.tmpl", &models.TemplateData{
				Error: errMsg,
			})
			return
		}
		session, _ := m.App.Session.Get(r, "posty")
		session.Values["userId"] = id
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{})
	}
}
func (m *Repository) OldBooks(w http.ResponseWriter, r *http.Request) {
	session, _ := m.App.Session.Get(r, "posty")
	id, ok := session.Values["userId"].(int)
	if !ok {
		http.Redirect(w, r, "login", http.StatusSeeOther)
	}
	var user models.User
	user, err := m.DB.FindUserByID(id)
	if err != nil {
		log.Println(err)
		return
	}
	var oldbook []models.BookDtls
	oldbook, err = m.DB.GetBooksByOld(user.Email, "Old Book")
	if err != nil {
		log.Println(err)
		return
	}
	stringMap := make(map[string]string)
	stringMap["Email"] = user.Email
	data := make(map[string]interface{})
	data["oldbooks"] = oldbook
	render.Template(w, r, "old_books.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})

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
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}
func (m *Repository) Search(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search.page.tmpl", &models.TemplateData{})
}
func (m *Repository) UserAddress(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "user_address.page.tmpl", &models.TemplateData{})
}
func (m *Repository) Setting(w http.ResponseWriter, r *http.Request) {
	session, _ := m.App.Session.Get(r, "posty")
	_, ok := session.Values["userId"].(int)
	if !ok {
		http.Redirect(w, r, "login", http.StatusSeeOther)
	}
	render.Template(w, r, "setting.page.tmpl", &models.TemplateData{
		Data: nil,
	})
}
func (m *Repository) Helpline(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "helpline.page.tmpl", &models.TemplateData{})
}
func (m *Repository) SellBook(w http.ResponseWriter, r *http.Request) {
	session, _ := m.App.Session.Get(r, "posty")
	_, ok := session.Values["userId"].(int)
	if !ok {
		http.Redirect(w, r, "login", http.StatusSeeOther)
	}
	render.Template(w, r, "sell_book.page.tmpl", &models.TemplateData{})
}
func (m *Repository) PostSellBook(w http.ResponseWriter, r *http.Request) {
	log.Println("hello from sellbook")

	// Đảm bảo gọi ParseMultipartForm trước
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Lấy dữ liệu từ form
	userEmail := r.FormValue("user")
	bookName := r.FormValue("name")
	author := r.FormValue("author")
	price := r.FormValue("price")

	// Lấy file từ form
	file, handler, err := r.FormFile("bookimg")
	if err != nil {
		log.Println("file là:", handler.Filename)
		http.Error(w, "Unable to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Kiểm tra dữ liệu bắt buộc
	form := forms.New(r.PostForm)
	form.Required("name", "author", "price")

	oldbook := models.BookDtls{
		BookName:     bookName,
		Author:       author,
		Price:        price,
		BookCategory: "Old Book",
		Status:       "Active",
		PhotoName:    handler.Filename,
		Email:        userEmail,
	}

	err = m.DB.AddBook(oldbook)
	if err != nil {
		log.Println(err)
	}

	var flash string
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
			flash = "Something went wrong on server"
			return
		}
		defer out.Close()

		// Ghi nội dung file
		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(w, "Unable to write file", http.StatusInternalServerError)
			flash = "Something went wrong on server"
			return
		}
		fmt.Println("File written to:", filePath)
		flash = "Sell book successfully"
	}

	render.Template(w, r, "sell_book.page.tmpl", &models.TemplateData{
		Flash: flash,
	})
}
func (m *Repository) DeleteOldBook(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	log.Println(email)
	bid, err := strconv.Atoi(r.URL.Query().Get("bid"))
	if err != nil {
		log.Println(err)
		return
	}
	err = m.DB.OldBookDelete(email, "Old Book", bid)
	if err != nil {
		log.Println(err)
		return
	}
	http.Redirect(w, r, "/old-books", http.StatusSeeOther)

}
func (m *Repository) GetOrderByUser(w http.ResponseWriter, r *http.Request) {
	session, _ := m.App.Session.Get(r, "posty")
	id, ok := session.Values["userId"].(int)
	if !ok {
		http.Redirect(w, r, "login", http.StatusSeeOther)
	}
	user, err := m.DB.FindUserByID(id)
	if err != nil {
		log.Println(err)
		return
	}
	orderbooks, err := m.DB.GetBookOrder(user.Email)
	if err != nil {
		log.Println(err)
		return
	}
	data := make(map[string]interface{})
	data["orderbooks"] = orderbooks
	render.Template(w, r, "order.page.tmpl", &models.TemplateData{
		Data: data,
	})

}
func (m *Repository) RemoveBook(w http.ResponseWriter, r *http.Request) {
	bid, err := strconv.Atoi(r.URL.Query().Get("bid"))
	if err != nil {
		log.Println(err)
		return
	}
	cid, err := strconv.Atoi(r.URL.Query().Get("cid"))
	if err != nil {
		log.Println(err)
		return
	}
	uid, err := strconv.Atoi(r.URL.Query().Get("uid"))
	if err != nil {
		log.Println(err)
		return
	}
	err = m.DB.DeleteBookC(bid, uid, cid)
	if err != nil {
		log.Println(err)
		return
	}
	http.Redirect(w, r, "/checkout", http.StatusSeeOther)
}
func (m *Repository) Order(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		log.Println(err)
		return
	}
	username := r.FormValue("username")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	address := r.FormValue("address")
	city := r.FormValue("city")
	state := r.FormValue("state")
	paymentMethod := r.FormValue("paymentmethod")
	if paymentMethod == "noselect" {
		// Thiết lập thông báo lỗi vào session (nếu sử dụng session)
		// Ví dụ sử dụng gorilla/sessions
		// session.Values["failedMsg"] = "please choose your payment method"
		// session.Save(r, w)

		// Chuyển hướng về checkout.jsp
		session, _ := m.App.Session.Get(r, "posty")
		session.Values["Error"] = "Please choose your payment method"
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "checkout.jsp", http.StatusSeeOther)
		return
	}
	fullAdd := address + ", " + city + ", " + state
	list, _, err := m.DB.GetBookByUser(int(userID))
	if err != nil {
		log.Println(err)
		return
	}
	if list == nil {
		session, _ := m.App.Session.Get(r, "posty")
		session.Values["Error"] = "Add Item"
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "checkout", http.StatusSeeOther)
	} else {
		var orderList []models.BookOrder
		for _, c := range list {
			o := models.BookOrder{
				Orderid:       fmt.Sprintf("BOOK-ORDER-00%d", rand.Intn(1000)),
				UserName:      username,
				Email:         email,
				Phone_no:      phone,
				FullAddress:   fullAdd,
				BookName:      c.BookName,
				Author:        c.Author,
				Price:         fmt.Sprintf("%.2f", c.Price),
				PaymentMethod: paymentMethod,
			}
			orderList = append(orderList, o)
		}
		err := m.DB.SaveOrder(orderList)
		if err != nil {
			log.Println(err)
			return
		}
		err = m.DB.DeleteAllBookC(userID)
		if err != nil {
			log.Println(err)
			return
		}
		http.Redirect(w, r, "order-success", http.StatusSeeOther)
	}
}
func (m *Repository) SearchBook(w http.ResponseWriter, r *http.Request) {
	searchString := r.FormValue("search")
	log.Println("asdasdad", searchString)
	var bookSearch []models.BookDtls
	bookSearch, err := m.DB.GetBookSearch(searchString)
	log.Println(bookSearch)
	if err != nil {
		log.Println(err)
		return
	}
	data := make(map[string]interface{})
	data["bookSearch"] = bookSearch
	render.Template(w, r, "search.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
func (m *Repository) AdminAllBooks(w http.ResponseWriter, r *http.Request) {
	var books []models.BookDtls
	books, err := m.DB.GetAllBooks()
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
func (m *Repository) AdminHome(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-home.page.tmpl", &models.TemplateData{})
}
func (m *Repository) PostEditBook(w http.ResponseWriter, r *http.Request) {
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
	err = m.DB.UpdateEditBook(bookName, authorName, bookStatus, float32(price), bookID)
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
func (m *Repository) EditBook(w http.ResponseWriter, r *http.Request) {
	bid, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Println(err)
		return
	}
	book, err := m.DB.GetBookById(bid)
	if err != nil {
		log.Println(err)
	}
	data := make(map[string]interface{})
	data["book"] = book
	render.Template(w, r, "admin-editbook.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
func (m *Repository) BookDelete(w http.ResponseWriter, r *http.Request) {
	bid, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Println(err)
		return
	}
	err = m.DB.DeleteBook(bid)
	if err != nil {
		log.Println(err)
		return
	}
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}
func (m *Repository) AdminOrders(w http.ResponseWriter, r *http.Request) {
	session, _ := m.App.Session.Get(r, "posty")
	_, ok := session.Values["userId"].(int)
	if !ok {
		http.Redirect(w, r, "login", http.StatusSeeOther)
	}
	books, err := m.DB.GetAllOrder()
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
func (m *Repository) AdminAddBook(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-addbook.page.tmpl", &models.TemplateData{})
}
func (m *Repository) PostAdminAddBook(w http.ResponseWriter, r *http.Request) {
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
	err = m.DB.AddBook(book)
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
