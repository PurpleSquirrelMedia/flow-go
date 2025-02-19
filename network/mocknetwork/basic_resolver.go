// Code generated by mockery v2.12.1. DO NOT EDIT.

package mocknetwork

import (
	context "context"
	net "net"

	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// BasicResolver is an autogenerated mock type for the BasicResolver type
type BasicResolver struct {
	mock.Mock
}

// LookupIPAddr provides a mock function with given fields: _a0, _a1
func (_m *BasicResolver) LookupIPAddr(_a0 context.Context, _a1 string) ([]net.IPAddr, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []net.IPAddr
	if rf, ok := ret.Get(0).(func(context.Context, string) []net.IPAddr); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]net.IPAddr)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LookupTXT provides a mock function with given fields: _a0, _a1
func (_m *BasicResolver) LookupTXT(_a0 context.Context, _a1 string) ([]string, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []string
	if rf, ok := ret.Get(0).(func(context.Context, string) []string); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewBasicResolver creates a new instance of BasicResolver. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewBasicResolver(t testing.TB) *BasicResolver {
	mock := &BasicResolver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
