// Code generated by mockery v2.12.1. DO NOT EDIT.

package mock

import (
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// EngineMetrics is an autogenerated mock type for the EngineMetrics type
type EngineMetrics struct {
	mock.Mock
}

// MessageHandled provides a mock function with given fields: engine, messages
func (_m *EngineMetrics) MessageHandled(engine string, messages string) {
	_m.Called(engine, messages)
}

// MessageReceived provides a mock function with given fields: engine, message
func (_m *EngineMetrics) MessageReceived(engine string, message string) {
	_m.Called(engine, message)
}

// MessageSent provides a mock function with given fields: engine, message
func (_m *EngineMetrics) MessageSent(engine string, message string) {
	_m.Called(engine, message)
}

// NewEngineMetrics creates a new instance of EngineMetrics. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewEngineMetrics(t testing.TB) *EngineMetrics {
	mock := &EngineMetrics{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
