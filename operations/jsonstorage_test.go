package operations

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupFile(t testing.TB, data []byte) (*os.File, func()) {
	f, err := os.CreateTemp("", "store-*.json")
	assert.NoError(t, err)

	_, err = f.Write(data)
	assert.NoError(t, err)

	return f, func() {
		assert.NoError(t, f.Close())
		assert.NoError(t, os.Remove(f.Name()), "failed to remove test file %s", f.Name())
	}
}

func TestSetupFile(t *testing.T) {
	content := "content"
	file, teardown := setupFile(t, []byte(content))

	assert.FileExists(t, file.Name())

	data, err := os.ReadFile(file.Name())
	assert.NoError(t, err)
	assert.Equal(t, content, string(data))

	teardown()
	assert.NoFileExists(t, file.Name())
}

func TestJSONStorage(t *testing.T) {
	t.Run("NewJSONStorage", func(t *testing.T) {
		data := []byte(`[{"amount":250,"nws":"needs","domain":"Groceries","name":"Groceries - supermarket"}]`)
		file, teardown := setupFile(t, data)
		defer teardown()

		js, err := NewJSONStorage(file)
		assert.NoError(t, err)
		assert.Len(t, js.Expenses, 1)
		assert.Equal(t, 250, js.Expenses[0].Amount)
	})

	t.Run("NewJSONStorage initializes with empty state when file is empty", func(t *testing.T) {
		file, teardown := setupFile(t, []byte(""))
		defer teardown()

		js, err := NewJSONStorage(file)
		assert.NoError(t, err)
		assert.Len(t, js.Expenses, 0)
	})

	t.Run("Add", func(t *testing.T) {
		file, teardown := setupFile(t, []byte("[]"))
		defer teardown()

		js, err := NewJSONStorage(file)
		assert.NoError(t, err)

		addSampleExpense(t, js)
		assert.Len(t, js.Expenses, 1)

		data, err := os.ReadFile(file.Name())
		assert.NoError(t, err)
		assert.Contains(t, string(data), `"amount":250`)
	})

	t.Run("List", func(t *testing.T) {
		data := []byte(`[{"amount":250,"nws":"needs","domain":"Groceries","name":"test"}]`)
		file, teardown := setupFile(t, data)
		defer teardown()

		js, err := NewJSONStorage(file)
		assert.NoError(t, err)

		assert.Len(t, js.List(), 1)
		assert.Equal(t, "test", js.List()[0].Name)
	})

	t.Run("Summary", func(t *testing.T) {
		file, teardown := setupFile(t, []byte("[]"))
		defer teardown()

		js, err := NewJSONStorage(file)
		assert.NoError(t, err)

		expenses := []Expense{
			{
				Amount: 200,
				NWS:    NWS_Needs,
			},
			{
				Amount: 150,
				NWS:    NWS_Needs,
			},
			{
				Amount: 500,
				NWS:    NWS_Wants,
			},
			{
				Amount: 50,
				NWS:    NWS_Savings,
			},
		}

		for _, e := range expenses {
			assert.NoError(t, js.Add(e))
		}

		summary := js.Summary()
		assert.Equal(t, 350, summary[NWS_Needs])
		assert.Equal(t, 500, summary[NWS_Wants])
		assert.Equal(t, 50, summary[NWS_Savings])
	})
}

func addSampleExpense(t testing.TB, es ExpenseStorage) Expense {
	expense := Expense{
		Amount: 250,
		NWS:    NWS_Needs,
		Domain: "Groceries",
		Name:   "Groceries - supermarket",
	}
	assert.NoError(t, es.Add(expense))
	return expense
}
