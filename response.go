package main

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func respondJSON(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Code:    status,
		Message: message,
		Data:    data,
	})
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, message, nil)
}
