package httpapi

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/create-user", h.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/delete-user", h.DeleteUser).Methods(http.MethodDelete)
	r.HandleFunc("/update-user", h.UpdateUser).Methods(http.MethodPut)
	r.HandleFunc("/get-user-by-email/{email}", h.GetUserByEmail).Methods(http.MethodGet)
	r.HandleFunc("/get-user-by-id/{id}", h.GetUserByID).Methods(http.MethodGet)
	//TODO: session auth with middleware token check
	return r
}
