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
	passwordHasher := util.NewPasswordHasher()
	tokenGenerator := util.NewTokenGenerator()

	RegisteredAuthRouter(e, db, userRepository, resendConfig, passwordHasher, tokenGenerator, uuidGenerator, timeProvider, txManager)
}

func RegisteredAuthRouter(
	e *echo.Echo,
	db *gorm.DB,
	userRepository *repository.UserRepository,
	resendConfig *config.ResendConfig,
	passwordHasher util.PasswordHasher,
	tokenGenerator util.TokenGenerator,
	uuidGenerator util.UUIDGenerator,
	timeProvider util.TimeProvider,
	txManager tx.Manager,
) {
	auth := e.Group("/auth")

	sessionRepository := repository.NewSessionRepository(db)
	repository := repository.NewAuthRepository(db)
	apiRepository := apirepository.NewSendVerificationMailRepository(resendConfig)
	interactor := interactor.NewAuthInteractor(repository, userRepository, apiRepository, sessionRepository, passwordHasher, tokenGenerator, uuidGenerator, timeProvider, txManager)
	controller := controller.NewAuthController(interactor)

	auth.POST("/signup", controller.SignUp)
	auth.POST("/verify", controller.Verify)
	auth.POST("/resend", controller.Resend)
	auth.POST("/login", controller.Login)
}
