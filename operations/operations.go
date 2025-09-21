package operations

import (
	"encoding/csv"
	"fmt"
	"slices"
	"strings"
)

const (
	NWS_Needs   = "needs"
	NWS_Wants   = "wants"
	NWS_Savings = "savings"
)

var nwss = []string{NWS_Needs, NWS_Wants, NWS_Savings}

func ValidateNWS(s string) error {
	if !slices.Contains(nwss, s) {
		return fmt.Errorf("supported values: %v", nwss)
	}
	return nil
}

type Expense struct {
	Amount int    `json:"amount"`
	NWS    string `json:"nws"`
	Domain string `json:"domain"`
	Name   string `json:"name"`
}

type Summary map[string]int

type ExpenseStorage interface {
	Add(Expense) error
	List() []Expense
	SummaryByNWS() Summary
	SummaryByDomain() Summary
}

func FormatCSVList(es ExpenseStorage) string {
	var b strings.Builder
	writer := csv.NewWriter(&b)

	writer.Write([]string{"name", "amount", "nws", "domain"})
	for _, e := range es.List() {
		writer.Write([]string{
			e.Name,
			fmt.Sprintf("%d", e.Amount),
			string(e.NWS),
			e.Domain,
		})
	}
	writer.Flush()
	return b.String()
}

func FormatCSVSummaryByNWS(es ExpenseStorage) string {
	summary := es.SummaryByNWS()
	return writeSummaryCSV("nws", nwss, summary)
}

func FormatCSVSummaryByDomain(es ExpenseStorage) string {
	summary := es.SummaryByDomain()

	domains := make([]string, 0, len(summary))
	for k := range summary {
		domains = append(domains, k)
	}

	slices.SortFunc(domains, func(a, b string) int {
		return summary[b] - summary[a]
	})

	return writeSummaryCSV("domain", domains, summary)
}

func writeSummaryCSV(name string, keys []string, values map[string]int) string {
	var b strings.Builder
	w := csv.NewWriter(&b)
	w.Write([]string{name, "amount"})

	total := 0
	for _, key := range keys {
		w.Write([]string{
			key,
			fmt.Sprintf("%d", values[key]),
		})
		total += values[key]
	}
	w.Write([]string{
		"total",
		fmt.Sprintf("%d", total),
	})

	w.Flush()
	return b.String()
}
