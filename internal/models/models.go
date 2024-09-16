package models

import "html/template"

type Cart struct {
	Cid         int
	Bid         int
	Uid         int
	BookName    string
	Author      string
	Price       float64
	Total_price float64
}
type BookOrder struct {
	ID            int
	Orderid       string
	UserName      string
	Email         string
	Phone_no      string
	FullAddress   string
	PaymentMethod string
	BookName      string
	Author        string
	Price         string
}
type BookDtls struct {
	BookID       int
	BookName     string
	Author       string
	Price        string
	BookCategory string
	Status       string
	PhotoName    string
	Email        string
}
type User struct {
	ID       int
	Name     string
	Email    string
	Password string
	Phone_no string
	Address  string
	City     string
}
type MailData struct {
	To       string
	From     string
	Subject  string
	Content  template.HTML
	Template string
}
