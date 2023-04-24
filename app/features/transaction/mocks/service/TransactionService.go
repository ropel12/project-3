// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	transaction "github.com/ropel12/project-3/app/features/transaction"
)

// TransactionService is an autogenerated mock type for the TransactionService type
type TransactionService struct {
	mock.Mock
}

// CreateCart provides a mock function with given fields: ctx, req
func (_m *TransactionService) CreateCart(ctx context.Context, req transaction.ReqCart) error {
	ret := _m.Called(ctx, req)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, transaction.ReqCart) error); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewTransactionService interface {
	mock.TestingT
	Cleanup(func())
}

// NewTransactionService creates a new instance of TransactionService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTransactionService(t mockConstructorTestingTNewTransactionService) *TransactionService {
	mock := &TransactionService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
