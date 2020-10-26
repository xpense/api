package model

import (
	"testing"
	"time"
)

type arguments struct {
	timestamp       time.Time
	amount          uint64
	transactionType TransactionType
}

func TestNewWithInvalidArguments(t *testing.T) {
	testCases := []struct {
		desc string
		args arguments
		err  error
	}{
		{
			desc: "Invalid Amount",
			args: arguments{
				timestamp:       time.Unix(0, 0),
				amount:          0,
				transactionType: "expense",
			},
			err: ErrorAmount,
		},
		{
			desc: "Invalid Transaction Type",
			args: arguments{
				timestamp:       time.Unix(0, 0),
				amount:          1,
				transactionType: "profit",
			},
			err: ErrorType,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			args := tC.args
			err := TransactionValidateCreateBody(args.timestamp, args.amount, args.transactionType)

			if err != tC.err {
				t.Errorf("Returned error message should be: '" + tC.err.Error() +
					"', but it's instead: '" + err.Error() + "'")
			}
		})
	}
}

func TestNewWithValidArguments(t *testing.T) {
	now := time.Now()
	args := arguments{
		timestamp:       now,
		amount:          1,
		transactionType: "expense",
	}

	err := TransactionValidateCreateBody(args.timestamp, args.amount, args.transactionType)
	if err != nil {
		t.Errorf("New Transaction shouldn return an error when provided with valid arguments")
	}
}
