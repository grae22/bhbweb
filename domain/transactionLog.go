package domain

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

type transactionLog struct {
	PeriodId     string
	Transactions []Transaction
}

func loadTransactionLog(bookName string, periodId string) (tlog transactionLog, err error) {

	filename := filenameForPeriodId(bookName, periodId)

	log.Printf("Loading transaction log \"%s\"...", filename)

	_, err = os.Stat(filename)
	if errors.Is(err, os.ErrNotExist) {
		log.Println("Transaction log not found, creating a new one.")
		tlog = transactionLog{
			PeriodId:     periodId,
			Transactions: make([]Transaction, 0),
		}
		return tlog, nil
	}

	bytes, err := os.ReadFile(filename)
	if err != nil {
		log.Println("Error reading file.")
		return tlog, err
	}

	err = json.Unmarshal(bytes, &tlog)
	if err != nil {
		log.Println("Error deserialising file.")
		return tlog, err
	}

	return tlog, nil
}

func loadTransactionLogFile(path string) (transactionLog, error) {

	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Println("Error reading file.")
		return transactionLog{}, err
	}

	var tlog transactionLog

	err = json.Unmarshal(bytes, &tlog)
	if err != nil {
		log.Println("Error deserialising file.")
		return transactionLog{}, err
	}

	return tlog, nil
}

func (tl *transactionLog) save(bookName string) error {

	log.Printf("Saving transaction log for book \"%s\" and period \"%s\"...", bookName, tl.PeriodId)

	if tl.PeriodId == "" {
		log.Println("Transaction log has no period id.")
		return errors.New("no period id")
	}

	bytes, err := json.MarshalIndent(tl, "", "  ")
	if err != nil {
		log.Println("Error serialising object.")
		return errors.New("serialisation failed")
	}

	filename := filenameForPeriodId(bookName, tl.PeriodId)

	err = os.WriteFile(filename, bytes, os.ModePerm)
	if err != nil {
		log.Println("Error writing file.")
		return err
	}

	return nil
}

func (tl *transactionLog) addTransactions(transactions []Transaction) {
	tl.Transactions = append(tl.Transactions, transactions...)
}

func filenameForPeriodId(bookName string, periodId string) string {
	return bookName + "_" + periodId + ".transactions"
}
