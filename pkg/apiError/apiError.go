package apiError

import (
	"net/http"
)

type APIError struct {
	err  string
	code int
}

func (e APIError) Error() string {
	return e.err
}

func (e APIError) Code() int {
	return e.code
}

func New(err string, code int) APIError {
	return APIError{
		err:  err,
		code: code,
	}
}

var (
	ErrInvalidJSON       = New("invalid json", http.StatusBadRequest)
	ErrInternal          = New("internal error", http.StatusInternalServerError)
	ErrUserAlreadyExists = New("user already exist", http.StatusBadRequest)
	ErrUserNotExists     = New("user not exist", http.StatusBadRequest)
	ErrWrongPassword     = New("wrong password", http.StatusBadRequest)
	ErrNotSessionHeader  = New("absent auth header", http.StatusBadRequest)
	ErrNotAuthService    = New("not auth service", http.StatusUnauthorized)
	ErrSessionExpired    = New("session expired", http.StatusBadRequest)
)
