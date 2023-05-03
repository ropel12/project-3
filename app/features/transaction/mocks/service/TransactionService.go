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

// CreateTransaction provides a mock function with given fields: ctx, req
func (_m *TransactionService) CreateTransaction(ctx context.Context, req transaction.ReqCheckout) (*transaction.Transaction, error) {
	ret := _m.Called(ctx, req)

	var r0 *transaction.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, transaction.ReqCheckout) (*transaction.Transaction, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, transaction.ReqCheckout) *transaction.Transaction); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*transaction.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, transaction.ReqCheckout) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByStatus provides a mock function with given fields: ctx, uid, status
func (_m *TransactionService) GetByStatus(ctx context.Context, uid int, status string) (*transaction.Response, error) {
	ret := _m.Called(ctx, uid, status)

	var r0 *transaction.Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, string) (*transaction.Response, error)); ok {
		return rf(ctx, uid, status)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, string) *transaction.Response); ok {
		r0 = rf(ctx, uid, status)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*transaction.Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, string) error); ok {
		r1 = rf(ctx, uid, status)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCart provides a mock function with given fields: ctx, uid
func (_m *TransactionService) GetCart(ctx context.Context, uid int) (*transaction.Response, error) {
	ret := _m.Called(ctx, uid)

	var r0 *transaction.Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (*transaction.Response, error)); ok {
		return rf(ctx, uid)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) *transaction.Response); ok {
		r0 = rf(ctx, uid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*transaction.Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, uid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDetail provides a mock function with given fields: ctx, invoice, uid
func (_m *TransactionService) GetDetail(ctx context.Context, invoice string, uid int) (*transaction.Response, error) {
	ret := _m.Called(ctx, invoice, uid)

	var r0 *transaction.Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int) (*transaction.Response, error)); ok {
		return rf(ctx, invoice, uid)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int) *transaction.Response); ok {
		r0 = rf(ctx, invoice, uid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*transaction.Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int) error); ok {
		r1 = rf(ctx, invoice, uid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetHistoryByuid provides a mock function with given fields: ctx, uid, page, limit
func (_m *TransactionService) GetHistoryByuid(ctx context.Context, uid int, page int, limit int) (*transaction.Response, error) {
	ret := _m.Called(ctx, uid, page, limit)

	var r0 *transaction.Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int, int) (*transaction.Response, error)); ok {
		return rf(ctx, uid, page, limit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, int, int) *transaction.Response); ok {
		r0 = rf(ctx, uid, page, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*transaction.Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, int, int) error); ok {
		r1 = rf(ctx, uid, page, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTickets provides a mock function with given fields: ctx, invoice, uid
func (_m *TransactionService) GetTickets(ctx context.Context, invoice string, uid int) (*transaction.Response, error) {
	ret := _m.Called(ctx, invoice, uid)

	var r0 *transaction.Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int) (*transaction.Response, error)); ok {
		return rf(ctx, invoice, uid)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int) *transaction.Response); ok {
		r0 = rf(ctx, invoice, uid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*transaction.Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int) error); ok {
		r1 = rf(ctx, invoice, uid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateStatus provides a mock function with given fields: ctx, status, invoice
func (_m *TransactionService) UpdateStatus(ctx context.Context, status string, invoice string) error {
	ret := _m.Called(ctx, status, invoice)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, status, invoice)
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
