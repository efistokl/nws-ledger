package operations

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONStorage(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		jsonStorage := NewJSONStorage(nil)
		var expenseStorage ExpenseStorage = jsonStorage

		addSampleExpense(t, expenseStorage)
		assert.Len(t, jsonStorage.Expenses, 1)
	})

	t.Run("Save", func(t *testing.T) {
		storage := &bytes.Buffer{}

		var expenseStorage ExpenseStorage = NewJSONStorage(storage)

		addSampleExpense(t, expenseStorage)
		assert.NoError(t, expenseStorage.Save())
		assert.True(t, json.Valid(storage.Bytes()), "source file backing JSONStorage should contain valid json: %s", storage.String())
	})

	t.Run("List", func(t *testing.T) {
		var expenseStorage ExpenseStorage = NewJSONStorage(nil)
		exp := addSampleExpense(t, expenseStorage)
		assert.Equal(t, []Expense{exp}, expenseStorage.List())
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
		var expenseStorage ExpenseStorage = NewJSONStorage(nil)
		addSampleExpense(t, expenseStorage)
		csv := FormatCSV(expenseStorage)

		assert.Equal(t, `name,amount,nws
Groceries - supermarket,250,needs
`, csv)
	})
}
