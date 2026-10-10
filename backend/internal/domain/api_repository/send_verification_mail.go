package apirepository

import (
	"crypto/rand"
	"math/big"
	"time"

	"backend/internal/config"
	"backend/internal/domain/apperror"
	"backend/internal/emailtemplate"
	logger "backend/internal/log"

	"github.com/resend/resend-go/v4"
)

type SendVerificationMailRepository struct {
	resendConfig *config.ResendConfig
}

func NewSendVerificationMailRepository(resendConfig *config.ResendConfig) *SendVerificationMailRepository {
	return &SendVerificationMailRepository{
		resendConfig: resendConfig,
	}
}

const verificationCodeLifetime = time.Minute

func (repository *SendVerificationMailRepository) SendEmailWithVerificationCode(email string) (resendSendEmailId string, verificationCode string, err error) {
	logger.Info("SendVerificationMailRepository: SendEmailWithVerificationCode")
	client := resend.NewClient(repository.resendConfig.APIKey)
	verificationCode, err = repository.createVerificationCode()
	if err != nil {
		return
	}

	content, err := emailtemplate.VerificationCode(emailtemplate.VerificationCodeData{
		Code:         verificationCode,
		ValidMinutes: int(verificationCodeLifetime.Minutes()),
	})
	if err != nil {
		return "", "", apperror.Wrap(apperror.InternalError, err)
	}

	params := &resend.SendEmailRequest{
		From:    repository.resendConfig.EmailFrom,
		To:      []string{email},
		Subject: content.Subject,
		Html:    content.HTML,
		Text:    content.Text,
		ReplyTo: repository.resendConfig.ReplyTo,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		return "", "", apperror.Wrap(apperror.ServiceUnavailable, err)
	}

	resendSendEmailId = sent.Id
	return
}

func (repository *SendVerificationMailRepository) createVerificationCode() (verificationCode string, err error) {
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
