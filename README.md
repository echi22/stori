# Stori Transaction Processor

This project processes a CSV file of debit and credit transactions for an account, stores the data in a SQLite database, and sends a summary email to the account's email address. The summary includes total balance, monthly transaction counts, and average debit/credit amounts, styled for easy reading.

## Project Structure

- `main.go` — Entry point; orchestrates DB setup, account loading, and transaction processing.
- `db/` — Database logic (SQLite), account and transaction storage, idempotency.
- `transaction/` — Transaction CSV parsing, validation, and grouping.
- `email/` — Email sending logic and summary HTML generation.
- `models/` — Shared data models.
- `accounts.csv` — Account metadata (AccountID, Name, Email).
- `transactions.csv` — Transaction data (AccountID, Id, Timestamp, Transaction).
- `config.json` — SMTP and app configuration.
- `stori_logo.png` — Logo for email branding.

## Design Decisions

- **Why `accounts.csv`?**
  - Separating account metadata (ID, name, email) into its own file allows for clean normalization, easier updates, and future extensibility (e.g., supporting multiple accounts, richer metadata, or account lookups).

- **Why `AccountID` and `Timestamp` in `transactions.csv`?**
  - `AccountID` links each transaction to its account, enabling support for multiple accounts and ensuring data integrity.
  - `Timestamp` (in ISO 8601 format) provides precise timing for each transaction, supports real-world scenarios (multiple transactions per day), and enables robust uniqueness and reporting (e.g., monthly summaries).

## Prerequisites

- Go 1.21+
- [SQLite3](https://www.sqlite.org/index.html) (for inspecting the DB, optional)

## Setup & Execution

1. **Clone the repository**
2. **Install dependencies:**
   ```sh
   go mod tidy
   ```
3. **Configure SMTP and account info:**
   - This project uses [Ethereal Email](https://ethereal.email/) for SMTP testing. Credentials are provided in `config.json`.
   - **Note:** This is insecure and for challenge/demo purposes only. See below for production advice.
   - Edit `accounts.csv` and `transactions.csv` as needed.
4. **Run the application:**
   ```sh
   go run main.go
   ```
5. **Check the output:**
   - The summary email will be sent to the account's email address (from `accounts.csv`).
   - The console will print a confirmation.
   - **To view the email:**
     1. Go to [https://ethereal.email/](https://ethereal.email/)
     2. Login with the SMTP credentials provided in your `config.json` (`smtp_user` and `smtp_pass`).
     3. Open the inbox to view the sent email.

## Security Note: SMTP Credentials

- **For this challenge, SMTP credentials are stored in `config.json` for simplicity.**
- **This is insecure and should never be done in production.**
- In a real-world application, secrets should be managed using environment variables or a secret manager (e.g., AWS Secrets Manager, Azure Key Vault, Google Secret Manager), and never committed to version control.
- For more on Ethereal Email, see [https://ethereal.email/](https://ethereal.email/).

## File Formats

### accounts.csv
```
AccountID,Name,Email
A123,Julien Terry,julien@example.com
```

### transactions.csv
```
AccountID,Id,Timestamp,Transaction
A123,0,2024-05-10T09:15:00Z,+60.5
A123,1,2024-05-15T15:00:00Z,-20.46
...
```
- `Timestamp` must be in ISO 8601 format (e.g., `2024-07-08T14:23:00Z`).
- Only one account per file is allowed; the program will error if multiple accounts are found.

## Testing

Run all tests with:
```sh
go test ./...
```
This covers CSV parsing, DB idempotency, summary generation, and error cases.

## Notes
- The database (`stori.db`) is idempotent: duplicate accounts and transactions are not created on repeated runs.
- The summary email is styled and includes a logo.
- If an account in the transactions file is not found in `accounts.csv`, the program will error.

## Extending
- To support multiple accounts per file, update the logic in `transaction/validate.go` and `main.go`.
- To add more summary stats, edit `email/summary.go`.


