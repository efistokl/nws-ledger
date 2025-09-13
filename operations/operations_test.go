package operations

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupFile(t testing.TB, data []byte) (string, func()) {
	f, err := os.CreateTemp("", "store-*.json")
	assert.NoError(t, err)
	defer f.Close()

	_, err = f.Write(data)
	assert.NoError(t, err)

	fileName := f.Name()

	return fileName, func() {
		assert.NoError(t, os.Remove(fileName), "failed to remove test file %s", fileName)
	}
}

func TestSetupFile(t *testing.T) {
	content := "content"
	fileName, teardown := setupFile(t, []byte(content))

	assert.FileExists(t, fileName)

	data, err := os.ReadFile(fileName)
	assert.NoError(t, err)
	assert.Equal(t, content, string(data))

	teardown()
	assert.NoFileExists(t, fileName)
}

func TestJSONStorage(t *testing.T) {
	t.Run("NewJSONStorage", func(t *testing.T) {
		data := []byte(`[{"amount":250,"nws":"needs","domain":"Groceries","name":"Groceries - supermarket"}]`)
		fileName, teardown := setupFile(t, data)
		defer teardown()

		js, err := NewJSONStorage(fileName)
		assert.NoError(t, err)
		assert.Len(t, js.Expenses, 1)
		assert.Equal(t, 250, js.Expenses[0].Amount)
	})

	t.Run("NewJSONStorage initializes with empty state when file doesn't exist", func(t *testing.T) {
		fileName, teardown := setupFile(t, []byte(""))
		teardown()

		js, err := NewJSONStorage(fileName)
		assert.NoError(t, err)
		assert.Len(t, js.Expenses, 0)
	})

	t.Run("Add creates file even if it didn't exist", func(t *testing.T) {
		fileName, teardown := setupFile(t, []byte(""))
		teardown()

		js, err := NewJSONStorage(fileName)
		assert.NoError(t, err)
		addSampleExpense(t, js)
		assert.Len(t, js.Expenses, 1)
		assert.FileExists(t, fileName)
		assert.NoError(t, os.Remove(fileName))
	})

	t.Run("Add", func(t *testing.T) {
		fileName, teardown := setupFile(t, []byte("[]"))
		defer teardown()

		js, err := NewJSONStorage(fileName)
		assert.NoError(t, err)

		addSampleExpense(t, js)
		assert.Len(t, js.Expenses, 1)

		data, err := os.ReadFile(fileName)
		assert.Contains(t, string(data), `"amount":250`)
	})

	t.Run("List", func(t *testing.T) {
		data := []byte(`[{"amount":250,"nws":"needs","domain":"Groceries","name":"test"}]`)
		fileName, teardown := setupFile(t, data)
		defer teardown()

		js, err := NewJSONStorage(fileName)
		assert.NoError(t, err)

		assert.Len(t, js.List(), 1)
		assert.Equal(t, "test", js.List()[0].Name)
	})
}

func addSampleExpense(t testing.TB, es ExpenseStorage) Expense {
	expense := Expense{
		Amount: 250,
		NWS:    NWS_Need,
		Domain: "Groceries",
		Name:   "Groceries - supermarket",
	}
	assert.NoError(t, es.Add(expense))
	return expense
}

func TestFormat(t *testing.T) {
	t.Run("csv", func(t *testing.T) {
		t.Run("JSONStorage", func(t *testing.T) {
			data := []byte(`[{"amount":250,"nws":"needs","domain":"Groceries","name":"Groceries - supermarket"}]`)
			fileName, teardown := setupFile(t, data)
			defer teardown()

			jsonStorage, err := NewJSONStorage(fileName)
			assert.NoError(t, err)
			csv := FormatCSV(jsonStorage)

			assert.Equal(t, "name,amount,nws\nGroceries - supermarket,250,needs\n", csv)
		})
	})
}
