package utils

import (
	"backend/dto"
	"encoding/json"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func Success(w http.ResponseWriter, status int, message string, data interface{}) {
	JSONResponse(w, status, dto.SuccessResponse{
		Success: true,
		Status:  status,
		Message: message,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, status int, message string) {
	JSONResponse(w, status, dto.ErrorResponse{
		Success: false,
		Status:  status,
		Message: message,
	})
}
