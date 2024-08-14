package domain

import (
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type transactionService struct {
	bookName string
}

func newTransactionSerivce(bookName string) transactionService {
	return transactionService{
		bookName: bookName,
	}
}

func (ts *transactionService) createTransaction(
	debitAccountId int,
	creditAccountId int,
	date time.Time,
	value Money,
	description string,
) (transaction Transaction, err error) {

	transaction = newTransaction(debitAccountId, creditAccountId, date, value, description)

	periodId := PeriodIdForDate(date)

	transactionLog, err := loadTransactionLog(ts.bookName, periodId)
	if err != nil {
		return
	}

	transactionLog.addTransactions([]Transaction{transaction})

	err = transactionLog.save(ts.bookName)

	return
}

func (ts *transactionService) transactionsForPeriod(periodId string) ([]Transaction, error) {

	transactionLog, err := loadTransactionLog(ts.bookName, periodId)
	if err != nil {
		return []Transaction{}, err
	}

	return transactionLog.Transactions, nil
}

func (ts *transactionService) moveTransactionsToAccount(
	fromtAccId int,
	toAccId int,
) error {

	err := filepath.Walk(
		".",
		func(
			path string,
			info fs.FileInfo,
			err error) error {

			if err != nil {
				return err
			}

			if !strings.HasPrefix(path, ts.bookName) ||
				!strings.HasSuffix(path, ".transactions") {
				return nil
			}

			tlog, err := loadTransactionLogFile(path)
			if err != nil {
				return err
			}

			var transactions []Transaction = make([]Transaction, 0, len(tlog.Transactions))

			for _, t := range tlog.Transactions {
				if t.DebitAccountId == fromtAccId {
					t.DebitAccountId = toAccId
				} else if t.CreditAccountId == fromtAccId {
					t.CreditAccountId = toAccId
				}

				transactions = append(transactions, t)
			}

			tlog.Transactions = transactions

			return tlog.save(ts.bookName)
		})

	return err
}

func (ts *transactionService) getAllTransactionLogsOldestFirst() ([]transactionLog, error) {
	tlogs := []transactionLog{}

	err := filepath.Walk(
		".",
		func(path string, info fs.FileInfo, err error) error {

			if err != nil {
				return err
			}

			if !strings.HasPrefix(path, ts.bookName) ||
				!strings.HasSuffix(path, ".transactions") {
				return nil
			}

			tlog, err := loadTransactionLogFile(path)
			if err != nil {
				return err
			}

			tlogs = append(tlogs, tlog)

			return nil
		})

	if err != nil {
		return []transactionLog{}, err
	}

	sort.Slice(
		tlogs,
		func(i, j int) bool {
			return tlogs[i].PeriodId < tlogs[j].PeriodId
		})

	return tlogs, err
}
