package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MateusBezerraa/pismo-technical-case/internal/adapter/httpx/dto"
	"github.com/MateusBezerraa/pismo-technical-case/internal/usecase"
)

type TransactionHandler struct {
	uc *usecase.TransactionUseCase
}

func NewTransactionHandler(uc *usecase.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{uc: uc}
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	tx, err := h.uc.Create(r.Context(), req.AccountID, req.OperationTypeID, req.Amount)
	if err != nil {
		writeError(w, statusFromError(err), messageFromError(err))
		return
	}
	writeJSON(w, http.StatusCreated, dto.TransactionResponse{
		TransactionID:   tx.ID,
		AccountID:       tx.AccountID,
		OperationTypeID: int(tx.OperationTypeID),
		Amount:          tx.Amount,
		EventDate:       tx.EventDate.Format(time.RFC3339),
	})
}
