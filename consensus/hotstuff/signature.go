package hotstuff

import (
	"github.com/onflow/flow-go/crypto"
	"github.com/onflow/flow-go/model/flow"
)

type RandomBeaconReconstructor interface {
	// Add a share
	TrustedAdd(signerID flow.Identifier, sig crypto.Signature) (bool, error)
	HasSufficientShares() bool
	// assume it has sufficient shares
	Reconstruct() (crypto.Signature, error)
}

// ThresholdSigAggregator aggregates the threshold signatures
type ThresholdSigAggregator interface {
	// TrustedAdd adds an already verified signature.
	// return (true, nil) means the signature has been added
	// return (false, nil) means the signature is a duplication
	TrustedAdd(signerID flow.Identifier, sig crypto.Signature) (bool, error)
	// Aggregate assumes enough shares have been collected, it aggregates the signatures
	// and return the aggregated signature.
	// if called concurrently, only one threshold will be running the aggregation.
	Aggregate() ([]byte, error)
}

// StakingSigAggregator aggregates the staking signatures
type StakingSigAggregator interface {
	// TrustedAdd adds an already verified signature.
	// return (true, nil) means the signature has been added
	// return (false, nil) means the signature is a duplication
	TrustedAdd(signerID flow.Identifier, sig crypto.Signature) (bool, error)

	// Aggregate assumes enough shares have been collected, it aggregates the signatures
	// and return the aggregated signature.
	// if called concurrently, only one threshold will be running the aggregation.
	Aggregate() ([]byte, error)
}
