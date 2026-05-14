package email

import (
	"fmt"
	"net/smtp"
)

type SMTPClient struct {
	Host     string
	Port     string
	Username string
	Password string
}

func NewSMTPClient(host, port, username, password string) *SMTPClient {
	return &SMTPClient{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

func (c *SMTPClient) Send(to, subject, htmlBody string) error {
	addr := fmt.Sprintf("%s:%s", c.Host, c.Port)
	auth := smtp.PlainAuth("", c.Username, c.Password, c.Host)

	// Gmail advertises STARTTLS, which smtp.SendMail will automatically use.
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
		htmlBody)

	return smtp.SendMail(addr, auth, c.Username, []string{to}, msg)
}
