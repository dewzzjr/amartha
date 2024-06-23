package action

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

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
