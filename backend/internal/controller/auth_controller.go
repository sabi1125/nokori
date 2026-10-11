package controller

import (
	"net/http"
	"strings"

	"backend/internal/domain/apperror"
	"backend/internal/domain/entities"
	"backend/internal/domain/interactor/inputport"
	logger "backend/internal/log"
	"backend/internal/validator"

	"github.com/labstack/echo/v4"
)

type AuthController struct {
	authInteractor inputport.AuthInteractorInputPort
}

func NewAuthController(authInteractor inputport.AuthInteractorInputPort) *AuthController {
	return &AuthController{
		authInteractor: authInteractor,
	}
}

// Signup
func (controller *AuthController) SignUp(c echo.Context) error {
	logger.Info("authController: Signup")

	ctx := c.Request().Context()

	var signUpParams entities.SignUp
	if err := c.Bind(&signUpParams); err != nil {
		return apperror.Wrap(apperror.InvalidParameter, err)
	}

	// trim white spaces and change email to lowercase
	trimmedStruct := trimSpaces(&signUpParams)

	if err := validator.ValidateStruct(trimmedStruct); err != nil {
		return apperror.Wrap(apperror.InvalidParameter, err)
	}

	err := controller.authInteractor.SignUp(ctx, signUpParams)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func trimSpaces(param *entities.SignUp) entities.SignUp {
	param.Email = strings.ToLower(strings.TrimSpace(param.Email))
	param.FirstName = strings.TrimSpace(param.FirstName)
	param.LastName = strings.TrimSpace(param.LastName)

	return *param
}
