package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/giicoo/go-auth-service/internal/config"
	"github.com/giicoo/go-auth-service/internal/delivery/httpapi"
	"github.com/giicoo/go-auth-service/internal/repository/sqlite"
	"github.com/giicoo/go-auth-service/internal/server"
	"github.com/giicoo/go-auth-service/internal/services"
	"github.com/giicoo/go-auth-service/pkg/prettylog"
)

func RunApp() {

	// init logger
	prettyHandler := prettylog.NewHandler(&slog.HandlerOptions{
		Level:       slog.LevelInfo,
		AddSource:   false,
		ReplaceAttr: nil,
	})
	logger := slog.New(prettyHandler)
	slog.SetDefault(logger)

	// Load Config
	cfg, err := config.LoadConfig("./config/config.yaml")
	if err != nil {
		slog.Error("", "err", err.Error())
		return
	}
	// Init layers
	repo := sqlite.NewRepo(cfg)
	services := services.NewServices(cfg, repo)
	handler := httpapi.NewHandler(services)

	// Start Server
	srv := server.NewServer(handler.Router)

	go func() {

		err := srv.StartServer()
		if err != nil {
			switch err {
			case http.ErrServerClosed:
				fmt.Println()
			default:
				log.Fatal(err)
			}
		}
	}()
	slog.Info("Server Start")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	// ShutDown Server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.ShutdownServer(ctx); err != nil {
		log.Fatal(err)
	} else {
		log.Print("Server stop.")
	}
}
