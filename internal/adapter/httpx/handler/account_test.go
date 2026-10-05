package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccountEndpoints(t *testing.T) {
	t.Run("POST /accounts creates account", func(t *testing.T) {
		router := setupServer(t)

		body := bytes.NewBufferString(`{"document_number":"12345678900"}`)
		req := httptest.NewRequest(http.MethodPost, "/accounts", body)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d — body: %s", rec.Code, rec.Body.String())
		}

		var resp map[string]any
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatal(err)
		}
		if resp["document_number"] != "12345678900" {
			t.Errorf("unexpected document: %v", resp["document_number"])
		}
		if resp["account_id"] == nil {
			t.Error("expected account_id in response")
		}
	})

	t.Run("POST /accounts with empty document returns 400", func(t *testing.T) {
		router := setupServer(t)

		body := bytes.NewBufferString(`{"document_number":""}`)
		req := httptest.NewRequest(http.MethodPost, "/accounts", body)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("POST /accounts with invalid JSON returns 400", func(t *testing.T) {
		router := setupServer(t)

		body := bytes.NewBufferString(`{invalid json}`)
		req := httptest.NewRequest(http.MethodPost, "/accounts", body)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("POST /accounts with duplicate document returns 409", func(t *testing.T) {
		router := setupServer(t)

		body1 := bytes.NewBufferString(`{"document_number":"12345678900"}`)
		req1 := httptest.NewRequest(http.MethodPost, "/accounts", body1)
		rec1 := httptest.NewRecorder()
		router.ServeHTTP(rec1, req1)
		if rec1.Code != http.StatusCreated {
			t.Fatalf("first create: expected 201, got %d", rec1.Code)
		}

		body2 := bytes.NewBufferString(`{"document_number":"12345678900"}`)
		req2 := httptest.NewRequest(http.MethodPost, "/accounts", body2)
		rec2 := httptest.NewRecorder()
		router.ServeHTTP(rec2, req2)
		if rec2.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d", rec2.Code)
		}
	})

	t.Run("GET /accounts/{id} returns account", func(t *testing.T) {
		router := setupServer(t)

		createBody := bytes.NewBufferString(`{"document_number":"12345678900"}`)
		createReq := httptest.NewRequest(http.MethodPost, "/accounts", createBody)
		createRec := httptest.NewRecorder()
		router.ServeHTTP(createRec, createReq)

		var created map[string]any
		if err := json.NewDecoder(createRec.Body).Decode(&created); err != nil {
			t.Fatal(err)
		}
		accountID := int64(created["account_id"].(float64))

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/accounts/%d", accountID), nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("GET /accounts/{id} returns 404 for non-existent", func(t *testing.T) {
		router := setupServer(t)

		req := httptest.NewRequest(http.MethodGet, "/accounts/99999", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", rec.Code)
		}
	})

	t.Run("GET /accounts/{id} returns 400 for invalid id", func(t *testing.T) {
		router := setupServer(t)

		req := httptest.NewRequest(http.MethodGet, "/accounts/abc", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
	})
}
