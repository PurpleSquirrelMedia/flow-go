package irrecoverable

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"

	"go.uber.org/atomic"
)

// Signaler sends the error out.
type Signaler struct {
	errChan   chan error
	errThrown *atomic.Bool
}

func NewSignaler() *Signaler {
	return &Signaler{
		errChan:   make(chan error, 1),
		errThrown: atomic.NewBool(false),
	}
}

// Error returns the Signaler's error channel.
func (s *Signaler) Error() <-chan error {
	return s.errChan
}

// Throw is a narrow drop-in replacement for panic, log.Fatal, log.Panic, etc
// anywhere there's something connected to the error channel. It only sends
// the first error it is called with to the error channel, and logs subsequent
// errors as unhandled.
func (s *Signaler) Throw(err error) {
	defer runtime.Goexit()

	if s.errThrown.CAS(false, true) {
		s.errChan <- err
		close(s.errChan)
	} else {
		// TODO: we simply log the unhandled irrecoverable for now, but we should probably
		// find a way to actually handle them.
		log.New(os.Stderr, "", log.LstdFlags).Println(fmt.Errorf("unhandled irrecoverable: %w", err))
	}
}

// We define a constrained interface to provide a drop-in replacement for context.Context
// including in interfaces that compose it.
type SignalerContext interface {
	context.Context
	Throw(err error) // delegates to the signaler
	sealed()         // private, to constrain builder to using WithSignaler
}

// private, to force context derivation / WithSignaler
type signalerCtx struct {
	context.Context
	signaler *Signaler
}

func (sc signalerCtx) sealed() {}

// Drop-in replacement for panic, log.Fatal, log.Panic, etc
// to use when we are able to get an SignalerContext and thread it down in the component
func (sc signalerCtx) Throw(err error) {
	sc.signaler.Throw(err)
}

// the One True Way of getting a SignalerContext
func WithSignaler(ctx context.Context, sig *Signaler) SignalerContext {
	return signalerCtx{ctx, sig}
}

// If we have an SignalerContext, we can directly ctx.Throw.
//
// But a lot of library methods expect context.Context, & we want to pass the same w/o boilerplate
// Moreover, we could have built with: context.WithCancel(irrecoverable.WithSignaler(ctx, sig)),
// "downcasting" to context.Context. Yet, we can still type-assert and recover.
//
// Throw can be a drop-in replacement anywhere we have a context.Context likely
// to support Irrecoverables. Note: this is not a method
func Throw(ctx context.Context, err error) {
	signalerAbleContext, ok := ctx.(SignalerContext)
	if ok {
		signalerAbleContext.Throw(err)
	}
	// Be spectacular on how this does not -but should- handle irrecoverables:
	log.Fatalf("irrecoverable error signaler not found for context, please implement! Unhandled irrecoverable error %v", err)
}
