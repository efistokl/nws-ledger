package operations

import (
	"encoding/csv"
	"fmt"
	"slices"
	"strings"
)

type NWS string

const (
	NWS_Needs   NWS = "needs"
	NWS_Wants   NWS = "wants"
	NWS_Savings NWS = "savings"
)

var nwss = []NWS{NWS_Needs, NWS_Wants, NWS_Savings}

func ValidateNWS(s string) error {
	if !slices.Contains(nwss, NWS(s)) {
		return fmt.Errorf("supported values: %v", nwss)
	}
	return nil
}

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

func FormatCSVSummary(es ExpenseStorage) string {
	summary := es.Summary()

	var b strings.Builder
	writer := csv.NewWriter(&b)

	writer.Write([]string{"nws", "amount"})
	total := 0
	for _, nws := range nwss {
		writer.Write([]string{
			string(nws),
			fmt.Sprintf("%d", summary[nws]),
		})
		total += summary[nws]
	}

	writer.Write([]string{
		"total",
		fmt.Sprintf("%d", total),
	})

	writer.Flush()
	return b.String()
}
