package operations

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONStorage(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		js, err := NewJSONStorage(nil, nil)
		assert.NoError(t, err)

		addSampleExpense(t, js)
		assert.Len(t, js.Expenses, 1)
	})

	t.Run("Save", func(t *testing.T) {
		storage := &bytes.Buffer{}

		js, err := NewJSONStorage(nil, storage)
		assert.NoError(t, err)

		addSampleExpense(t, js)
		assert.NoError(t, js.Save())
		assert.True(t, json.Valid(storage.Bytes()), "source file backing JSONStorage should contain valid json: %s", storage.String())
	})

	t.Run("List", func(t *testing.T) {
		js, err := NewJSONStorage(nil, nil)
		assert.NoError(t, err)

		exp := addSampleExpense(t, js)
		assert.Equal(t, []Expense{exp}, js.List())
	})

	t.Run("NewJSONStorage", func(t *testing.T) {
		storage := bytes.NewBuffer([]byte(`[{"amount":250,"nws":"needs","domain":"Groceries","name":"Groceries - supermarket"}]`))

		js, err := NewJSONStorage(storage, nil)
		assert.NoError(t, err)

		assert.Len(t, js.List(), 1)
		assert.Equal(t, 250, js.List()[0].Amount)
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
		jsonStorage, err := NewJSONStorage(nil, nil)
		assert.NoError(t, err)
		addSampleExpense(t, jsonStorage)
		csv := FormatCSV(jsonStorage)

		assert.Equal(t, `name,amount,nws
Groceries - supermarket,250,needs
`, csv)
	})
}
