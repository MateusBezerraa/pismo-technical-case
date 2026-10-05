package handler

import (
	"errors"
	"net/http"

	"github.com/MateusBezerraa/pismo-technical-case/internal/domain"
)

func statusFromError(err error) int {
	switch {
	case errors.Is(err, domain.ErrAccountNotFound):
		return http.StatusNotFound // 404
	case errors.Is(err, domain.ErrInvalidDocument),
		errors.Is(err, domain.ErrInvalidOperation),
		errors.Is(err, domain.ErrInvalidAmount):
		return http.StatusBadRequest //400
	case errors.Is(err, domain.ErrDocumentAlreadyExists):
		return http.StatusConflict // 409
	default:
		return http.StatusInternalServerError // 500
	}
}

func messageFromError(err error) string {
	switch {
	case errors.Is(err, domain.ErrAccountNotFound),
		errors.Is(err, domain.ErrInvalidDocument),
		errors.Is(err, domain.ErrInvalidOperation),
		errors.Is(err, domain.ErrInvalidAmount),
		errors.Is(err, domain.ErrDocumentAlreadyExists):
		return err.Error()
	default:
		return "internal server error"
	}
}
