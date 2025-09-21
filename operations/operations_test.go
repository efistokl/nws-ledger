package operations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type StubExpenseStorage struct {
	expenses []Expense
	summary  SummaryByNWS
}

func (s *StubExpenseStorage) Add(_ Expense) error {
	return nil
}

func (s *StubExpenseStorage) List() []Expense {
	return s.expenses
}

func (s *StubExpenseStorage) Summary() SummaryByNWS {
	return s.summary
}

func TestFormat(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		store := &StubExpenseStorage{
			[]Expense{{
				Amount: 250,
				NWS:    NWS_Needs,
				Name:   "Groceries - supermarket",
				Domain: "groceries",
			}},
			nil,
		}

		csv := FormatCSVList(store)
		assert.Equal(t, "name,amount,nws,domain\nGroceries - supermarket,250,needs,groceries\n", csv)
	})

	t.Run("Summary", func(t *testing.T) {
		store := &StubExpenseStorage{
			nil,
			SummaryByNWS{
				NWS_Needs:   300,
				NWS_Wants:   100,
				NWS_Savings: 50,
			},
		}

		csv := FormatCSVSummary(store)
		assert.Equal(t, "nws,amount\nneeds,300\nwants,100\nsavings,50\ntotal,450\n", csv)
	})
}

func TestValidateNWS(t *testing.T) {
	assert.NoError(t, ValidateNWS("wants"))
	assert.NoError(t, ValidateNWS("needs"))
	assert.NoError(t, ValidateNWS("savings"))
	assert.Error(t, ValidateNWS("deeds"))
}
