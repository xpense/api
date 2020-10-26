package transaction_test

import (
	"expense-api/transaction"
	"testing"
	"time"
)

type arguments struct {
	id              uint64
	date            *time.Time
	amount          uint64
	transactionType transaction.Type
}

func TestNewWithInvalidArguments(t *testing.T) {
	testCases := []struct {
		desc         string
		args         arguments
		errorMessage string
	}{
		{
			desc: "Invalid ID",
			args: arguments{
				id:              0,
				date:            nil,
				amount:          0,
				transactionType: 20,
			},
			errorMessage: "Cannot create new Transaction with a non positive id.",
		},
		{
			desc: "Invalid Amount",
			args: arguments{
				id:              1,
				date:            nil,
				amount:          0,
				transactionType: 20,
			},
			errorMessage: "Cannot create new Transaction with an amount of 0.",
		},
		{
			desc: "Invalid Transaction Type",
			args: arguments{
				id:              1,
				date:            nil,
				amount:          1,
				transactionType: 20,
			},
			errorMessage: "Invalid Transaction Type. Please use either 'Income' or 'Expense'.",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			args := tC.args
			got, err := transaction.New(args.id, args.date, args.amount, args.transactionType)

			if got != nil {
				t.Errorf("Invalid New Transaction should be nil")
			}

			if err.Error() != tC.errorMessage {
				t.Errorf("Returned error message should be: '" + tC.errorMessage +
					"', but it's instead: '" + err.Error() + "'")
			}
		})
	}
}

func TestNewWithValidArguments(t *testing.T) {
	now := time.Now()
	args := arguments{
		id:              1,
		date:            &now,
		amount:          1,
		transactionType: 1,
	}

	got, err := transaction.New(args.id, args.date, args.amount, args.transactionType)
	if err != nil {
		t.Errorf("New Transaction shouldn return an error when provided with valid arguments")
	}

	switch true {
	case got == nil:
		t.Errorf("New Transaction shouldn't be nil")
	case got.ID != args.id:
		t.Errorf("New Transaction 'ID' field should be set to provided id")
	case got.Date != now:
		t.Errorf("New Transaction 'Date' field should be set to provided date")
	case got.Type != args.transactionType:
		t.Errorf("New Transaction 'TransactionType' field should be set to provided transaction type")
	}
}
