// Package dto holds the request/response shapes for the HTTP layer.
//
// DTOs are kept separate from domain types so JSON and domain can evolve
// independently — the domain never needs `json:` tags, and the API never
// leaks internal representation.
package dto

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number"`
}

type CreateTransactionRequest struct {
	AccountID       int64   `json:"account_id"`
	OperationTypeID int     `json:"operation_type_id"`
	Amount          float64 `json:"amount"`
}

type AccountResponse struct {
	AccountID      int64  `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}

type TransactionResponse struct {
	TransactionID   int64   `json:"transaction_id"`
	AccountID       int64   `json:"account_id"`
	OperationTypeID int     `json:"operation_type_id"`
	Amount          float64 `json:"amount"`
	EventDate       string  `json:"event_date"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
