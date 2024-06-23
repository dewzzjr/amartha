package action

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/dewzzjr/amartha/reconcile/files"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

type FlowType int

func (f FlowType) Float() float64 {
	return float64(f)
}

func (f FlowType) String() string {
	switch f {
	case Debit:
		return "DEBIT"
	case Credit:
		return "CREDIT"
	default:
		return ""
	}
}

const (
	Debit FlowType = iota - 1
	_
	Credit
)

func ParseFlowType(s string) FlowType {
	switch strings.ToUpper(s) {
	case "DEBIT":
		return Debit
	case "CREDIT":
		return Credit
	default:
		return FlowType(0)
	}
}

type Param struct {
	SourceFile     string
	StatementFiles []string
	Start          cli.Timestamp
	End            cli.Timestamp
}

type Output struct {
	TransactionID, BankID string
	NotFound              *TransactionModel
}

type Result struct {
	Total, Matched, Unmatched int
	TransactionMissing        []TransactionModel
	BankMissing               map[string][]BankModel
}

func (r Result) TotalDisrepancies() float64 {
	// TODO calculate total disrepancy
	return 0
}

type TransactionModel struct {
	ID        string
	Amount    float64
	Type      FlowType
	Timestamp time.Time

	next []chan *TransactionModel
}

func ToTransactionModel(record []string, filter func(time.Time) bool) *TransactionModel {
	if len(record) != 4 {
		return nil
	}
	timestamp, err := time.Parse("2006/01/02 15:04:05", record[3])
	if err != nil {
		return nil
	}
	if !filter(timestamp) {
		return nil
	}
	amount, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		return nil
	}
	return &TransactionModel{
		ID:        record[0],
		Amount:    amount,
		Type:      ParseFlowType(record[2]),
		Timestamp: timestamp,
		next:      make([]chan *TransactionModel, 0),
	}
}

func (m *TransactionModel) SetChannel(chnls ...chan *TransactionModel) *TransactionModel {
	if m == nil {
		return nil
	}
	m.next = append(m.next, chnls...)
	return m
}

func (m *TransactionModel) Next() bool {
	if m == nil {
		return false
	}
	if len(m.next) == 0 {
		return false
	}
	index := rand.Intn(len(m.next))
	ch := DeleteAndPop(&m.next, index)
	defer func() { ch <- m }()
	return true
}

func (m TransactionModel) String() string {
	return fmt.Sprintf("[%s\t%.2f\t%s\t%s]", m.ID, m.Amount, m.Type, m.Timestamp.Format("2006/01/02 15:04:05"))
}

type BankModel struct {
	ID     string
	Amount float64
	Date   time.Time
}

func ToBankModel(record []string, filter func(time.Time) bool) *BankModel {
	if len(record) != 3 {
		return nil
	}

	date, err := time.Parse("2006/01/02", record[2])
	if err != nil {
		return nil
	}
	if !filter(date) {
		return nil
	}
	amount, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		return nil
	}

	return &BankModel{
		ID:     record[0],
		Amount: amount,
		Date:   date,
	}
}

func (m BankModel) String() string {
	return fmt.Sprintf("[%s\t%.2f\t%s]", m.ID, m.Amount, m.Date.Format("2006/01/02"))
}

func Filter(start, end *time.Time, dateOnly bool) func(time.Time) bool {
	return func(timestamp time.Time) bool {
		afterStart := true
		matchStart := false
		if start != nil {
			afterStart = timestamp.After(*start)
			matchStart = timestamp.Equal(*start)
		}
		beforeEnd := true
		matchEnd := false
		if end != nil {
			beforeEnd = timestamp.Before(*end)
			matchEnd = timestamp.Equal(*end)
		}
		between := afterStart && beforeEnd
		if !dateOnly {
			return between
		}
		return between || matchStart || matchEnd
	}
}

func DeleteAndPop[T any](slice *[]T, s int) (result T) {
	result = (*slice)[s]
	*slice = append((*slice)[:s], (*slice)[s+1:]...)
	return
}

func Run(param Param) (*Result, error) {
	scanner, err := files.Reader(param.SourceFile)
	if err != nil {
		return nil, err
	}
	inputCh := []chan *TransactionModel{}
	outputCh := make(chan Output)
	missingCh := make(map[string]chan []BankModel)
	for i, file := range param.StatementFiles {
		inputCh = append(inputCh, make(chan *TransactionModel, 5))
		missingCh[file] = make(chan []BankModel)
		go Worker(file,
			inputCh[i],
			outputCh,
			missingCh[file],
			Filter(param.Start.Value(), param.End.Value(), true))
	}

	var counter int
	for record := range scanner {
		log.Debug().Msgf("%+v", record)
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
		BankMissing:        make(map[string][]BankModel),
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
		result.BankMissing[file] = <-ch
	}
	return &result, nil
}

func Worker(path string, inputCh <-chan *TransactionModel, outputCh chan<- Output, missingCh chan<- []BankModel, filter func(time.Time) bool) {
	scanner, err := files.Reader(path)
	if err != nil {
		return
	}
	models := make([]BankModel, 0)
	for record := range scanner {
		log.Debug().Msgf("%+v", record)
		m := ToBankModel(record, filter)
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
