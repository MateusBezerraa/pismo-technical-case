package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/MateusBezerraa/pismo-technical-case/internal/domain"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Save(ctx context.Context, a *domain.Account) error {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO accounts (document_number) VALUES (?)`,
		a.DocumentNumber,
	)
	if err != nil {
		// modernc.org/sqlite returns a generic error with the SQLite message.
		// We check the message because the driver does not expose a stable typed code.
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return domain.ErrDocumentAlreadyExists
		}
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	a.ID = id
	return nil
}

func (r *AccountRepository) FindByID(ctx context.Context, id int64) (*domain.Account, error) {
	var a domain.Account
	err := r.db.QueryRowContext(ctx,
		`SELECT id, document_number FROM accounts WHERE id = ?`, id,
	).Scan(&a.ID, &a.DocumentNumber)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrAccountNotFound
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}
