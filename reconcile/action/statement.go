package action

import (
	"fmt"
	"strconv"
	"time"
)

type StatementModel struct {
	ID     string
	Amount float64
	Date   time.Time
}

func ToStatementModel(record []string, filter func(time.Time) bool) *StatementModel {
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

	return &StatementModel{
		ID:     record[0],
		Amount: amount,
		Date:   date,
	}
}

func (m StatementModel) String() string {
	return fmt.Sprintf("[%s\t%.2f\t%s]", m.ID, m.Amount, m.Date.Format("2006/01/02"))
}
