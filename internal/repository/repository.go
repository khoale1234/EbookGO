package repository

import "Ebook/internal/models"

type DatabaseRepo interface {
	Register(user models.User) error
	Login(email, testPassword string) (int, string, error)
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
	GetBookByUser(id int) ([]models.Cart, float64, error)
	GetBookSearch(search string) ([]models.BookDtls, error)
	DeleteBookC(bid, uid, cid int) error
	CheckUser(email string) bool
	FindUserByID(id int) (models.User, error)
	GetBooksByOld(email string, category string) ([]models.BookDtls, error)
	OldBookDelete(email string, category string, bid int) error
	GetBookOrder(email string) ([]models.BookOrder, error)
	AddCart(c models.Cart) error
	SaveOrder(orderlist []models.BookOrder) error
	DeleteAllBookC(uid int) error
	CheckPassword(uid int) string
	UpdateProfile(name, email, phone_no string, uid int) error
	GetAllOrder() ([]models.BookOrder, error)
}
