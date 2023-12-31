// Code generated by mockery v2.20.0. DO NOT EDIT.

package banking

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockTransactionStorage is an autogenerated mock type for the TransactionStorage type
type MockTransactionStorage struct {
	mock.Mock
}

// Set provides a mock function with given fields: ctx, t
func (_m *MockTransactionStorage) Set(ctx context.Context, t Transaction) error {
	ret := _m.Called(ctx, t)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, Transaction) error); ok {
		r0 = rf(ctx, t)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewMockTransactionStorage interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockTransactionStorage creates a new instance of MockTransactionStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockTransactionStorage(t mockConstructorTestingTNewMockTransactionStorage) *MockTransactionStorage {
	mock := &MockTransactionStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
