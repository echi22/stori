package transaction

import (
	"errors"
	"stori/models"
)

func GroupTransactionsByAccountID(txs []models.Transaction) map[string][]models.Transaction {
	grouped := make(map[string][]models.Transaction)
	for _, tx := range txs {
		grouped[tx.AccountID] = append(grouped[tx.AccountID], tx)
	}
	return grouped
}

func ValidateSingleAccount(grouped map[string][]models.Transaction) (string, error) {
	if len(grouped) == 0 {
		return "", errors.New("no transactions found in file")
	}
	if len(grouped) > 1 {
		return "", errors.New("the file contains more than one account. Only one account per file is allowed")
	}
	for accountID := range grouped {
		return accountID, nil
	}
	return "", errors.New("unexpected error in account grouping")
}
