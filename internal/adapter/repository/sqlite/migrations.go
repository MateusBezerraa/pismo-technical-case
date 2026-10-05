package sqlite

import "database/sql"

func migrate(db *sql.DB) error {
	schema := `
    CREATE TABLE IF NOT EXISTS accounts (
        id              INTEGER PRIMARY KEY AUTOINCREMENT,
        document_number TEXT NOT NULL
    );

    CREATE TABLE IF NOT EXISTS transactions (
        id                INTEGER PRIMARY KEY AUTOINCREMENT,
        account_id        INTEGER NOT NULL REFERENCES accounts(id),
        operation_type_id INTEGER NOT NULL,
        amount            REAL    NOT NULL,
        event_date        DATETIME NOT NULL
    );

    CREATE UNIQUE INDEX IF NOT EXISTS idx_accounts_document
    ON accounts(document_number);
    `
	_, err := db.Exec(schema)
	return err
}
