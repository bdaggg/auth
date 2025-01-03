package service

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	dialer *gomail.Dialer
	from   string
}

func NewEmailService(host string, port int, username, password, from string) *EmailService {
	return &EmailService{
		dialer: gomail.NewDialer(host, port, username, password),
		from:   from,
	}
}

func (s *EmailService) SendVerificationEmail(to, token string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Email Doğrulama")
	m.SetBody("text/html", fmt.Sprintf(`
		<h1>Email Adresinizi Doğrulayın</h1>
		<p>Aşağıdaki linke tıklayarak email adresinizi doğrulayabilirsiniz:</p>
		<a href="http://localhost:8080/api/v1/auth/verify-email?token=%s">Email Doğrula</a>
	`, token))

	return s.dialer.DialAndSend(m)
}

func (s *EmailService) SendPasswordResetEmail(to, token string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Şifre Sıfırlama")
	m.SetBody("text/html", fmt.Sprintf(`
		<h1>Şifrenizi Sıfırlayın</h1>
		<p>Aşağıdaki linke tıklayarak şifrenizi sıfırlayabilirsiniz:</p>
		<a href="http://localhost:8080/api/v1/auth/reset-password?token=%s">Şifre Sıfırla</a>
	`, token))

	return s.dialer.DialAndSend(m)
}
