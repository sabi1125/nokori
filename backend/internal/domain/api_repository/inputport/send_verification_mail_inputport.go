//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
package inputport

type SendVerificationMailInputPort interface {
	SendEmailWithVerificationCode(email string) (resendSendEmailId string, verificationCode string, err error)
}
