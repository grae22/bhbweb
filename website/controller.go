package website

import (
	"bhbweb/domain"
)

type Controller struct {
	Book *domain.Book
}

func NewController(book *domain.Book) Controller {
	return Controller{
		Book: book,
	}
}
