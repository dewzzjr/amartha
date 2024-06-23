package action

import (
	"time"

	"github.com/dewzzjr/amartha/reconcile/files"
	"github.com/rs/zerolog/log"
)

func Run(param Param) (*Result, error) {
	scanner, err := files.Reader(param.SourceFile)
	if err != nil {
		return nil, err
	}
	inputCh := []chan *TransactionModel{}
	outputCh := make(chan Output)
	missingCh := make(map[string]chan []StatementModel)
	for i, file := range param.StatementFiles {
		inputCh = append(inputCh, make(chan *TransactionModel, 100))
		missingCh[file] = make(chan []StatementModel)
		go Worker(file,
			inputCh[i],
			outputCh,
			missingCh[file],
			Filter(param.Start.Value(), param.End.Value(), true))
	}

	var counter int
	for record := range scanner {
		log.Debug().Any("file", param.SourceFile).Msgf("%+v", record)
		if ToTransactionModel(
			record,
			Filter(param.Start.Value(), param.End.Value(), false),
		).
			SetChannel(inputCh...).
			Next() {
			counter++
			continue
		}
	}

	result := Result{
		TransactionMissing: []TransactionModel{},
		StatementMissing:   make(map[string][]StatementModel),
	}
	for res := range outputCh {
		if counter--; counter == 0 {
			close(outputCh)
		}

		result.Total++
		if res.NotFound == nil {
			result.Matched++
			continue
		}

		result.Unmatched++
		result.TransactionMissing = append(
			result.TransactionMissing,
			*res.NotFound,
		)
	}
	for _, ch := range inputCh {
		close(ch)
	}
	for file, ch := range missingCh {
		result.StatementMissing[file] = <-ch
	}
	return &result, nil
}

func Worker(path string, inputCh <-chan *TransactionModel, outputCh chan<- Output, missingCh chan<- []StatementModel, filter func(time.Time) bool) {
	scanner, err := files.Reader(path)
	if err != nil {
		return
	}
	models := make([]StatementModel, 0)
	for record := range scanner {
		log.Debug().Any("file", path).Msgf("%+v", record)
		m := ToStatementModel(record, filter)
		if m != nil {
			models = append(models, *m)
		}
	}

	for transaction := range inputCh {
		var found bool
		for i, bank := range models {
			if !CompareTimestampWithDate(transaction.Timestamp, bank.Date) {
				continue
			}
			if transaction.Amount*transaction.Type.Float() != bank.Amount {
				continue
			}
			outputCh <- Output{
				TransactionID: transaction.ID,
				BankID:        bank.ID,
			}
			found = true
			DeleteAndPop(&models, i)
			break
		}
		if !found && !transaction.Next() {
			outputCh <- Output{
				TransactionID: transaction.ID,
				NotFound:      transaction,
			}
		}
	}
	missingCh <- models
}

func CompareTimestampWithDate(timestamp, date time.Time) bool {
	return time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, timestamp.Location()).Equal(date)
}
