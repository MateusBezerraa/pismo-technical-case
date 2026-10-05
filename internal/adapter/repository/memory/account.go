package memory

import (
	"context"
	"sync"

	"github.com/MateusBezerraa/pismo-technical-case/internal/domain"
)

// AccountRepository is the in-memory implementation of usecase.AccountRepository.
//
// It exists as an alternative to SQLite — the use case wires against the
// interface, so swapping adapters requires no changes above. Uses RWMutex
// because reads (FindByID) outnumber writes (Save).
type AccountRepository struct {
	mu       sync.RWMutex
	accounts map[int64]*domain.Account
	nextID   int64
}

func NewAccountRepository() *AccountRepository {
	return &AccountRepository{
		accounts: make(map[int64]*domain.Account),
		nextID:   1,
	}
}

func (r *AccountRepository) Save(ctx context.Context, a *domain.Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	a.ID = r.nextID
	r.nextID++
	r.accounts[a.ID] = a
	return nil
}

func (r *AccountRepository) FindByID(ctx context.Context, id int64) (*domain.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.accounts[id]
	if !ok {
		return nil, domain.ErrAccountNotFound
	}
	return a, nil
}
