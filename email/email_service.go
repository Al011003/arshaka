package service

type EmailService interface {
	Send(to, subject, body string) error
	SendOTP(to, userName, otp string) error
}
