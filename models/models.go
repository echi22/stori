package models

type Config struct {
	SMTPHost       string `json:"smtp_host"`
	SMTPPort       int    `json:"smtp_port"`
	SMTPUser       string `json:"smtp_user"`
	SMTPPass       string `json:"smtp_pass"`
	FromEmail      string `json:"from_email"`
	AccountName    string `json:"account_name"`
	RecipientEmail string `json:"recipient_email"`
	LogoValue      string `json:"logo_value"`
}

type Account struct {
	ID        int
	AccountID string
	Name      string
	Email     string
}

type Transaction struct {
	ID        int
	Timestamp string
	Amount    float64
	AccountID string
}
