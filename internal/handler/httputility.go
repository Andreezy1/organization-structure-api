package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"org_struct_api/internal/models"
	"strconv"
)

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if body == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Println("encode response", "error", err)
	}
}

func writeServiceError(w http.ResponseWriter, err error) {
	writeError(w, statusFromErr(err), err.Error())
}

func statusFromErr(err error) int {
	switch {
	case errors.Is(err, models.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, models.ErrConflict), errors.Is(err, models.ErrCycle):
		return http.StatusConflict
	case errors.Is(err, models.ErrValidation):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func parsePathID(r *http.Request, key string) (uint, error) {
	s := r.PathValue(key)
	if s == "" {
		return 0, fmt.Errorf("%s is required", key)
	}
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return uint(n), nil
}

func queryIntDefault(r *http.Request, key string, def int) (int, error) {
	s := r.URL.Query().Get(key)
	if s == "" {
		return def, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return n, nil
}
