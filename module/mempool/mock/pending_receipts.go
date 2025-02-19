// Code generated by mockery v2.12.1. DO NOT EDIT.

package mempool

import (
	flow "github.com/onflow/flow-go/model/flow"

	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// PendingReceipts is an autogenerated mock type for the PendingReceipts type
type PendingReceipts struct {
	mock.Mock
}

// Add provides a mock function with given fields: receipt
func (_m *PendingReceipts) Add(receipt *flow.ExecutionReceipt) bool {
	ret := _m.Called(receipt)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*flow.ExecutionReceipt) bool); ok {
		r0 = rf(receipt)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ByPreviousResultID provides a mock function with given fields: previousReusltID
func (_m *PendingReceipts) ByPreviousResultID(previousReusltID flow.Identifier) []*flow.ExecutionReceipt {
	ret := _m.Called(previousReusltID)

	var r0 []*flow.ExecutionReceipt
	if rf, ok := ret.Get(0).(func(flow.Identifier) []*flow.ExecutionReceipt); ok {
		r0 = rf(previousReusltID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*flow.ExecutionReceipt)
		}
	}

	return r0
}

// PruneUpToHeight provides a mock function with given fields: height
func (_m *PendingReceipts) PruneUpToHeight(height uint64) error {
	ret := _m.Called(height)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint64) error); ok {
		r0 = rf(height)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Rem provides a mock function with given fields: receiptID
func (_m *PendingReceipts) Rem(receiptID flow.Identifier) bool {
	ret := _m.Called(receiptID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(flow.Identifier) bool); ok {
		r0 = rf(receiptID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// NewPendingReceipts creates a new instance of PendingReceipts. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewPendingReceipts(t testing.TB) *PendingReceipts {
	mock := &PendingReceipts{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
