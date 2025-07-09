# Stori Transaction Processor

This project processes a CSV file of debit and credit transactions for an account, stores the data in a SQLite database, and sends a summary email to the account's email address using [Brevo](https://brevo.com/) for email delivery. The summary includes total balance, monthly transaction counts, and average debit/credit amounts, styled for easy reading.

## Project Structure

- `main.go` — Entry point for local/Docker execution; orchestrates DB setup, account loading, and transaction processing.
- `lambda_main.go` — Entry point for AWS Lambda deployment; triggers the Lambda handler on S3 events.
- `db/` — Database logic (SQLite); handles account and transaction storage, idempotency, and updates account emails if the account exists.
- `transaction/` — Transaction CSV parsing, validation, and grouping.
- `email/` — Email sending logic and summary HTML generation.
- `models/` — Shared data models, including config options for logo type and value.
- `accounts.csv` — Account metadata (AccountID, Name, Email).
- `transactions.csv` — Transaction data (AccountID, Id, Timestamp, Transaction).
- `config.json` — SMTP and app configuration.
- `Dockerfile` — For building and running the app in a container.
  
## Test the Lambda Function with the Public Uploader

You can test the full workflow (S3 upload → Lambda → email) using the public uploader:

1. **Go to the Public Uploader Page**  
   [https://echi22.github.io/stori/](https://echi22.github.io/stori/)

2. **Configure Account Info**  
   - Use the "Update Account Info" form to set your Account ID, Name, and Email. This will generate and upload a new `accounts.csv` to S3.
   - **Important:** The Account ID must match the one used in your transactions file.

3. **Upload a Transactions File**  
   - Click "Select file to upload" and choose your transaction file (check below for correct format).
   - Click "Upload."

4. **Lambda Trigger**  
   - The upload to the `storifiles` S3 bucket will automatically trigger your AWS Lambda function (as configured).

5. **Check the Email**  
   - The summary email will be sent to the address you entered in the Account Info form.
   - Check your inbox for the summary email. (If you don't see it, check your spam folder.)
   - Sample output:
   
   ![Screenshot 2025-07-08 at 1 15 37 AM](https://github.com/user-attachments/assets/336d8d06-d505-4702-88fc-6685294abb6b)

**Notes:**
- Files uploaded via the web page are public in the S3 bucket. Do not upload sensitive data.
- This is for testing purposes only.

---

## Try it Online

You can launch this project in a GitHub Codespace (cloud-based dev environment) with all dependencies pre-installed:

[![Open in Codespaces](https://github.com/codespaces/badge.svg)](https://github.com/codespaces/new/echi22/stori)

- No local setup required!
- Just open a terminal and run:
  ```sh
  go run main.go
  ```

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
1.  **Update the email address in `accounts.csv` to your own email before running.**

2. **Build the Docker image:**
   ```sh
   docker build -t stori-app .
   ```
3. **Run the app in a container:**
   ```sh
   docker run --rm stori-app
   ```
4. **View the email:**
   - The summary email will be sent to the address specified in `accounts.csv` using Brevo.
   - Check your inbox for the summary email.
   - If you can't find the email, check your spam folder

### Run Locally (Go)

1.  **Update the email address in `accounts.csv` to your own email before running.**

2. **Install dependencies:**
   ```sh
   go mod tidy
   ```
3. **Run the application:**
   ```sh
   go run main.go
   ```
4. **View the email:**
   - The summary email will be sent to the address specified in `accounts.csv` using Brevo.
   - Check your inbox for the summary email.
   - If you can't find the email, check your spam folder

### Try it Online (Codespaces)

See the section above for Codespaces instructions.

### AWS Lambda

See section above for Lambda instructions.

### Note on Database Persistence in Lambda

The Lambda function uses a SQLite database stored in the `/tmp` directory during each invocation. This database is **ephemeral**—it is created fresh for each Lambda run and is not persisted between invocations.

- **Implication:** Data is not shared or saved across Lambda executions. Each run starts with an empty database.
- **Production alternative:** For a real-world, persistent solution, use a managed database such as Amazon RDS (PostgreSQL/MySQL), DynamoDB, or another cloud database service.
- **Why SQLite here?** This project uses SQLite in `/tmp` for simplicity and ease of local testing/demo purposes.

## Security Note: SMTP Credentials and Brevo

- **This project uses [Brevo](https://www.brevo.com/) for email delivery.**
- SMTP credentials are stored in `config.json` for simplicity. **This is insecure and should never be done in production.**
- In a real-world application, secrets should be managed using environment variables or a secret manager (e.g., AWS Secrets Manager, Azure Key Vault, Google Secret Manager), and never committed to version control.
- For more on Brevo, see [https://www.brevo.com/](https://www.brevo.com/).

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



