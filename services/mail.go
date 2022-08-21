package services

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
)

const (
	senderName = "Rooms 💬"
)

type Mail struct {
	To      []string
	Subject string
	Body    string
}

type IMailService interface {
	Send(mail Mail)
}

type MailService struct {
	from string
	addr string
	auth smtp.Auth
}

var (
	mailOnce    sync.Once
	mailService IMailService
)

func GetMailService() IMailService {
	mailOnce.Do(func() {
		log.Info().Msg("Initializing mail service")

		email := os.Getenv("EMAIL")
		password := os.Getenv("EMAIL_PASSWORD")
		smtpHost := os.Getenv("SMTP_HOST")
		smtpPort := os.Getenv("SMTP_PORT")
		addr := smtpHost + ":" + smtpPort
		auth := smtp.PlainAuth("", email, password, smtpHost)

		mailService = &MailService{
			from: email,
			addr: addr,
			auth: auth,
		}
	})
	return mailService
}

func (s *MailService) Send(mail Mail) {
	log.Debug().Msg("Sending mail")
	msg := s.buildMail(mail)

	err := smtp.SendMail(s.addr, s.auth, s.from, mail.To, msg)
	if err != nil {
		log.Error().Err(err).Msg("Error sending mail")
	}
}

func (s *MailService) buildMail(mail Mail) []byte {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", senderName)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return []byte(msg)
}
