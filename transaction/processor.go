package transaction

import (
	"encoding/csv"
	"io"
	"os"
	"stori/models"
	"strconv"
)

type Processor interface {
	Process(filePath string) ([]models.Transaction, float64, float64, float64, int, error)
}

type CSVTransactionProcessor struct{}

func (p *CSVTransactionProcessor) Process(filePath string) ([]models.Transaction, float64, float64, float64, int, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	_, err = r.Read() // skip header
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	var txs []models.Transaction
	var credits, debits, balance float64
	count := 0
	for {
		rec, errRead := r.Read()
		if errRead == io.EOF {
			break
		}
		if errRead != nil {
			return nil, 0, 0, 0, 0, errRead
		}
		if len(rec) < 4 {
			continue
		}
		amt, errParse := strconv.ParseFloat(rec[3], 64)
		if errParse != nil {
			continue
		}
		tx := models.Transaction{
			ID:        count,
			Timestamp: rec[2],
			Amount:    amt,
			AccountID: rec[0],
		}
		txs = append(txs, tx)
		if amt > 0 {
			credits += amt
		} else {
			debits += amt
		}
		balance += amt
		count++
	}
	return txs, credits, debits, balance, count, nil
}
