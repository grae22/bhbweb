package domain

import (
	"errors"
	"log"
	"os"
	"slices"
	"time"
)

type Book struct {
	Name         string
	accounts     accountService
	transactions transactionService
	balances     periodBalancesService
}

func NewBook(name string) (*Book, error) {
	accounts, err := newAccountService(name, false)
	if err != nil {
		return nil, err
	}

	newBook := Book{
		Name:         name,
		accounts:     accounts,
		transactions: newTransactionSerivce(name),
		balances:     newPeriodBalancesService(name, &accounts),
	}

	newBook.accounts.addRootAccount("Bank", TypeBank)
	newBook.accounts.addRootAccount("Income", TypeIncome)
	newBook.accounts.addRootAccount("Expenses", TypeExpense)
	newBook.accounts.addRootAccount("Debtors", TypeDebtor)
	newBook.accounts.addRootAccount("Creditors", TypeCreditor)

	return &newBook, nil
}

func LoadBook(name string) (*Book, error) {
	accounts, err := newAccountService(name, true)
	if errors.Is(err, os.ErrNotExist) {
		return NewBook(name)
	} else if err != nil {
		return nil, err
	}

	newBook := Book{
		Name:         name,
		accounts:     accounts,
		transactions: newTransactionSerivce(name),
		balances:     newPeriodBalancesService(name, &accounts),
	}

	return &newBook, nil
}

func (b Book) AccountForId(id int) (Account, bool) {
	return b.accounts.accountForId(id)
}

func (b *Book) AddAccount(
	parentId int,
	name string,
) error {

	acc, err := b.accounts.addAccount(parentId, name)
	if err != nil {
		return err
	}

	err = b.transactions.moveTransactionsToAccount(parentId, acc.Id)
	if err != nil {
		return errors.New("error checking for parent account transactions")
	}

	err = b.balances.recalculateBalances(b.accounts, b.transactions)

	return err
}

func (b *Book) GetAccounts() []Account {
	return b.accounts.getAccounts()
}

func (b *Book) AddTransaction(
	debitAccId int,
	creditAccId int,
	date time.Time,
	value Money,
	description string,
) error {

	debitAcc, ok := b.accounts.accountForId(debitAccId)
	if !ok {
		return errors.New("debit account not found")
	}
	if len(debitAcc.ChildrenId) > 0 {
		return errors.New("debit account has children")
	}

	creditAcc, ok := b.accounts.accountForId(creditAccId)
	if !ok {
		return errors.New("credit account not found")
	}
	if len(creditAcc.ChildrenId) > 0 {
		return errors.New("credit account has children")
	}

	if value < 0 {
		return errors.New("negative values are not allowed")
	}

	transaction, err := b.transactions.createTransaction(
		debitAccId,
		creditAccId,
		date,
		value,
		description)

	if err != nil {
		return err
	}

	return b.balances.updateBalances(transaction, b.accounts)
}

func (b *Book) BalancesForDate(date time.Time) (PeriodBalances, error) {
	return b.balances.balancesForDate(date)
}

func (b *Book) AccountTransactionsForPeriod(
	accountId int,
	periodId string,
) ([]Transaction, error) {

	transactions, err := b.transactions.transactionsForPeriod(periodId)
	if err != nil {
		return []Transaction{}, err
	}

	account, ok := b.accounts.accountForId(accountId)

	if !ok {
		return []Transaction{}, errors.New("account not found")
	}

	accountIdToInclude := []int{accountId}
	childrenId := b.getChildrenIdRecursive(account)
	accountIdToInclude = append(accountIdToInclude, childrenId...)

	accountTransactions := []Transaction{}

	for _, t := range transactions {
		if !slices.Contains(accountIdToInclude, t.DebitAccountId) &&
			!slices.Contains(accountIdToInclude, t.CreditAccountId) {
			continue
		}

		accountTransactions = append(accountTransactions, t)
	}

	return accountTransactions, nil
}

func PeriodIdForDate(t time.Time) string {
	return t.Format("0601")
}

func DateForPeriodId(id string) time.Time {
	date, err := time.Parse("0601", id)
	if err != nil {
		log.Printf("Failed to get date from period id \"%s\".\n", id)
	}
	return date
}

func (b *Book) getChildrenIdRecursive(account Account) []int {

	childrenId := []int{}

	for _, childId := range account.ChildrenId {
		childrenId = append(childrenId, childId)

		child, ok := b.accounts.accountForId(childId)
		if !ok {
			log.Println("Account not found.")
			return childrenId
		}

		foundChildrenId := b.getChildrenIdRecursive(child)

		childrenId = append(childrenId, foundChildrenId...)
	}

	return childrenId
}
