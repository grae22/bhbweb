package website

import (
	"bhbweb/domain"
	"fmt"
	"net/http"
	"sort"
	"strconv"
)

func (c Controller) ViewAccount(
	response http.ResponseWriter,
	request *http.Request,
) {
	accId, err := strconv.Atoi(request.URL.Query().Get("accId"))
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	periodId := request.URL.Query().Get("period")
	if periodId == "" {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	responseHtml := pageHeader
	responseHtml += "<span><sup><a href='/home'>Home</a></sup></span>"

	periodData, err := c.Book.BalancesForDate(domain.DateForPeriodId(periodId))
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}

	if accData, ok := periodData.BalanceByAccountId[accId]; ok {
		responseHtml += fmt.Sprintf("<p><center>Opening: %s</center></p>", formatMoney(accData.Opening))
		responseHtml += "<p><center><table>" +
			"<th>Date</th>" +
			"<th>Value</th>" +
			"<th>Account</th>" +
			"<th>Contra</th>" +
			"<th>Description</th>"

		transactions, err := c.Book.AccountTransactionsForPeriod(accId, periodId)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(err.Error()))
			return
		}

		sort.Slice(
			transactions,
			func(i, j int) bool {
				return transactions[i].Date.After(transactions[j].Date)
			})

		for _, t := range transactions {
			value := t.Value

			if t.DebitAccountId == accId {
				value = -value
			}

			drAcc, ok := c.Book.AccountForId(t.DebitAccountId)
			if !ok {
				drAcc = domain.Account{
					Name: "(UNKNOWN)",
				}
			}

			crAcc, ok := c.Book.AccountForId(t.CreditAccountId)
			if !ok {
				crAcc = domain.Account{
					Name: "(UNKNOWN)",
				}
			}

			var acc domain.Account
			var contraAcc domain.Account

			if t.DebitAccountId == accId {
				acc = drAcc
				contraAcc = crAcc
			} else {
				acc = crAcc
				contraAcc = drAcc
			}

			responseHtml += fmt.Sprintf(
				"<tr><td style='padding-left:10px;padding-right:10px;'>%s</td>"+
					"<td style='padding-left:10px;padding-right:10px;'>%s</td>"+
					"<td style='padding-left:10px;padding-right:10px;'>%s</td>"+
					"<td style='padding-left:10px;padding-right:10px;'>%s</td>"+
					"<td style='padding-left:10px;padding-right:10px;'>%s</td></tr>",
				t.Date.Format("2006/01/02"),
				formatMoney(value),
				acc.Name,
				contraAcc.Name,
				t.Description)
		}

		responseHtml += "</table></center></p>"
		responseHtml += fmt.Sprintf("<p><center>Closing: %s</center></p>", formatMoney(accData.Closing))
	}

	responseHtml += pageFooter

	response.Write([]byte(responseHtml))
}

func formatMoney(value domain.Money) string {
	valueAbs := value
	if valueAbs < 0 {
		valueAbs *= -1
	}

	s := fmt.Sprintf("%d.%02d", valueAbs/100, valueAbs%100)

	if value < 0 {
		s = "(" + s + ")"
	}

	return s
}
