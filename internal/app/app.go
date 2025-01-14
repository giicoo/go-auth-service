package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/giicoo/go-auth-service/internal/config"
	"github.com/giicoo/go-auth-service/internal/delivery/httpapi"
	"github.com/giicoo/go-auth-service/internal/repository/sqlite"
	"github.com/giicoo/go-auth-service/internal/server"
	"github.com/giicoo/go-auth-service/internal/services"
	"github.com/giicoo/go-auth-service/pkg/beauti_json_formatter"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

func RunApp() {
	// init jwt
	//jwt.GenerateServiceToken() // debug

	// init logger
	formatter := beauti_json_formatter.NewFormatter(true)
	logrus.SetFormatter(formatter)
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	// Load Config
	cfg, err := config.LoadConfig("./config/config.yaml")
	if err != nil {
		logrus.Errorf("load config file: %s", err)
		return
	}
	logrus.Info("Logger and Config init")

	// Load DB
	db, err := sql.Open("sqlite3", cfg.DB.Path)
	if err != nil {
		logrus.Errorf("open db: %s", err)
		return
	}
	logrus.Info("DB connect")

	// Init layers
	repo := sqlite.NewRepo(cfg, db)
	services := services.NewServices(cfg, repo)
	handler := httpapi.NewHandler(services)

	// Init repo
	if err := repo.InitRepo(); err != nil {
		logrus.Errorf("init repo: %s", err)
		return
	}

	// Create router
	r := handler.NewRouter()
	h := cors.Default().Handler(r)
	// Start Server
	srv := server.NewServer(h)

	go func() {

		err := srv.StartServer()
		if err != nil {
			switch err {
			case http.ErrServerClosed:

			default:
				logrus.Errorf("start server: %s", err)
				return
			}
		}
	}()
	logrus.Info("Server Start")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
	fmt.Println()
	err = db.Close()
	if err != nil {
		logrus.Errorf("err with db close: %s", err)
		return
	}
	logrus.Info("DB close")
	// ShutDown Server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.ShutdownServer(ctx); err != nil {
		logrus.Errorf("shutdown server: %s", err)
		return
	} else {
		logrus.Info("Server stop")
	}
}
