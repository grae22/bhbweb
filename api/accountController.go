package api

import (
	"bhbweb/domain"
	"net/http"
	"strconv"
	"strings"
)

type AccountController struct {
	book *domain.Book
}

func NewAccountController(book *domain.Book) AccountController {
	return AccountController{
		book: book,
	}
}

func (ac AccountController) AddAccount(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	parentIdStr := request.FormValue("parentId")
	name := strings.Trim(request.FormValue("name"), "")

	if len(parentIdStr) == 0 {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte("Parent Id is required"))
		return
	}

	parentId, err := strconv.Atoi(parentIdStr)

	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte("Parent Id is an invalid id"))
		return
	}

	if len(name) == 0 {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte("Name is required"))
		return
	}

	ac.book.AddAccount(parentId, name)

	http.Redirect(
		responseWriter,
		request,
		"/home",
		http.StatusSeeOther)
}
