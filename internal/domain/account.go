package domain

type Account struct {
	ID             int64
	DocumentNumber string
}

func NewAccount(documentNumber string) (*Account, error) {
	if documentNumber == "" {
		return nil, ErrInvalidDocument
	}
	return &Account{DocumentNumber: documentNumber}, nil
}
