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
	// TODO calculate total disrepancy
	return 0
}
