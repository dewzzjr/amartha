package action

import "strings"

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
