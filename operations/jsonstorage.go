package operations

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// JSONStorage stores expenses in the filesystem in JSON format
type JSONStorage struct {
	database *json.Encoder
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
	return storage.database.Encode(storage.Expenses)
}

func (storage *JSONStorage) List() []Expense {
	return storage.Expenses
}

func (storage *JSONStorage) SummaryByNWS() Summary {
	summary := make(Summary)
	for _, e := range storage.Expenses {
		summary[e.NWS] += e.Amount
	}
	return summary
}

func (storage *JSONStorage) SummaryByDomain() Summary {
	summary := make(Summary)
	for _, e := range storage.Expenses {
		summary[e.Domain] += e.Amount
	}
	return summary
}
