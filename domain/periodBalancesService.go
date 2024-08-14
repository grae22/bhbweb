package domain

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type periodBalancesService struct {
	bookName string
	accounts *accountService
}

func newPeriodBalancesService(
	bookName string,
	accounts *accountService,
) periodBalancesService {

	return periodBalancesService{
		bookName: bookName,
		accounts: accounts,
	}
}

func (pbs periodBalancesService) updateBalances(
	transaction Transaction,
	accounts accountService,
) error {

	log.Println("Updating period balances...")

	periodId := PeriodIdForDate(transaction.Date)

	balances, subsequentBalances, err := loadPeriodBalances(
		pbs.bookName,
		periodId,
		*pbs.accounts)

	if err != nil {
		log.Println("Error while loading period balances.")
		return err
	}

	allBalances := []PeriodBalances{balances}
	allBalances = append(allBalances, subsequentBalances...)

	err = updateAccountBalanceRecursive(
		transaction.DebitAccountId,
		-transaction.Value,
		accounts,
		allBalances)

	if err != nil {
		log.Println("Error udpating debit account balances.")
		return err
	}

	err = updateAccountBalanceRecursive(
		transaction.CreditAccountId,
		transaction.Value,
		accounts,
		allBalances)

	if err != nil {
		log.Println("Error udpating credit account balances.")
		return err
	}

	for _, b := range allBalances {
		err = b.save(pbs.bookName)
		if err != nil {
			log.Println("Error while writing updated balances.")
			return err
		}
	}

	return nil
}

func (pbs *periodBalancesService) recalculateBalances(
	accounts accountService,
	transactions transactionService,
) error {

	log.Printf("Recalculating balances for \"%s\"...\n", pbs.bookName)

	err := pbs.deleteAllBalances()
	if err != nil {
		return err
	}

	tlogs, err := transactions.getAllTransactionLogsOldestFirst()
	if err != nil {
		return err
	}

	for _, tlog := range tlogs {
		for _, t := range tlog.Transactions {
			err = pbs.updateBalances(t, accounts)

			if err != nil {
				return err
			}
		}
	}

	log.Println("Recalculating balances completed.")

	return nil
}

func (pbs periodBalancesService) balancesForDate(date time.Time) (PeriodBalances, error) {

	log.Printf("Finding balances for date \"%s\"...\n", date.Format("2006-01-02"))

	periodId := PeriodIdForDate(date)
	log.Printf("Finding balances for period \"%s\"...", periodId)

	pbals, _, err := loadPeriodBalances(
		pbs.bookName,
		periodId,
		*pbs.accounts)

	return pbals, err
}

func (pbs periodBalancesService) deleteAllBalances() error {
	log.Printf("Deleting existing balances for \"%s\"...\n", pbs.bookName)

	err := filepath.Walk(
		".",
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !strings.HasPrefix(path, pbs.bookName) ||
				!strings.HasSuffix(path, ".balances") {
				return nil
			}

			err = os.Remove(path)

			return err
		})

	return err
}

func updateAccountBalanceRecursive(
	accountId int,
	value Money,
	accounts accountService,
	balances []PeriodBalances,
) error {

	for i, b := range balances {
		accBalances := b.BalanceByAccountId[accountId]
		accBalances.Closing += value

		if i > 0 {
			accBalances.Opening += value
		}

		b.BalanceByAccountId[accountId] = accBalances
	}

	acc, ok := accounts.accountForId(accountId)
	if !ok {
		return errors.New("account not found for id")
	}

	if acc.ParentId == NoParentId {
		return nil
	}

	return updateAccountBalanceRecursive(
		acc.ParentId,
		value,
		accounts,
		balances)
}
