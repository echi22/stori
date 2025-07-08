package db

import (
	"database/sql"
	"encoding/csv"
	"os"
	"stori/models"

	_ "github.com/mattn/go-sqlite3"
)

type DB interface {
	Init(dbPath string) error
	SaveAccount(account models.Account) (int, error)
	SaveTransaction(tx models.Transaction) error
	Close() error
	LoadAccountsFromCSV(csvPath string) error
	GetAccountIDByAccountID(accountID string) (int, error)
}

type SQLiteDB struct {
	db *sql.DB
}

func (s *SQLiteDB) Init(dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS account (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			account_id TEXT UNIQUE NOT NULL,
			name TEXT NOT NULL,
			email TEXT
		);
		CREATE TABLE IF NOT EXISTS transactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp TEXT,
			amount REAL NOT NULL,
			account_id INTEGER,
			FOREIGN KEY(account_id) REFERENCES account(id),
			UNIQUE(account_id, timestamp, amount)
		);
	`)
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *SQLiteDB) SaveAccount(account models.Account) (int, error) {
	res, err := s.db.Exec("INSERT OR IGNORE INTO account (account_id, name, email) VALUES (?, ?, ?)", account.AccountID, account.Name, account.Email)
	if err != nil {
		return 0, err
	}
	var id int64
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		// Already exists, fetch id
		row := s.db.QueryRow("SELECT id FROM account WHERE account_id = ?", account.AccountID)
		row.Scan(&id)
	} else {
		id, err = res.LastInsertId()
	}
	return int(id), err
}

func (s *SQLiteDB) SaveTransaction(tx models.Transaction) error {
	_, err := s.db.Exec("INSERT OR IGNORE INTO transactions (timestamp, amount, account_id) VALUES (?, ?, ?)", tx.Timestamp, tx.Amount, tx.AccountID)
	return err
}

func (s *SQLiteDB) Close() error {
	return s.db.Close()
}

func (s *SQLiteDB) LoadAccountsFromCSV(csvPath string) error {
	f, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer f.Close()

	r := csv.NewReader(f)
	_, err = r.Read() // skip header
	if err != nil {
		return err
	}
	for {
		rec, errRead := r.Read()
		if errRead != nil {
			break
		}
		if len(rec) < 3 {
			continue
		}
		account := models.Account{
			AccountID: rec[0],
			Name:      rec[1],
			Email:     rec[2],
		}
		_, err = s.SaveAccount(account)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SQLiteDB) GetAccountIDByAccountID(accountID string) (int, error) {
	row := s.db.QueryRow("SELECT id FROM account WHERE account_id = ?", accountID)
	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (s *SQLiteDB) DB() *sql.DB {
	return s.db
}
