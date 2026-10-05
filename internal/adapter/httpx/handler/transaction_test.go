package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTransactionEndpoints(t *testing.T) {
	router := setupServer(t)

	// Create an account to use in the tests
	createBody := bytes.NewBufferString(`{"document_number":"12345678900"}`)
	createReq := httptest.NewRequest(http.MethodPost, "/accounts", createBody)
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, createReq)
	var created map[string]any
	json.NewDecoder(createRec.Body).Decode(&created)
	accountID := int64(created["account_id"].(float64))

	t.Run("POST /transactions creates debit with negative amount", func(t *testing.T) {
		body := bytes.NewBufferString(fmt.Sprintf(
			`{"account_id":%d,"operation_type_id":1,"amount":50.0}`, accountID))
		req := httptest.NewRequest(http.MethodPost, "/transactions", body)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d — body: %s", rec.Code, rec.Body.String())
		}
		var resp map[string]any
		json.NewDecoder(rec.Body).Decode(&resp)
		if resp["amount"] != -50.0 {
			t.Errorf("expected -50.0, got %v", resp["amount"])
		}
	})

	t.Run("POST /transactions creates credit with positive amount", func(t *testing.T) {
		body := bytes.NewBufferString(fmt.Sprintf(
			`{"account_id":%d,"operation_type_id":4,"amount":60.0}`, accountID))
		req := httptest.NewRequest(http.MethodPost, "/transactions", body)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d", rec.Code)
		}
		var resp map[string]any
		json.NewDecoder(rec.Body).Decode(&resp)
		if resp["amount"] != 60.0 {
			t.Errorf("expected 60.0, got %v", resp["amount"])
		}
	})

	t.Run("returns 404 for non-existent account", func(t *testing.T) {
		body := bytes.NewBufferString(`{"account_id":99999,"operation_type_id":1,"amount":10.0}`)
		req := httptest.NewRequest(http.MethodPost, "/transactions", body)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", rec.Code)
		}
	})

	t.Run("returns 400 for invalid operation type", func(t *testing.T) {
		body := bytes.NewBufferString(fmt.Sprintf(
			`{"account_id":%d,"operation_type_id":99,"amount":10.0}`, accountID))
		req := httptest.NewRequest(http.MethodPost, "/transactions", body)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("returns 400 for zero amount", func(t *testing.T) {
		body := bytes.NewBufferString(fmt.Sprintf(
			`{"account_id":%d,"operation_type_id":1,"amount":0}`, accountID))
		req := httptest.NewRequest(http.MethodPost, "/transactions", body)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
	})
}
