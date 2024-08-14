package domain

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sort"
)

type accountService struct {
	bookName      string
	accountsById  map[int]*Account
	nextAccountId int
}

func newAccountService(bookName string) (accountService, error) {
	accounts := accountService{
		bookName:     bookName,
		accountsById: make(map[int]*Account),
	}

	err := accounts.load()

	return accounts, err
}

func (as *accountService) save() error {
	log.Printf("Saving accounts for \"%s\"...", as.bookName)

	f, err := os.OpenFile(
		as.bookName+".accounts",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		os.ModePerm)

	if err != nil {
		log.Println("Error opening file.")
		return err
	}

	defer f.Close()

	accounts := make([]Account, 0, len(as.accountsById))
	for _, acc := range as.accountsById {
		accounts = append(accounts, *acc)
	}

	bytes, err := json.MarshalIndent(accounts, "", "  ")

	if err != nil {
		log.Println("Error serialising accounts.")
		return err
	}

	_, err = f.Write(bytes)

	if err != nil {
		log.Println("Error writing file.")
		return err
	}

	return nil
}

func (as *accountService) accountForId(id int) (Account, bool) {
	a, ok := as.accountsById[id]

	if !ok {
		zero := Account{}
		return zero, false
	}

	return *a, ok
}

func (as *accountService) addRootAccount(
	name string,
	accountType AccountType,
) {
	account := newAccount(
		as.nextAccountId,
		accountType,
		name,
		NoParentId)

	as.nextAccountId++

	as.accountsById[account.Id] = &account
}

func (as *accountService) addAccount(
	parentId int,
	name string,
) (Account, error) {

	parent, ok := as.accountsById[parentId]

	if !ok {
		return Account{}, errors.New("no account found for parent id")
	}

	child := parent.createChild(as.nextAccountId, name)

	as.nextAccountId++

	as.accountsById[child.Id] = &child

	err := as.save()

	return child, err
}

func (as *accountService) getAccounts() []Account {
	accounts := make([]Account, 0, len(as.accountsById))

	for _, acc := range as.accountsById {
		accounts = append(accounts, *acc)
	}

	sort.Slice(
		accounts,
		func(i, j int) bool {
			return accounts[i].Id < accounts[j].Id
		})

	return accounts
}

func (as *accountService) load() error {
	log.Printf("Loading accounts for \"%s\"...", as.bookName)

	if _, err := os.Stat(as.bookName + ".accounts"); errors.Is(err, os.ErrNotExist) {
		log.Printf("File not found, creating new accounts...")
		return nil
	}

	bytes, err := os.ReadFile(as.bookName + ".accounts")

	if err != nil {
		log.Println("Error while reading file.")
		return err
	}

	accounts := []Account{}

	err = json.Unmarshal(bytes, &accounts)

	if err != nil {
		log.Println("Error while deserialising file.")
		return err
	}

	for _, a := range accounts {
		as.accountsById[a.Id] = &a

		if as.nextAccountId <= a.Id {
			as.nextAccountId = a.Id + 1
		}
	}

	return nil
}
