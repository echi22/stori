package email

import (
	"encoding/base64"
	"fmt"
	"net/smtp"
	"os"
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

	// for limitations in the email service, we need to use an external logo
	logoHTML := ""
	if s.Config.LogoType == "url" && s.Config.LogoValue != "" {
		logoHTML = fmt.Sprintf(`<img src="%s" alt="Stori Logo" style="width:120px;">`, s.Config.LogoValue)
	} else {
		logoBase64 := ""
		if logoPath != "" {
			data, err := os.ReadFile(logoPath)
			if err == nil {
				logoBase64 = base64.StdEncoding.EncodeToString(data)
			}
		}
		if logoBase64 != "" {
			logoHTML = fmt.Sprintf(`<img src="data:image/png;base64,%s" alt="Stori Logo" style="width:120px;">`, logoBase64)
		}
	}

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + logoHTML + htmlBody

	addr := fmt.Sprintf("%s:%d", s.Config.SMTPHost, s.Config.SMTPPort)
	return smtp.SendMail(addr, auth, from, to, []byte(message))
}
