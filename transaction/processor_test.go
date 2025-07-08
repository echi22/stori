package transaction

import (
	"os"
	"stori/models"
	"testing"
)

func TestCSVTransactionProcessor(t *testing.T) {
	csvContent := "AccountID,Id,Timestamp,Transaction\nA123,0,2024-07-08T14:23:00Z,+60.5\nA123,1,2024-07-08T15:00:00Z,-20.46\n"
	f, err := os.CreateTemp("", "txs-*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	f.WriteString(csvContent)
	f.Close()

	p := &CSVTransactionProcessor{}
	txs, credits, debits, balance, count, err := p.Process(f.Name())
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}
	if len(txs) != 2 || count != 2 {
		t.Errorf("Expected 2 transactions, got %d", len(txs))
	}
	if credits != 60.5 {
		t.Errorf("Expected credits 60.5, got %f", credits)
	}
	if debits != -20.46 {
		t.Errorf("Expected debits -20.46, got %f", debits)
	}
	if balance != 40.04 {
		t.Errorf("Expected balance 40.04, got %f", balance)
	}
}

func TestValidateSingleAccount_MultipleAccounts(t *testing.T) {
	csvContent := "AccountID,Id,Timestamp,Transaction\nA123,0,2024-07-08T14:23:00Z,+60.5\nB456,1,2024-07-08T15:00:00Z,-20.46\n"
	f, err := os.CreateTemp("", "txs-*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	f.WriteString(csvContent)
	f.Close()

	p := &CSVTransactionProcessor{}
	txs, _, _, _, _, err := p.Process(f.Name())
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}
	grouped := GroupTransactionsByAccountID(txs)
	_, err = ValidateSingleAccount(grouped)
	if err == nil {
		t.Error("Expected error for multiple accounts, got nil")
	}
}

func TestValidateSingleAccount_Empty(t *testing.T) {
	grouped := GroupTransactionsByAccountID([]models.Transaction{})
	_, err := ValidateSingleAccount(grouped)
	if err == nil {
		t.Error("Expected error for empty transactions, got nil")
	}
}

func TestCSVTransactionProcessor_InvalidAmount(t *testing.T) {
	csvContent := "AccountID,Id,Timestamp,Transaction\nA123,0,2024-07-08T14:23:00Z,notanumber\nA123,1,2024-07-08T15:00:00Z,+20.00\n"
	f, err := os.CreateTemp("", "txs-*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	f.WriteString(csvContent)
	f.Close()

	p := &CSVTransactionProcessor{}
	txs, _, _, _, count, err := p.Process(f.Name())
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}
	if len(txs) != 1 || count != 1 {
		t.Errorf("Expected 1 valid transaction, got %d", len(txs))
	}
}
