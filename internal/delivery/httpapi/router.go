package httpapi

import (
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (h *Handler) NewRouter() *mux.Router {
	rout := mux.NewRouter()
	r := rout.Methods("GET", "POST", "PUT", "DELETE").Subrouter()
	r.Use(h.MiddlewareGetSessionHeader)

	r.HandleFunc("/create-user", h.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/get-user-by-id/{id}", h.GetUserByID).Methods(http.MethodGet)
	r.HandleFunc("/delete-user", h.DeleteUser).Methods(http.MethodDelete)
	r.HandleFunc("/update-email", h.UpdateEmailUser).Methods(http.MethodPut)
	r.HandleFunc("/update-password", h.UpdatePasswordUser).Methods(http.MethodPut)
	r.HandleFunc("/check-user", h.CheckUser).Methods(http.MethodPost)

	r.HandleFunc("/create-session", h.CreateSession).Methods(http.MethodPost)
	r.HandleFunc("/get-session/{session_id}", h.GetSession).Methods(http.MethodGet)
	r.HandleFunc("/get-sessions/{user_id}", h.GetListSession).Methods(http.MethodGet)
	r.HandleFunc("/delete-session", h.DeleteSession).Methods(http.MethodDelete)
	r.HandleFunc("/delete-sessions", h.DeleteListSession).Methods(http.MethodDelete)

	// swagger docs
	rout.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)
	return rout
}
