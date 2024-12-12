package httpapi

import "github.com/gorilla/mux"

func (h *Handler) NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/create-user", h.CreateUser).Methods("POST")
	return r
}
