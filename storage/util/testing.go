package util

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dgraph-io/badger/v2"
	"github.com/stretchr/testify/require"

	"github.com/dapperlabs/flow-go/module/metrics"
	storage "github.com/dapperlabs/flow-go/storage/badger"
	"github.com/dapperlabs/flow-go/utils/unittest"
)

func StorageLayer(t testing.TB, db *badger.DB) (*storage.Headers, *storage.Guarantees, *storage.Seals, *storage.Index, *storage.Payloads, *storage.Blocks, *storage.EpochSetups, *storage.EpochCommits) {
	metrics := metrics.NewNoopCollector()
	headers := storage.NewHeaders(metrics, db)
	guarantees := storage.NewGuarantees(metrics, db)
	seals := storage.NewSeals(metrics, db)
	index := storage.NewIndex(metrics, db)
	payloads := storage.NewPayloads(db, index, guarantees, seals)
	blocks := storage.NewBlocks(db, headers, payloads)
	setups := storage.NewEpochSetups(metrics, db)
	commits := storage.NewEpochCommits(metrics, db)
	return headers, guarantees, seals, index, payloads, blocks, setups, commits
}

func RunWithStorageLayer(t testing.TB, f func(*badger.DB, *storage.Headers, *storage.Guarantees, *storage.Seals, *storage.Index, *storage.Payloads, *storage.Blocks, *storage.EpochSetups, *storage.EpochCommits)) {
	unittest.RunWithBadgerDB(t, func(db *badger.DB) {
		headers, guarantees, seals, index, payloads, blocks, setups, commits := StorageLayer(t, db)
		f(db, headers, guarantees, seals, index, payloads, blocks, setups, commits)
	})
}

func CreateFiles(t *testing.T, dir string, names ...string) {
	for _, name := range names {
		file, err := os.Create(filepath.Join(dir, name))
		require.NoError(t, err)
		err = file.Close()
		require.NoError(t, err)
	}
}
