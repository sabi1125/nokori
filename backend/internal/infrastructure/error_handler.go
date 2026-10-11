package infrastructure

import (
	"errors"
	"net/http"

	"backend/internal/domain/apperror"
	logger "backend/internal/log"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	if appErr, ok := errors.AsType[*apperror.AppError](err); ok {
		if appErr.Err != nil {
			logger.Error(appErr.Err)
		}
		c.JSON(appErr.Status, map[string]string{"code": appErr.Code, "message": appErr.Message})
		return
	}

	if echoErr, ok := errors.AsType[*echo.HTTPError](err); ok {
		switch internalEchoErr := echoErr.Code; internalEchoErr {
		case 400:
			c.JSON(echoErr.Code, map[string]any{"code": "bad_request", "message": echoErr.Message})
		case 404:
			c.JSON(echoErr.Code, map[string]any{"code": "route_not_found", "message": echoErr.Message})
		case 405:
			c.JSON(echoErr.Code, map[string]any{"code": "not_allowed", "message": echoErr.Message})
		default:
			c.JSON(echoErr.Code, map[string]any{"code": "unknown", "message": echoErr.Message})
		}
		return
	}

	logger.Error(err)
	c.JSON(http.StatusInternalServerError, map[string]string{"code": apperror.InternalError.Code, "message": apperror.InternalError.Message})
}
