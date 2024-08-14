package api

import (
	"bhbweb/domain"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type TransactionController struct {
	book *domain.Book
}

func NewTransactionController(book *domain.Book) TransactionController {
	return TransactionController{
		book: book,
	}
}

func (c TransactionController) AddTransaction(
	response http.ResponseWriter,
	request *http.Request,
) {

	debitAccIdStr := request.FormValue("debitAccId")
	debitAccId, err := strconv.Atoi(debitAccIdStr)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Missing or invalid debit account id."))
		return
	}

	creditAccIdStr := request.FormValue("creditAccId")
	creditAccId, err := strconv.Atoi(creditAccIdStr)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Missing or invalid credit account id."))
		return
	}

	dateStr := request.FormValue("date")
	date, err := time.Parse("2006/1/2", dateStr)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Missing date or invalid format, use yyyy/m/d"))
		return
	}

	valuesStr := strings.Split(request.FormValue("value"), ".")
	if len(valuesStr) > 2 {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Invalid value."))
		return
	}
	integerValue, err := strconv.Atoi(valuesStr[0])
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Missing or invalid value."))
		return
	}
	fractionalValue := 0
	if len(valuesStr) > 1 {
		fractionalValue, err = strconv.Atoi(valuesStr[1])
		if err != nil || fractionalValue > 99 {
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte("Invalid value."))
			return
		}
	}

	description := request.FormValue("description")

	err = c.book.AddTransaction(
		debitAccId,
		creditAccId,
		date,
		domain.Money((integerValue*100)+fractionalValue),
		description)

	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	}

	http.Redirect(
		response,
		request,
		"/home",
		http.StatusSeeOther)
}
