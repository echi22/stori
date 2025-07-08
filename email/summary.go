package email

import (
	"fmt"
	"sort"
	"stori/models"
	"time"
)

func GenerateAccountSummaryHTML(accountID string, txs []models.Transaction) string {
	var credits, debits, balance float64
	var creditCount, debitCount int
	monthTxCount := make(map[string]int)
	monthOrder := make(map[string]time.Time)

	for _, tx := range txs {
		// Parse month and year from timestamp
		monthYear := "Unknown"
		if t, err := time.Parse(time.RFC3339, tx.Timestamp); err == nil {
			monthYear = t.Format("January 2006")
			monthOrder[monthYear] = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
		}
		monthTxCount[monthYear]++

		if tx.Amount > 0 {
			credits += tx.Amount
			creditCount++
		} else {
			debits += tx.Amount
			debitCount++
		}
		balance += tx.Amount
	}

	avgCredit := 0.0
	if creditCount > 0 {
		avgCredit = credits / float64(creditCount)
	}
	avgDebit := 0.0
	if debitCount > 0 {
		avgDebit = debits / float64(debitCount)
	}

	// Sort months for consistent output
	months := make([]string, 0, len(monthTxCount))
	for m := range monthTxCount {
		months = append(months, m)
	}
	sort.Slice(months, func(i, j int) bool {
		return monthOrder[months[i]].Before(monthOrder[months[j]])
	})

	// Left column: total balance and transactions per month
	leftCol := fmt.Sprintf("Total balance is %.2f<br>", balance)
	for _, month := range months {
		leftCol += fmt.Sprintf("Number of transactions in %s: %d<br>", month, monthTxCount[month])
	}

	// Right column: average debit and credit
	rightCol := fmt.Sprintf("Average debit amount: %.2f<br>Average credit amount: %.2f", avgDebit, avgCredit)

	return fmt.Sprintf(`
		<table style="width:100%%; max-width:600px; border-collapse:collapse;">
			<tr>
				<td style="vertical-align:top; width:50%%; padding:8px; border-right:1px solid #eee;">%s</td>
				<td style="vertical-align:top; width:50%%; padding:8px;">%s</td>
			</tr>
		</table>
	`, leftCol, rightCol)
}

// SendAccountSummary generates the summary HTML and sends the email using the provided sender.
func SendAccountSummary(sender Sender, accountID, recipient string, txs []models.Transaction, logoPath string) error {
	htmlBody := GenerateAccountSummaryHTML(accountID, txs)
	return sender.Send(recipient, "Stori Account Summary", htmlBody, logoPath)
}
