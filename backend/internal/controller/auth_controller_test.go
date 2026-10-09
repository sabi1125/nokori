package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
