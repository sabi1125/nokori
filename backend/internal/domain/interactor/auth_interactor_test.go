package interactor

import (
	"context"
	"errors"
	"testing"
	"time"

	apiInputportMock "backend/internal/domain/api_repository/inputport/mock"
	"backend/internal/domain/apperror"
	"backend/internal/domain/entities"
	inputportMock "backend/internal/domain/repository/inputport/mock"
	logger "backend/internal/log"
	txMock "backend/internal/tx/mock"
	utilMock "backend/internal/util/mock"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestMain(m *testing.M) {
	logger.InitForTest()
	m.Run()
}

func mockWithinTransaction(tx *txMock.MockManager) {
	tx.EXPECT().
		WithinTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		}).
		AnyTimes()
}

func TestAuthInteractor_SignUp(t *testing.T) {
	ctx := context.Background()

	testSignupParam := entities.SignUp{
		FirstName: "sabir",
		LastName:  "barahi",
		Email:     "example@gmail.com",
		Password:  "testpassword12345",
	}

	testUserId := "f79a7721-f7a4-441c-9136-a4db768c9bc3"
	testVerificationCodeId := "5fb7436d-b021-43cc-9dbf-231b2169164a"
	testResendId := "9d950124-25f8-4b41-b88b-5aea038684f4"

	testNewUser := entities.Users{
		UserID:       testUserId,
		FirstName:    "sabir",
		LastName:     "barahi",
		Email:        "example@gmail.com",
		PasswordHash: "",
		Verified:     false,
	}
	testVerifiedUser := entities.Users{
		UserID:       testUserId,
		FirstName:    "sabir",
		LastName:     "barahi",
		Email:        "example@gmail.com",
		PasswordHash: "",
		Verified:     true,
	}
	testUnVerifiedUser := entities.Users{
		UserID:       testUserId,
		FirstName:    "sabir",
		LastName:     "barahi",
		Email:        "example@gmail.com",
		PasswordHash: "",
		Verified:     false,
	}

	testVerificationCode := "testVerificationCode"
	fixedNow := time.Date(2026, 10, 10, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		prepareFunc func(
			mar *inputportMock.MockAuthRepositoryInputPort,
			mur *inputportMock.MockUserRepositoryInputPort,
			msve *apiInputportMock.MockSendVerificationMailInputPort,
			mUUID *utilMock.MockUUIDGenerator,
			mTime *utilMock.MockTimeProvider,
		)
		wantedError string
	}{
		{
			name: "Success",
			prepareFunc: func(
				mar *inputportMock.MockAuthRepositoryInputPort,
				mur *inputportMock.MockUserRepositoryInputPort,
				msve *apiInputportMock.MockSendVerificationMailInputPort,
				mUUID *utilMock.MockUUIDGenerator,
				mTime *utilMock.MockTimeProvider,
			) {
				mur.EXPECT().FindUserByEmail(ctx, testSignupParam.Email).Return(entities.Users{}, nil)
				mUUID.EXPECT().NewV7().Return(testUserId, nil)
				mUUID.EXPECT().NewV7().Return(testVerificationCodeId, nil)
				msve.EXPECT().SendEmailWithVerificationCode(testSignupParam.Email).Return(testResendId, testVerificationCode, nil)
				mTime.EXPECT().Now().Return(fixedNow)
				mur.EXPECT().CreateUser(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, user entities.Users) error {
						err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash),
							[]byte(testSignupParam.Password))
						assert.NoError(t, err)

						user.PasswordHash = ""
						assert.Equal(t, testNewUser, user)
						return nil
					})
				mar.EXPECT().InsertVerificationCode(ctx, entities.VerificationCodes{
					VerificationCodeId: testVerificationCodeId,
					UserId:             testUserId,
					Code:               testVerificationCode,
					ResendEmailId:      testResendId,
					ExpiresAt:          fixedNow.Add(verificationCodeLifetime),
				})
			},
			wantedError: "",
		},
		{
			name: "fail with internal error while finding user",
			prepareFunc: func(
				mar *inputportMock.MockAuthRepositoryInputPort,
				mur *inputportMock.MockUserRepositoryInputPort,
				msve *apiInputportMock.MockSendVerificationMailInputPort,
				mUUID *utilMock.MockUUIDGenerator,
				mTime *utilMock.MockTimeProvider,
			) {
				mur.EXPECT().FindUserByEmail(ctx, testSignupParam.Email).Return(entities.Users{}, apperror.Wrap(apperror.InternalError, errors.New("internal error")))
			},
			wantedError: apperror.InternalError.Code,
		},
		{
			name: "when user id generation fails",
			prepareFunc: func(
				mar *inputportMock.MockAuthRepositoryInputPort,
				mur *inputportMock.MockUserRepositoryInputPort,
				msve *apiInputportMock.MockSendVerificationMailInputPort,
				mUUID *utilMock.MockUUIDGenerator,
				mTime *utilMock.MockTimeProvider,
			) {
				mur.EXPECT().FindUserByEmail(ctx, testSignupParam.Email).Return(entities.Users{}, nil)
				mUUID.EXPECT().NewV7().Return("", apperror.Wrap(apperror.InternalError, errors.New("uuid generation error")))
			},
			wantedError: apperror.InternalError.Code,
		},
		{
			name: "when verification code id generation fails",
			prepareFunc: func(
				mar *inputportMock.MockAuthRepositoryInputPort,
				mur *inputportMock.MockUserRepositoryInputPort,
				msve *apiInputportMock.MockSendVerificationMailInputPort,
				mUUID *utilMock.MockUUIDGenerator,
				mTime *utilMock.MockTimeProvider,
			) {
				mur.EXPECT().FindUserByEmail(ctx, testSignupParam.Email).Return(entities.Users{}, nil)
				mUUID.EXPECT().NewV7().Return(testUserId, nil)
				mUUID.EXPECT().NewV7().Return("", apperror.Wrap(apperror.InternalError, errors.New("uuid generation error")))
			},
			wantedError: apperror.InternalError.Code,
		},
		{
			name: "when sending email with verification code fails to create verification code",
			prepareFunc: func(
				mar *inputportMock.MockAuthRepositoryInputPort,
				mur *inputportMock.MockUserRepositoryInputPort,
				msve *apiInputportMock.MockSendVerificationMailInputPort,
				mUUID *utilMock.MockUUIDGenerator,
				mTime *utilMock.MockTimeProvider,
			) {
				mur.EXPECT().FindUserByEmail(ctx, testSignupParam.Email).Return(entities.Users{}, nil)
				mUUID.EXPECT().NewV7().Return(testUserId, nil)
				mUUID.EXPECT().NewV7().Return(testVerificationCodeId, nil)
				msve.EXPECT().SendEmailWithVerificationCode(testSignupParam.Email).Return("", "", apperror.Wrap(apperror.InternalError, errors.New("failed to create verification code")))
			},
			wantedError: apperror.InternalError.Code,
		},
		{
			name: "when sending email with verification code fails to create mail content",
			prepareFunc: func(
				mar *inputportMock.MockAuthRepositoryInputPort,
				mur *inputportMock.MockUserRepositoryInputPort,
				msve *apiInputportMock.MockSendVerificationMailInputPort,
				mUUID *utilMock.MockUUIDGenerator,
				mTime *utilMock.MockTimeProvider,
			) {
				mur.EXPECT().FindUserByEmail(ctx, testSignupParam.Email).Return(entities.Users{}, nil)
				mUUID.EXPECT().NewV7().Return(testUserId, nil)
				mUUID.EXPECT().NewV7().Return(testVerificationCodeId, nil)
				msve.EXPECT().SendEmailWithVerificationCode(testSignupParam.Email).Return("", "", apperror.Wrap(apperror.InternalError, errors.New("failed to create verification code")))
			},
			wantedError: apperror.InternalError.Code,
		},
		{
			name: "when sending email with verification code fails to send",
			prepareFunc: func(
				mar *inputportMock.MockAuthRepositoryInputPort,
				mur *inputportMock.MockUserRepositoryInputPort,
				msve *apiInputportMock.MockSendVerificationMailInputPort,
				mUUID *utilMock.MockUUIDGenerator,
				mTime *utilMock.MockTimeProvider,
			) {
				mur.EXPECT().FindUserByEmail(ctx, testSignupParam.Email).Return(entities.Users{}, nil)
				mUUID.EXPECT().NewV7().Return(testUserId, nil)
				mUUID.EXPECT().NewV7().Return(testVerificationCodeId, nil)
				msve.EXPECT().SendEmailWithVerificationCode(testSignupParam.Email).Return("", "", apperror.Wrap(apperror.ServiceUnavailable, errors.New("failed to send verification email")))
			},
			wantedError: apperror.ServiceUnavailable.Code,
		},
		{
			name: "when failed to create user",
			prepareFunc: func(
				mar *inputportMock.MockAuthRepositoryInputPort,
				mur *inputportMock.MockUserRepositoryInputPort,
				msve *apiInputportMock.MockSendVerificationMailInputPort,
				mUUID *utilMock.MockUUIDGenerator,
				mTime *utilMock.MockTimeProvider,
			) {
				mur.EXPECT().FindUserByEmail(ctx, testSignupParam.Email).Return(entities.Users{}, nil)
				mUUID.EXPECT().NewV7().Return(testUserId, nil)
				mUUID.EXPECT().NewV7().Return(testVerificationCodeId, nil)
				msve.EXPECT().SendEmailWithVerificationCode(testSignupParam.Email).Return(testResendId, testVerificationCode, nil)
				mTime.EXPECT().Now().Return(fixedNow)
				mur.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(apperror.Wrap(apperror.InternalError, errors.New("failed to create user")))
			},
			wantedError: apperror.InternalError.Code,
		},
		{
			name: "failed to create verification code record",
			prepareFunc: func(
				mar *inputportMock.MockAuthRepositoryInputPort,
				mur *inputportMock.MockUserRepositoryInputPort,
				msve *apiInputportMock.MockSendVerificationMailInputPort,
				mUUID *utilMock.MockUUIDGenerator,
				mTime *utilMock.MockTimeProvider,
			) {
				mur.EXPECT().FindUserByEmail(ctx, testSignupParam.Email).Return(entities.Users{}, nil)
				mUUID.EXPECT().NewV7().Return(testUserId, nil)
				mUUID.EXPECT().NewV7().Return(testVerificationCodeId, nil)
				msve.EXPECT().SendEmailWithVerificationCode(testSignupParam.Email).Return(testResendId, testVerificationCode, nil)
				mTime.EXPECT().Now().Return(fixedNow)
				mur.EXPECT().CreateUser(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, user entities.Users) error {
						err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash),
							[]byte(testSignupParam.Password))
						assert.NoError(t, err)

						user.PasswordHash = ""
						assert.Equal(t, testNewUser, user)
						return nil
					})
				mar.EXPECT().InsertVerificationCode(ctx, entities.VerificationCodes{
					VerificationCodeId: testVerificationCodeId,
					UserId:             testUserId,
					Code:               testVerificationCode,
					ResendEmailId:      testResendId,
					ExpiresAt:          fixedNow.Add(verificationCodeLifetime),
				}).Return(apperror.Wrap(apperror.InternalError, errors.New("failed to record verification code")))
			},
			wantedError: apperror.InternalError.Code,
		},
		{
			name: "when user exists and verified",
			prepareFunc: func(
				mar *inputportMock.MockAuthRepositoryInputPort,
				mur *inputportMock.MockUserRepositoryInputPort,
				msve *apiInputportMock.MockSendVerificationMailInputPort,
				mUUID *utilMock.MockUUIDGenerator,
				mTime *utilMock.MockTimeProvider,
			) {
				mur.EXPECT().FindUserByEmail(ctx, testSignupParam.Email).Return(testVerifiedUser, nil)
			},
			wantedError: apperror.UserAlreadyExists.Code,
		},
		{
			name: "when user exists but not verified",
			prepareFunc: func(
				mar *inputportMock.MockAuthRepositoryInputPort,
				mur *inputportMock.MockUserRepositoryInputPort,
				msve *apiInputportMock.MockSendVerificationMailInputPort,
				mUUID *utilMock.MockUUIDGenerator,
				mTime *utilMock.MockTimeProvider,
			) {
				mur.EXPECT().FindUserByEmail(ctx, testSignupParam.Email).Return(testUnVerifiedUser, nil)
			},
			wantedError: apperror.EmailNotVerified.Code,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthRepository := inputportMock.NewMockAuthRepositoryInputPort(ctrl)
			mockUserRepository := inputportMock.NewMockUserRepositoryInputPort(ctrl)
			mockSendVerificationMailRepository := apiInputportMock.NewMockSendVerificationMailInputPort(ctrl)
			mockUUID := utilMock.NewMockUUIDGenerator(ctrl)
			mockTimeProvider := utilMock.NewMockTimeProvider(ctrl)
			mockTx := txMock.NewMockManager(ctrl)
			mockWithinTransaction(mockTx)

			tt.prepareFunc(mockAuthRepository, mockUserRepository, mockSendVerificationMailRepository, mockUUID, mockTimeProvider)

			interactor := NewAuthInteractor(mockAuthRepository, mockUserRepository, mockSendVerificationMailRepository, mockUUID, mockTimeProvider, mockTx)

			err := interactor.SignUp(ctx, testSignupParam)

			if tt.wantedError == "" {
				assert.NoError(t, err)
			} else {
				appErr, ok := errors.AsType[*apperror.AppError](err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantedError, appErr.Code)
			}
		})
	}
}
