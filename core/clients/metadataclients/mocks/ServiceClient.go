// Code generated by mockery v1.0.0
package mocks

import mock "github.com/stretchr/testify/mock"
import models "github.com/edgexfoundry/edgex-go/core/domain/models"

// ServiceClient is an autogenerated mock type for the ServiceClient type
type MockServiceClient struct {
	mock.Mock
}

// Add provides a mock function with given fields: ds
func (_m *MockServiceClient) Add(ds *models.DeviceService) (string, error) {
	ret := _m.Called(ds)

	var r0 string
	if rf, ok := ret.Get(0).(func(*models.DeviceService) string); ok {
		r0 = rf(ds)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.DeviceService) error); ok {
		r1 = rf(ds)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeviceServiceForName provides a mock function with given fields: name
func (_m *MockServiceClient) DeviceServiceForName(name string) (models.DeviceService, error) {
	ret := _m.Called(name)

	var r0 models.DeviceService
	if rf, ok := ret.Get(0).(func(string) models.DeviceService); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(models.DeviceService)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateLastConnected provides a mock function with given fields: id, time
func (_m *MockServiceClient) UpdateLastConnected(id string, time int64) error {
	ret := _m.Called(id, time)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, int64) error); ok {
		r0 = rf(id, time)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateLastReported provides a mock function with given fields: id, time
func (_m *MockServiceClient) UpdateLastReported(id string, time int64) error {
	ret := _m.Called(id, time)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, int64) error); ok {
		r0 = rf(id, time)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
