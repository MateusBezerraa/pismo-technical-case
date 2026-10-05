package domain_test

import (
	"testing"
	"time"

	"github.com/MateusBezerraa/pismo-technical-case/internal/domain"
)

func TestNewTransaction(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name       string
		opType     domain.OperationType
		amount     float64
		wantAmount float64
		wantErr    error
	}{
		{"normal purchase is negative", domain.NormalPurchase, 50.0, -50.0, nil},
		{"installments is negative", domain.PurchaseWithInstallments, 23.5, -23.5, nil},
		{"withdrawal is negative", domain.Withdrawal, 18.7, -18.7, nil},
		{"credit voucher is positive", domain.CreditVoucher, 60.0, 60.0, nil},
		{"invalid operation type", domain.OperationType(99), 10.0, 0, domain.ErrInvalidOperation},
		{"zero amount", domain.NormalPurchase, 0, 0, domain.ErrInvalidAmount},
		{"negative amount", domain.NormalPurchase, -10.0, 0, domain.ErrInvalidAmount},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := domain.NewTransaction(1, tt.opType, tt.amount, now)
			if err != tt.wantErr {
				t.Fatalf("expected err %v, got %v", tt.wantErr, err)
			}
			if tt.wantErr == nil && tx.Amount != tt.wantAmount {
				t.Errorf("expected amount %v, got %v", tt.wantAmount, tx.Amount)
			}
		})
	}
}
