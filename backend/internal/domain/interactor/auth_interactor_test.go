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
	"backend/internal/util"
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
		PasswordHash: "testPasswordHash",
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

	testPasswordHash := "testPasswordHash"
	testVerificationCode := "testVerificationCode"
	fixedNow := time.Date(2026, 10, 10, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		prepareFunc func(
			mar *inputportMock.MockAuthRepositoryInputPort,
			mur *inputportMock.MockUserRepositoryInputPort,
			msve *apiInputportMock.MockSendVerificationMailInputPort,
			mHasher *utilMock.MockPasswordHasher,
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
				mHasher *utilMock.MockPasswordHasher,
				mUUID *utilMock.MockUUIDGenerator,
				mTime *utilMock.MockTimeProvider,
			) {
				mur.EXPECT().FindUserByEmail(ctx, testSignupParam.Email).Return(entities.Users{}, nil)
				mUUID.EXPECT().NewV7().Return(testUserId, nil)
				mHasher.EXPECT().Hash(testSignupParam.Password).Return(testPasswordHash, nil)
				mUUID.EXPECT().NewV7().Return(testVerificationCodeId, nil)
				msve.EXPECT().SendEmailWithVerificationCode(testSignupParam.Email).Return(testResendId, testVerificationCode, nil)
				mTime.EXPECT().Now().Return(fixedNow)
				mur.EXPECT().CreateUser(gomock.Any(), testNewUser).Return(nil)
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
				mHasher *utilMock.MockPasswordHasher,
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
				mHasher *utilMock.MockPasswordHasher,
				mUUID *utilMock.MockUUIDGenerator,
				mTime *utilMock.MockTimeProvider,
			) {
				mur.EXPECT().FindUserByEmail(ctx, testSignupParam.Email).Return(entities.Users{}, nil)
				mUUID.EXPECT().NewV7().Return("", apperror.Wrap(apperror.InternalError, errors.New("uuid generation error")))
			},
			wantedError: apperror.InternalError.Code,
		},
		{
			name: "when password hashing fails",
			prepareFunc: func(
				mar *inputportMock.MockAuthRepositoryInputPort,
				mur *inputportMock.MockUserRepositoryInputPort,
				msve *apiInputportMock.MockSendVerificationMailInputPort,
				mHasher *utilMock.MockPasswordHasher,
				mUUID *utilMock.MockUUIDGenerator,
				mTime *utilMock.MockTimeProvider,
			) {
				mur.EXPECT().FindUserByEmail(ctx, testSignupParam.Email).Return(entities.Users{}, nil)
				mUUID.EXPECT().NewV7().Return(testUserId, nil)
				mHasher.EXPECT().Hash(testSignupParam.Password).Return("", errors.New("hashing failed"))
			},
			wantedError: apperror.InternalError.Code,
		},
		{
			name: "when verification code id generation fails",
			prepareFunc: func(
				mar *inputportMock.MockAuthRepositoryInputPort,
				mur *inputportMock.MockUserRepositoryInputPort,
				msve *apiInputportMock.MockSendVerificationMailInputPort,
				mHasher *utilMock.MockPasswordHasher,
				mUUID *utilMock.MockUUIDGenerator,
				mTime *utilMock.MockTimeProvider,
			) {
				mur.EXPECT().FindUserByEmail(ctx, testSignupParam.Email).Return(entities.Users{}, nil)
				mUUID.EXPECT().NewV7().Return(testUserId, nil)
				mHasher.EXPECT().Hash(testSignupParam.Password).Return(testPasswordHash, nil)
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
				mHasher *utilMock.MockPasswordHasher,
				mUUID *utilMock.MockUUIDGenerator,
				mTime *utilMock.MockTimeProvider,
			) {
				mur.EXPECT().FindUserByEmail(ctx, testSignupParam.Email).Return(entities.Users{}, nil)
				mUUID.EXPECT().NewV7().Return(testUserId, nil)
				mHasher.EXPECT().Hash(testSignupParam.Password).Return(testPasswordHash, nil)
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
				mHasher *utilMock.MockPasswordHasher,
				mUUID *utilMock.MockUUIDGenerator,
				mTime *utilMock.MockTimeProvider,
			) {
				mur.EXPECT().FindUserByEmail(ctx, testSignupParam.Email).Return(entities.Users{}, nil)
				mUUID.EXPECT().NewV7().Return(testUserId, nil)
				mHasher.EXPECT().Hash(testSignupParam.Password).Return(testPasswordHash, nil)
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
				mHasher *utilMock.MockPasswordHasher,
				mUUID *utilMock.MockUUIDGenerator,
				mTime *utilMock.MockTimeProvider,
			) {
				mur.EXPECT().FindUserByEmail(ctx, testSignupParam.Email).Return(entities.Users{}, nil)
				mUUID.EXPECT().NewV7().Return(testUserId, nil)
				mHasher.EXPECT().Hash(testSignupParam.Password).Return(testPasswordHash, nil)
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
				mHasher *utilMock.MockPasswordHasher,
				mUUID *utilMock.MockUUIDGenerator,
				mTime *utilMock.MockTimeProvider,
			) {
				mur.EXPECT().FindUserByEmail(ctx, testSignupParam.Email).Return(entities.Users{}, nil)
				mUUID.EXPECT().NewV7().Return(testUserId, nil)
				mHasher.EXPECT().Hash(testSignupParam.Password).Return(testPasswordHash, nil)
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
				mHasher *utilMock.MockPasswordHasher,
				mUUID *utilMock.MockUUIDGenerator,
				mTime *utilMock.MockTimeProvider,
			) {
				mur.EXPECT().FindUserByEmail(ctx, testSignupParam.Email).Return(entities.Users{}, nil)
				mUUID.EXPECT().NewV7().Return(testUserId, nil)
				mHasher.EXPECT().Hash(testSignupParam.Password).Return(testPasswordHash, nil)
				mUUID.EXPECT().NewV7().Return(testVerificationCodeId, nil)
				msve.EXPECT().SendEmailWithVerificationCode(testSignupParam.Email).Return(testResendId, testVerificationCode, nil)
				mTime.EXPECT().Now().Return(fixedNow)
				mur.EXPECT().CreateUser(gomock.Any(), testNewUser).Return(nil)
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
				mHasher *utilMock.MockPasswordHasher,
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
				mHasher *utilMock.MockPasswordHasher,
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
			mockPasswordHasher := utilMock.NewMockPasswordHasher(ctrl)
			mockTokenGenerator := utilMock.NewMockTokenGenerator(ctrl)
			mockSessionRepository := inputportMock.NewMockSessionRepositoryInputPort(ctrl)
			mockTx := txMock.NewMockManager(ctrl)
			mockWithinTransaction(mockTx)

			tt.prepareFunc(mockAuthRepository, mockUserRepository, mockSendVerificationMailRepository, mockPasswordHasher, mockUUID, mockTimeProvider)

			interactor := NewAuthInteractor(mockAuthRepository, mockUserRepository, mockSendVerificationMailRepository, mockSessionRepository, mockPasswordHasher, mockTokenGenerator, mockUUID, mockTimeProvider, mockTx)

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

func TestAuthInteractor_VerifyCode(t *testing.T) {
	ctx := context.Background()
	testFixedTime := time.Date(2026, 10, 10, 12, 0, 0, 0, time.UTC).Add(10 * time.Minute)
	testVerificationCode := entities.VerificationCodes{
		VerificationCodeId: "f3354e2b-8205-4de3-a708-f230b3fba168",
		Code:               "123456",
		ResendEmailId:      "621a0b3c-b04a-40d3-ba65-dfd0b8f8c62a",
		Attempts:           0,
		ExpiresAt:          testFixedTime,
	}
	testInvalidVerificationCode := entities.VerificationCodes{
		VerificationCodeId: "f3354e2b-8205-4de3-a708-f230b3fba168",
		Code:               "133456",
		ResendEmailId:      "621a0b3c-b04a-40d3-ba65-dfd0b8f8c62a",
		Attempts:           0,
		ExpiresAt:          testFixedTime,
	}

	testVerificationRequest := entities.VerificationRequest{
		Email:            "example@email.com",
		VerificationCode: "123456",
	}

	testUser := entities.Users{
		UserID:       "771f9087-6a77-4405-84ea-b1c34ef7a952",
		FirstName:    "sabir",
		LastName:     "barahi",
		Email:        "example@email.com",
		PasswordHash: "hashedpassword",
		Verified:     false,
	}

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
			name: "success",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(ctx, testVerificationRequest.Email).Return(testUser, nil)
				mar.EXPECT().GetVerifiationCodeByEmail(ctx, testUser.UserID).Return(testVerificationCode, nil)
				mTime.EXPECT().Now().Return(testFixedTime)
				mur.EXPECT().UpdateUserVerification(ctx, testUser.UserID, true).Return(nil)
			},
		},
		{
			name: "when FindUserByEmail returns error",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(ctx, testVerificationRequest.Email).Return(entities.Users{}, apperror.Wrap(apperror.InternalError, errors.New("internal_error")))
			},
			wantedError: apperror.InternalError.Code,
		},
		{
			name: "when email is empty",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(ctx, testVerificationRequest.Email).Return(entities.Users{
					Email: "",
				}, nil)
			},
			wantedError: apperror.EmailNotRegistered.Code,
		},
		{
			name: "when verified user already exists",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(ctx, testVerificationRequest.Email).Return(entities.Users{
					Email:    "example@email.com",
					Verified: true,
				}, nil)
			},
			wantedError: apperror.UserAlreadyExists.Code,
		},
		{
			name: "when code already expired",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(ctx, testVerificationRequest.Email).Return(testUser, nil)
				mar.EXPECT().GetVerifiationCodeByEmail(ctx, testUser.UserID).Return(testVerificationCode, nil)
				mTime.EXPECT().Now().Return(testFixedTime.Add(time.Second))
			},
			wantedError: apperror.TokenAlreadyExpired.Code,
		},
		{
			name: "when attempts is 5 times",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(ctx, testVerificationRequest.Email).Return(testUser, nil)
				mar.EXPECT().GetVerifiationCodeByEmail(ctx, testUser.UserID).Return(entities.VerificationCodes{
					Attempts: 5,
				}, nil)
			},
			wantedError: apperror.TooManyAttempts.Code,
		},
		{
			name: "when code already expired",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(ctx, testVerificationRequest.Email).Return(testUser, nil)
				mar.EXPECT().GetVerifiationCodeByEmail(ctx, testUser.UserID).Return(testInvalidVerificationCode, nil)
				mTime.EXPECT().Now().Return(testFixedTime)
				mar.EXPECT().UpdateVerificationCodeAttempts(ctx, testInvalidVerificationCode.VerificationCodeId, testInvalidVerificationCode.Attempts+1).Return(nil)
			},
			wantedError: apperror.InvalidVerificationCode.Code,
		},
		{
			name: "when GetVerifiationCodeByEmail returns error",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(ctx, testVerificationRequest.Email).Return(testUser, nil)
				mar.EXPECT().GetVerifiationCodeByEmail(ctx, testUser.UserID).Return(entities.VerificationCodes{}, apperror.Wrap(apperror.InternalError, errors.New("internal_error")))
			},
			wantedError: apperror.InternalError.Code,
		},
		{
			name: "when UpdateVerificationCodeAttempts returns error",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(ctx, testVerificationRequest.Email).Return(testUser, nil)
				mar.EXPECT().GetVerifiationCodeByEmail(ctx, testUser.UserID).Return(testInvalidVerificationCode, nil)
				mTime.EXPECT().Now().Return(testFixedTime)
				mar.EXPECT().UpdateVerificationCodeAttempts(ctx, testInvalidVerificationCode.VerificationCodeId, testInvalidVerificationCode.Attempts+1).Return(apperror.Wrap(apperror.InternalError, errors.New("internal_error")))
			},
			wantedError: apperror.InternalError.Code,
		},
		{
			name: "when UpdateUserVerification returns error",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(ctx, testVerificationRequest.Email).Return(testUser, nil)
				mar.EXPECT().GetVerifiationCodeByEmail(ctx, testUser.UserID).Return(testVerificationCode, nil)
				mTime.EXPECT().Now().Return(testFixedTime)
				mur.EXPECT().UpdateUserVerification(ctx, testUser.UserID, true).Return(apperror.Wrap(apperror.InternalError, errors.New("internal_error")))
			},
			wantedError: apperror.InternalError.Code,
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
			mockPasswordHasher := utilMock.NewMockPasswordHasher(ctrl)
			mockTokenGenerator := utilMock.NewMockTokenGenerator(ctrl)
			mockSessionRepository := inputportMock.NewMockSessionRepositoryInputPort(ctrl)
			mockTx := txMock.NewMockManager(ctrl)
			mockWithinTransaction(mockTx)

			tt.prepareFunc(mockAuthRepository, mockUserRepository, mockSendVerificationMailRepository, mockUUID, mockTimeProvider)

			interactor := NewAuthInteractor(mockAuthRepository, mockUserRepository, mockSendVerificationMailRepository, mockSessionRepository, mockPasswordHasher, mockTokenGenerator, mockUUID, mockTimeProvider, mockTx)

			err := interactor.VerifyCode(ctx, testVerificationRequest)

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

func TestAuthInteractor_ResendVerificationCode(t *testing.T) {
	ctx := context.Background()

	testUserId := "f79a7721-f7a4-441c-9136-a4db768c9bc3"
	testVerificationCodeId := "5fb7436d-b021-43cc-9dbf-231b2169164a"
	testResendId := "9d950124-25f8-4b41-b88b-5aea038684f4"
	testNewResendId := "c5fd626d-4b1d-4092-a651-e367174a8edd"
	testFixedTime := time.Date(2026, 10, 10, 12, 0, 0, 0, time.UTC)
	testNewVerificationCode := "121245"
	// testVerificationCode below expires at testFixedTime, so it was sent lifetime earlier
	testSentAt := testFixedTime.Add(-verificationCodeLifetime)

	testNewUser := entities.Users{
		UserID:       testUserId,
		FirstName:    "sabir",
		LastName:     "barahi",
		Email:        "example@email.com",
		PasswordHash: "",
		Verified:     false,
	}

	testVerificationCode := entities.VerificationCodes{
		VerificationCodeId: testVerificationCodeId,
		UserId:             testUserId,
		Code:               "123456",
		ResendEmailId:      testResendId,
		Attempts:           1,
		ExpiresAt:          testFixedTime,
	}

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
			name: "success",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(ctx, testNewUser.Email).Return(testNewUser, nil)
				mar.EXPECT().GetVerifiationCodeByEmail(ctx, testNewUser.UserID).Return(testVerificationCode, nil)
				mTime.EXPECT().Now().Return(testFixedTime)
				msve.EXPECT().SendEmailWithVerificationCode(testNewUser.Email).Return(testNewResendId, testNewVerificationCode, nil)
				mTime.EXPECT().Now().Return(testFixedTime)
				mar.EXPECT().UpdateVerificationCodeByVerificationCodeId(ctx, entities.VerificationCodes{
					VerificationCodeId: testVerificationCode.VerificationCodeId,
					UserId:             testVerificationCode.UserId,
					Code:               testNewVerificationCode,
					ResendEmailId:      testNewResendId,
					Attempts:           0,
					ExpiresAt:          testFixedTime.Add(verificationCodeLifetime),
				}).Return(nil)
			},
		},
		{
			name: "when internal error occurrs",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(ctx, testNewUser.Email).Return(entities.Users{}, apperror.Wrap(apperror.InternalError, errors.New("internal_error")))
			},
			wantedError: apperror.InternalError.Code,
		},
		{
			name: "when email is not a registered email",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(ctx, testNewUser.Email).Return(entities.Users{}, nil)
			},
			wantedError: apperror.EmailNotRegistered.Code,
		},
		{
			name: "when email is already verified",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(ctx, testNewUser.Email).Return(entities.Users{
					Email:    "example@email.com",
					Verified: true,
				}, nil)
			},
			wantedError: apperror.UserAlreadyExists.Code,
		},
		{
			name: "when resend is too soon",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(ctx, testNewUser.Email).Return(testNewUser, nil)
				mar.EXPECT().GetVerifiationCodeByEmail(ctx, testNewUser.UserID).Return(testVerificationCode, nil)
				mTime.EXPECT().Now().Return(testFixedTime.Add(20 * -time.Minute))
			},
			wantedError: apperror.ResendTooSoon.Code,
		},
		{
			name: "when resend is one second before the cooldown ends",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(ctx, testNewUser.Email).Return(testNewUser, nil)
				mar.EXPECT().GetVerifiationCodeByEmail(ctx, testNewUser.UserID).Return(testVerificationCode, nil)
				mTime.EXPECT().Now().Return(testSentAt.Add(resendCooldown - time.Second))
			},
			wantedError: apperror.ResendTooSoon.Code,
		},
		{
			name: "success when the cooldown has just ended",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(ctx, testNewUser.Email).Return(testNewUser, nil)
				mar.EXPECT().GetVerifiationCodeByEmail(ctx, testNewUser.UserID).Return(testVerificationCode, nil)
				mTime.EXPECT().Now().Return(testSentAt.Add(resendCooldown))
				msve.EXPECT().SendEmailWithVerificationCode(testNewUser.Email).Return(testNewResendId, testNewVerificationCode, nil)
				mTime.EXPECT().Now().Return(testSentAt.Add(resendCooldown))
				mar.EXPECT().UpdateVerificationCodeByVerificationCodeId(ctx, entities.VerificationCodes{
					VerificationCodeId: testVerificationCode.VerificationCodeId,
					UserId:             testVerificationCode.UserId,
					Code:               testNewVerificationCode,
					ResendEmailId:      testNewResendId,
					Attempts:           0,
					ExpiresAt:          testSentAt.Add(resendCooldown).Add(verificationCodeLifetime),
				}).Return(nil)
			},
		},
		{
			name: "success while the old code is still valid",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(ctx, testNewUser.Email).Return(testNewUser, nil)
				mar.EXPECT().GetVerifiationCodeByEmail(ctx, testNewUser.UserID).Return(testVerificationCode, nil)
				mTime.EXPECT().Now().Return(testSentAt.Add(5 * time.Minute))
				msve.EXPECT().SendEmailWithVerificationCode(testNewUser.Email).Return(testNewResendId, testNewVerificationCode, nil)
				mTime.EXPECT().Now().Return(testSentAt.Add(5 * time.Minute))
				mar.EXPECT().UpdateVerificationCodeByVerificationCodeId(ctx, entities.VerificationCodes{
					VerificationCodeId: testVerificationCode.VerificationCodeId,
					UserId:             testVerificationCode.UserId,
					Code:               testNewVerificationCode,
					ResendEmailId:      testNewResendId,
					Attempts:           0,
					ExpiresAt:          testSentAt.Add(5 * time.Minute).Add(verificationCodeLifetime),
				}).Return(nil)
			},
		},
		{
			name: "when getting the verification code fails",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(ctx, testNewUser.Email).Return(testNewUser, nil)
				mar.EXPECT().GetVerifiationCodeByEmail(ctx, testNewUser.UserID).Return(entities.VerificationCodes{}, apperror.Wrap(apperror.InternalError, errors.New("internal_error")))
			},
			wantedError: apperror.InternalError.Code,
		},
		{
			name: "when sending the verification email fails",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(ctx, testNewUser.Email).Return(testNewUser, nil)
				mar.EXPECT().GetVerifiationCodeByEmail(ctx, testNewUser.UserID).Return(testVerificationCode, nil)
				mTime.EXPECT().Now().Return(testFixedTime)
				msve.EXPECT().SendEmailWithVerificationCode(testNewUser.Email).Return("", "", new(apperror.ServiceUnavailable))
			},
			wantedError: apperror.ServiceUnavailable.Code,
		},
		{
			name: "when updating the verification code fails",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(ctx, testNewUser.Email).Return(testNewUser, nil)
				mar.EXPECT().GetVerifiationCodeByEmail(ctx, testNewUser.UserID).Return(testVerificationCode, nil)
				mTime.EXPECT().Now().Return(testFixedTime)
				msve.EXPECT().SendEmailWithVerificationCode(testNewUser.Email).Return(testNewResendId, testNewVerificationCode, nil)
				mTime.EXPECT().Now().Return(testFixedTime)
				mar.EXPECT().UpdateVerificationCodeByVerificationCodeId(ctx, gomock.Any()).Return(apperror.Wrap(apperror.InternalError, errors.New("internal_error")))
			},
			wantedError: apperror.InternalError.Code,
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
			mockPasswordHasher := utilMock.NewMockPasswordHasher(ctrl)
			mockTokenGenerator := utilMock.NewMockTokenGenerator(ctrl)
			mockSessionRepository := inputportMock.NewMockSessionRepositoryInputPort(ctrl)
			mockTx := txMock.NewMockManager(ctrl)
			mockWithinTransaction(mockTx)

			tt.prepareFunc(mockAuthRepository, mockUserRepository, mockSendVerificationMailRepository, mockUUID, mockTimeProvider)

			interactor := NewAuthInteractor(mockAuthRepository, mockUserRepository, mockSendVerificationMailRepository, mockSessionRepository, mockPasswordHasher, mockTokenGenerator, mockUUID, mockTimeProvider, mockTx)

			err := interactor.ResendVerificationCode(ctx, testNewUser.Email)

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

func TestAuthInteractor_Login(t *testing.T) {
	ctx := context.Background()

	testToken := "testToken"
	testExpiresAt := time.Date(2026, 10, 10, 12, 0, 0, 0, time.UTC)
	testUserId := "d1271cfd-c311-428a-9a90-03f16c3a23c7"
	testFirstName := "testFirstName"
	testLastName := "testLastName"
	testEmailAddress := "example@email.com"
	testPassword := "testpassword12345"

	testLoginParameter := entities.LoginParameter{
		Email:    testEmailAddress,
		Password: testPassword,
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(testLoginParameter.Password), bcrypt.MinCost)
	testPasswordHash := string(hash)
	testHashedToken := "testToken"

	testLoginUser := entities.LoginUser{
		UserId:    testUserId,
		FirstName: testFirstName,
		LastName:  testLastName,
		Email:     testEmailAddress,
	}

	testLoginResponse := entities.LoginResponse{
		Token:     testToken,
		ExpiresAt: testExpiresAt.Add(720 * time.Hour),
		User:      testLoginUser,
	}

	testUser := entities.Users{
		UserID:       testUserId,
		FirstName:    testFirstName,
		LastName:     testLastName,
		Email:        testEmailAddress,
		PasswordHash: testPasswordHash,
		Verified:     true,
	}

	testSessionId := "acfdef27-4b7c-4714-87e9-afb4940748da"

	testSession := entities.Sessions{
		SessionId: testSessionId,
		UserId:    testUserId,
		TokenHash: testHashedToken,
		ExpiresAt: testExpiresAt.Add(720 * time.Hour),
	}

	tests := []struct {
		name        string
		prepareFunc func(
			mar *inputportMock.MockAuthRepositoryInputPort,
			mur *inputportMock.MockUserRepositoryInputPort,
			msr *inputportMock.MockSessionRepositoryInputPort,
			msve *apiInputportMock.MockSendVerificationMailInputPort,
			mph *utilMock.MockPasswordHasher,
			mtg *utilMock.MockTokenGenerator,
			mUUID *utilMock.MockUUIDGenerator,
			mTime *utilMock.MockTimeProvider,
		)
		wantedLoginResponse entities.LoginResponse
		wantedError         string
	}{
		{
			name: "success",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msr *inputportMock.MockSessionRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mph *utilMock.MockPasswordHasher, mtg *utilMock.MockTokenGenerator, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(gomock.Any(), testLoginParameter.Email).Return(testUser, nil)
				mph.EXPECT().Compare(testUser.PasswordHash, testLoginParameter.Password).Return(nil)
				mtg.EXPECT().GenerateToken().Return(testToken, testHashedToken, nil)
				mUUID.EXPECT().NewV7().Return(testSessionId, nil)
				mTime.EXPECT().Now().Return(testExpiresAt)
				msr.EXPECT().CreateSession(gomock.Any(), testSession).Return(nil)
			},
			wantedLoginResponse: testLoginResponse,
			wantedError:         "",
		},
		{
			name: "when user is not found",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msr *inputportMock.MockSessionRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mph *utilMock.MockPasswordHasher, mtg *utilMock.MockTokenGenerator, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(gomock.Any(), testLoginParameter.Email).Return(entities.Users{}, nil)
			},
			wantedLoginResponse: entities.LoginResponse{},
			wantedError:         apperror.InvalidCredentials.Code,
		},
		{
			name: "when user is not verified",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msr *inputportMock.MockSessionRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mph *utilMock.MockPasswordHasher, mtg *utilMock.MockTokenGenerator, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(gomock.Any(), testLoginParameter.Email).Return(entities.Users{
					Email:    testEmailAddress,
					Verified: false,
				}, nil)
			},
			wantedLoginResponse: entities.LoginResponse{},
			wantedError:         apperror.EmailNotVerified.Code,
		},
		{
			name: "when invalid password",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msr *inputportMock.MockSessionRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mph *utilMock.MockPasswordHasher, mtg *utilMock.MockTokenGenerator, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(gomock.Any(), testLoginParameter.Email).Return(testUser, nil)
				mph.EXPECT().Compare(testUser.PasswordHash, testLoginParameter.Password).Return(util.ErrPasswordMismatch)
			},
			wantedLoginResponse: entities.LoginResponse{},
			wantedError:         apperror.InvalidCredentials.Code,
		},
		{
			name: "when the stored password hash is broken",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msr *inputportMock.MockSessionRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mph *utilMock.MockPasswordHasher, mtg *utilMock.MockTokenGenerator, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(gomock.Any(), testLoginParameter.Email).Return(testUser, nil)
				mph.EXPECT().Compare(testUser.PasswordHash, testLoginParameter.Password).Return(errors.New("hash too short"))
			},
			wantedLoginResponse: entities.LoginResponse{},
			wantedError:         apperror.InternalError.Code,
		},
		{
			name: "when finding the user fails",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msr *inputportMock.MockSessionRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mph *utilMock.MockPasswordHasher, mtg *utilMock.MockTokenGenerator, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(gomock.Any(), testLoginParameter.Email).Return(entities.Users{}, apperror.Wrap(apperror.InternalError, errors.New("db error")))
			},
			wantedLoginResponse: entities.LoginResponse{},
			wantedError:         apperror.InternalError.Code,
		},
		{
			name: "when generating the token fails",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msr *inputportMock.MockSessionRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mph *utilMock.MockPasswordHasher, mtg *utilMock.MockTokenGenerator, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(gomock.Any(), testLoginParameter.Email).Return(testUser, nil)
				mph.EXPECT().Compare(testUser.PasswordHash, testLoginParameter.Password).Return(nil)
				mtg.EXPECT().GenerateToken().Return("", "", apperror.Wrap(apperror.InternalError, errors.New("rand error")))
			},
			wantedLoginResponse: entities.LoginResponse{},
			wantedError:         apperror.InternalError.Code,
		},
		{
			name: "when generating the session id fails",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msr *inputportMock.MockSessionRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mph *utilMock.MockPasswordHasher, mtg *utilMock.MockTokenGenerator, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(gomock.Any(), testLoginParameter.Email).Return(testUser, nil)
				mph.EXPECT().Compare(testUser.PasswordHash, testLoginParameter.Password).Return(nil)
				mtg.EXPECT().GenerateToken().Return(testToken, testHashedToken, nil)
				mUUID.EXPECT().NewV7().Return("", errors.New("uuid error"))
			},
			wantedLoginResponse: entities.LoginResponse{},
			wantedError:         apperror.InternalError.Code,
		},
		{
			name: "when creating the session fails",
			prepareFunc: func(mar *inputportMock.MockAuthRepositoryInputPort, mur *inputportMock.MockUserRepositoryInputPort, msr *inputportMock.MockSessionRepositoryInputPort, msve *apiInputportMock.MockSendVerificationMailInputPort, mph *utilMock.MockPasswordHasher, mtg *utilMock.MockTokenGenerator, mUUID *utilMock.MockUUIDGenerator, mTime *utilMock.MockTimeProvider) {
				mur.EXPECT().FindUserByEmail(gomock.Any(), testLoginParameter.Email).Return(testUser, nil)
				mph.EXPECT().Compare(testUser.PasswordHash, testLoginParameter.Password).Return(nil)
				mtg.EXPECT().GenerateToken().Return(testToken, testHashedToken, nil)
				mUUID.EXPECT().NewV7().Return(testSessionId, nil)
				mTime.EXPECT().Now().Return(testExpiresAt)
				msr.EXPECT().CreateSession(gomock.Any(), testSession).Return(apperror.Wrap(apperror.InternalError, errors.New("db error")))
			},
			wantedLoginResponse: entities.LoginResponse{},
			wantedError:         apperror.InternalError.Code,
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
			mockPasswordHasher := utilMock.NewMockPasswordHasher(ctrl)
			mockTokenGenerator := utilMock.NewMockTokenGenerator(ctrl)
			mockSessionRepository := inputportMock.NewMockSessionRepositoryInputPort(ctrl)
			mockTx := txMock.NewMockManager(ctrl)
			mockWithinTransaction(mockTx)

			tt.prepareFunc(mockAuthRepository, mockUserRepository, mockSessionRepository, mockSendVerificationMailRepository, mockPasswordHasher, mockTokenGenerator, mockUUID, mockTimeProvider)

			interactor := NewAuthInteractor(mockAuthRepository, mockUserRepository, mockSendVerificationMailRepository, mockSessionRepository, mockPasswordHasher, mockTokenGenerator, mockUUID, mockTimeProvider, mockTx)

			loginResponse, err := interactor.Login(ctx, testLoginParameter)

			if tt.wantedError == "" {
				assert.NoError(t, err)
			} else {
				appErr, ok := errors.AsType[*apperror.AppError](err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantedError, appErr.Code)
			}

			assert.Equal(t, tt.wantedLoginResponse, loginResponse)
		})
	}
}
