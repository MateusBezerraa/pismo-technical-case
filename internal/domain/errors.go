package domain

import "errors"

var (
	ErrAccountNotFound       = errors.New("account not found")
	ErrInvalidOperation      = errors.New("invalid operation type")
	ErrInvalidAmount         = errors.New("amount must be greater than zero")
	ErrInvalidDocument       = errors.New("document number is required")
	ErrDocumentAlreadyExists = errors.New("document number already exists")
)
