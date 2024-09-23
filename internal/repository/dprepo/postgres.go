package dbrepo

import (
	"Ebook/internal/models"
	"context"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (r *UserRepo) Register(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return err
	}
	query := `
	INSERT INTO users (name, email, phone_no, password)
	VALUES ($1, $2, $3, $4)
`
	_, err = r.DB.ExecContext(ctx, query,
		user.Name,
		user.Email,
		user.Phone_no,
		hashedPassword,
	)
	if err != nil {
		return err
	}

	return nil
}
func (r *UserRepo) Login(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var id int
	var hashedPassword string
	row := r.DB.QueryRowContext(ctx, "select uid, password from users where email= $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, " ", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("Incorrect Password")
	} else if err != nil {
		return 0, "", err
	}
	return id, hashedPassword, nil
}
func (r *UserRepo) CheckUser(email string) bool {
	f := true
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `select * from users where email=$1`
	rows, _ := r.DB.QueryContext(ctx, query, email)
	for rows.Next() {
		f = true
	}
	return f
}
func (r *UserRepo) FindUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var user models.User
	query := `select * from users where uid=$1`
	row := r.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone_no,
		&user.Password,
	)
	if err != nil {
		return user, err
	}
	return user, nil
}
func (r *UserRepo) UpdateProfile(name, email, phone_no string, uid int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `update users set name= $1,email= $2,phone_no= $3 where uid=$4`
	_, err := r.DB.ExecContext(ctx, query, name, email, phone_no, uid)
	if err != nil {
		return err
	}
	return nil
}
func (r *UserRepo) CheckPassword(uid int) string {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `select password from users where uid=$1`
	row := r.DB.QueryRowContext(ctx, query, uid)
	var password string
	err := row.Scan(&password)
	if err != nil {
		log.Println(err)
		return ""
	}
	return password
}
func (m *BookRepo) GetAllBooks() ([]models.BookDtls, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var books []models.BookDtls
	query := `select * from bookdtsl`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return books, err
	}
	for rows.Next() {
		var book models.BookDtls
		err := rows.Scan(
			&book.BookID,
			&book.BookName,
			&book.Author,
			&book.Price,
			&book.BookCategory,
			&book.Status,
			&book.PhotoName,
			&book.Email,
		)
		if err != nil {
			return books, err
		}
		books = append(books, book)
	}
	if err = rows.Err(); err != nil {
		return books, nil
	}
	return books, nil
}
func (m *BookRepo) GetNewBooks() ([]models.BookDtls, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var books []models.BookDtls
	query := `Select * from bookdtsl where bookCategory=$1 and status=$2 order by bookId DESC`
	rows, err := m.DB.QueryContext(ctx, query, "New Book", "Active")
	if err != nil {
		return books, err
	}
	for rows.Next() {
		var book models.BookDtls
		err := rows.Scan(
			&book.BookID,
			&book.BookName,
			&book.Author,
			&book.Price,
			&book.BookCategory,
			&book.Status,
			&book.PhotoName,
			&book.Email,
		)
		if err != nil {
			return books, err
		}
		books = append(books, book)
	}
	if err = rows.Err(); err != nil {
		return books, nil
	}
	return books, nil
}
func (m *BookRepo) GetOldBooks() ([]models.BookDtls, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var books []models.BookDtls
	query := `Select * from bookdtsl where bookCategory=$1 and status=$2 order by bookId DESC`
	rows, err := m.DB.QueryContext(ctx, query, "Old Book", "Active")
	if err != nil {
		return books, err
	}
	for rows.Next() {
		var book models.BookDtls
		err := rows.Scan(
			&book.BookID,
			&book.BookName,
			&book.Author,
			&book.Price,
			&book.BookCategory,
			&book.Status,
			&book.PhotoName,
			&book.Email,
		)
		if err != nil {
			return books, err
		}
		books = append(books, book)
	}
	if err = rows.Err(); err != nil {
		return books, nil
	}
	return books, nil
}
func (m *BookRepo) GetRecentBooks() ([]models.BookDtls, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var books []models.BookDtls
	query := `Select * from bookdtsl where status=$1 order by bookId DESC`
	rows, err := m.DB.QueryContext(ctx, query, "Active")
	if err != nil {
		return books, err
	}
	for rows.Next() {
		var book models.BookDtls
		err := rows.Scan(
			&book.BookID,
			&book.BookName,
			&book.Author,
			&book.Price,
			&book.BookCategory,
			&book.Status,
			&book.PhotoName,
			&book.Email,
		)
		if err != nil {
			return books, err
		}
		books = append(books, book)
	}
	if err = rows.Err(); err != nil {
		return books, nil
	}
	return books, nil
}
func (m *BookRepo) GetSomeNewBooks() ([]models.BookDtls, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var books []models.BookDtls
	query := `SELECT * FROM bookdtsl WHERE bookCategory = $1 AND status = $2 ORDER BY bookId DESC LIMIT 4`
	rows, err := m.DB.QueryContext(ctx, query, "New Book", "Active")
	if err != nil {
		return books, err
	}
	for rows.Next() {
		var book models.BookDtls
		err := rows.Scan(
			&book.BookID,
			&book.BookName,
			&book.Author,
			&book.Price,
			&book.BookCategory,
			&book.Status,
			&book.PhotoName,
			&book.Email,
		)
		if err != nil {
			return books, err
		}
		books = append(books, book)
	}
	if err = rows.Err(); err != nil {
		return books, nil
	}
	return books, nil
}
func (m *BookRepo) GetSomeOldBooks() ([]models.BookDtls, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var books []models.BookDtls
	query := `Select * from bookdtsl where bookCategory = $1 and status = $2 order by bookId DESC LIMIT 4`
	rows, err := m.DB.QueryContext(ctx, query, "Old Book", "Active")
	if err != nil {
		return books, err
	}
	for rows.Next() {
		var book models.BookDtls
		err := rows.Scan(
			&book.BookID,
			&book.BookName,
			&book.Author,
			&book.Price,
			&book.BookCategory,
			&book.Status,
			&book.PhotoName,
			&book.Email,
		)
		if err != nil {
			return books, err
		}
		books = append(books, book)
	}
	if err = rows.Err(); err != nil {
		return books, nil
	}
	return books, nil
}
func (m *BookRepo) GetSomeRecentBooks() ([]models.BookDtls, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var books []models.BookDtls
	query := `Select * from bookdtsl where status=$1 order by bookId DESC LIMIT 4`
	rows, err := m.DB.QueryContext(ctx, query, "Active")
	if err != nil {
		return books, err
	}
	for rows.Next() {
		var book models.BookDtls
		err := rows.Scan(
			&book.BookID,
			&book.BookName,
			&book.Author,
			&book.Price,
			&book.BookCategory,
			&book.Status,
			&book.PhotoName,
			&book.Email,
		)
		if err != nil {
			return books, err
		}
		books = append(books, book)
	}
	if err = rows.Err(); err != nil {
		return books, nil
	}
	return books, nil
}
func (m *BookRepo) AddBook(b models.BookDtls) error {
	log.Println("hello from repo")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `INSERT INTO bookdtsl (bookname, author, price, bookCategory, status, photo,email) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := m.DB.ExecContext(ctx, query, b.BookName, b.Author, b.Price, b.BookCategory, b.Status, b.PhotoName, b.Email)
	if err != nil {
		return err
	}

	return nil
}
func (m *BookRepo) GetBookById(id int) (models.BookDtls, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `select * from bookdtsl where bookId=$1`
	row := m.DB.QueryRowContext(ctx, query, id)
	var book models.BookDtls
	err := row.Scan(
		&book.BookID,
		&book.BookName,
		&book.Author,
		&book.Price,
		&book.BookCategory,
		&book.Status,
		&book.PhotoName,
		&book.Email,
	)
	if err != nil {
		return book, err
	}
	return book, nil
}
func (m *BookRepo) UpdateEditBook(bookName, author, status string, price float32, bookID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `update bookdtsl set bookname=$1, author=$2, price=$3, status=$4 where bookId=$5`
	_, err := m.DB.ExecContext(ctx, query, bookName, author, price, status, bookID)
	if err != nil {
		return err
	}
	return nil
}
func (m *BookRepo) DeleteBook(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `delete from bookdtsl where bookId=$1`
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
func (m *BookRepo) GetBookSearch(search string) ([]models.BookDtls, error) {
	var list []models.BookDtls

	query := `
	SELECT * 
	FROM bookdtsl 
	WHERE (bookname LIKE $1 OR author LIKE $2 OR bookCategory LIKE $3) AND status = $4
`

	rows, err := m.DB.Query(query, "%"+search+"%", "%"+search+"%", "%"+search+"%", "Active")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var b models.BookDtls
		if err := rows.Scan(&b.BookID, &b.BookName, &b.Author, &b.Price, &b.BookCategory, &b.Status, &b.PhotoName, &b.Email); err != nil {
			return nil, err
		}
		list = append(list, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}
func (m *CartRepo) GetBookByUserC(id int) ([]models.Cart, float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var carts []models.Cart
	var totalPrice float64

	query := `SELECT * FROM cart WHERE uid=$1`
	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return carts, totalPrice, err
	}
	defer rows.Close() // Đảm bảo đóng rows sau khi sử dụng

	for rows.Next() {
		var cart models.Cart
		err := rows.Scan(
			&cart.Cid,
			&cart.Bid,
			&cart.Uid,
			&cart.BookName,
			&cart.Author,
			&cart.Price,
		)
		if err != nil {
			return carts, totalPrice, err
		}

		// Cộng dồn giá trị total price
		totalPrice += cart.Price
		carts = append(carts, cart)
	}

	return carts, totalPrice, nil
}
func (m *CartRepo) DeleteBookC(bid, uid, cid int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `delete from cart where bid=$1 and uid=$2 and cid=$3 `
	_, err := m.DB.ExecContext(ctx, query, bid, uid, cid)
	if err != nil {
		return err
	}
	return nil
}

func (m *BookRepo) GetBooksByOld(email string, category string) ([]models.BookDtls, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var books []models.BookDtls
	query := `Select * from bookdtsl where email=$1 and bookCategory=$2`
	rows, err := m.DB.QueryContext(ctx, query, email, category)
	if err != nil {
		return books, err
	}
	for rows.Next() {
		var book models.BookDtls
		err := rows.Scan(
			&book.BookID,
			&book.BookName,
			&book.Author,
			&book.Price,
			&book.BookCategory,
			&book.Status,
			&book.PhotoName,
			&book.Email,
		)
		if err != nil {
			return books, err
		}
		books = append(books, book)
	}
	if err = rows.Err(); err != nil {
		return books, nil
	}
	return books, nil
}
func (m *BookRepo) OldBookDelete(email string, category string, bid int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `delete from bookdtsl where bookCategory=$1 and email=$2 and bookId=$3 `
	_, err := m.DB.ExecContext(ctx, query, category, email, bid)
	if err != nil {
		return err
	}
	return nil
}
func (m *OrderRepo) GetBookOrder(email string) ([]models.BookOrder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `select * from orders where email =$1`
	var booksOrder []models.BookOrder
	rows, err := m.DB.QueryContext(ctx, query, email)
	if err != nil {
		return booksOrder, err
	}
	for rows.Next() {
		var book models.BookOrder
		err = rows.Scan(
			&book.Orderid,
			&book.UserName,
			&book.Email,
			&book.FullAddress,
			&book.Phone_no,
			&book.BookName,
			&book.Author,
			&book.Price,
			&book.PaymentMethod,
		)
		if err != nil {
			return booksOrder, err
		}
		booksOrder = append(booksOrder, book)
	}
	return booksOrder, err
}
func (m *CartRepo) AddCart(c models.Cart) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `insert into cart(bid,uid,bookName,author,price) values ($1,$2,$3,$4,$5)`
	_, err := m.DB.ExecContext(ctx, query, c.Bid, c.Uid, c.BookName, c.Author, c.Price)
	if err != nil {
		log.Println(err)
	}
	return nil
}
func (m *OrderRepo) SaveOrder(orderlist []models.BookOrder) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `insert into orders(orderid,user_name,email,address,phone,book_name,author,price,payment)
	values($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`
	for _, c := range orderlist {
		_, err := m.DB.ExecContext(ctx, query,
			c.Orderid,
			c.UserName,
			c.Email,
			c.FullAddress,
			c.Phone_no,
			c.BookName,
			c.Author,
			c.Price,
			c.PaymentMethod)
		if err != nil {
			return err
		}
	}
	return nil
}
func (m *CartRepo) DeleteAllBookC(uid int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `delete from cart where uid=$1`
	_, err := m.DB.ExecContext(ctx, query, uid)
	if err != nil {
		return err
	}
	return nil
}

func (m *OrderRepo) GetAllOrder() ([]models.BookOrder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `select * from orders`
	var booksOrder []models.BookOrder
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return booksOrder, err
	}
	for rows.Next() {
		var book models.BookOrder
		err = rows.Scan(
			&book.Orderid,
			&book.UserName,
			&book.Email,
			&book.FullAddress,
			&book.Phone_no,
			&book.BookName,
			&book.Author,
			&book.Price,
			&book.PaymentMethod,
		)
		if err != nil {
			return booksOrder, err
		}
		booksOrder = append(booksOrder, book)
	}
	return booksOrder, err
}
