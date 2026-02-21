package interfaces

import "net/http"

// AppError adalah struct agar Usecase bisa kirim status code & message sekaligus
type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

// Helper untuk membuat error baru
func NewAppError(code int, msg string) error {
	return &AppError{Code: code, Message: msg}
}

// --- Errors dengan Status Code ---
var (
	ErrInvalidCredentials = NewAppError(http.StatusUnauthorized, "invalid username or password")
	ErrInvalidJSONFormat  = NewAppError(http.StatusBadRequest, "invalid json format")
	ErrInternalServer     = NewAppError(http.StatusInternalServerError, "something went wrong on our side")
)
