package usecase

import (
	"context"

	"github.com/MateusBezerraa/pismo-technical-case/internal/domain"
)

type AccountRepository interface {
	Save(ctx context.Context, account *domain.Account) error
	FindByID(ctx context.Context, id int64) (*domain.Account, error)
}

type AccountUseCase struct {
	repo AccountRepository
}

func NewAccountUseCase(repo AccountRepository) *AccountUseCase {
	return &AccountUseCase{repo: repo}
}

func (uc *AccountUseCase) Create(ctx context.Context, documentNumber string) (*domain.Account, error) {
	account, err := domain.NewAccount(documentNumber)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.Save(ctx, account); err != nil {
		return nil, err
	}
	return account, nil
}

func (uc *AccountUseCase) GetByID(ctx context.Context, id int64) (*domain.Account, error) {
	return uc.repo.FindByID(ctx, id)
}
