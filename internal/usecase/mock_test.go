package usecase_test

import (
	"context"

	"github.com/MateusBezerraa/pismo-technical-case/internal/domain"
)

// ===== AccountRepository =====

type mockAccountRepo struct {
	accounts map[int64]*domain.Account
	nextID   int64
}

func newMockAccountRepo() *mockAccountRepo {
	return &mockAccountRepo{accounts: map[int64]*domain.Account{}, nextID: 1}
}

func (m *mockAccountRepo) Save(_ context.Context, a *domain.Account) error {
	a.ID = m.nextID
	m.nextID++
	m.accounts[a.ID] = a
	return nil
}

func (m *mockAccountRepo) FindByID(_ context.Context, id int64) (*domain.Account, error) {
	a, ok := m.accounts[id]
	if !ok {
		return nil, domain.ErrAccountNotFound
	}
	return a, nil
}

// ===== TransactionRepository =====

type mockTransactionRepo struct {
	saved  []*domain.Transaction
	nextID int64
}

func newMockTransactionRepo() *mockTransactionRepo {
	return &mockTransactionRepo{nextID: 1}
}

func (m *mockTransactionRepo) Save(_ context.Context, tx *domain.Transaction) error {
	tx.ID = m.nextID
	m.nextID++
	m.saved = append(m.saved, tx)
	return nil
}
