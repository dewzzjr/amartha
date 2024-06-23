package action

import "github.com/urfave/cli/v2"

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
	StatementMissing          map[string][]StatementModel
}

func (r Result) TotalDisrepancies() float64 {
	var transaction, statement float64
	for _, t := range r.TransactionMissing {
		transaction += t.Amount * t.Type.Float()
	}

	for _, statements := range r.StatementMissing {
		for _, s := range statements {
			statement += s.Amount
		}
	}
	if transaction > statement {
		return transaction - statement
	}
	return statement - transaction
}
