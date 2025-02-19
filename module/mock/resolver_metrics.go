// Code generated by mockery v2.12.1. DO NOT EDIT.

package mock

import (
	mock "github.com/stretchr/testify/mock"

	testing "testing"

	time "time"
)

// ResolverMetrics is an autogenerated mock type for the ResolverMetrics type
type ResolverMetrics struct {
	mock.Mock
}

// DNSLookupDuration provides a mock function with given fields: duration
func (_m *ResolverMetrics) DNSLookupDuration(duration time.Duration) {
	_m.Called(duration)
}

// OnDNSCacheHit provides a mock function with given fields:
func (_m *ResolverMetrics) OnDNSCacheHit() {
	_m.Called()
}

// OnDNSCacheInvalidated provides a mock function with given fields:
func (_m *ResolverMetrics) OnDNSCacheInvalidated() {
	_m.Called()
}

// OnDNSCacheMiss provides a mock function with given fields:
func (_m *ResolverMetrics) OnDNSCacheMiss() {
	_m.Called()
}

// OnDNSLookupRequestDropped provides a mock function with given fields:
func (_m *ResolverMetrics) OnDNSLookupRequestDropped() {
	_m.Called()
}

// NewResolverMetrics creates a new instance of ResolverMetrics. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewResolverMetrics(t testing.TB) *ResolverMetrics {
	mock := &ResolverMetrics{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
