package response

import (
	"api-guardian/internal/domain/auth/interfaces"
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

// Error sekarang pintar, dia bongkar status code dari Usecase sendiri
func Error(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError
	message := err.Error()

	// Cek apakah ini AppError dari Usecase?
	if appErr, ok := err.(*interfaces.AppError); ok {
		statusCode = appErr.Code
		message = appErr.Message
	}

	JSON(w, statusCode, map[string]string{"error": message})
}
