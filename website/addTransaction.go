package website

import (
	"bhbweb/domain"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"
)

func (c Controller) AddTransaction(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	responseHtml := pageHeader

	debitAccIdStr := request.URL.Query().Get("debitAccId")
	debitAccId, err := strconv.Atoi(debitAccIdStr)

	if err != nil {
		debitAccId = -1
	}

	creditAccIdStr := request.URL.Query().Get("creditAccId")
	creditAccId, err := strconv.Atoi(creditAccIdStr)

	if err != nil {
		creditAccId = -1
	}

	responseHtml += "<sup><a href='/home'>Home</a></sup>"

	if request.URL.Query().Has("error") {
		responseHtml += "<span style=''>" + request.URL.Query().Get("error") + "</span>"
	}

	dateNow := time.Now().Format("2006/1/2")

	responseHtml +=
		"<form action='/api/transaction' method='post'>" +
			"<table><caption><b>New Transaction</b></caption>" +
			"<tr><td>Debit:</td><td>" + buildDropdown("debitAccId", c.Book.GetAccounts(), debitAccId) + "</td></tr>" +
			"<tr><td>Credit:</td><td>" + buildDropdown("creditAccId", c.Book.GetAccounts(), creditAccId) + "</td></tr>" +
			"<tr><td>Date:</td><td><input name='date' id='date' type='text' value='" + dateNow + "' /> <span style='color:#a0a0a0'><i>yyyy/m/d</i></span></td></tr>" +
			"<tr><td>Value:</td><td><input name='value' id='value' type='text' /></td></tr>" +
			"<tr><td>Description:</td><td><input name='description' id='description' type='text' /></td></tr>" +
			"<tr><td></td><td><input type='submit' value='Create' /></td></tr>" +
			"</table>" +
			"</form>" +
			pageFooter

	responseWriter.Write([]byte(responseHtml))
}

func buildDropdown(
	name string,
	accounts []domain.Account,
	defaultId int,
) string {

	var qualifiedNames []string = make([]string, 0, len(accounts))
	var accByQn map[string]*domain.Account = make(map[string]*domain.Account, len(accounts))

	for _, a := range accounts {
		qn := buildQualifiedAccountName(a.Id, accounts)
		qualifiedNames = append(qualifiedNames, qn)
		accByQn[qn] = &a
	}

	sort.Slice(qualifiedNames, func(i, j int) bool {
		return qualifiedNames[i] < qualifiedNames[j]
	})

	s := "<select name='" + name + "' id='" + name + "'>"

	for _, qn := range qualifiedNames {
		a := accByQn[qn]

		var selected string
		if a.Id == defaultId {
			selected = "selected"
		}

		qualifiedName := buildQualifiedAccountName(a.Id, accounts)

		s += fmt.Sprintf("<option value='%d' %s>%s</option>", a.Id, selected, qualifiedName)
	}

	s += "</select>"

	return s
}

func buildQualifiedAccountName(
	accId int,
	accounts []domain.Account,
) (qualifiedName string) {

	var acc *domain.Account
	var accountsById map[int]*domain.Account = make(map[int]*domain.Account, len(accounts))

	for _, a := range accounts {
		if a.Id == accId {
			acc = &a
		}

		accountsById[a.Id] = &a
	}

	if acc == nil {
		return "error"
	}

	qualifiedName = acc.Name

	for parentId := acc.ParentId; parentId >= 0; {
		parentAcc, ok := accountsById[parentId]

		if !ok {
			return "error"
		}

		qualifiedName = parentAcc.Name + " > " + qualifiedName

		parentId = parentAcc.ParentId
	}

	return qualifiedName
}
