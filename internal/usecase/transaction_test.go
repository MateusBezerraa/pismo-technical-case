package usecase_test

import (
	"context"
	"testing"

	"github.com/MateusBezerraa/pismo-technical-case/internal/domain"
	"github.com/MateusBezerraa/pismo-technical-case/internal/usecase"
)

func TestTransactionUseCase_Create(t *testing.T) {
	accountRepo := newMockAccountRepo()
	txRepo := newMockTransactionRepo()
	uc := usecase.NewTransactionUseCase(txRepo, accountRepo)

	acc, _ := usecase.NewAccountUseCase(accountRepo).Create(context.Background(), "12345678900")

	t.Run("creates debit transaction with negative amount", func(t *testing.T) {
		tx, err := uc.Create(context.Background(), acc.ID, int(domain.NormalPurchase), 50.0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tx.Amount != -50.0 {
			t.Errorf("expected -50.0, got %v", tx.Amount)
		}
	})

	t.Run("creates credit transaction with positive amount", func(t *testing.T) {
		tx, err := uc.Create(context.Background(), acc.ID, int(domain.CreditVoucher), 60.0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tx.Amount != 60.0 {
			t.Errorf("expected 60.0, got %v", tx.Amount)
		}
	})

	t.Run("returns error when account does not exist", func(t *testing.T) {
		_, err := uc.Create(context.Background(), 9999, int(domain.NormalPurchase), 10.0)
		if err != domain.ErrAccountNotFound {
			t.Errorf("expected ErrAccountNotFound, got %v", err)
		}
	})

	t.Run("returns error for invalid operation type", func(t *testing.T) {
		_, err := uc.Create(context.Background(), acc.ID, 99, 10.0)
		if err != domain.ErrInvalidOperation {
			t.Errorf("expected ErrInvalidOperation, got %v", err)
		}
	})

	t.Run("returns error for non-positive amount", func(t *testing.T) {
		_, err := uc.Create(context.Background(), acc.ID, int(domain.NormalPurchase), 0)
		if err != domain.ErrInvalidAmount {
			t.Errorf("expected ErrInvalidAmount, got %v", err)
		}
	})
}
