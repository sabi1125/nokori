package interactor

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"backend/internal/config"
	"backend/internal/domain/apperror"
	"backend/internal/domain/entities"
	"backend/internal/domain/repository/inputport"
	"backend/internal/emailtemplate"
	logger "backend/internal/log"
	"backend/internal/tx"
	"backend/internal/util"

	"github.com/resend/resend-go/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthInteractor struct {
	authRepository inputport.AuthRepositoryInputPort
	userRepository inputport.UserRepositoryInputPort
	resendConfig   *config.ResendConfig
	uuidGenerator  util.UUIDGenerator
	timeProvider   util.TimeProvider
	txManager      tx.Manager
}

func NewAuthInteractor(
	authRepository inputport.AuthRepositoryInputPort,
	userRepository inputport.UserRepositoryInputPort,
	resendConfig *config.ResendConfig,
	uuidGenerator util.UUIDGenerator,
	timeProvider util.TimeProvider,
	txManager tx.Manager,
) *AuthInteractor {
	return &AuthInteractor{
		authRepository: authRepository,
		userRepository: userRepository,
		resendConfig:   resendConfig,
		uuidGenerator:  uuidGenerator,
		timeProvider:   timeProvider,
		txManager:      txManager,
	}
}

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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupParams.Password), bcrypt.DefaultCost)
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

	verificationCode, err := interactor.createVerificationCode()
	if err != nil {
		return
	}

	resendId, err := interactor.sendEmailWithVerificationCode(newUser.Email, verificationCode)
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

const verificationCodeLifetime = time.Minute

func (interactor *AuthInteractor) sendEmailWithVerificationCode(email string, verificationCode string) (resendSendEmailId string, err error) {
	logger.Info("AuthInteractor: sendEmailWithVerificationCode")
	client := resend.NewClient(interactor.resendConfig.APIKey)

	content, err := emailtemplate.VerificationCode(emailtemplate.VerificationCodeData{
		Code:         verificationCode,
		ValidMinutes: int(verificationCodeLifetime.Minutes()),
	})
	if err != nil {
		return "", apperror.Wrap(apperror.InternalError, err)
	}

	params := &resend.SendEmailRequest{
		From:    interactor.resendConfig.EmailFrom,
		To:      []string{email},
		Subject: content.Subject,
		Html:    content.HTML,
		Text:    content.Text,
		ReplyTo: interactor.resendConfig.ReplyTo,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		return "", apperror.Wrap(apperror.ServiceUnavailable, err)
	}

	resendSendEmailId = sent.Id

	return
}

func (interactor *AuthInteractor) createVerificationCode() (verificationCode string, err error) {
	max := big.NewInt(1000000)

	nBig, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", apperror.Wrap(apperror.InternalError, err)
	}

	num := nBig.Int64()
	buf := make([]byte, 6)

	for i := 5; i >= 0; i-- {
		buf[i] = '0' + byte(num%10) // '0' is ASCII 48. Adding the remainder gives the digit character.
		num /= 10
	}

	return string(buf), nil
}
