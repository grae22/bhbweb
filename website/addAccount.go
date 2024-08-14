package website

import (
	"net/http"
	"strconv"
)

func (c Controller) AddAccount(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	responseHtml := pageHeader

	parentIdStr := request.URL.Query().Get("parentId")
	parentId, err := strconv.Atoi(parentIdStr)

	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	parentAccount, ok := c.Book.AccountForId(parentId)

	if !ok {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	responseHtml += "<sup><a href='/home'>Home</a></sup>"

	if request.URL.Query().Has("error") {
		responseHtml += "<span style=''>" + request.URL.Query().Get("error") + "</span>"
	}

	responseHtml += "<p><center><form action='/api/account' method='post'>"
	responseHtml += "<table><caption><b>Add Account</b></caption><tr><td>Parent:</td>"
	responseHtml += "<td><input name='parentId' id='parentId' type='hidden' value='" + parentIdStr + "' />" + parentAccount.Name + "</td></tr>"
	responseHtml += "<tr><td>Name:</td><td><input name='name' id='name' type='text' /></td></tr>"
	responseHtml += "<tr><td></td><td><input type='submit' value='Add' /></td></tr>"
	responseHtml += "</table>"
	responseHtml += "</form></center></p>"

	responseHtml += "<p><center><b>Warning!</b> If the parent account contains transactions they will be moved to the new account.</center></p>"

	responseHtml += pageFooter

	responseWriter.Write([]byte(responseHtml))
}
