package domain

import "time"

type OperationType int

const (
	NormalPurchase           OperationType = 1
	PurchaseWithInstallments OperationType = 2
	Withdrawal               OperationType = 3
	CreditVoucher            OperationType = 4
)

func (o OperationType) isValid() bool {
	switch o {
	case NormalPurchase, PurchaseWithInstallments, Withdrawal, CreditVoucher:
		return true
	}
	return false
}

func (o OperationType) isDebit() bool {
	return o == NormalPurchase || o == PurchaseWithInstallments || o == Withdrawal
}

type Transaction struct {
	ID              int64
	AccountID       int64
	OperationTypeID OperationType
	Amount          float64
	EventDate       time.Time
}

func NewTransaction(accountID int64, opType OperationType, amount float64, eventDate time.Time) (*Transaction, error) {
	if !opType.isValid() {
		return nil, ErrInvalidOperation
	}
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}
	if opType.isDebit() {
		amount = -amount
	}
	return &Transaction{
		AccountID:       accountID,
		OperationTypeID: opType,
		Amount:          amount,
		EventDate:       eventDate,
	}, nil
}
