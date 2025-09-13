package operations

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
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
	fileName string
	Expenses []Expense
}

// NewJSONStorage initializes JSONStorage by parsing the file "fileName".
// On every "Add" it rewrites the file.
func NewJSONStorage(fileName string) (*JSONStorage, error) {
	_, err := os.Stat(fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &JSONStorage{fileName, []Expense{}}, nil
		}
		return nil, fmt.Errorf("failed to get FileInfo for %s: %v", fileName, err)
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed opening file %s: %v", fileName, err)
	}
	defer file.Close()

	var expenses []Expense
	if err := json.NewDecoder(file).Decode(&expenses); err != nil {
		return nil, fmt.Errorf("failed decoding JSON from file %s: %v", fileName, err)
	}

	return &JSONStorage{fileName, expenses}, nil
}

func (storage *JSONStorage) Add(e Expense) error {
	storage.Expenses = append(storage.Expenses, e)
	return storage.flush()
}

func (storage *JSONStorage) flush() error {
	file, err := os.OpenFile(storage.fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(storage.Expenses)
}

func (storage *JSONStorage) List() []Expense {
	return storage.Expenses
}
