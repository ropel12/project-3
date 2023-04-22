// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"

	event "github.com/ropel12/project-3/app/features/event"
	mock "github.com/stretchr/testify/mock"

	multipart "mime/multipart"
)

// EventService is an autogenerated mock type for the EventService type
type EventService struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, req, file
func (_m *EventService) Create(ctx context.Context, req event.ReqCreate, file multipart.File) (int, error) {
	ret := _m.Called(ctx, req, file)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, event.ReqCreate, multipart.File) (int, error)); ok {
		return rf(ctx, req, file)
	}
	if rf, ok := ret.Get(0).(func(context.Context, event.ReqCreate, multipart.File) int); ok {
		r0 = rf(ctx, req, file)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, event.ReqCreate, multipart.File) error); ok {
		r1 = rf(ctx, req, file)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MyEvent provides a mock function with given fields: ctx, uid, limit, page
func (_m *EventService) MyEvent(ctx context.Context, uid int, limit int, page int) (*event.Response, error) {
	ret := _m.Called(ctx, uid, limit, page)

	var r0 *event.Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int, int) (*event.Response, error)); ok {
		return rf(ctx, uid, limit, page)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, int, int) *event.Response); ok {
		r0 = rf(ctx, uid, limit, page)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*event.Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, int, int) error); ok {
		r1 = rf(ctx, uid, limit, page)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewEventService interface {
	mock.TestingT
	Cleanup(func())
}

// NewEventService creates a new instance of EventService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewEventService(t mockConstructorTestingTNewEventService) *EventService {
	mock := &EventService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
