package email

import (
	"stori/models"
	"strings"
	"testing"
)

type MockSender struct {
	Called    bool
	Recipient string
	Subject   string
	Body      string
	LogoPath  string
}

func (m *MockSender) Send(recipient, subject, htmlBody, logoPath string) error {
	m.Called = true
	m.Recipient = recipient
	m.Subject = subject
	m.Body = htmlBody
	m.LogoPath = logoPath
	return nil
}

func TestSMTPSenderInterface(t *testing.T) {
	mock := &MockSender{}
	sender := Sender(mock)
	err := sender.Send("test@example.com", "Test Subject", "<b>Body</b>", "logo.png")
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}
	if !mock.Called {
		t.Error("Send was not called")
	}
	if mock.Recipient != "test@example.com" {
		t.Errorf("Recipient mismatch: %s", mock.Recipient)
	}
	if mock.Subject != "Test Subject" {
		t.Errorf("Subject mismatch: %s", mock.Subject)
	}
	if mock.Body != "<b>Body</b>" {
		t.Errorf("Body mismatch: %s", mock.Body)
	}
	if mock.LogoPath != "logo.png" {
		t.Errorf("LogoPath mismatch: %s", mock.LogoPath)
	}
}

func TestGenerateAccountSummaryHTML(t *testing.T) {
	txs := []models.Transaction{
		{Timestamp: "2024-07-08T14:23:00Z", Amount: 100.0, AccountID: "A123"},
		{Timestamp: "2024-07-08T15:00:00Z", Amount: -20.0, AccountID: "A123"},
		{Timestamp: "2024-08-01T10:00:00Z", Amount: 50.0, AccountID: "A123"},
		{Timestamp: "2025-07-08T14:23:00Z", Amount: 30.0, AccountID: "A123"},
	}
	html := GenerateAccountSummaryHTML("A123", txs)
	if !strings.Contains(html, "Total balance is 160.00") {
		t.Error("Expected total balance in summary")
	}
	if !strings.Contains(html, "Number of transactions in July 2024: 2") {
		t.Error("Expected July 2024 count in summary")
	}
	if !strings.Contains(html, "Number of transactions in August 2024: 1") {
		t.Error("Expected August 2024 count in summary")
	}
	if !strings.Contains(html, "Number of transactions in July 2025: 1") {
		t.Error("Expected July 2025 count in summary")
	}
	if !strings.Contains(html, "Average debit amount: -20.00") {
		t.Error("Expected average debit in summary")
	}
	if !strings.Contains(html, "Average credit amount: 60.00") {
		t.Error("Expected average credit in summary")
	}
}

func TestGenerateAccountSummaryHTML_Empty(t *testing.T) {
	html := GenerateAccountSummaryHTML("A123", []models.Transaction{})
	if !strings.Contains(html, "Total balance is 0.00") {
		t.Error("Expected total balance 0.00 in summary")
	}
}
