// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	entities "github.com/ropel12/project-3/app/entities"
	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"

	redis "github.com/go-redis/redis/v8"
)

// EventRepo is an autogenerated mock type for the EventRepo type
type EventRepo struct {
	mock.Mock
}

// Create provides a mock function with given fields: db, data
func (_m *EventRepo) Create(db *gorm.DB, data entities.Event) (*int, error) {
	ret := _m.Called(db, data)

	var r0 *int
	var r1 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, entities.Event) (*int, error)); ok {
		return rf(db, data)
	}
	if rf, ok := ret.Get(0).(func(*gorm.DB, entities.Event) *int); ok {
		r0 = rf(db, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*int)
		}
	}

	if rf, ok := ret.Get(1).(func(*gorm.DB, entities.Event) error); ok {
		r1 = rf(db, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateComment provides a mock function with given fields: db, comment
func (_m *EventRepo) CreateComment(db *gorm.DB, comment entities.UserComments) (*entities.UserComments, error) {
	ret := _m.Called(db, comment)

	var r0 *entities.UserComments
	var r1 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, entities.UserComments) (*entities.UserComments, error)); ok {
		return rf(db, comment)
	}
	if rf, ok := ret.Get(0).(func(*gorm.DB, entities.UserComments) *entities.UserComments); ok {
		r0 = rf(db, comment)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.UserComments)
		}
	}

	if rf, ok := ret.Get(1).(func(*gorm.DB, entities.UserComments) error); ok {
		r1 = rf(db, comment)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateTicket provides a mock function with given fields: db, comment
func (_m *EventRepo) CreateTicket(db *gorm.DB, comment entities.Type) (*entities.Type, error) {
	ret := _m.Called(db, comment)

	var r0 *entities.Type
	var r1 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, entities.Type) (*entities.Type, error)); ok {
		return rf(db, comment)
	}
	if rf, ok := ret.Get(0).(func(*gorm.DB, entities.Type) *entities.Type); ok {
		r0 = rf(db, comment)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Type)
		}
	}

	if rf, ok := ret.Get(1).(func(*gorm.DB, entities.Type) error); ok {
		r1 = rf(db, comment)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: db, id, uid
func (_m *EventRepo) Delete(db *gorm.DB, id int, uid int) error {
	ret := _m.Called(db, id, uid)

	var r0 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, int, int) error); ok {
		r0 = rf(db, id, uid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteTicket provides a mock function with given fields: db, id
func (_m *EventRepo) DeleteTicket(db *gorm.DB, id int) (*entities.Type, error) {
	ret := _m.Called(db, id)

	var r0 *entities.Type
	var r1 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, int) (*entities.Type, error)); ok {
		return rf(db, id)
	}
	if rf, ok := ret.Get(0).(func(*gorm.DB, int) *entities.Type); ok {
		r0 = rf(db, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Type)
		}
	}

	if rf, ok := ret.Get(1).(func(*gorm.DB, int) error); ok {
		r1 = rf(db, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields: db, rds, limit, offset, search
func (_m *EventRepo) GetAll(db *gorm.DB, rds *redis.Client, limit int, offset int, search string) ([]*entities.Event, int, error) {
	ret := _m.Called(db, rds, limit, offset, search)

	var r0 []*entities.Event
	var r1 int
	var r2 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, *redis.Client, int, int, string) ([]*entities.Event, int, error)); ok {
		return rf(db, rds, limit, offset, search)
	}
	if rf, ok := ret.Get(0).(func(*gorm.DB, *redis.Client, int, int, string) []*entities.Event); ok {
		r0 = rf(db, rds, limit, offset, search)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entities.Event)
		}
	}

	if rf, ok := ret.Get(1).(func(*gorm.DB, *redis.Client, int, int, string) int); ok {
		r1 = rf(db, rds, limit, offset, search)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(*gorm.DB, *redis.Client, int, int, string) error); ok {
		r2 = rf(db, rds, limit, offset, search)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetById provides a mock function with given fields: db, id
func (_m *EventRepo) GetById(db *gorm.DB, id int) (*entities.Event, error) {
	ret := _m.Called(db, id)

	var r0 *entities.Event
	var r1 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, int) (*entities.Event, error)); ok {
		return rf(db, id)
	}
	if rf, ok := ret.Get(0).(func(*gorm.DB, int) *entities.Event); ok {
		r0 = rf(db, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Event)
		}
	}

	if rf, ok := ret.Get(1).(func(*gorm.DB, int) error); ok {
		r1 = rf(db, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByUid provides a mock function with given fields: db, rds, uid, limit, offset
func (_m *EventRepo) GetByUid(db *gorm.DB, rds *redis.Client, uid int, limit int, offset int) ([]*entities.Event, int, error) {
	ret := _m.Called(db, rds, uid, limit, offset)

	var r0 []*entities.Event
	var r1 int
	var r2 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, *redis.Client, int, int, int) ([]*entities.Event, int, error)); ok {
		return rf(db, rds, uid, limit, offset)
	}
	if rf, ok := ret.Get(0).(func(*gorm.DB, *redis.Client, int, int, int) []*entities.Event); ok {
		r0 = rf(db, rds, uid, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entities.Event)
		}
	}

	if rf, ok := ret.Get(1).(func(*gorm.DB, *redis.Client, int, int, int) int); ok {
		r1 = rf(db, rds, uid, limit, offset)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(*gorm.DB, *redis.Client, int, int, int) error); ok {
		r2 = rf(db, rds, uid, limit, offset)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// JoinEvent provides a mock function with given fields: db, participant
func (_m *EventRepo) JoinEvent(db *gorm.DB, participant entities.Participants) (*entities.Participants, error) {
	ret := _m.Called(db, participant)

	var r0 *entities.Participants
	var r1 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, entities.Participants) (*entities.Participants, error)); ok {
		return rf(db, participant)
	}
	if rf, ok := ret.Get(0).(func(*gorm.DB, entities.Participants) *entities.Participants); ok {
		r0 = rf(db, participant)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Participants)
		}
	}

	if rf, ok := ret.Get(1).(func(*gorm.DB, entities.Participants) error); ok {
		r1 = rf(db, participant)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: db, event
func (_m *EventRepo) Update(db *gorm.DB, event entities.Event) (*entities.Event, error) {
	ret := _m.Called(db, event)

	var r0 *entities.Event
	var r1 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, entities.Event) (*entities.Event, error)); ok {
		return rf(db, event)
	}
	if rf, ok := ret.Get(0).(func(*gorm.DB, entities.Event) *entities.Event); ok {
		r0 = rf(db, event)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Event)
		}
	}

	if rf, ok := ret.Get(1).(func(*gorm.DB, entities.Event) error); ok {
		r1 = rf(db, event)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewEventRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewEventRepo creates a new instance of EventRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewEventRepo(t mockConstructorTestingTNewEventRepo) *EventRepo {
	mock := &EventRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
