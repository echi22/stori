package db

import (
	"stori/models"
	"strconv"
	"testing"
)

func TestSQLiteDB(t *testing.T) {
	db := &SQLiteDB{}
	err := db.Init(":memory:")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer db.Close()

	account := models.Account{AccountID: "A123", Name: "Test User", Email: "test@example.com"}
	accountID, err := db.SaveAccount(account)
	if err != nil {
		t.Fatalf("SaveAccount failed: %v", err)
	}

	tx := models.Transaction{Timestamp: "2024-07-08T14:23:00Z", Amount: 100.0, AccountID: strconv.Itoa(accountID)}
	err = db.SaveTransaction(tx)
	if err != nil {
		t.Fatalf("SaveTransaction failed: %v", err)
	}

	exists, err := AccountExists(db.DB(), "A123")
	if err != nil {
		t.Fatalf("AccountExists failed: %v", err)
	}
	if !exists {
		t.Error("Expected account to exist")
	}

	exists, err = AccountExists(db.DB(), "B999")
	if err != nil {
		t.Fatalf("AccountExists failed: %v", err)
	}
	if exists {
		t.Error("Expected account to not exist")
	}

	_, err = db.GetAccountIDByAccountID("B999")
	if err == nil {
		t.Error("Expected error for non-existent account, got nil")
	}
}

func TestSaveAccount_Duplicate(t *testing.T) {
	db := &SQLiteDB{}
	db.Init(":memory:")
	defer db.Close()
	account := models.Account{AccountID: "A123", Name: "Test User", Email: "test@example.com"}
	id1, err := db.SaveAccount(account)
	if err != nil {
		t.Fatalf("First SaveAccount failed: %v", err)
	}
	id2, err := db.SaveAccount(account)
	if err != nil {
		t.Fatalf("Second SaveAccount failed: %v", err)
	}
	if id1 != id2 {
		t.Error("Expected same ID for duplicate account")
	}
}

func TestSaveTransaction_Duplicate(t *testing.T) {
	db := &SQLiteDB{}
	db.Init(":memory:")
	defer db.Close()
	account := models.Account{AccountID: "A123", Name: "Test User", Email: "test@example.com"}
	accountID, _ := db.SaveAccount(account)

	tx := models.Transaction{Timestamp: "2024-07-08T14:23:00Z", Amount: 100.0, AccountID: strconv.Itoa(accountID)}
	err := db.SaveTransaction(tx)
	if err != nil {
		t.Fatalf("First SaveTransaction failed: %v", err)
	}
	err = db.SaveTransaction(tx)
	if err != nil {
		t.Fatalf("Second SaveTransaction failed: %v", err)
	}
	// Count transactions in DB
	row := db.db.QueryRow("SELECT COUNT(*) FROM transactions WHERE account_id = ?", strconv.Itoa(accountID))
	var count int
	row.Scan(&count)
	if count != 1 {
		t.Errorf("Expected 1 transaction, got %d", count)
	}
}
