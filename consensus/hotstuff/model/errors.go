package model

import (
	"errors"
	"fmt"

	"github.com/onflow/flow-go/model/flow"
)

var (
	ErrUnverifiableBlock = errors.New("block proposal can't be verified, because its view is above the finalized view, but its QC is below the finalized view")
	ErrInvalidFormat     = errors.New("invalid signature format")
	ErrInvalidSignature  = errors.New("invalid signature")
)

/* ****************************** NoVoteError ****************************** */

// NoVoteError contains the reason of why the voter didn't vote for a block proposal.
type NoVoteError struct {
	Msg string
}

func (e NoVoteError) Error() string { return e.Msg }

// IsNoVoteError returns whether an error is NoVoteError
func IsNoVoteError(err error) bool {
	var e NoVoteError
	return errors.As(err, &e)
}

/* *************************** ConfigurationError ************************** */

type ConfigurationError struct {
	Msg string
}

func (e ConfigurationError) Error() string { return e.Msg }

type MissingBlockError struct {
	View    uint64
	BlockID flow.Identifier
}

/* *************************** MissingBlockError *************************** */

func (e MissingBlockError) Error() string {
	return fmt.Sprintf("missing Block at view %d with ID %v", e.View, e.BlockID)
}

// IsMissingBlockError returns whether an error is MissingBlockError
func IsMissingBlockError(err error) bool {
	var e MissingBlockError
	return errors.As(err, &e)
}

/* *************************** InvalidBlockError *************************** */

type InvalidBlockError struct {
	BlockID flow.Identifier
	View    uint64
	Err     error
}

func (e InvalidBlockError) Error() string {
	return fmt.Sprintf("invalid block %x at view %d: %s", e.BlockID, e.View, e.Err.Error())
}

// IsInvalidBlockError returns whether an error is InvalidBlockError
func IsInvalidBlockError(err error) bool {
	var e InvalidBlockError
	return errors.As(err, &e)
}

func (e InvalidBlockError) Unwrap() error {
	return e.Err
}

/* *************************** InvalidVoteError **************************** */

type InvalidVoteError struct {
	VoteID flow.Identifier
	View   uint64
	Err    error
}

func (e InvalidVoteError) Error() string {
	return fmt.Sprintf("invalid vote %x for view %d: %s", e.VoteID, e.View, e.Err.Error())
}

// IsInvalidVoteError returns whether an error is InvalidVoteError
func IsInvalidVoteError(err error) bool {
	var e InvalidVoteError
	return errors.As(err, &e)
}

func (e InvalidVoteError) Unwrap() error {
	return e.Err
}

func NewInvalidVoteErrorf(vote *Vote, msg string, args ...interface{}) error {
	return InvalidVoteError{
		VoteID: vote.ID(),
		View:   vote.View,
		Err:    fmt.Errorf(msg, args...),
	}
}

/* ****************** ByzantineThresholdExceededError ********************** */

// ByzantineThresholdExceededError is raised if HotStuff detects malicious conditions which
// prove a Byzantine threshold of consensus replicas has been exceeded.
// Per definition, the byzantine threshold is exceeded is there are byzantine consensus
// replicas with _at least_ 1/3 stake.
type ByzantineThresholdExceededError struct {
	Evidence string
}

func (e ByzantineThresholdExceededError) Error() string {
	return e.Evidence
}

/* **************************** DoubleVoteError **************************** */

type DoubleVoteError struct {
	FirstVote       *Vote
	ConflictingVote *Vote
	err             error
}

func (e DoubleVoteError) Error() string {
	return e.err.Error()
}

// IsDoubleVoteError returns whether an error is DoubleVoteError
func IsDoubleVoteError(err error) bool {
	var e DoubleVoteError
	return errors.As(err, &e)
}

// AsDoubleVoteError determines whether the given error is a DoubleVoteError
// (potentially wrapped). It follows the same semantics as a checked type cast.
func AsDoubleVoteError(err error) (*DoubleVoteError, bool) {
	var e DoubleVoteError
	ok := errors.As(err, &e)
	if ok {
		return &e, true
	}
	return nil, false
}

func (e DoubleVoteError) Unwrap() error {
	return e.err
}

func NewDoubleVoteErrorf(firstVote, conflictingVote *Vote, msg string, args ...interface{}) error {
	return DoubleVoteError{
		FirstVote:       firstVote,
		ConflictingVote: conflictingVote,
		err:             fmt.Errorf(msg, args...),
	}
}

/* ************************* DuplicatedSignerError ************************* */

// DuplicatedSignerError indicates that a signature from the same node ID has already been added
type DuplicatedSignerError struct {
	err error
}

func NewDuplicatedSignerError(err error) error {
	return DuplicatedSignerError{err}
}

func NewDuplicatedSignerErrorf(msg string, args ...interface{}) error {
	return DuplicatedSignerError{err: fmt.Errorf(msg, args...)}
}

func (e DuplicatedSignerError) Error() string { return e.err.Error() }
func (e DuplicatedSignerError) Unwrap() error { return e.err }

// IsDuplicatedSignerError returns whether err is an DuplicatedSignerError
func IsDuplicatedSignerError(err error) bool {
	var e DuplicatedSignerError
	return errors.As(err, &e)
}

/* ********************* InvalidSignatureIncludedError ********************* */

// InvalidSignatureIncludedError indicates that some signatures, included via TrustedAdd, are invalid
type InvalidSignatureIncludedError struct {
	err error
}

func NewInvalidSignatureIncludedError(err error) error {
	return InvalidSignatureIncludedError{err}
}

func NewInvalidSignatureIncludedErrorf(msg string, args ...interface{}) error {
	return InvalidSignatureIncludedError{fmt.Errorf(msg, args...)}
}

func (e InvalidSignatureIncludedError) Error() string { return e.err.Error() }
func (e InvalidSignatureIncludedError) Unwrap() error { return e.err }

// IsInvalidSignatureIncludedError returns whether err is an InvalidSignatureIncludedError
func IsInvalidSignatureIncludedError(err error) bool {
	var e InvalidSignatureIncludedError
	return errors.As(err, &e)
}

/* ********************** InsufficientSignaturesError ********************** */

// InsufficientSignaturesError indicates that not enough signatures have been stored to complete the operation.
type InsufficientSignaturesError struct {
	err error
}

func NewInsufficientSignaturesError(err error) error {
	return InsufficientSignaturesError{err}
}

func NewInsufficientSignaturesErrorf(msg string, args ...interface{}) error {
	return InsufficientSignaturesError{fmt.Errorf(msg, args...)}
}

func (e InsufficientSignaturesError) Error() string { return e.err.Error() }
func (e InsufficientSignaturesError) Unwrap() error { return e.err }

// IsInsufficientSignaturesError returns whether err is an InsufficientSignaturesError
func IsInsufficientSignaturesError(err error) bool {
	var e InsufficientSignaturesError
	return errors.As(err, &e)
}

/* ********************** InsufficientSignaturesError ********************** */

// InvalidSignerError indicates that the signer is not authorized or unknown
type InvalidSignerError struct {
	err error
}

func NewInvalidSignerError(err error) error {
	return InvalidSignerError{err}
}

func NewInvalidSignerErrorf(msg string, args ...interface{}) error {
	return InvalidSignerError{fmt.Errorf(msg, args...)}
}

func (e InvalidSignerError) Error() string { return e.err.Error() }
func (e InvalidSignerError) Unwrap() error { return e.err }

// IsInvalidSignerError returns whether err is an InvalidSignerError
func IsInvalidSignerError(err error) bool {
	var e InvalidSignerError
	return errors.As(err, &e)
}
