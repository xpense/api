package model

// type arguments struct {
// 	timestamp       time.Time
// 	amount          uint64
// 	transactionType TransactionType
// }

// func TestNewWithInvalidArguments(t *testing.T) {
// 	testCases := []struct {
// 		desc string
// 		args arguments
// 		err  error
// 	}{
// 		{
// 			desc: "Invalid Amount",
// 			args: arguments{
// 				timestamp:       time.Unix(0, 0),
// 				amount:          0,
// 				transactionType: "expense",
// 			},
// 			err: ErrorAmount,
// 		},
// 		{
// 			desc: "Invalid Transaction Type",
// 			args: arguments{
// 				timestamp:       time.Unix(0, 0),
// 				amount:          1,
// 				transactionType: "profit",
// 			},
// 			err: ErrorType,
// 		},
// 	}

// 	for _, tC := range testCases {
// 		t.Run(tC.desc, func(t *testing.T) {
// 			args := tC.args
// 			if err := TransactionValidateBody(args.timestamp, args.amount, args.transactionType); err != nil {
// 				t.Errorf("Invalid New Transaction should be nil")
// 			}

// 			if err != tC.err {
// 				t.Errorf("Returned error message should be: '" + tC.err.Error() +
// 					"', but it's instead: '" + err.Error() + "'")
// 			}
// 		})
// 	}
// }

// func TestNewWithValidArguments(t *testing.T) {
// 	now := time.Now()
// 	args := arguments{
// 		timestamp:       now,
// 		amount:          1,
// 		transactionType: "expense",
// 	}

// 	got, err := TransactionValidateBody(args.timestamp, args.amount, args.transactionType)
// 	if err != nil {
// 		t.Errorf("New Transaction shouldn return an error when provided with valid arguments")
// 	}

// 	switch true {
// 	case got == nil:
// 		t.Errorf("New Transaction shouldn't be nil")
// 	case !got.Timestamp.Equal(now):
// 		t.Errorf("New Transaction 'Date' field should be set to provided date")
// 	case got.Type != args.transactionType:
// 		t.Errorf("New Transaction 'TransactionType' field should be set to provided transaction type")
// 	}
// }
