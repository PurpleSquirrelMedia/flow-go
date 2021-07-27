// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import (
	flow "github.com/onflow/flow-go/model/flow"
	mock "github.com/stretchr/testify/mock"
)

// Consumer is an autogenerated mock type for the Consumer type
type Consumer struct {
	mock.Mock
}

// BlockFinalized provides a mock function with given fields: block
func (_m *Consumer) BlockFinalized(block *flow.Header) {
	_m.Called(block)
}

// BlockProcessable provides a mock function with given fields: block
func (_m *Consumer) BlockProcessable(block *flow.Header) {
	_m.Called(block)
}

// EpochCommittedPhaseStarted provides a mock function with given fields: currentEpochCounter, first
func (_m *Consumer) EpochCommittedPhaseStarted(currentEpochCounter uint64, first *flow.Header) {
	_m.Called(currentEpochCounter, first)
}

// EpochSetupPhaseStarted provides a mock function with given fields: currentEpochCounter, first
func (_m *Consumer) EpochSetupPhaseStarted(currentEpochCounter uint64, first *flow.Header) {
	_m.Called(currentEpochCounter, first)
}

// EpochTransition provides a mock function with given fields: newEpochCounter, first
func (_m *Consumer) EpochTransition(newEpochCounter uint64, first *flow.Header) {
	_m.Called(newEpochCounter, first)
}
