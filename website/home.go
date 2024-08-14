package website

import (
	"bhbweb/domain"
	"fmt"
	"net/http"
	"time"
)

func (c Controller) Home(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	period := request.URL.Query().Get("period")

	if period == "" {
		period = time.Now().Format("0601")
	}

	showForDate, _ := time.Parse("0601", period)

	periodData, err := c.Book.BalancesForDate(showForDate)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte(err.Error()))
		return
	}

	previousPeriodId := showForDate.AddDate(0, -1, 0).Format("0601")
	nextPeriodId := showForDate.AddDate(0, 1, 0).Format("0601")

	responseHtml := pageHeader +
		"<span><center>" +
		fmt.Sprintf("<a href='/home?period=%s'>%s</a> ", previousPeriodId, showForDate.AddDate(0, -1, 0).Format("Jan 06")) +
		"<b>" + showForDate.Format("Jan 06") + "</b>" +
		fmt.Sprintf(" <a href='/home?period=%s'>%s</a>", nextPeriodId, showForDate.AddDate(0, 1, 0).Format("Jan 06")) +
		"</center></span>" +
		"<p><center><table>"

	for _, a := range c.Book.GetAccounts() {
		if a.ParentId != domain.NoParentId {
			continue
		}
		responseHtml += buildAccountTreeNode(a, *c.Book, periodData, 0)
	}

	responseHtml += "</table><center></p>"
	responseHtml += pageFooter

	responseWriter.Write([]byte(responseHtml))
}

func buildAccountTreeNode(
	a domain.Account,
	book domain.Book,
	periodData domain.PeriodBalances,
	level int,
) (resultHtml string) {

	var nameIndent string

	if level > 0 {
		for i := 0; i < level; i++ {
			nameIndent += "&nbsp;&nbsp;&nbsp;&nbsp;"
		}
	}

	balances, ok := periodData.BalanceByAccountId[a.Id]

	if !ok {
		balances = domain.AccountBalances{}
	}

	closingBalanceAbs := balances.Closing
	if closingBalanceAbs < 0 {
		closingBalanceAbs *= -1
	}

	balanceStr := fmt.Sprintf("%d.%02d", closingBalanceAbs/100, closingBalanceAbs%100)

	if balances.Closing < 0 {
		balanceStr = "(" + balanceStr + ")"
	}

	resultHtml += fmt.Sprintf(
		"<tr><td><sup><a href='/addAccount?parentId=%d'>[+]</a>"+
			" <a href='/addTransaction?debitAccId=%d'>[d]</a>"+
			" <a href='/addTransaction?creditAccId=%d'>[c]</a>"+
			"</sup></td><td style='padding-left:10px;padding-right:10px;'>%s<a href='/viewAccount?accId=%d&period=%s'>%s</a></td><td style='text-align:right;'>%s</td></tr>",
		a.Id,
		a.Id,
		a.Id,
		nameIndent,
		a.Id,
		periodData.PeriodId,
		a.Name,
		balanceStr)

	for _, childId := range a.ChildrenId {
		childAccount, ok := book.AccountForId(childId)

		if !ok {
			continue
		}

		resultHtml += buildAccountTreeNode(childAccount, book, periodData, level+1)
	}

	return resultHtml
}
