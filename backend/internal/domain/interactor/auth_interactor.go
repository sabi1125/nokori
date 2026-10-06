package interactor

import (
	"context"
	"crypto/rand"
	"math/big"

	"backend/internal/config"
	"backend/internal/domain/apperror"
	"backend/internal/domain/entities"
	"backend/internal/domain/repository/inputport"
	logger "backend/internal/log"
	"backend/internal/tx"
	"backend/internal/util"

	"github.com/resend/resend-go/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthInteractor struct {
	authRepository inputport.AuthRepositoryInputPort
	userRepository inputport.UserRepositoryInputPort

	resendConfig *config.ResendConfig
	txManager    tx.Manager
}

func NewAuthInteractor(
	authRepository inputport.AuthRepositoryInputPort,
	userRepository inputport.UserRepositoryInputPort,
	resendConfig *config.ResendConfig,
	txManager tx.Manager,
) *AuthInteractor {
	return &AuthInteractor{
		authRepository: authRepository,
		userRepository: userRepository,

		resendConfig: resendConfig,
		txManager:    txManager,
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

	uuid := util.NewUUIDGenerator()
	newUserId, err := uuid.NewV7()
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

	newVerificationCodeId, err := uuid.NewV7()
	if err != nil {
		return apperror.Wrap(apperror.InternalError, err)
	}

	verificationCode, err := interactor.createVerificationCode()
	if err != nil {
		return
	}

	resendSendEmailId, err := interactor.sendEmailWithVerificationCode(ctx, newUser.Email, newVerificationCodeId, verificationCode)
	if err != nil {
		return
	}

	// TODO: create datetime package in util because we need to create expires_at for the verification code.

	txErr := interactor.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := interactor.userRepository.CreateUser(ctx, newUser); err != nil {
			return err
		}

		// TODO: Verification code that is sent should be inserted to the database in the same transaction.
		// if saving of the verification code fails the whole transaction should rollback because this is the most import part to check if the user's email is valid

		// TODO: another problem is that we should save the the id returned from resend sdk in the database. so we need to fix the er diagram
		// fix the er diagram before creating the struct for the verification code

		return nil
	})

	if txErr != nil {
		err = txErr
		return
	}

	return
}

func (interactor *AuthInteractor) sendEmailWithVerificationCode(ctx context.Context, email string, uuid string, verificationCode string) (resendSendEmailId string, err error) {
	logger.Info("AuthInteractor: sendEmailWithVerificationCode")
	client := resend.NewClient(interactor.resendConfig.APIKey)

	// TODO: need to fix the message we are going to send to the user's over all
	params := &resend.SendEmailRequest{
		From:    interactor.resendConfig.EmailFrom,
		To:      []string{email},
		Html:    "<strong>Hi! This is Nokori.</strong>", // TODO: add the verification code here and in the subject
		Subject: "Hi this is nokori.",
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
