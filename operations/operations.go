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
	Flush() error
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
	store    io.Writer
	Expenses []Expense `json:"expenses"`
}

// TODO:
// - Remake JSONStorage to handle file saving.
// - Save on every add for now.
// - Revisit if I want to have Flush at all.
// - For TDD use testing/fstest
// JSONStorage as it is now is unusable

// NewJSONStorage initializes JSONStorage by parsing the "source".
// The "source" and "store" can point to one object
func NewJSONStorage(source io.Reader, store io.Writer) (*JSONStorage, error) {
	if source == nil {
		return &JSONStorage{store: store, Expenses: make([]Expense, 0)}, nil
	}

	var expenses []Expense
	err := json.NewDecoder(source).Decode(&expenses)
	if err != nil {
		return nil, err
	}
	return &JSONStorage{store: store, Expenses: expenses}, nil
}

func (storage *JSONStorage) Add(e Expense) error {
	storage.Expenses = append(storage.Expenses, e)
	return nil
}

func (storage *JSONStorage) Flush() error {
	return json.NewEncoder(storage.store).Encode(storage.Expenses)
}

func (storage *JSONStorage) List() []Expense {
	return storage.Expenses
}
