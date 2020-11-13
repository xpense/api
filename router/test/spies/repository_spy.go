// Code generated by mockery v2.3.0. DO NOT EDIT.

package spies

import (
	model "expense-api/model"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// RepositorySpy is an autogenerated mock type for the Repository type
type RepositorySpy struct {
	mock.Mock
}

// TransactionCreate provides a mock function with given fields: timestamp, amount, transactionType
func (_m *RepositorySpy) TransactionCreate(timestamp time.Time, amount uint64, transactionType model.TransactionType) (*model.Transaction, error) {
	ret := _m.Called(timestamp, amount, transactionType)

	var r0 *model.Transaction
	if rf, ok := ret.Get(0).(func(time.Time, uint64, model.TransactionType) *model.Transaction); ok {
		r0 = rf(timestamp, amount, transactionType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(time.Time, uint64, model.TransactionType) error); ok {
		r1 = rf(timestamp, amount, transactionType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransactionDelete provides a mock function with given fields: id
func (_m *RepositorySpy) TransactionDelete(id uint) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TransactionGet provides a mock function with given fields: id
func (_m *RepositorySpy) TransactionGet(id uint) (*model.Transaction, error) {
	ret := _m.Called(id)

	var r0 *model.Transaction
	if rf, ok := ret.Get(0).(func(uint) *model.Transaction); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransactionList provides a mock function with given fields:
func (_m *RepositorySpy) TransactionList() ([]*model.Transaction, error) {
	ret := _m.Called()

	var r0 []*model.Transaction
	if rf, ok := ret.Get(0).(func() []*model.Transaction); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransactionUpdate provides a mock function with given fields: id, timestamp, amount, transactionType
func (_m *RepositorySpy) TransactionUpdate(id uint, timestamp time.Time, amount uint64, transactionType model.TransactionType) (*model.Transaction, error) {
	ret := _m.Called(id, timestamp, amount, transactionType)

	var r0 *model.Transaction
	if rf, ok := ret.Get(0).(func(uint, time.Time, uint64, model.TransactionType) *model.Transaction); ok {
		r0 = rf(id, timestamp, amount, transactionType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint, time.Time, uint64, model.TransactionType) error); ok {
		r1 = rf(id, timestamp, amount, transactionType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserCreate provides a mock function with given fields: firstName, LastName, Email, Password, Salt
func (_m *RepositorySpy) UserCreate(firstName string, LastName string, Email string, Password string, Salt string) (*model.User, error) {
	ret := _m.Called(firstName, LastName, Email, Password, Salt)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(string, string, string, string, string) *model.User); ok {
		r0 = rf(firstName, LastName, Email, Password, Salt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, string, string) error); ok {
		r1 = rf(firstName, LastName, Email, Password, Salt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserDelete provides a mock function with given fields: id
func (_m *RepositorySpy) UserDelete(id uint) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserGet provides a mock function with given fields: id
func (_m *RepositorySpy) UserGet(id uint) (*model.User, error) {
	ret := _m.Called(id)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(uint) *model.User); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserGetWithEmail provides a mock function with given fields: email
func (_m *RepositorySpy) UserGetWithEmail(email string) (*model.User, error) {
	ret := _m.Called(email)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(string) *model.User); ok {
		r0 = rf(email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserUpdate provides a mock function with given fields: id, firstName, LastName, Email
func (_m *RepositorySpy) UserUpdate(id uint, firstName string, LastName string, Email string) (*model.User, error) {
	ret := _m.Called(id, firstName, LastName, Email)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(uint, string, string, string) *model.User); ok {
		r0 = rf(id, firstName, LastName, Email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint, string, string, string) error); ok {
		r1 = rf(id, firstName, LastName, Email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}