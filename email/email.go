package email

import (
	"fmt"
	"net/smtp"
	"stori/models"
)

type Sender interface {
	Send(recipient, subject, htmlBody, logoPath string) error
}

type SMTPSender struct {
	Config models.Config
}

func (s *SMTPSender) Send(recipient, subject, htmlBody, logoPath string) error {
	auth := smtp.PlainAuth("", s.Config.SMTPUser, s.Config.SMTPPass, s.Config.SMTPHost)
	from := s.Config.FromEmail
	to := []string{recipient}

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = recipient
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	// we'll use the logo from
	logoHTML := ""
	if s.Config.LogoValue != "" {
		logoHTML = fmt.Sprintf(`<img src="%s" alt="Stori Logo" style="width:120px;">`, s.Config.LogoValue)
	}
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + logoHTML + htmlBody

	addr := fmt.Sprintf("%s:%d", s.Config.SMTPHost, s.Config.SMTPPort)
	return smtp.SendMail(addr, auth, from, to, []byte(message))
}
