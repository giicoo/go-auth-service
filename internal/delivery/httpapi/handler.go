package httpapi

import (
	"github.com/giicoo/go-auth-service/internal/services"
	"github.com/gorilla/mux"
)

type Handler struct {
	Router   *mux.Router
	Services *services.Services
}

func NewHandler(services *services.Services) *Handler {
	return &Handler{
		Router:   NewRouter(),
		Services: services,
	}
}
