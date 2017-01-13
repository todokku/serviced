package mocks

import datastore "github.com/control-center/serviced/datastore"
import mock "github.com/stretchr/testify/mock"
import service "github.com/control-center/serviced/domain/service"
import servicestatemanager "github.com/control-center/serviced/scheduler/servicestatemanager"

// Facade is an autogenerated mock type for the Facade type
type Facade struct {
	mock.Mock
}

// GetTenantIDs provides a mock function with given fields: ctx
func (_m *Facade) GetTenantIDs(ctx datastore.Context) ([]string, error) {
	ret := _m.Called(ctx)

	var r0 []string
	if rf, ok := ret.Get(0).(func(datastore.Context) []string); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(datastore.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ScheduleServiceBatch provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *Facade) ScheduleServiceBatch(_a0 datastore.Context, _a1 []*service.Service, _a2 string, _a3 service.DesiredState) (int, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 int
	if rf, ok := ret.Get(0).(func(datastore.Context, []*service.Service, string, service.DesiredState) int); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(datastore.Context, []*service.Service, string, service.DesiredState) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateService provides a mock function with given fields: ctx, svc
func (_m *Facade) UpdateService(ctx datastore.Context, svc service.Service) error {
	ret := _m.Called(ctx, svc)

	var r0 error
	if rf, ok := ret.Get(0).(func(datastore.Context, service.Service) error); ok {
		r0 = rf(ctx, svc)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WaitSingleService provides a mock function with given fields: _a0, _a1, _a2
func (_m *Facade) WaitSingleService(_a0 *service.Service, _a1 service.DesiredState, _a2 <-chan interface{}) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(*service.Service, service.DesiredState, <-chan interface{}) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

var _ servicestatemanager.Facade = (*Facade)(nil)
