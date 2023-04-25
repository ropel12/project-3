// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	entities "github.com/ropel12/project-3/app/entities"
	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"
)

// TransactionRepo is an autogenerated mock type for the TransactionRepo type
type TransactionRepo struct {
	mock.Mock
}

// Create provides a mock function with given fields: db, cart
func (_m *TransactionRepo) Create(db *gorm.DB, cart entities.Carts) error {
	ret := _m.Called(db, cart)

	var r0 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, entities.Carts) error); ok {
		r0 = rf(db, cart)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateTransaction provides a mock function with given fields: db, data
func (_m *TransactionRepo) CreateTransaction(db *gorm.DB, data entities.Transaction) error {
	ret := _m.Called(db, data)

	var r0 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, entities.Transaction) error); ok {
		r0 = rf(db, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetCart provides a mock function with given fields: db, uid
func (_m *TransactionRepo) GetCart(db *gorm.DB, uid int) ([]entities.Carts, error) {
	ret := _m.Called(db, uid)

	var r0 []entities.Carts
	var r1 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, int) ([]entities.Carts, error)); ok {
		return rf(db, uid)
	}
	if rf, ok := ret.Get(0).(func(*gorm.DB, int) []entities.Carts); ok {
		r0 = rf(db, uid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entities.Carts)
		}
	}

	if rf, ok := ret.Get(1).(func(*gorm.DB, int) error); ok {
		r1 = rf(db, uid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDetailUser provides a mock function with given fields: db, uid
func (_m *TransactionRepo) GetDetailUser(db *gorm.DB, uid int) *entities.User {
	ret := _m.Called(db, uid)

	var r0 *entities.User
	if rf, ok := ret.Get(0).(func(*gorm.DB, int) *entities.User); ok {
		r0 = rf(db, uid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.User)
		}
	}

	return r0
}

type mockConstructorTestingTNewTransactionRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewTransactionRepo creates a new instance of TransactionRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTransactionRepo(t mockConstructorTestingTNewTransactionRepo) *TransactionRepo {
	mock := &TransactionRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
