package dbrepo

import (
	"Ebook/internal/models"
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) Register(user models.User) error {
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
	_, err = m.DB.ExecContext(ctx, query,
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
func (m *postgresDBRepo) Login(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var id int
	var hashedPassword string
	row := m.DB.QueryRowContext(ctx, "select uid, password from users where email= $1", email)
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
func (m *postgresDBRepo) GetAllBooks() ([]models.BookDtls, error) {
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
			&book.BookID,
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
func (m *postgresDBRepo) GetNewBooks() ([]models.BookDtls, error) {
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
func (m *postgresDBRepo) GetOldBooks() ([]models.BookDtls, error) {
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
func (m *postgresDBRepo) GetRecentBooks() ([]models.BookDtls, error) {
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
func (m *postgresDBRepo) GetSomeNewBooks() ([]models.BookDtls, error) {
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
func (m *postgresDBRepo) GetSomeOldBooks() ([]models.BookDtls, error) {
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
func (m *postgresDBRepo) GetSomeRecentBooks() ([]models.BookDtls, error) {
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
func (m *postgresDBRepo) AddBook(b models.BookDtls) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `INSERT INTO bookdtls (bookname, author, price, bookCategory, status, photo,email) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := m.DB.ExecContext(ctx, query, b.BookName, b.Author, b.Price, b.BookCategory, b.Status, b.PhotoName, b.Email)
	if err != nil {
		return err
	}

	return nil
}
func (m *postgresDBRepo) GetBookById(id int) (models.BookDtls, error) {
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
func (m *postgresDBRepo) UpdateEditBook(book models.BookDtls) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `update bookdtls set bookname=$1, author=$2, price=$3, status=$4 where bookId=$5`
	_, err := m.DB.ExecContext(ctx, query, book.BookName, book.Author, book.Price, book.Status, book.BookID)
	if err != nil {
		return err
	}
	return nil
}
func (m *postgresDBRepo) DeleteBook(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `delete from bookdtls where bookId=$1`
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
func (m *postgresDBRepo) GetBookSearch(search string) ([]models.BookDtls, error) {
	var list []models.BookDtls

	query := `
		SELECT * 
		FROM bookdtls 
		WHERE (bookname LIKE ? OR author LIKE ? OR bookCategory LIKE ?) AND status = ?
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
func (m *postgresDBRepo) CartAdd(c models.Cart) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `insert into cart (bid, uid, bookName, author, price, total_price) values ($1,$2,$3,$4,$5,$6)`
	_, err := m.DB.ExecContext(ctx, query, c.Bid, c.Uid, c.BookName, c.Author, c.Price, c.Total_price)
	if err != nil {
		return err
	}
	return nil
}
func (m *postgresDBRepo) GetBookByUser(id int) ([]models.BookDtls, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var books []models.BookDtls
	query := `select * from cart where uid=$1`
	rows, err := m.DB.QueryContext(ctx, query, id)
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
	return books, nil
}
func (m *postgresDBRepo) DeleteBookC(bid, uid, cid int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `delete from cart where bid=$1 and uid=$2 and cid=$3 `
	_, err := m.DB.ExecContext(ctx, query, bid, uid, cid)
	if err != nil {
		return err
	}
	return nil
}
func (m *postgresDBRepo) UpdateProfile(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `update user set name=$1,email=$2,phone_no=$3 where uid=$4`
	_, err := m.DB.ExecContext(ctx, query, u.Name, u.Email, u.Phone_no, u.ID)
	if err != nil {
		return err
	}
	return nil
}
func (m *postgresDBRepo) CheckUser(email string) bool {
	f := true
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `select * from user where email=$1`
	rows, _ := m.DB.QueryContext(ctx, query, email)
	for rows.Next() {
		f = true
	}
	return f
}
func (m *postgresDBRepo) FindUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var user models.User
	query := `select * from user where id=$1`
	row := m.DB.QueryRowContext(ctx, query, id)
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
