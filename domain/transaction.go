package domain

import "time"

type Transaction struct {
	DebitAccountId  int
	CreditAccountId int
	Date            time.Time
	Value           Money
	Description     string
	Timestamp       time.Time
}

func newTransaction(
	debitAccountId int,
	creditAccountId int,
	date time.Time,
	value Money,
	description string,
) Transaction {

	return Transaction{
		DebitAccountId:  debitAccountId,
		CreditAccountId: creditAccountId,
		Date:            date,
		Value:           value,
		Description:     description,
		Timestamp:       time.Now(),
	}
}
