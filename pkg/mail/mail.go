package mail

import (
	"github.com/conan194351/BTL-KTPM/internal/config"
	"gopkg.in/gomail.v2"
)

type MailService interface {
	SendEmail(to string, subject string, body string) error
}

type MailServiceImpl struct {
}

func NewMailService() MailService {
	return &MailServiceImpl{}
}

func (ms *MailServiceImpl) SendEmail(to string, subject string, body string) error {
	mailConfig := config.GetConfig().Mail
	smtpHost := mailConfig.Host
	smtpPort := mailConfig.Port
	senderEmail := mailConfig.SenderEmail
	senderPassword := mailConfig.Password

	// Tạo email
	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Thiết lập SMTP client
	d := gomail.NewDialer(smtpHost, smtpPort, senderEmail, senderPassword)

	// Gửi email
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
