package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"stori/db"
	"stori/email"
	"stori/models"
	"stori/transaction"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func downloadFromS3(svc *s3.S3, bucket, key, dest string) error {
	out, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}
	defer out.Body.Close()
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, out.Body)
	return err
}

func HandleRequest(ctx context.Context, s3Event events.S3Event) error {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)

	// Read config from env
	cfg := &models.Config{
		SMTPHost:       getEnv("SMTP_HOST", ""),
		SMTPPort:       func() int { p, _ := strconv.Atoi(getEnv("SMTP_PORT", "587")); return p }(),
		SMTPUser:       getEnv("SMTP_USER", ""),
		SMTPPass:       getEnv("SMTP_PASS", ""),
		AccountName:    getEnv("ACCOUNT_NAME", ""),
		RecipientEmail: getEnv("RECIPIENT_EMAIL", ""),
	}

	// Download accounts.csv (assume in same bucket, or use bundled file)
	bucket := s3Event.Records[0].S3.Bucket.Name
	key := s3Event.Records[0].S3.Object.Key
	localTxPath := "/tmp/transactions.csv"
	if err := downloadFromS3(svc, bucket, key, localTxPath); err != nil {
		return fmt.Errorf("failed to download transactions: %w", err)
	}
	accountsKey := "accounts.csv"
	localAccPath := "/tmp/accounts.csv"
	if err := downloadFromS3(svc, bucket, accountsKey, localAccPath); err != nil {
		return fmt.Errorf("failed to download accounts: %w", err)
	}

	dbPath := "/tmp/stori.db"
	dbImpl := &db.SQLiteDB{}
	if err := dbImpl.Init(dbPath); err != nil {
		return fmt.Errorf("failed to init DB: %w", err)
	}
	defer dbImpl.Close()

	if err := dbImpl.LoadAccountsFromCSV(localAccPath); err != nil {
		return fmt.Errorf("failed to load accounts: %w", err)
	}

	processor := &transaction.CSVTransactionProcessor{}
	txs, _, _, _, _, err := processor.Process(localTxPath)
	if err != nil {
		return fmt.Errorf("failed to process transactions: %w", err)
	}
	grouped := transaction.GroupTransactionsByAccountID(txs)
	accountID, err := transaction.ValidateSingleAccount(grouped)
	if err != nil {
		return err
	}

	exists, err := db.AccountExists(dbImpl.DB(), accountID)
	if err != nil {
		return fmt.Errorf("DB error: %w", err)
	}
	if !exists {
		return fmt.Errorf("account %s not found in DB", accountID)
	}

	emailAddr, err := db.GetAccountEmailByID(dbImpl.DB(), accountID)
	if err != nil {
		return fmt.Errorf("could not get email for account %s: %w", accountID, err)
	}

	accountDbID, err := dbImpl.GetAccountIDByAccountID(accountID)
	if err != nil {
		return fmt.Errorf("could not get DB id for account %s: %w", accountID, err)
	}

	for _, tx := range grouped[accountID] {
		_ = dbImpl.SaveTransaction(models.Transaction{
			Timestamp: tx.Timestamp,
			Amount:    tx.Amount,
			AccountID: strconv.Itoa(accountDbID),
		})
	}

	sender := &email.SMTPSender{Config: *cfg}
	logoPath := "/tmp/stori_logo.png"
	// Try to download logo from S3, fallback to bundled if not found
	logoKey := "stori_logo.png"
	if err := downloadFromS3(svc, bucket, logoKey, logoPath); err != nil {
		logoPath = "stori_logo.png" // fallback
	}
	if err := email.SendAccountSummary(sender, accountID, emailAddr, grouped[accountID], logoPath); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Summary email sent to %s!", emailAddr)
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
