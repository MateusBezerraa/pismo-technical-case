package domain_test

import (
	"testing"

	"github.com/MateusBezerraa/pismo-technical-case/internal/domain"
)

func TestNewAccount(t *testing.T) {
	t.Run("valid document", func(t *testing.T) {
		acc, err := domain.NewAccount("12345678900")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if acc.DocumentNumber != "12345678900" {
			t.Errorf("expected document 12345678900, got %s", acc.DocumentNumber)
		}
	})

	t.Run("empty document returns error", func(t *testing.T) {
		_, err := domain.NewAccount("")
		if err != domain.ErrInvalidDocument {
			t.Errorf("expected ErrInvalidDocument, got %v", err)
		}
	})
}
