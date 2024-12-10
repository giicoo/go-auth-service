package server

import (
	"context"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(r http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    "localhost:8080",
			Handler: r,
		},
	}
}

func (srv *Server) StartServer() error {
	return srv.httpServer.ListenAndServe()
}

func (srv *Server) ShutdownServer(ctx context.Context) error {
	return srv.httpServer.Shutdown(ctx)
}
