package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"backend/internal/domain/apperror"
	"backend/internal/domain/entities"
	interactorMock "backend/internal/domain/interactor/inputport/mock"
	logger "backend/internal/log"
	"backend/internal/validator"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestMain(m *testing.M) {
	logger.InitForTest()
	validator.Init()
	m.Run()
}

func TestAuthController_SignUp(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		mockSetup    func(m *interactorMock.MockAuthInteractorInputPort)
		expectedCode string
	}{
		{
			name: "Success",
			body: `{"first_name":"sabir", "last_name":"barahi", "email":"example@gmail.com", "password":"testpassword12345"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().SignUp(gomock.Any(), entities.SignUp{
					FirstName: "sabir",
					LastName:  "barahi",
					Email:     "example@gmail.com",
					Password:  "testpassword12345",
				}).Return(nil)
			},
			expectedCode: "",
		},
		{
			name: "Success with white spaces in first name, last name and email and change to lower case for email",
			body: `{"first_name":" sabir ", "last_name":" barahi ", "email":" EXAMPLE@gmail.com", "password":"testpassword12345"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().SignUp(gomock.Any(), entities.SignUp{
					FirstName: "sabir",
					LastName:  "barahi",
					Email:     "example@gmail.com",
					Password:  "testpassword12345",
				}).Return(nil)
			},
			expectedCode: "",
		},
		{
			name: "Success when email has multiple dots",
			body: `{"first_name":"sabir", "last_name":"barahi", "email":"example.jo.hn@gmail.com", "password":"testpassword12345"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().SignUp(gomock.Any(), entities.SignUp{
					FirstName: "sabir",
					LastName:  "barahi",
					Email:     "example.jo.hn@gmail.com",
					Password:  "testpassword12345",
				}).Return(nil)
			},
			expectedCode: "",
		},
		{
			name: "Success when first name is short",
			body: `{"first_name":"Yu", "last_name":"barahi", "email":"example@gmail.com", "password":"testpassword12345"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().SignUp(gomock.Any(), entities.SignUp{
					FirstName: "Yu",
					LastName:  "barahi",
					Email:     "example@gmail.com",
					Password:  "testpassword12345",
				}).Return(nil)
			},
			expectedCode: "",
		},
		{
			name: "Success when last name is short",
			body: `{"first_name":"sabir", "last_name":"li", "email":"example@gmail.com", "password":"testpassword12345"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().SignUp(gomock.Any(), entities.SignUp{
					FirstName: "sabir",
					LastName:  "li",
					Email:     "example@gmail.com",
					Password:  "testpassword12345",
				}).Return(nil)
			},
			expectedCode: "",
		},
		{
			name: "Success when first name 255chars",
			body: `{"first_name":"sabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabir", "last_name":"barahi", "email":"example@gmail.com", "password":"testpassword12345"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().SignUp(gomock.Any(), entities.SignUp{
					FirstName: "sabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabir",
					LastName:  "barahi",
					Email:     "example@gmail.com",
					Password:  "testpassword12345",
				}).Return(nil)
			},
			expectedCode: "",
		},
		{
			name: "Success when last name 255chars",
			body: `{"first_name":"sabir", "last_name":"sabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabir", "email":"example@gmail.com", "password":"testpassword12345"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().SignUp(gomock.Any(), entities.SignUp{
					FirstName: "sabir",
					LastName:  "sabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabir",
					Email:     "example@gmail.com",
					Password:  "testpassword12345",
				}).Return(nil)
			},
			expectedCode: "",
		},
		{
			name: "Success when password is 72chars long",
			body: `{"first_name":"sabir", "last_name":"barahi", "email":"example@gmail.com", "password":"test123456test123456test123456test123456test123456test123456test12345671"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().SignUp(gomock.Any(), entities.SignUp{
					FirstName: "sabir",
					LastName:  "barahi",
					Email:     "example@gmail.com",
					Password:  "test123456test123456test123456test123456test123456test123456test12345671",
				}).Return(nil)
			},
			expectedCode: "",
		},
		{
			name: "Success when email is 255chars long",
			body: `{"first_name":"sabir", "last_name":"barahi", "email":"sabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabir@gmail.com", "password":"testpassword12345"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().SignUp(gomock.Any(), entities.SignUp{
					FirstName: "sabir",
					LastName:  "barahi",
					Email:     "sabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabir@gmail.com",
					Password:  "testpassword12345",
				}).Return(nil)
			},
			expectedCode: "",
		},
		{
			name:         "fail when first name exceeds 255chars",
			body:         `{"first_name":"sabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabir", "last_name":"barahi", "email":"example@gmail.com", "password":"testpassword12345"}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name:         "fail when last name exceeds 255chars",
			body:         `{"first_name":"sabir", "last_name":"sabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabir", "email":"example@gmail.com", "password":"testpassword12345"}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name:         "when password exceeds 72chars",
			body:         `{"first_name":"sabir", "last_name":"barahi", "email":"example@gmail.com", "password":"testtest123456test123456test123456test123456test123456test123456test12345671"}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name:         "when email exceeds 255chars",
			body:         `{"first_name":"sabir", "last_name":"barahi", "email":"sabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabir@gmail.com", "password":"testpassword12345"}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name:         "when firstname is not provided",
			body:         `{"first_name":"", "last_name":"barahi", "email":"example@gmail.com", "password":"testpassword12345"}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name:         "when lastname is not provided",
			body:         `{"first_name":"sabir", "last_name":"", "email":"example@gmail.com", "password":"testpassword12345"}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name:         "when email is not provided",
			body:         `{"first_name":"sabir", "last_name":"barahi", "email":"", "password":"testpassword12345"}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name:         "when password is not provided",
			body:         `{"first_name":"sabir", "last_name":"barahi", "email":"email@gmail.com", "password":""}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name:         "when email is not valid format",
			body:         `{"first_name":"sabir", "last_name":"barahi", "email":"emailgmail.com", "password":"testpassword12345"}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name:         "when first name is only white space",
			body:         `{"first_name":"      ", "last_name":"barahi", "email":"email@gmail.com", "password":"testpassword12345"}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name:         "when last name is only white space",
			body:         `{"first_name":"sabir", "last_name":"      ", "email":"email@gmail.com", "password":"testpassword12345"}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name:         "when pasword is not ASCII",
			body:         `{"first_name":"sabir", "last_name":"barahi", "email":"email@gmail.com", "password":"テストパスワード"}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name:         "when request is broken",
			body:         `{"first_name""sabir", "last_name":"barahi", "email":"email@gmail.com", "password":"testpassword12345"}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name: "when user already registered",
			body: `{"first_name":"sabir", "last_name":"barahi", "email":"example@gmail.com", "password":"testpassword12345"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().SignUp(gomock.Any(), entities.SignUp{
					FirstName: "sabir",
					LastName:  "barahi",
					Email:     "example@gmail.com",
					Password:  "testpassword12345",
				}).Return(new(apperror.UserAlreadyExists))
			},
			expectedCode: apperror.UserAlreadyExists.Code,
		},
		{
			name: "when user already registered but not verified",
			body: `{"first_name":"sabir", "last_name":"barahi", "email":"example@gmail.com", "password":"testpassword12345"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().SignUp(gomock.Any(), entities.SignUp{
					FirstName: "sabir",
					LastName:  "barahi",
					Email:     "example@gmail.com",
					Password:  "testpassword12345",
				}).Return(new(apperror.EmailNotVerified))
			},
			expectedCode: apperror.EmailNotVerified.Code,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthInteractor := interactorMock.NewMockAuthInteractorInputPort(ctrl)
			tt.mockSetup(mockAuthInteractor)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/auth/signup", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			controller := NewAuthController(mockAuthInteractor)

			err := controller.SignUp(c)

			if tt.expectedCode == "" {
				assert.NoError(t, err)
			} else {
				appErr, ok := errors.AsType[*apperror.AppError](err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedCode, appErr.Code)
			}
		})
	}
}

func TestAuthController_Verify(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		mockSetup    func(m *interactorMock.MockAuthInteractorInputPort)
		expectedCode string
	}{
		{
			name: "success",
			body: `{"email":"example@email.com", "verification_code":"123456"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().VerifyCode(gomock.Any(), entities.VerificationRequest{
					Email:            "example@email.com",
					VerificationCode: "123456",
				}).Return(nil)
			},
			expectedCode: "",
		},
		{
			name:         "fail when verification code is not numbers",
			body:         `{"email":"example@email.com", "verification_code":"abcdef"}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name:         "fail when verification code is shorter that 6",
			body:         `{"email":"example@email.com", "verification_code":"12345"}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name:         "fail when verification code is longer that 6",
			body:         `{"email":"example@email.com", "verification_code":"1234567"}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name: "success when email is 255chars",
			body: `{"email":"sabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabir@gmail.com", "verification_code":"123456"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().VerifyCode(gomock.Any(), entities.VerificationRequest{
					Email:            "sabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabir@gmail.com",
					VerificationCode: "123456",
				}).Return(nil)
			},
			expectedCode: "",
		},
		{
			name:         "fail when email is longer than 255chars",
			body:         `{"email":"sabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabir@gmail.com", "verification_code":"123456"}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name: "fails with internal error",
			body: `{"email":"example@email.com", "verification_code":"123456"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().VerifyCode(gomock.Any(), entities.VerificationRequest{
					Email:            "example@email.com",
					VerificationCode: "123456",
				}).Return(apperror.Wrap(apperror.InternalError, errors.New("internal_error")))
			},
			expectedCode: apperror.InternalError.Code,
		},
		{
			name: "fails when email is not a registered email",
			body: `{"email":"example@email.com", "verification_code":"123456"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().VerifyCode(gomock.Any(), entities.VerificationRequest{
					Email:            "example@email.com",
					VerificationCode: "123456",
				}).Return(new(apperror.EmailNotRegistered))
			},
			expectedCode: apperror.EmailNotRegistered.Code,
		},
		{
			name: "fails when user already registered and verified",
			body: `{"email":"example@email.com", "verification_code":"123456"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().VerifyCode(gomock.Any(), entities.VerificationRequest{
					Email:            "example@email.com",
					VerificationCode: "123456",
				}).Return(new(apperror.UserAlreadyExists))
			},
			expectedCode: apperror.UserAlreadyExists.Code,
		},
		{
			name: "fails when verification code is expired",
			body: `{"email":"example@email.com", "verification_code":"123456"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().VerifyCode(gomock.Any(), entities.VerificationRequest{
					Email:            "example@email.com",
					VerificationCode: "123456",
				}).Return(new(apperror.TokenAlreadyExpired))
			},
			expectedCode: apperror.TokenAlreadyExpired.Code,
		},
		{
			name: "fails when verification code is invalid",
			body: `{"email":"example@email.com", "verification_code":"123456"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().VerifyCode(gomock.Any(), entities.VerificationRequest{
					Email:            "example@email.com",
					VerificationCode: "123456",
				}).Return(new(apperror.InvalidVerificationCode))
			},
			expectedCode: apperror.InvalidVerificationCode.Code,
		},
		{
			name: "fails when too many attempts",
			body: `{"email":"example@email.com", "verification_code":"123456"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().VerifyCode(gomock.Any(), entities.VerificationRequest{
					Email:            "example@email.com",
					VerificationCode: "123456",
				}).Return(new(apperror.TooManyAttempts))
			},
			expectedCode: apperror.TooManyAttempts.Code,
		},
		{
			name: "success with email trimmed and lowercased",
			body: `{"email":"  Example@Email.com ", "verification_code":"123456"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().VerifyCode(gomock.Any(), entities.VerificationRequest{
					Email:            "example@email.com",
					VerificationCode: "123456",
				}).Return(nil)
			},
			expectedCode: "",
		},
		{
			name:         "fail when email is not a valid email",
			body:         `{"email":"not-an-email", "verification_code":"123456"}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name:         "fail when email is missing",
			body:         `{"verification_code":"123456"}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name:         "fail when verification code is missing",
			body:         `{"email":"example@email.com"}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name:         "fail when body is malformed json",
			body:         `{"email":`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthInteractor := interactorMock.NewMockAuthInteractorInputPort(ctrl)
			tt.mockSetup(mockAuthInteractor)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/auth/verify", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			controller := NewAuthController(mockAuthInteractor)

			err := controller.Verify(c)

			if tt.expectedCode == "" {
				assert.NoError(t, err)
			} else {
				appErr, ok := errors.AsType[*apperror.AppError](err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedCode, appErr.Code)
			}
		})
	}
}

func TestAuthController_Resend(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		mockSetup    func(m *interactorMock.MockAuthInteractorInputPort)
		expectedCode string
	}{
		{
			name: "success",
			body: `{"email":"example@email.com"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().ResendVerificationCode(gomock.Any(), "example@email.com").Return(nil)
			},
		},
		{
			name: "success when email has whitespace",
			body: `{"email":" example@email.com "}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().ResendVerificationCode(gomock.Any(), "example@email.com").Return(nil)
			},
		},
		{
			name:         "when email is not valid",
			body:         `{"email":" example"}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name:         "when email is more than 255chars long",
			body:         `{"email":"sabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabir@gmail.com"}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name: "success when email is 255chars long",
			body: `{"email":"sabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabir@gmail.com"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().ResendVerificationCode(gomock.Any(), "sabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabirsabir@gmail.com").Return(nil)
			},
			expectedCode: "",
		},
		{
			name: "when ResendVerificationCode returns internal error",
			body: `{"email":"example@email.com"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().ResendVerificationCode(gomock.Any(), "example@email.com").Return(apperror.Wrap(apperror.InternalError, errors.New("internal_error")))
			},
			expectedCode: apperror.InternalError.Code,
		},
		{
			name: "when email is not a registered email",
			body: `{"email":"example@email.com"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().ResendVerificationCode(gomock.Any(), "example@email.com").Return(new(apperror.EmailNotRegistered))
			},
			expectedCode: apperror.EmailNotRegistered.Code,
		},
		{
			name: "when user for the email address already exists and is already a valid user",
			body: `{"email":"example@email.com"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().ResendVerificationCode(gomock.Any(), "example@email.com").Return(new(apperror.UserAlreadyExists))
			},
			expectedCode: apperror.UserAlreadyExists.Code,
		},
		{
			name: "when resend code request is too soon and previous code is still valid",
			body: `{"email":"example@email.com"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().ResendVerificationCode(gomock.Any(), "example@email.com").Return(new(apperror.ResendTooSoon))
			},
			expectedCode: apperror.ResendTooSoon.Code,
		},
		{
			name: "when sending the verification email fails",
			body: `{"email":"example@email.com"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().ResendVerificationCode(gomock.Any(), "example@email.com").Return(new(apperror.ServiceUnavailable))
			},
			expectedCode: apperror.ServiceUnavailable.Code,
		},
		{
			name: "success with email trimmed and lowercased",
			body: `{"email":" Example@EMAIL.com "}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().ResendVerificationCode(gomock.Any(), "example@email.com").Return(nil)
			},
		},
		{
			name:         "fail when email is missing",
			body:         `{}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name:         "fail when body is malformed json",
			body:         `{"email":`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthInteractor := interactorMock.NewMockAuthInteractorInputPort(ctrl)
			tt.mockSetup(mockAuthInteractor)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/auth/resend", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			controller := NewAuthController(mockAuthInteractor)

			err := controller.Resend(c)

			if tt.expectedCode == "" {
				assert.NoError(t, err)
			} else {
				appErr, ok := errors.AsType[*apperror.AppError](err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedCode, appErr.Code)
			}
		})
	}
}

func TestAuthController_Login(t *testing.T) {
	testResponse := entities.LoginResponse{
		Token:     "test-session-token",
		ExpiresAt: time.Date(2026, 11, 21, 12, 0, 0, 0, time.UTC),
		User: entities.LoginUser{
			UserId:    "f79a7721-f7a4-441c-9136-a4db768c9bc3",
			FirstName: "sabir",
			LastName:  "barahi",
			Email:     "example@email.com",
		},
	}
	validParams := entities.LoginParameter{Email: "example@email.com", Password: "password123"}

	tests := []struct {
		name         string
		body         string
		mockSetup    func(m *interactorMock.MockAuthInteractorInputPort)
		expectedCode string
	}{
		{
			name: "success",
			body: `{"email":"example@email.com","password":"password123"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().Login(gomock.Any(), validParams).Return(testResponse, nil)
			},
		},
		{
			name: "success with email trimmed and lowercased",
			body: `{"email":" Example@EMAIL.com ","password":"password123"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().Login(gomock.Any(), validParams).Return(testResponse, nil)
			},
		},
		{
			name:         "fail when email is not valid",
			body:         `{"email":"example","password":"password123"}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name:         "fail when email is longer than 255chars",
			body:         `{"email":"` + strings.Repeat("sabir", 50) + `@gmail.com","password":"password123"}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name:         "fail when email is missing",
			body:         `{"password":"password123"}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name:         "fail when password is missing",
			body:         `{"email":"example@email.com"}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name:         "fail when password is shorter than 8",
			body:         `{"email":"example@email.com","password":"pass123"}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name:         "fail when password is longer than 72",
			body:         `{"email":"example@email.com","password":"` + strings.Repeat("a", 73) + `"}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name:         "fail when password is not ASCII",
			body:         `{"email":"example@email.com","password":"のこりpassword"}`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name:         "fail when body is malformed json",
			body:         `{"email":`,
			mockSetup:    func(m *interactorMock.MockAuthInteractorInputPort) {},
			expectedCode: apperror.InvalidParameter.Code,
		},
		{
			name: "fail when credentials are invalid",
			body: `{"email":"example@email.com","password":"password123"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().Login(gomock.Any(), validParams).Return(entities.LoginResponse{}, new(apperror.InvalidCredentials))
			},
			expectedCode: apperror.InvalidCredentials.Code,
		},
		{
			name: "fail when email is not verified",
			body: `{"email":"example@email.com","password":"password123"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().Login(gomock.Any(), validParams).Return(entities.LoginResponse{}, new(apperror.EmailNotVerified))
			},
			expectedCode: apperror.EmailNotVerified.Code,
		},
		{
			name: "fail with internal error",
			body: `{"email":"example@email.com","password":"password123"}`,
			mockSetup: func(m *interactorMock.MockAuthInteractorInputPort) {
				m.EXPECT().Login(gomock.Any(), validParams).Return(entities.LoginResponse{}, apperror.Wrap(apperror.InternalError, errors.New("internal_error")))
			},
			expectedCode: apperror.InternalError.Code,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthInteractor := interactorMock.NewMockAuthInteractorInputPort(ctrl)
			tt.mockSetup(mockAuthInteractor)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			controller := NewAuthController(mockAuthInteractor)

			err := controller.Login(c)

			if tt.expectedCode == "" {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, rec.Code)
				// the body must match the spec's field names
				assert.JSONEq(t, `{
					"token": "test-session-token",
					"expires_at": "2026-11-21T12:00:00Z",
					"user": {
						"user_id": "f79a7721-f7a4-441c-9136-a4db768c9bc3",
						"first_name": "sabir",
						"last_name": "barahi",
						"email": "example@email.com"
					}
				}`, rec.Body.String())
			} else {
				appErr, ok := errors.AsType[*apperror.AppError](err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedCode, appErr.Code)
			}
		})
	}
}
