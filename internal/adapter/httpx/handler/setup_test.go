package handler_test

import (
	"net/http"
	"testing"

	"github.com/MateusBezerraa/pismo-technical-case/internal/adapter/httpx"
	"github.com/MateusBezerraa/pismo-technical-case/internal/adapter/httpx/handler"
	"github.com/MateusBezerraa/pismo-technical-case/internal/adapter/repository/sqlite"
	"github.com/MateusBezerraa/pismo-technical-case/internal/usecase"
)

// setupServer creates a router with SQLite :memory: isolated for the tests.
func setupServer(t *testing.T) http.Handler {
	t.Helper()

	db, err := sqlite.Open(":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	accountRepo := sqlite.NewAccountRepository(db)
	txRepo := sqlite.NewTransactionRepository(db)

	accountUC := usecase.NewAccountUseCase(accountRepo)
	txUC := usecase.NewTransactionUseCase(txRepo, accountRepo)

	accountH := handler.NewAccountHandler(accountUC)
	txH := handler.NewTransactionHandler(txUC)
	healthH := handler.NewHealthHandler()

	router := httpx.NewRouter(accountH, txH, healthH)
	return router
}
