package operations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type StubExpenseStorage struct {
	expenses []Expense
}

func (s *StubExpenseStorage) Add(_ Expense) error {
	return nil
}

func (s *StubExpenseStorage) List() []Expense {
	return s.expenses
}

func (s *StubExpenseStorage) Summary() SummaryByNWS {
	return nil
}

func TestFormat(t *testing.T) {
	store := &StubExpenseStorage{
		[]Expense{{
			Amount: 250,
			NWS:    NWS_Needs,
			Name:   "Groceries - supermarket",
		}},
	}

	csv := FormatCSV(store)
	assert.Equal(t, "name,amount,nws\nGroceries - supermarket,250,needs\n", csv)
}
