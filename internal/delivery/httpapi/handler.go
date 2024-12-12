package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/giicoo/go-auth-service/internal/services"
	"github.com/giicoo/go-auth-service/pkg/apiError"
)

type Handler struct {
	Services *services.Services
}

func NewHandler(services *services.Services) *Handler {
	return &Handler{
		Services: services,
	}
}

func httpError(w http.ResponseWriter, err error) {
	var apiErr apiError.APIError
	if errors.As(err, &apiErr) {
		http.Error(w, err.Error(), apiErr.Code())
	} else {
		http.Error(w, apiError.ErrInternal.Error(), apiError.ErrInternal.Code())
	}

}

func httpResponse(w http.ResponseWriter, r interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(r); err != nil {
		return err
	}
	return nil
}
