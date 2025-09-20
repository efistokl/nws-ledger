package operations

import (
	"encoding/csv"
	"fmt"
	"strings"
)

type NWS string

const (
	NWS_Needs   NWS = "needs"
	NWS_Wants   NWS = "wants"
	NWS_Savings NWS = "savings"
)

type Expense struct {
	Amount int    `json:"amount"`
	NWS    NWS    `json:"nws"`
	Domain string `json:"domain"`
	Name   string `json:"name"`
}

type SummaryByNWS map[NWS]int

type ExpenseStorage interface {
	Add(Expense) error
	List() []Expense
	Summary() SummaryByNWS
}

func FormatCSVList(es ExpenseStorage) string {
	var b strings.Builder
	writer := csv.NewWriter(&b)

	writer.Write([]string{"name", "amount", "nws"})
	for _, e := range es.List() {
		writer.Write([]string{
			e.Name,
			fmt.Sprintf("%d", e.Amount),
			string(e.NWS),
		})
	}
	writer.Flush()
	return b.String()
}
