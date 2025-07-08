package db

import (
	"database/sql"
)

func GetAccountEmailByID(db *sql.DB, accountID string) (string, error) {
	row := db.QueryRow("SELECT email FROM account WHERE account_id = ?", accountID)
	var email string
	if err := row.Scan(&email); err != nil {
		return "", err
	}
	return email, nil
}

func AccountExists(db *sql.DB, accountID string) (bool, error) {
	row := db.QueryRow("SELECT 1 FROM account WHERE account_id = ?", accountID)
	var exists int
	if err := row.Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
