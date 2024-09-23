package repository

import "Ebook/internal/models"

type UserRepo interface {
	Register(user models.User) error
	Login(email, testPassword string) (int, string, error)
	CheckUser(email string) bool
	FindUserByID(id int) (models.User, error)
	UpdateProfile(name, email, phone_no string, uid int) error
	CheckPassword(uid int) string
}
type BookRepo interface {
	GetAllBooks() ([]models.BookDtls, error)
	GetOldBooks() ([]models.BookDtls, error)
	GetNewBooks() ([]models.BookDtls, error)
	GetRecentBooks() ([]models.BookDtls, error)
	GetSomeOldBooks() ([]models.BookDtls, error)
	GetSomeNewBooks() ([]models.BookDtls, error)
	GetSomeRecentBooks() ([]models.BookDtls, error)
	AddBook(b models.BookDtls) error
	GetBookById(id int) (models.BookDtls, error)
	UpdateEditBook(bookName, author, status string, price float32, bookID int) error
	DeleteBook(id int) error
	GetBookSearch(search string) ([]models.BookDtls, error)
	GetBooksByOld(email string, category string) ([]models.BookDtls, error)
	OldBookDelete(email string, category string, bid int) error
}
type OrderRepo interface {
	SaveOrder(orderlist []models.BookOrder) error
	GetBookOrder(email string) ([]models.BookOrder, error)
	GetAllOrder() ([]models.BookOrder, error)
}
type CartRepo interface {
	AddCart(c models.Cart) error
	DeleteBookC(bid, uid, cid int) error
	GetBookByUserC(id int) ([]models.Cart, float64, error)
	DeleteAllBookC(uid int) error
}
type DatabaseRepo interface {
	UserRepo() UserRepo   // Trả về repository cho User
	BookRepo() BookRepo   // Trả về repository cho Book
	OrderRepo() OrderRepo // Trả về repository cho Order
	CartRepo() CartRepo   // Trả về repository cho Cart
}
