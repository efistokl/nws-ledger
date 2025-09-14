package operations

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
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
	database *json.Encoder
	// fileName string
	Expenses []Expense
}

// NewJSONStorage initializes JSONStorage by parsing the file.
// On every "Add" it rewrites the file.
func NewJSONStorage(file *os.File) (*JSONStorage, error) {
	err := initializeJSONStorageFile(file)
	if err != nil {
		return nil, fmt.Errorf("problem initializing JSONStorage file, %v", err)
	}

	var expenses []Expense
	if err := json.NewDecoder(file).Decode(&expenses); err != nil {
		return nil, fmt.Errorf("failed decoding JSON from file %s: %v", file.Name(), err)
	}

	return &JSONStorage{
		database: json.NewEncoder(&tape{file}),
		Expenses: expenses,
	}, nil
}

func initializeJSONStorageFile(file *os.File) error {
	file.Seek(0, io.SeekStart)

	info, err := file.Stat()

	if err != nil {
		return fmt.Errorf("problem getting file info from file %s: %v", file.Name(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, io.SeekStart)
	}

	return nil
}

func (storage *JSONStorage) Add(e Expense) error {
	storage.Expenses = append(storage.Expenses, e)
	return storage.flush()
}

func (storage *JSONStorage) flush() error {
	return storage.database.Encode(storage.Expenses)
}

func (storage *JSONStorage) List() []Expense {
	return storage.Expenses
}
