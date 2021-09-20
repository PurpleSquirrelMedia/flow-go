// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import (
	mock "github.com/stretchr/testify/mock"

	network "github.com/onflow/flow-go/network"
)

// ReadyDoneAwareNetwork is an autogenerated mock type for the ReadyDoneAwareNetwork type
type ReadyDoneAwareNetwork struct {
	mock.Mock
}

// Done provides a mock function with given fields:
func (_m *ReadyDoneAwareNetwork) Done() <-chan struct{} {
	ret := _m.Called()

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// Ready provides a mock function with given fields:
func (_m *ReadyDoneAwareNetwork) Ready() <-chan struct{} {
	ret := _m.Called()

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// Register provides a mock function with given fields: channel, engine
func (_m *ReadyDoneAwareNetwork) Register(channel network.Channel, engine network.Engine) (network.Conduit, error) {
	ret := _m.Called(channel, engine)

	var r0 network.Conduit
	if rf, ok := ret.Get(0).(func(network.Channel, network.Engine) network.Conduit); ok {
		r0 = rf(channel, engine)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(network.Conduit)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(network.Channel, network.Engine) error); ok {
		r1 = rf(channel, engine)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
