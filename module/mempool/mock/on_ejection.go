// Code generated by mockery v2.12.1. DO NOT EDIT.

package mempool

import (
	flow "github.com/onflow/flow-go/model/flow"

	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// OnEjection is an autogenerated mock type for the OnEjection type
type OnEjection struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0
func (_m *OnEjection) Execute(_a0 flow.Entity) {
	_m.Called(_a0)
}

// NewOnEjection creates a new instance of OnEjection. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewOnEjection(t testing.TB) *OnEjection {
	mock := &OnEjection{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
