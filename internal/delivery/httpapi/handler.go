package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/giicoo/go-auth-service/internal/entity"
	"github.com/giicoo/go-auth-service/internal/services"
)

type Handler struct {
	Services *services.Services
}

func NewHandler(services *services.Services) *Handler {
	return &Handler{
		Services: services,
	}
}

func httpError(w http.ResponseWriter, err error, code int) {
	switch err.(type) {
	case entity.PublicError:
		http.Error(w, err.Error(), code)
	case entity.PrivateError:
		http.Error(w, entity.ErrInternal.Error(), code)
	default:
		http.Error(w, entity.ErrInternal.Error(), code)
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
