package mail

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/go-mail/mail"
)

type Mailer struct {
	dialer *mail.Dialer
	from   string
}

func NewMailer(host string, port int, username, password, from string) *Mailer {
	d := mail.NewDialer(host, port, username, password)
	d.TLSConfig = &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: false,
	}

	return &Mailer{
		dialer: d,
		from:   from,
	}
}

func (m *Mailer) Send(to, subject, body string) error {
	msg := mail.NewMessage()
	msg.SetHeader("From", m.from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	if err := m.dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent to %s", to)
	return nil
}
