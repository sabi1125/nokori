package infrastructure

import (
	"backend/internal/config"
	"backend/internal/controller"
	"backend/internal/domain/interactor"
	"backend/internal/domain/repository"
	"backend/internal/tx"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Router(e *echo.Echo, db *gorm.DB) {
	txManager := NewTransactionManager(db)

	// shared repositories
	userRepository := repository.NewUserRepository(db)

	// config
	resendConfig := config.LoadResendConfigFromEnv()

	RegisteredHealthRouter(e, db)
	RegisteredAuthRouter(e, db, userRepository, resendConfig, txManager)
}

func RegisteredHealthRouter(e *echo.Echo, db *gorm.DB) {
	health := e.Group("/health")
	repository := repository.NewHealthRepository(db)
	interactor := interactor.NewHealthInteractor(repository)
	controller := controller.NewHealthController(interactor)

	health.GET("", controller.Health)
}

func RegisteredAuthRouter(
	e *echo.Echo,
	db *gorm.DB,
	userRepository *repository.UserRepository,
	resendConfig *config.ResendConfig,
	txManager tx.Manager,
) {
	auth := e.Group("/auth")

	repository := repository.NewAuthRepository(db)
	interactor := interactor.NewAuthInteractor(repository, userRepository, resendConfig, txManager)
	controller := controller.NewAuthController(interactor)

	auth.POST("/signup", controller.SignUp)
}
