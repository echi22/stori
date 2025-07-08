# Stori Transaction Processor

This project processes a CSV file of debit and credit transactions for an account, stores the data in a SQLite database, and sends a summary email to the account's email address. The summary includes total balance, monthly transaction counts, and average debit/credit amounts, styled for easy reading.

## Test the Lambda Function with the Public Uploader

You can test the full workflow (S3 upload → Lambda → email) using the public uploader:

1. **Go to the Public Uploader Page**  
   [https://echi22.github.io/stori/](https://echi22.github.io/stori/)

2. **Upload a File**  
   - Click "Select file to upload" and choose your transaction file (check below for correct format).
   - Click "Upload."
   - You should see a success message and a link to the uploaded file if the upload is successful.

3. **Lambda Trigger**  
   - The upload to the `storifiles` S3 bucket will automatically trigger your AWS Lambda function (as configured).

4. **Check the Email**  
   - Go to [https://ethereal.email/](https://ethereal.email/)
   - Login with the SMTP credentials provided in your `config.json` (`smtp_user` and `smtp_pass`).
   - Open the inbox to view the sent email.
   - Sample output:
   
   ![Screenshot 2025-07-08 at 1 15 37 AM](https://github.com/user-attachments/assets/336d8d06-d505-4702-88fc-6685294abb6b)


**Notes:**
- Files uploaded via the web page are public in the S3 bucket. Do not upload sensitive data.
- I know this is insecure, it's just for testing purpouses.

---

## Try it Online

You can launch this project in a GitHub Codespace (cloud-based dev environment) with all dependencies pre-installed:

[![Open in Codespaces](https://github.com/codespaces/badge.svg)](https://github.com/codespaces/new/echi22/stori)

- No local setup required!
- Just open a terminal and run:
  ```sh
  go run main.go
  ```

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
- `Dockerfile` — For building and running the app in a container.

## Design Decisions

- **Why `accounts.csv`?**
  - Separating account metadata (ID, name, email) into its own file allows for clean normalization, easier updates, and future extensibility (e.g., supporting multiple accounts, richer metadata, or account lookups).

- **Why `AccountID` and `Timestamp` in `transactions.csv`?**
  - `AccountID` links each transaction to its account, enabling support for multiple accounts and ensuring data integrity.
  - `Timestamp` (in ISO 8601 format) provides precise timing for each transaction, supports real-world scenarios (multiple transactions per day), and enables robust uniqueness and reporting (e.g., monthly summaries).

## Prerequisites

- Go 1.21+
- [Docker](https://www.docker.com/) (for containerized testing)
- [SQLite3](https://www.sqlite.org/index.html) (for inspecting the DB, optional)

## Setup & Execution

### Run with Docker

1. **Build the Docker image:**
   ```sh
   docker build -t stori-app .
   ```
2. **Run the app in a container:**
   ```sh
   docker run --rm stori-app
   ```
3. **View the email:**
   - Go to [https://ethereal.email/](https://ethereal.email/)
   - Login with the SMTP credentials provided in your `config.json` (`smtp_user` and `smtp_pass`).
   - Open the inbox to view the sent email.

### Run Locally (Go)

1. **Install dependencies:**
   ```sh
   go mod tidy
   ```
2. **Run the application:**
   ```sh
   go run main.go
   ```

### Try it Online (Codespaces)

See the section above for Codespaces instructions.

### AWS Lambda

See section above for Lambda instructions.

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



