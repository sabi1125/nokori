package interactor

import (
	"context"
	"crypto/subtle"
	"errors"
	"time"

	api_inputport "backend/internal/domain/api_repository/inputport"
	"backend/internal/domain/apperror"
	"backend/internal/domain/entities"
	"backend/internal/domain/repository/inputport"
	logger "backend/internal/log"
	"backend/internal/tx"
	"backend/internal/util"
)

type AuthInteractor struct {
	authRepository                 inputport.AuthRepositoryInputPort
	userRepository                 inputport.UserRepositoryInputPort
	sendVerificationMailRepository api_inputport.SendVerificationMailInputPort
	sessionRepository              inputport.SessionRepositoryInputPort
	passwordHasher                 util.PasswordHasher
	tokenGenerator                 util.TokenGenerator
	uuidGenerator                  util.UUIDGenerator
	timeProvider                   util.TimeProvider
	txManager                      tx.Manager
}

func NewAuthInteractor(
	authRepository inputport.AuthRepositoryInputPort,
	userRepository inputport.UserRepositoryInputPort,
	sendVerificationMailRepository api_inputport.SendVerificationMailInputPort,
	sessionRepository inputport.SessionRepositoryInputPort,
	passwordHasher util.PasswordHasher,
	tokenGenerator util.TokenGenerator,
	uuidGenerator util.UUIDGenerator,
	timeProvider util.TimeProvider,
	txManager tx.Manager,
) *AuthInteractor {
	return &AuthInteractor{
		authRepository:                 authRepository,
		userRepository:                 userRepository,
		sendVerificationMailRepository: sendVerificationMailRepository,
		passwordHasher:                 passwordHasher,
		tokenGenerator:                 tokenGenerator,
		sessionRepository:              sessionRepository,
		uuidGenerator:                  uuidGenerator,
		timeProvider:                   timeProvider,
		txManager:                      txManager,
	}
}

const (
	verificationCodeLifetime = 10 * time.Minute
	// a new code can be requested this long after the last one was sent,
	// even if that code is still valid
	resendCooldown = time.Minute
	// a session lasts this long and is extended by this much while it is used
	sessionLifetime = 30 * 24 * time.Hour
)

func (interactor *AuthInteractor) SignUp(ctx context.Context, signupParams entities.SignUp) (err error) {
	logger.Info("AuthInteractor: SignUp")

	user, err := interactor.userRepository.FindUserByEmail(ctx, signupParams.Email)
	if err != nil {
		return
	}

	if user.Email != "" && user.Verified {
		return new(apperror.UserAlreadyExists)
	}

	if user.Email != "" && !user.Verified {
		return new(apperror.EmailNotVerified)
	}

	newUserId, err := interactor.uuidGenerator.NewV7()
	if err != nil {
		return apperror.Wrap(apperror.InternalError, err)
	}

	hashedPassword, err := interactor.passwordHasher.Hash(signupParams.Password)
	if err != nil {
		return apperror.Wrap(apperror.InternalError, err)
	}

	newUser := entities.Users{
		UserID:       newUserId,
		FirstName:    signupParams.FirstName,
		LastName:     signupParams.LastName,
		Email:        signupParams.Email,
		PasswordHash: string(hashedPassword),
		Verified:     false,
	}

	newVerificationCodeId, err := interactor.uuidGenerator.NewV7()
	if err != nil {
		return apperror.Wrap(apperror.InternalError, err)
	}

	resendId, verificationCode, err := interactor.sendVerificationMailRepository.SendEmailWithVerificationCode(newUser.Email)
	if err != nil {
		return
	}

	expiryTime := interactor.timeProvider.Now().Add(verificationCodeLifetime)

	verificationCodeParam := entities.VerificationCodes{
		VerificationCodeId: newVerificationCodeId,
		UserId:             newUserId,
		Code:               verificationCode,
		ResendEmailId:      resendId,
		ExpiresAt:          expiryTime,
	}

	txErr := interactor.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := interactor.userRepository.CreateUser(ctx, newUser); err != nil {
			return err
		}

		if err := interactor.authRepository.InsertVerificationCode(ctx, verificationCodeParam); err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		err = txErr
		return
	}

	return
}

func (interactor *AuthInteractor) VerifyCode(ctx context.Context, verificationRequestParams entities.VerificationRequest) (err error) {
	logger.Info("AuthInteractor: VerifyCode")

	user, err := interactor.userRepository.FindUserByEmail(ctx, verificationRequestParams.Email)
	if err != nil {
		return
	}

	if user.Email == "" {
		return new(apperror.EmailNotRegistered)
	}

	if user.Verified {
		return new(apperror.UserAlreadyExists)
	}

	verificationCodes, err := interactor.authRepository.GetVerifiationCodeByEmail(ctx, user.UserID)
	if err != nil {
		return
	}

	// TODO: parallel requests for the same user can all read the same attempts value
	// before any of them increments it, so a burst can get more than 5 guesses.
	// fix later by locking the row (SELECT ... FOR UPDATE) or by incrementing
	// with a single conditional UPDATE (WHERE attempts < 5) before comparing.
	if verificationCodes.Attempts >= 5 {
		return new(apperror.TooManyAttempts)
	}

	beforeExpiry := verificationCodes.ExpiresAt.Before(interactor.timeProvider.Now())
	if beforeExpiry {
		return new(apperror.TokenAlreadyExpired)
	}

	// comparing hashed verification code and the verification code sent by the user
	if subtle.ConstantTimeCompare([]byte(verificationCodes.Code), []byte(verificationRequestParams.VerificationCode)) != 1 {
		// increment the attempts by 1
		if err := interactor.authRepository.UpdateVerificationCodeAttempts(ctx, verificationCodes.VerificationCodeId, verificationCodes.Attempts+1); err != nil {
			return err
		}
		return new(apperror.InvalidVerificationCode)
	}

	// changing the verification status of the user to true
	err = interactor.userRepository.UpdateUserVerification(ctx, user.UserID, true)
	if err != nil {
		return err
	}

	return
}

func (interactor *AuthInteractor) ResendVerificationCode(ctx context.Context, email string) (err error) {
	logger.Info("AuthInteractor: ResendVerificationCode")

	user, err := interactor.userRepository.FindUserByEmail(ctx, email)
	if err != nil {
		return
	}

	if user.Email == "" {
		return new(apperror.EmailNotRegistered)
	}

	if user.Verified {
		return new(apperror.UserAlreadyExists)
	}

	verificationCode, err := interactor.authRepository.GetVerifiationCodeByEmail(ctx, user.UserID)
	if err != nil {
		return
	}

	// the code was sent lifetime before it expires
	sentAt := verificationCode.ExpiresAt.Add(-verificationCodeLifetime)

	// TODO: parallel requests for the same user can all read the old expires_at
	// before any of them saves the new one, so each sends an email. the frontend
	// disables the resend button to prevent it (see decisions.md, verification-limits).
	// fix later by locking the row (SELECT ... FOR UPDATE) inside a transaction.
	if interactor.timeProvider.Now().Before(sentAt.Add(resendCooldown)) {
		return new(apperror.ResendTooSoon)
	}

	resendId, newVerificationCode, err := interactor.sendVerificationMailRepository.SendEmailWithVerificationCode(email)
	if err != nil {
		return
	}

	verificationCode.Code = newVerificationCode
	verificationCode.ExpiresAt = interactor.timeProvider.Now().Add(verificationCodeLifetime)
	verificationCode.Attempts = 0
	verificationCode.ResendEmailId = resendId

	err = interactor.authRepository.UpdateVerificationCodeByVerificationCodeId(ctx, verificationCode)
	if err != nil {
		return
	}

	return
}

func (interactor *AuthInteractor) Login(ctx context.Context, loginParameter entities.LoginParameter) (loginResponse entities.LoginResponse, err error) {
	logger.Info("AuthInteractor: Login")

	user, err := interactor.userRepository.FindUserByEmail(ctx, loginParameter.Email)
	if err != nil {
		return
	}

	if user.Email == "" {
		err = new(apperror.InvalidCredentials)
		return
	}

	if !user.Verified {
		err = new(apperror.EmailNotVerified)
		return
	}

	err = interactor.passwordHasher.Compare(user.PasswordHash, loginParameter.Password)
	if errors.Is(err, util.ErrPasswordMismatch) {
		err = new(apperror.InvalidCredentials)
		return
	}
	if err != nil {
		err = apperror.Wrap(apperror.InternalError, err)
		return
	}

	token, hash, err := interactor.tokenGenerator.GenerateToken()
	if err != nil {
		return
	}

	sessionId, err := interactor.uuidGenerator.NewV7()
	if err != nil {
		err = apperror.Wrap(apperror.InternalError, err)
		return
	}
	sessionExpiresAt := interactor.timeProvider.Now().Add(sessionLifetime)

	session := entities.Sessions{
		SessionId: sessionId,
		UserId:    user.UserID,
		TokenHash: hash,
		ExpiresAt: sessionExpiresAt,
	}

	err = interactor.sessionRepository.CreateSession(ctx, session)
	if err != nil {
		return
	}

	loginResponse = entities.LoginResponse{
		Token:     token,
		ExpiresAt: sessionExpiresAt,
		User: entities.LoginUser{
			UserId:    user.UserID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
	}

	return
}
