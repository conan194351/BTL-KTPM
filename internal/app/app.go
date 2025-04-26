package app

import (
	"fmt"
	"github.com/conan194351/BTL-KTPM/internal/config"
	"github.com/conan194351/BTL-KTPM/internal/routes"
	"github.com/conan194351/BTL-KTPM/pkg/logger"
	"net/http"
)

type App struct {
	logger logger.Logger
}

func New() *App {
	return &App{
		logger: logger.NewZapLogger("App", true),
	}
}

func (app *App) Start() {
	fmt.Println("\033[31m\r\n" + `
	 | |/ /__   __|  __ \|  \/  |
	 | ' /   | |  | |__) | \  / |
	 |  <    | |  |  ___/| |\/| |
	 | . \   | |  | |    | |  | |
	 |_|\_\  |_|  |_|    |_|  |_|
	` + "\033[m")
	app.logger.Info("App is running...", nil)
	config.InitDatabase()
	r := routes.InitRoutes()
	addr := config.GetConfig().Server.GetAddr()
	server := &http.Server{
		Addr:        addr,
		Handler:     r,
		ReadTimeout: config.GetConfig().Server.GetReadTimeout(),
	}
	app.logger.Info(fmt.Sprintf("Server is running on %s", addr), nil)
	app.logger.Fatal(server.ListenAndServe(), "Server is shutting down...", nil)
}
