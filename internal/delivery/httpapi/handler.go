package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/giicoo/go-auth-service/internal/entity/models"
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
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	var apiErr apiError.APIError
	if errors.As(err, &apiErr) {
		response := models.ErrorResponse{
			Error: err.Error(),
		}
		w.WriteHeader(apiErr.Code())
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, apiError.ErrInternal.Error(), apiError.ErrInternal.Code())
		}

	} else {
		response := models.ErrorResponse{
			Error: apiError.ErrInternal.Error(),
		}
		w.WriteHeader(apiError.ErrInternal.Code())
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, apiError.ErrInternal.Error(), apiError.ErrInternal.Code())
		}

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
