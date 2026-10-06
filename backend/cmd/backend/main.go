package main

import (
	"backend/internal/config"
	"backend/internal/infrastructure"
	"backend/internal/log"
	"backend/internal/validator"

	"github.com/labstack/echo/v4"
)

func main() {
	zapCfg := config.LoadZapConfig()
	logger.Init(zapCfg)
	defer logger.Sync()
	validator.Init()

	dbCfg := config.LoadDbConfig()
	db := infrastructure.Connection(dbCfg)

	// Create an Echo instance
	e := echo.New()
	e.Use(logger.MiddlewareLogger(logger.Get()))
	e.HTTPErrorHandler = infrastructure.ErrorHandler
	infrastructure.Router(e, db)

	// Start the server
	logger.Info("starting server")
	e.Logger.Fatal(e.Start(":8080"))
}
