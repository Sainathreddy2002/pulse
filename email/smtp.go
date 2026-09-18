package email

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

type Sender struct {
	host string
	port string
	user string
	pass string
	from string
}

func NewSenderFromEnv() (*Sender, error) {
	host := os.Getenv("SMTP_HOST")
	if host == "" {
		host = "smtp.gmail.com"
	}
	port := os.Getenv("SMTP_PORT")
	if port == "" {
		port = "587"
	}
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")
	if from == "" {
		from = user
	}
	if user == "" || pass == "" {
		return nil, fmt.Errorf("SMTP_USER and SMTP_PASSWORD are required")
	}
	return &Sender{
		host: host,
		port: port,
		user: user,
		pass: pass,
		from: from,
	}, nil
}

func (s *Sender) Send(to, subject, body string) error {
	addr := s.host + ":" + s.port
	msg := strings.Join([]string{
		"From: " + s.from,
		"To: " + to,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	auth := smtp.PlainAuth("", s.user, s.pass, s.host)
	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg))
}
