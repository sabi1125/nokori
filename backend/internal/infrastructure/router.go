package infrastructure

import (
	"backend/internal/config"
	"backend/internal/controller"
	apirepository "backend/internal/domain/api_repository"
	"backend/internal/domain/interactor"
	"backend/internal/domain/repository"
	"backend/internal/tx"
	"backend/internal/util"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Router(e *echo.Echo, db *gorm.DB) {
	txManager := NewTransactionManager(db)

	// shared repositories
	userRepository := repository.NewUserRepository(db)

	// config
	resendConfig := config.LoadResendConfigFromEnv()

	// utils
	timeProvider := util.NewTimeProvider()
	uuidGenerator := util.NewUUIDGenerator()

	RegisteredHealthRouter(e, db)
	RegisteredAuthRouter(e, db, userRepository, resendConfig, uuidGenerator, timeProvider, txManager)
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
	uuidGenerator util.UUIDGenerator,
	timeProvider util.TimeProvider,
	txManager tx.Manager,
) {
	auth := e.Group("/auth")

	repository := repository.NewAuthRepository(db)
	apiRepository := apirepository.NewSendVerificationMailRepository(resendConfig)
	interactor := interactor.NewAuthInteractor(repository, userRepository, apiRepository, uuidGenerator, timeProvider, txManager)
	controller := controller.NewAuthController(interactor)

	auth.POST("/signup", controller.SignUp)
}
