package account

type AccountRepository interface {
	Save(account *Account) error
	Update(account *Account) error
	FindByID(accountID AccountID) (*Account, error)
}
