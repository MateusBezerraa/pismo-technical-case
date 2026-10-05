package sqlite

import (
	"context"
	"database/sql"

	"github.com/MateusBezerraa/pismo-technical-case/internal/domain"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Save(ctx context.Context, tx *domain.Transaction) error {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO transactions
         (account_id, operation_type_id, amount, event_date)
         VALUES (?, ?, ?, ?)`,
		tx.AccountID, int(tx.OperationTypeID), tx.Amount, tx.EventDate,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	tx.ID = id
	return nil
}
