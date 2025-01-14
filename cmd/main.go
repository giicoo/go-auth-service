package main

import (
	_ "github.com/giicoo/go-auth-service/docs"
	"github.com/giicoo/go-auth-service/internal/app"
)

// @title           go-auth-service
// @version         1.0
// @description     Сервис для ведения БД пользователей и контроля их сессий. Используются принципы чистой архитектуры. Доступ к сервису осуществляется только через запросы подписанные JWT токеном которые генерируется при запуске и сохраняется в ENV.

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

// @host      localhost:8080
// @BasePath  /
func main() {
	app.RunApp()
}
