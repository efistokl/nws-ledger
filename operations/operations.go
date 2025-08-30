package operations

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type NWS string

const (
	NWS_Need   NWS = "needs"
	NWS_Want   NWS = "wants"
	NWS_Saving NWS = "savings"
)

type Expense struct {
	Amount int    `json:"amount"`
	NWS    NWS    `json:"nws"`
	Domain string `json:"domain"`
	Name   string `json:"name"`
}

type ExpenseStorage interface {
	Add(Expense) error
	Save() error
	List() []Expense
}

func FormatCSV(es ExpenseStorage) string {
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

type JSONStorage struct {
	file     io.ReadWriter
	Expenses []Expense `json:"expenses"`
}

func NewJSONStorage(source io.ReadWriter) *JSONStorage {
	return &JSONStorage{file: source}
}

func (storage *JSONStorage) Add(e Expense) error {
	storage.Expenses = append(storage.Expenses, e)
	return nil
}

func (storage *JSONStorage) Save() error {
	return json.NewEncoder(storage.file).Encode(storage.Expenses)
}

func (storage *JSONStorage) List() []Expense {
	return storage.Expenses
}
