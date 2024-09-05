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
	UpdateEditBook(book models.BookDtls) error
	DeleteBook(id int) error
	CartAdd(c models.Cart) error
	GetBookByUser(id int) ([]models.BookDtls, error)
	GetBookSearch(search string) ([]models.BookDtls, error)
	DeleteBookC(bid, uid, cid int) error
	CheckUser(email string) bool
}
