package dbrepo

import (
	"Ebook/internal/config"
	"Ebook/internal/repository"
	"database/sql"
)

type PostgresDBRepo struct {
	App       *config.AppConfig
	DB        *sql.DB
	userRepo  repository.UserRepo
	bookRepo  repository.BookRepo
	cartRepo  repository.CartRepo
	orderRepo repository.OrderRepo
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &PostgresDBRepo{
		App:       a,
		DB:        conn,
		userRepo:  NewUserRepo(conn),
		bookRepo:  NewBookRepo(conn),
		cartRepo:  NewCartRepo(conn),
		orderRepo: NewOrderRepo(conn),
	}
}

func (p *PostgresDBRepo) BookRepo() repository.BookRepo {
	return p.bookRepo
}

func (p *PostgresDBRepo) CartRepo() repository.CartRepo {
	return p.cartRepo
}

func (p *PostgresDBRepo) OrderRepo() repository.OrderRepo {
	return p.orderRepo
}

func (p *PostgresDBRepo) UserRepo() repository.UserRepo {
	return p.userRepo
}

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) repository.UserRepo {
	return &UserRepo{DB: db}
}

type BookRepo struct {
	DB *sql.DB
}

func NewBookRepo(db *sql.DB) repository.BookRepo {
	return &BookRepo{DB: db}
}

type CartRepo struct {
	DB *sql.DB
}

func NewCartRepo(db *sql.DB) repository.CartRepo {
	return &CartRepo{DB: db}
}

type OrderRepo struct {
	DB *sql.DB
}

func NewOrderRepo(db *sql.DB) repository.OrderRepo {
	return &OrderRepo{DB: db}
}
