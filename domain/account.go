package domain

type AccountType int
type Money int

type Account struct {
	Id          int
	AccountType AccountType
	Name        string
	ParentId    int
	ChildrenId  []int
}

const (
	TypeBank AccountType = iota
	TypeIncome
	TypeExpense
	TypeCreditor
	TypeDebtor
)

const NoParentId int = -1

func newAccount(
	id int,
	accountType AccountType,
	name string,
	parentId int,
) Account {

	newAccount := Account{
		Id:          id,
		AccountType: accountType,
		Name:        name,
		ParentId:    parentId,
		ChildrenId:  make([]int, 0),
	}

	return newAccount
}

func (a *Account) createChild(
	id int,
	name string,
) Account {

	child := newAccount(
		id,
		a.AccountType,
		name,
		a.Id)

	a.ChildrenId = append(a.ChildrenId, child.Id)

	return child
}
