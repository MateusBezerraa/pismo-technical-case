package memory

import (
	"context"
	"sync"

	"github.com/MateusBezerraa/pismo-technical-case/internal/domain"
)

// TransactionRepository is the in-memory implementation for transactions.
//
// Uses plain Mutex (not RWMutex) because the only operation today is Save —
// no concurrent reads to optimize for. If FindByID is added, switch to RWMutex.
type TransactionRepository struct {
	mu           sync.Mutex
	transactions map[int64]*domain.Transaction
	nextID       int64
}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{
		transactions: make(map[int64]*domain.Transaction),
		nextID:       1,
	}
}

func (r *TransactionRepository) Save(ctx context.Context, tx *domain.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	tx.ID = r.nextID
	r.nextID++
	r.transactions[tx.ID] = tx
	return nil
}
