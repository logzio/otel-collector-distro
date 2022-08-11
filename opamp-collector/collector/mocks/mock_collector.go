// Code generated by mockery v2.12.2. DO NOT EDIT.

package mocks

import (
	context "context"

	collector "github.com/yotamloe/otel-collector-distro/opamp-collector/collector"

	mock "github.com/stretchr/testify/mock"

	testing "testing"

	zap "go.uber.org/zap"
)

// MockCollector is an autogenerated mock type for the Collector type
type MockCollector struct {
	mock.Mock
}

// GetLoggingOpts provides a mock function with given fields:
func (_m *MockCollector) GetLoggingOpts() []zap.Option {
	ret := _m.Called()

	var r0 []zap.Option
	if rf, ok := ret.Get(0).(func() []zap.Option); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]zap.Option)
		}
	}

	return r0
}

// Restart provides a mock function with given fields: _a0
func (_m *MockCollector) Restart(_a0 context.Context) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Run provides a mock function with given fields: _a0
func (_m *MockCollector) Run(_a0 context.Context) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetLoggingOpts provides a mock function with given fields: _a0
func (_m *MockCollector) SetLoggingOpts(_a0 []zap.Option) {
	_m.Called(_a0)
}

// Status provides a mock function with given fields:
func (_m *MockCollector) Status() <-chan *collector.Status {
	ret := _m.Called()

	var r0 <-chan *collector.Status
	if rf, ok := ret.Get(0).(func() <-chan *collector.Status); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan *collector.Status)
		}
	}

	return r0
}

// Stop provides a mock function with given fields:
func (_m *MockCollector) Stop() {
	_m.Called()
}

// NewMockCollector creates a new instance of MockCollector. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockCollector(t testing.TB) *MockCollector {
	mock := &MockCollector{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
