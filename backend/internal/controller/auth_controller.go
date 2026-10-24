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

func (controller *AuthController) Verify(c echo.Context) (err error) {
	logger.Info("AuthController: Verify")

	ctx := c.Request().Context()

	var verificationRequestParam entities.VerificationRequest
	if err := c.Bind(&verificationRequestParam); err != nil {
		return apperror.Wrap(apperror.InvalidParameter, err)
	}

	// trim white spaces and change email to lowercase
	verificationRequestParam.Email = strings.ToLower(strings.TrimSpace(verificationRequestParam.Email))

	if err := validator.ValidateStruct(verificationRequestParam); err != nil {
		return apperror.Wrap(apperror.InvalidParameter, err)
	}

	err = controller.authInteractor.VerifyCode(ctx, verificationRequestParam)
	if err != nil {
		return
	}

	return c.NoContent(http.StatusOK)
}

func (controller *AuthController) Resend(c echo.Context) (err error) {
	logger.Info("AuthController: Resend")

	ctx := c.Request().Context()

	var resendCodeRequestParam entities.ResendCodeRequest
	if err := c.Bind(&resendCodeRequestParam); err != nil {
		return apperror.Wrap(apperror.InvalidParameter, err)
	}

	resendCodeRequestParam.Email = strings.ToLower(strings.TrimSpace(resendCodeRequestParam.Email))

	if err := validator.ValidateStruct(resendCodeRequestParam); err != nil {
		return apperror.Wrap(apperror.InvalidParameter, err)
	}

	err = controller.authInteractor.ResendVerificationCode(ctx, resendCodeRequestParam.Email)
	if err != nil {
		return
	}

	return c.NoContent(http.StatusOK)
}

func (controller *AuthController) Login(c echo.Context) (err error) {
	logger.Info("AuthController: Login")

	ctx := c.Request().Context()

	var loginParameter entities.LoginParameter
	if err := c.Bind(&loginParameter); err != nil {
		return apperror.Wrap(apperror.InvalidParameter, err)
	}

	loginParameter.Email = strings.ToLower(strings.TrimSpace(loginParameter.Email))

	if err := validator.ValidateStruct(loginParameter); err != nil {
		return apperror.Wrap(apperror.InvalidParameter, err)
	}

	res, err := controller.authInteractor.Login(ctx, loginParameter)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, res)
}
