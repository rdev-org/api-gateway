package utils

import (
	"encoding/json"
	"net/http"
)

func setHeader(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func ResponseOk(w http.ResponseWriter, data any) {
	setHeader(w, http.StatusOK, data)
}

func ResponseCreated(w http.ResponseWriter, data any) {
	setHeader(w, http.StatusCreated, data)
}

func ResponseError(w http.ResponseWriter, statusCode int, message string) {
	setHeader(w, statusCode, map[string]string{"error": message})
}
