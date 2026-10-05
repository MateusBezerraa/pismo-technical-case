package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MateusBezerraa/pismo-technical-case/internal/adapter/httpx"
	"github.com/MateusBezerraa/pismo-technical-case/internal/adapter/httpx/handler"
	"github.com/MateusBezerraa/pismo-technical-case/internal/adapter/repository/sqlite"
	"github.com/MateusBezerraa/pismo-technical-case/internal/usecase"
)

// main wires the application together and starts the HTTP server.
//
// This is the only place that knows which concrete adapter is in use.
// Swapping SQLite for Postgres or in-memory means changing only this file.
func main() {
	// 1. Persistence (outbound adapter)
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "pismo.db"
	}
	db, err := sqlite.Open(dbPath)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	// 2. Repositories
	accountRepo := sqlite.NewAccountRepository(db)
	txRepo := sqlite.NewTransactionRepository(db)

	// 3. Use cases
	accountUC := usecase.NewAccountUseCase(accountRepo)
	txUC := usecase.NewTransactionUseCase(txRepo, accountRepo)

	// 4. HTTP handlers (inbound adapter)
	accountH := handler.NewAccountHandler(accountUC)
	txH := handler.NewTransactionHandler(txUC)
	healthH := handler.NewHealthHandler()

	// 5. Router
	router := httpx.NewRouter(accountH, txH, healthH)

	// 6. HTTP server with timeouts
	//
	// Timeouts are non-negotiable in production. Without them, a slow
	// or malicious client can hold a connection open indefinitely,
	// exhausting the server's file descriptors (Slowloris attack).
	//
	//   - ReadHeaderTimeout: max time to read request headers
	//   - ReadTimeout:       max time to read the whole request body
	//   - WriteTimeout:      max time to write the response
	//   - IdleTimeout:       max keep-alive idle time between requests
	//
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// 7. Start server in a goroutine
	//
	// ListenAndServe blocks until the server stops. Running it in a
	// goroutine lets main continue and wait for shutdown signals.
	// Errors from a normal shutdown (ErrServerClosed) are not real errors.
	//
	go func() {
		log.Printf("listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	// 8. Wait for SIGINT / SIGTERM
	//
	// SIGINT  → Ctrl+C during local development
	// SIGTERM → what Docker/Kubernetes send when stopping a container
	//
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // blocks until a signal arrives
	log.Println("shutdown signal received")

	// 9. Graceful shutdown with timeout
	//
	// Shutdown stops accepting new connections and waits up to 10s for
	// in-flight requests to finish. If the timeout expires, remaining
	// requests are aborted — better to die than hang forever.
	//
	// Once Shutdown returns, the deferred db.Close() runs safely, with
	// no queries still in progress.
	//
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
	log.Println("server stopped gracefully")
}
