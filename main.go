package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"stori/db"
	"stori/email"
	"stori/models"
	"stori/transaction"
	"strconv"
)

func readConfig(path string) (*models.Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg models.Config
	err = json.Unmarshal(data, &cfg)
	return &cfg, err
}

func setupDB() *db.SQLiteDB {
	dbImpl := &db.SQLiteDB{}
	err := dbImpl.Init("stori.db")
	if err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}
	return dbImpl
}

func loadAccountsData(dbImpl *db.SQLiteDB) {
	err := dbImpl.LoadAccountsFromCSV("accounts.csv")
	if err != nil {
		log.Fatalf("Failed to load accounts: %v", err)
	}
}

func processFile(dbImpl *db.SQLiteDB) {
	cfg, err := readConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	processor := &transaction.CSVTransactionProcessor{}
	txs, _, _, _, _, err := processor.Process("transactions.csv")
	if err != nil {
		log.Fatalf("Failed to process transactions: %v", err)
	}

	grouped := transaction.GroupTransactionsByAccountID(txs)
	accountID, err := transaction.ValidateSingleAccount(grouped)
	if err != nil {
		log.Fatalf("%v", err)
	}

	exists, err := db.AccountExists(dbImpl.DB(), accountID)
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}
	if !exists {
		log.Fatalf("Error: Account %s not found in DB. Please check accounts.csv.", accountID)
	}

	emailAddr, err := db.GetAccountEmailByID(dbImpl.DB(), accountID)
	if err != nil {
		log.Fatalf("Error: Could not get email for account %s.", accountID)
	}

	accountDbID, err := dbImpl.GetAccountIDByAccountID(accountID)
	if err != nil {
		log.Fatalf("Error: Could not get DB id for account %s.", accountID)
	}

	for _, tx := range grouped[accountID] {
		_ = dbImpl.SaveTransaction(models.Transaction{
			Timestamp: tx.Timestamp,
			Amount:    tx.Amount,
			AccountID: strconv.Itoa(accountDbID),
		})
	}

	sender := &email.SMTPSender{Config: *cfg}
	err = email.SendAccountSummary(sender, accountID, emailAddr, grouped[accountID], "stori_logo.png")
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	fmt.Printf("Summary email sent to %s!\n", emailAddr)
}

func main() {
	dbImpl := setupDB()
	defer dbImpl.Close()

	loadAccountsData(dbImpl)
	processFile(dbImpl)
}
