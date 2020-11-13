// Code generated by mockery v2.3.0. DO NOT EDIT.

package spies

import (
	utils "expense-api/utils"

	mock "github.com/stretchr/testify/mock"
)

// JWTServiceSpy is an autogenerated mock type for the JWTService type
type JWTServiceSpy struct {
	mock.Mock
}

// CreateJWT provides a mock function with given fields: email
func (_m *JWTServiceSpy) CreateJWT(email string) (string, error) {
	ret := _m.Called(email)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(email)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ValidateJWT provides a mock function with given fields: tokenString
func (_m *JWTServiceSpy) ValidateJWT(tokenString string) (*utils.CustomClaims, error) {
	ret := _m.Called(tokenString)

	var r0 *utils.CustomClaims
	if rf, ok := ret.Get(0).(func(string) *utils.CustomClaims); ok {
		r0 = rf(tokenString)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*utils.CustomClaims)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(tokenString)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
