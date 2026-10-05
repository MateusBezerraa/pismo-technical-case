package usecase_test

import (
	"context"
	"testing"

	"github.com/MateusBezerraa/pismo-technical-case/internal/domain"
	"github.com/MateusBezerraa/pismo-technical-case/internal/usecase"
)

func TestAccountUseCase_Create(t *testing.T) {
	uc := usecase.NewAccountUseCase(newMockAccountRepo())

	t.Run("creates account successfully", func(t *testing.T) {
		acc, err := uc.Create(context.Background(), "12345678900")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if acc.ID != 1 {
			t.Errorf("expected ID 1, got %d", acc.ID)
		}
		if acc.DocumentNumber != "12345678900" {
			t.Errorf("unexpected document: %s", acc.DocumentNumber)
		}
	})

	t.Run("returns domain error for empty document", func(t *testing.T) {
		_, err := uc.Create(context.Background(), "")
		if err != domain.ErrInvalidDocument {
			t.Errorf("expected ErrInvalidDocument, got %v", err)
		}
	})
}

func TestAccountUseCase_GetByID(t *testing.T) {
	uc := usecase.NewAccountUseCase(newMockAccountRepo())

	created, _ := uc.Create(context.Background(), "12345678900")

	t.Run("returns existing account", func(t *testing.T) {
		acc, err := uc.GetByID(context.Background(), created.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if acc.ID != created.ID {
			t.Errorf("expected ID %d, got %d", created.ID, acc.ID)
		}
	})

	t.Run("returns error for non-existent account", func(t *testing.T) {
		_, err := uc.GetByID(context.Background(), 999)
		if err != domain.ErrAccountNotFound {
			t.Errorf("expected ErrAccountNotFound, got %v", err)
		}
	})
}
