package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MateusBezerraa/pismo-technical-case/internal/adapter/httpx/dto"
	"github.com/MateusBezerraa/pismo-technical-case/internal/usecase"
)

type AccountHandler struct {
	uc *usecase.AccountUseCase
}

func NewAccountHandler(uc *usecase.AccountUseCase) *AccountHandler {
	return &AccountHandler{uc: uc}
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	account, err := h.uc.Create(r.Context(), req.DocumentNumber)
	if err != nil {
		writeError(w, statusFromError(err), messageFromError(err))
		return
	}
	writeJSON(w, http.StatusCreated, dto.AccountResponse{
		AccountID:      account.ID,
		DocumentNumber: account.DocumentNumber,
	})
}

func (h *AccountHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("accountId"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid account id")
		return
	}
	account, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, statusFromError(err), messageFromError(err))
		return
	}
	writeJSON(w, http.StatusOK, dto.AccountResponse{
		AccountID:      account.ID,
		DocumentNumber: account.DocumentNumber,
	})
}
