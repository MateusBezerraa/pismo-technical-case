package usecase

import (
	"context"
	"time"

	"github.com/MateusBezerraa/pismo-technical-case/internal/domain"
)

type TransactionRepository interface {
	Save(ctx context.Context, tx *domain.Transaction) error
}

type TransactionUseCase struct {
	txRepo      TransactionRepository
	accountRepo AccountRepository
}

func NewTransactionUseCase(txRepo TransactionRepository, accountRepo AccountRepository) *TransactionUseCase {
	return &TransactionUseCase{txRepo: txRepo, accountRepo: accountRepo}
}

func (uc *TransactionUseCase) Create(
	ctx context.Context,
	accountID int64,
	operationTypeID int,
	amount float64,
) (*domain.Transaction, error) {

	if _, err := uc.accountRepo.FindByID(ctx, accountID); err != nil {
		return nil, err
	}

	tx, err := domain.NewTransaction(
		accountID,
		domain.OperationType(operationTypeID),
		amount,
		time.Now().UTC(),
	)
	if err != nil {
		return nil, err
	}

	if err := uc.txRepo.Save(ctx, tx); err != nil {
		return nil, err
	}
	return tx, nil
}
