package userhandler

import (
	"Ebook/internal/config"
	"Ebook/internal/driver"
	"Ebook/internal/forms"
	"Ebook/internal/helpers"
	"Ebook/internal/models"
	"Ebook/internal/render"
	"Ebook/internal/repository"
	dbrepo "Ebook/internal/repository/dprepo"
	"context"
	"database/sql"
	"encoding/json"
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

var RepoUser *UserHandler

type UserHandler struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

func NewUserRepository(a *config.AppConfig, db *driver.DB) *UserHandler {
	return &UserHandler{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}
func NewUserHandlers(r *UserHandler) {
	RepoUser = r
}
func (m *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}
func (m *UserHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {

	url := config.AppConfigOauth.GoogleLoginConfig.AuthCodeURL("randomstate")

	http.Redirect(w, r, url, http.StatusSeeOther)
}
func (m *UserHandler) GoogleCallBack(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state != "randomstate" {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		log.Println("States don't Match!!")
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Authorization code not found", http.StatusBadRequest)
		return
	}

	googlecon := config.GoogleConfig()
	token, err := googlecon.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
		log.Println("Code-Token Exchange Failed")
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		http.Error(w, "Failed to fetch user data", http.StatusInternalServerError)
		log.Println("User Data Fetch Failed")
		return
	}
	defer resp.Body.Close()

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Không thể phân tích dữ liệu người dùng", http.StatusInternalServerError)
		log.Println("Phân tích JSON thất bại")
		return
	}

	// Phân tích dữ liệu JSON
	var userInfo map[string]interface{}
	if err := json.Unmarshal(userData, &userInfo); err != nil {
		http.Error(w, "Không thể phân tích dữ liệu JSON", http.StatusInternalServerError)
		log.Println("Phân tích JSON thất bại:", err)
		return
	}

	// Lấy thông tin cần thiết
	email, _ := userInfo["email"].(string)
	name, _ := userInfo["name"].(string)

	// Kiểm tra xem người dùng đã tồn tại chưa
	user, err := m.DB.UserRepo().FindUserByEmail(email)
	if err != nil {
		// Nếu người dùng chưa tồn tại, tạo mới
		if err == sql.ErrNoRows {
			// Nếu người dùng chưa tồn tại, tạo mới
			newUser := models.User{
				Email: email,
				Name:  name,
			}
			err = m.DB.UserRepo().Register(newUser)
			if err != nil {
				http.Error(w, "Không thể tạo người dùng mới", http.StatusInternalServerError)
				log.Println("Tạo người dùng thất bại:", err)
				return
			}
			user = newUser
		} else {
			http.Error(w, "Lỗi khi tìm kiếm người dùng", http.StatusInternalServerError)
			log.Println("Tìm kiếm người dùng thất bại:", err)
			return
		}
	}
	session, _ := m.App.Session.Get(r, "posty")
	session.Values["userId"] = user.ID
	log.Println(session.Values["userId"])
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Không thể lưu phiên", http.StatusInternalServerError)
		log.Println("Lưu phiên thất bại:", err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (m *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "register.page.tmpl", &models.TemplateData{})
}
func (m *UserHandler) PostRegister(w http.ResponseWriter, r *http.Request) {
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
	err = m.DB.UserRepo().Register(user)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
func (m *UserHandler) Home(w http.ResponseWriter, r *http.Request) {
	var newbooks []models.BookDtls

	newbooks, err := m.DB.BookRepo().GetSomeNewBooks()
	if err != nil {
		helpers.ServerError(w, err)
	}
	data := make(map[string]interface{})
	data["newbooks"] = newbooks
	var recentbooks []models.BookDtls
	recentbooks, err = m.DB.BookRepo().GetSomeRecentBooks()
	if err != nil {
		helpers.ServerError(w, err)
	}
	data["recentbooks"] = recentbooks
	var oldbooks []models.BookDtls
	oldbooks, err = m.DB.BookRepo().GetSomeOldBooks()
	if err != nil {
		helpers.ServerError(w, err)
	}
	data["oldbooks"] = oldbooks
	fmt.Println("render home")

	render.Template(w, r, "index.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
func (m *UserHandler) AllNewBooks(w http.ResponseWriter, r *http.Request) {
	var newbooks []models.BookDtls
	newbooks, err := m.DB.BookRepo().GetNewBooks()
	if err != nil {
		helpers.ServerError(w, err)
	}
	data := make(map[string]interface{})
	data["newbooks"] = newbooks
	render.Template(w, r, "all_new_books.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
func (m *UserHandler) AllOldBooks(w http.ResponseWriter, r *http.Request) {
	var oldbooks []models.BookDtls
	oldbooks, err := m.DB.BookRepo().GetOldBooks()
	if err != nil {
		helpers.ServerError(w, err)
	}
	data := make(map[string]interface{})
	data["oldbooks"] = oldbooks
	render.Template(w, r, "all_old_books.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
func (m *UserHandler) AllRecentBooks(w http.ResponseWriter, r *http.Request) {
	var recentbooks []models.BookDtls
	recentbooks, err := m.DB.BookRepo().GetRecentBooks()
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
func (m *UserHandler) BookDetail(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
	}
	var book models.BookDtls
	book, err = m.DB.BookRepo().GetBookById(ID)
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
func (m *UserHandler) AddCart(w http.ResponseWriter, r *http.Request) {
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
	book, err := m.DB.BookRepo().GetBookById(bid)
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
	err = m.DB.CartRepo().AddCart(cart)
	if err != nil {
		log.Println(err)
		return
	}
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}
func (m *UserHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	session, _ := m.App.Session.Get(r, "posty")
	id, ok := session.Values["userId"].(int)
	if !ok {
		http.Redirect(w, r, "login", http.StatusSeeOther)
	}
	var cart []models.Cart
	cart, totalPrice, err := m.DB.CartRepo().GetBookByUserC(id)
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
func (m *UserHandler) EditProfile(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "edit-profile.page.tmpl", &models.TemplateData{})
}
func (m *UserHandler) PostEditProfile(w http.ResponseWriter, r *http.Request) {
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
	err = bcrypt.CompareHashAndPassword([]byte(m.DB.UserRepo().CheckPassword(userID)), []byte(password))
	if err != nil {
		warning := "Wrong Password"
		render.Template(w, r, "edit-profile.page.tmpl", &models.TemplateData{
			Form:    nil,
			Warning: warning,
		})
	} else {
		err = m.DB.UserRepo().UpdateProfile(fullName, email, phone, userID)
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
func (m *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
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
func (m *UserHandler) PostLogin(w http.ResponseWriter, r *http.Request) {
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
		id, _, err := m.DB.UserRepo().Login(email, password)
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
func (m *UserHandler) OldBooks(w http.ResponseWriter, r *http.Request) {
	session, _ := m.App.Session.Get(r, "posty")
	id, ok := session.Values["userId"].(int)
	if !ok {
		http.Redirect(w, r, "login", http.StatusSeeOther)
	}
	var user models.User
	user, err := m.DB.UserRepo().FindUserByID(id)
	if err != nil {
		log.Println(err)
		return
	}
	var oldbook []models.BookDtls
	oldbook, err = m.DB.BookRepo().GetBooksByOld(user.Email, "Old Book")
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
func (m *UserHandler) OrderSuccess(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "order_success.page.tmpl", &models.TemplateData{})
}
func (m *UserHandler) GetOrderByUser(w http.ResponseWriter, r *http.Request) {
	session, _ := m.App.Session.Get(r, "posty")
	id, ok := session.Values["userId"].(int)
	if !ok {
		http.Redirect(w, r, "login", http.StatusSeeOther)
	}
	user, err := m.DB.UserRepo().FindUserByID(id)
	if err != nil {
		log.Println(err)
		return
	}
	orderbooks, err := m.DB.OrderRepo().GetBookOrder(user.Email)
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
func (m *UserHandler) Search(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search.page.tmpl", &models.TemplateData{})
}
func (m *UserHandler) UserAddress(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "user_address.page.tmpl", &models.TemplateData{})
}
func (m *UserHandler) Setting(w http.ResponseWriter, r *http.Request) {
	session, _ := m.App.Session.Get(r, "posty")
	_, ok := session.Values["userId"].(int)
	if !ok {
		http.Redirect(w, r, "login", http.StatusSeeOther)
	}
	render.Template(w, r, "setting.page.tmpl", &models.TemplateData{
		Data: nil,
	})
}
func (m *UserHandler) Helpline(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "helpline.page.tmpl", &models.TemplateData{})
}
func (m *UserHandler) SellBook(w http.ResponseWriter, r *http.Request) {
	session, _ := m.App.Session.Get(r, "posty")
	_, ok := session.Values["userId"].(int)
	if !ok {
		http.Redirect(w, r, "login", http.StatusSeeOther)
	}
	render.Template(w, r, "sell_book.page.tmpl", &models.TemplateData{})
}
func (m *UserHandler) PostSellBook(w http.ResponseWriter, r *http.Request) {
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

	err = m.DB.BookRepo().AddBook(oldbook)
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
func (m *UserHandler) DeleteOldBook(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	log.Println(email)
	bid, err := strconv.Atoi(r.URL.Query().Get("bid"))
	if err != nil {
		log.Println(err)
		return
	}
	err = m.DB.BookRepo().OldBookDelete(email, "Old Book", bid)
	if err != nil {
		log.Println(err)
		return
	}
	http.Redirect(w, r, "/old-books", http.StatusSeeOther)
}
func (m *UserHandler) RemoveBook(w http.ResponseWriter, r *http.Request) {
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
	err = m.DB.CartRepo().DeleteBookC(bid, uid, cid)
	if err != nil {
		log.Println(err)
		return
	}
	http.Redirect(w, r, "/checkout", http.StatusSeeOther)
}
func (m *UserHandler) Order(w http.ResponseWriter, r *http.Request) {
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
	list, _, err := m.DB.CartRepo().GetBookByUserC(int(userID))
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
		err := m.DB.OrderRepo().SaveOrder(orderList)
		if err != nil {
			log.Println(err)
			return
		}
		err = m.DB.CartRepo().DeleteAllBookC(userID)
		if err != nil {
			log.Println(err)
			return
		}
		http.Redirect(w, r, "order-success", http.StatusSeeOther)
	}
}
func (m *UserHandler) SearchBook(w http.ResponseWriter, r *http.Request) {
	searchString := r.FormValue("search")
	log.Println("asdasdad", searchString)
	var bookSearch []models.BookDtls
	bookSearch, err := m.DB.BookRepo().GetBookSearch(searchString)
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
