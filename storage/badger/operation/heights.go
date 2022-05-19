// (c) 2019 Dapper Labs - ALL RIGHTS RESERVED

package operation

import (
	badger "github.com/dgraph-io/badger/v2"
)

func InsertRootHeight(height uint64) func(*badger.Txn) error {
	return insert(makePrefix(codeRootHeight), height)
}

func RetrieveRootHeight(height *uint64) func(*badger.Txn) error {
	return retrieve(makePrefix(codeRootHeight), height)
}

func InsertFinalizedHeight(height uint64) func(*badger.Txn) error {
	return insert(makePrefix(codeFinalizedHeight), height)
}

func UpdateFinalizedHeight(height uint64) func(*badger.Txn) error {
	return update(makePrefix(codeFinalizedHeight), height)
}

func RetrieveFinalizedHeight(height *uint64) func(*badger.Txn) error {
	return retrieve(makePrefix(codeFinalizedHeight), height)
}

func InsertSealedHeight(height uint64) func(*badger.Txn) error {
	return insert(makePrefix(codeSealedHeight), height)
}

func UpdateSealedHeight(height uint64) func(*badger.Txn) error {
	return update(makePrefix(codeSealedHeight), height)
}

func RetrieveSealedHeight(height *uint64) func(*badger.Txn) error {
	return retrieve(makePrefix(codeSealedHeight), height)
}

func InsertLastCompleteBlockHeight(height uint64) func(*badger.Txn) error {
	return insert(makePrefix(codeLastCompleteBlockHeight), height)
}

func UpdateLastCompleteBlockHeight(height uint64) func(*badger.Txn) error {
	return update(makePrefix(codeLastCompleteBlockHeight), height)
}

func RetrieveLastCompleteBlockHeight(height *uint64) func(*badger.Txn) error {
	return retrieve(makePrefix(codeLastCompleteBlockHeight), height)
}
