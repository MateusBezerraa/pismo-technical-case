package httpx

import (
	"net/http"

	"github.com/MateusBezerraa/pismo-technical-case/internal/adapter/httpx/handler"
	"github.com/MateusBezerraa/pismo-technical-case/internal/adapter/httpx/middleware"
)

// NewRouter registers routes and wraps them with middleware.
func NewRouter(
	accountH *handler.AccountHandler,
	txH *handler.TransactionHandler,
	healthH *handler.HealthHandler,

) http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthH.Check)
	mux.HandleFunc("POST /accounts", accountH.Create)
	mux.HandleFunc("GET /accounts/{accountId}", accountH.GetByID)
	mux.HandleFunc("POST /transactions", txH.Create)

	// Order: Recovery (Extern) → Logging (Intern) → mux
	return middleware.Recovery(middleware.Logging(mux))
}
