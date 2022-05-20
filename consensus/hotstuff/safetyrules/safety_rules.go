package safetyrules

import (
	"fmt"

	"github.com/onflow/flow-go/consensus/hotstuff"
	"github.com/onflow/flow-go/consensus/hotstuff/model"
	"github.com/onflow/flow-go/model/flow"
)

// SafetyRules is a dedicated module that is responsible for persisting chain safety by producing
// votes and timeouts. It follows voting and timeout rules for creating votes and timeouts respectively.
// Caller can be sure that created vote or timeout doesn't break safety and can be used in consensus process.
// SafetyRules relies on hotstuff.Persister to store latest state of hotstuff.SafetyData.
//
// The voting rules implemented by SafetyRules are:
// 1. Replicas vote strictly in increasing rounds
// 2. Each block has to include a TC or a QC from the previous round.
//   a. [Happy path] If the previous round resulted in a QC then new QC should extend it.
//   b. [Recovery path] If the previous round did *not* result in a QC, the leader of the
//   subsequent round *must* include a valid TC for the previous round in its block.
//
// NOT safe for concurrent use.
type SafetyRules struct {
	signer     hotstuff.Signer
	persist    hotstuff.Persister
	committee  hotstuff.DynamicCommittee // only produce votes when we are valid committee members
	safetyData *hotstuff.SafetyData
}

var _ hotstuff.SafetyRules = (*SafetyRules)(nil)

// New creates a new SafetyRules instance
func New(
	signer hotstuff.Signer,
	persist hotstuff.Persister,
	committee hotstuff.DynamicCommittee,
	safetyData *hotstuff.SafetyData,
) *SafetyRules {
	return &SafetyRules{
		signer:     signer,
		persist:    persist,
		committee:  committee,
		safetyData: safetyData,
	}
}

// ProduceVote will make a decision on whether it will vote for the given proposal, the returned
// error indicates whether to vote or not.
// In order to ensure that only safe proposals are being voted on, by checking if proposer is a valid committee member and
// proposal complies with voting rules.
// We expect that only well-formed proposals with valid signatures are submitted for voting.
// The curView is taken as input to ensure SafetyRules will only vote for proposals at current view and prevent double voting.
// Returns:
//  * (vote, nil): On the _first_ block for the current view that is safe to vote for.
//    Subsequently, voter does _not_ vote for any other block with the same (or lower) view.
//  * (nil, model.NoVoteError): If the voter decides that it does not want to vote for the given block.
//    This is a sentinel error and _expected_ during normal operation.
// All other errors are unexpected and potential symptoms of uncovered edge cases or corrupted internal state (fatal).
func (r *SafetyRules) ProduceVote(proposal *model.Proposal, curView uint64) (*model.Vote, error) {
	block := proposal.Block
	// sanity checks:
	if curView != block.View {
		return nil, fmt.Errorf("expecting block for current view %d, but block's view is %d", curView, block.View)
	}

	err := r.IsSafeToVote(proposal)
	if err != nil {
		return nil, fmt.Errorf("not safe to vote for proposal %x: %w", proposal.Block.BlockID, err)
	}

	// we expect that only valid proposals are submitted for voting
	// we need to make sure that proposer is not ejected to decide to vote or not
	_, err = r.committee.IdentityByBlock(block.BlockID, block.ProposerID)
	if model.IsInvalidSignerError(err) {
		// the proposer must be ejected since the proposal has already been validated,
		// which ensures that the proposer was a valid committee member at the start of the epoch
		return nil, model.NewNoVoteErrorf("not voting - proposer ejected")
	}
	if err != nil {
		return nil, fmt.Errorf("could not retrieve proposer Identity %x at block %x: %w", block.ProposerID, block.BlockID, err)
	}

	// Do not produce a vote for blocks where we are not a valid committee member.
	// HotStuff will ask for a vote for the first block of the next epoch, even if we
	// have zero weight in the next epoch. Such vote can't be used to produce valid QCs.
	_, err = r.committee.IdentityByBlock(block.BlockID, r.committee.Self())
	if model.IsInvalidSignerError(err) {
		return nil, model.NewNoVoteErrorf("not voting committee member for block %x", block.BlockID)
	}
	if err != nil {
		return nil, fmt.Errorf("could not get self identity: %w", err)
	}

	vote, err := r.signer.CreateVote(block)
	if err != nil {
		return nil, fmt.Errorf("could not vote for block: %w", err)
	}

	// vote for the current view has been produced, update safetyData
	r.safetyData.HighestAcknowledgedView = curView
	if r.safetyData.LockedOneChainView < block.QC.View {
		r.safetyData.LockedOneChainView = block.QC.View
	}

	err = r.persist.PutSafetyData(r.safetyData)
	if err != nil {
		return nil, fmt.Errorf("could not persist safety data: %w", err)
	}

	return vote, nil
}

// ProduceTimeout takes current view, highest locally known QC and TC and decides whether to produce timeout for current view.
// Returns:
//  * (timeout, nil): On the _first_ block for the current view that is safe to vote for.
//    Subsequently, voter does _not_ vote for any other block with the same (or lower) view.
//  * (nil, model.NoTimeoutError): If the safety module decides that it is not safe to timeout under current conditions.
//    This is a sentinel error and _expected_ during normal operation.
// All other errors are unexpected and potential symptoms of uncovered edge cases or corrupted internal state (fatal).
func (r *SafetyRules) ProduceTimeout(curView uint64, highestQC *flow.QuorumCertificate, lastViewTC *flow.TimeoutCertificate) (*model.TimeoutObject, error) {
	lastTimeout := r.safetyData.LastTimeout
	if lastTimeout != nil && lastTimeout.View == curView {
		return lastTimeout, nil
	}

	if !r.IsSafeToTimeout(curView, highestQC, lastViewTC) {
		return nil, model.NewNoTimeoutErrorf("not safe to time out under current conditions")
	}

	timeout, err := r.signer.CreateTimeout(curView, highestQC, lastViewTC)
	if err != nil {
		return nil, fmt.Errorf("could not create timeout at view %d: %w", curView, err)
	}

	r.safetyData.HighestAcknowledgedView = curView
	r.safetyData.LastTimeout = timeout

	err = r.persist.PutSafetyData(r.safetyData)
	if err != nil {
		return nil, fmt.Errorf("could not persist safety data: %w", err)
	}

	return timeout, nil
}

// IsSafeToVote checks if this proposal is valid in terms of voting rules, if voting for this proposal won't break safety rules.
func (r *SafetyRules) IsSafeToVote(proposal *model.Proposal) error {
	blockView := proposal.Block.View
	qcView := proposal.Block.QC.View

	// block's view must be larger than the view of the included QC
	if blockView <= qcView {
		return fmt.Errorf("block's view %d must be larger than the view of the included QC %d", blockView, qcView)
	}

	// This check satisfies voting rule 1
	// 1. Replicas vote strictly in increasing rounds,
	// block's view must be greater than the view that we have voted for
	if blockView <= r.safetyData.HighestAcknowledgedView {
		return model.NewNoVoteErrorf("not safe to vote, we have already voted for this view: %d <= %d",
			blockView, r.safetyData.HighestAcknowledgedView)
	}

	// This check satisfies voting rule 2a:
	// 2a. [Happy path] If the previous round resulted in a QC then new QC should extend it.
	if blockView == qcView+1 {
		return nil
	}

	return r.IsSafeToExtend(blockView, qcView, proposal.LastViewTC)
}

// IsSafeToExtend performs safety checks if proposal can be extended in case of recovery path, we will call
// this function only if previous round resulted in TC, to know if it's safe to extend such proposal.
func (r *SafetyRules) IsSafeToExtend(blockView, qcView uint64, lastViewTC *flow.TimeoutCertificate) error {
	// These checks satisfy voting rule 2b:
	// [Recovery Path] If the previous round did *not* result in a QC, the leader of the
	// subsequent round *must* include a valid TC for the previous round in its block.
	if lastViewTC == nil {
		return fmt.Errorf("block's view %d is not sequential with included QC view %d, last view TC not included", blockView, qcView)
	}
	if blockView != lastViewTC.View+1 {
		return fmt.Errorf("last view TC %d is not sequential for block %d", lastViewTC.View, blockView)
	}
	if qcView < lastViewTC.TOHighestQC.View {
		return fmt.Errorf("QC's view %d should be at least %d", qcView, lastViewTC.TOHighestQC.View)
	}
	return nil
}

// IsSafeToTimeout checks if it's safe to timeout with proposed data, if timing out won't break safety rules.
// highestQC is the valid QC with the greatest view that we have observed.
// lastViewTC is the TC for the previous view (might be nil)
func (r *SafetyRules) IsSafeToTimeout(curView uint64, highestQC *flow.QuorumCertificate, lastViewTC *flow.TimeoutCertificate) bool {
	if highestQC.View < r.safetyData.LockedOneChainView ||
		curView+1 <= r.safetyData.HighestAcknowledgedView ||
		curView <= highestQC.View {
		return false
	}

	return (curView == highestQC.View+1) || (curView == lastViewTC.View+1)
}
