// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import flow "github.com/dapperlabs/flow-go/model/flow"
import mock "github.com/stretchr/testify/mock"

// Transactions is an autogenerated mock type for the Transactions type
type Transactions struct {
	mock.Mock
}

// ByFingerprint provides a mock function with given fields: fingerprint
func (_m *Transactions) ByFingerprint(fingerprint flow.Fingerprint) (*flow.Transaction, error) {
	ret := _m.Called(fingerprint)

	var r0 *flow.Transaction
	if rf, ok := ret.Get(0).(func(flow.Fingerprint) *flow.Transaction); ok {
		r0 = rf(fingerprint)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(flow.Fingerprint) error); ok {
		r1 = rf(fingerprint)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Insert provides a mock function with given fields: tx
func (_m *Transactions) Insert(tx *flow.Transaction) error {
	ret := _m.Called(tx)

	var r0 error
	if rf, ok := ret.Get(0).(func(*flow.Transaction) error); ok {
		r0 = rf(tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
