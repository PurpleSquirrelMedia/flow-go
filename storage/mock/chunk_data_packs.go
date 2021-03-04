// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import (
	flow "github.com/onflow/flow-go/model/flow"
	mock "github.com/stretchr/testify/mock"

	storage "github.com/onflow/flow-go/storage"
)

// ChunkDataPacks is an autogenerated mock type for the ChunkDataPacks type
type ChunkDataPacks struct {
	mock.Mock
}

// BatchStore provides a mock function with given fields: c, batch
func (_m *ChunkDataPacks) BatchStore(c *flow.ChunkDataPack, batch storage.BatchStorage) error {
	ret := _m.Called(c, batch)

	var r0 error
	if rf, ok := ret.Get(0).(func(*flow.ChunkDataPack, storage.BatchStorage) error); ok {
		r0 = rf(c, batch)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ByChunkID provides a mock function with given fields: chunkID
func (_m *ChunkDataPacks) ByChunkID(chunkID flow.Identifier) (*flow.ChunkDataPack, error) {
	ret := _m.Called(chunkID)

	var r0 *flow.ChunkDataPack
	if rf, ok := ret.Get(0).(func(flow.Identifier) *flow.ChunkDataPack); ok {
		r0 = rf(chunkID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.ChunkDataPack)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(chunkID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Remove provides a mock function with given fields: chunkID
func (_m *ChunkDataPacks) Remove(chunkID flow.Identifier) error {
	ret := _m.Called(chunkID)

	var r0 error
	if rf, ok := ret.Get(0).(func(flow.Identifier) error); ok {
		r0 = rf(chunkID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Store provides a mock function with given fields: c
func (_m *ChunkDataPacks) Store(c *flow.ChunkDataPack) error {
	ret := _m.Called(c)

	var r0 error
	if rf, ok := ret.Get(0).(func(*flow.ChunkDataPack) error); ok {
		r0 = rf(c)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
