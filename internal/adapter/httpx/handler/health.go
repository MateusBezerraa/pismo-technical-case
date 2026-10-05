package handler

import "net/http"

// HealthHandler responds to liveness probes.
// Kept dependency-free on purpose: liveness only checks that the process
// is running. A readiness probe (checking DB connectivity) would be a
// separate handler with its own route.
type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
