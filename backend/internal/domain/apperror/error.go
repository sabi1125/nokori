package apperror

import "net/http"

var (
	InvalidParameter        = AppError{Code: "invalid_parameters", Status: http.StatusBadRequest, Message: "Request includes invalid parameters."}
	EmailNotVerified        = AppError{Code: "email_not_verified", Status: http.StatusConflict, Message: "Account for the provided email address has not been verified."}
	UserAlreadyExists       = AppError{Code: "user_already_exists", Status: http.StatusConflict, Message: "Account for the provided email address is already registered."}
	InternalError           = AppError{Code: "internal_error", Status: http.StatusInternalServerError, Message: "Something went wrong. Please wait and try again."}
	ServiceUnavailable      = AppError{Code: "service_unavailable", Status: http.StatusServiceUnavailable, Message: "Some services are unreachable."}
	EmailNotRegistered      = AppError{Code: "email_not_registered", Status: http.StatusNotFound, Message: "Provided email address is not a registered email address."}
	TokenAlreadyExpired     = AppError{Code: "token_already_expired", Status: http.StatusUnauthorized, Message: "Verification code has expired. Please request a new one."}
	InvalidVerificationCode = AppError{Code: "invalid_verification_code", Status: http.StatusBadRequest, Message: "Verification code is incorrect."}
	TooManyAttempts         = AppError{Code: "too_many_attempts", Status: http.StatusTooManyRequests, Message: "Too many wrong attempts. Please request a new code."}
)
