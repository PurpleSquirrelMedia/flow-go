package computation

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/rs/zerolog"

	"github.com/onflow/flow-go/engine/execution"
	"github.com/onflow/flow-go/engine/execution/computation/computer"
	"github.com/onflow/flow-go/engine/execution/state/delta"
	"github.com/onflow/flow-go/fvm"
	"github.com/onflow/flow-go/fvm/state"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module"
	"github.com/onflow/flow-go/module/mempool/entity"
	"github.com/onflow/flow-go/state/protocol"
	"github.com/onflow/flow-go/utils/logging"
)

type VirtualMachine interface {
	Run(fvm.Context, fvm.Procedure, state.Ledger) error
	GetAccount(fvm.Context, flow.Address, state.Ledger) (*flow.Account, error)
}

type ComputationManager interface {
	ExecuteScript([]byte, [][]byte, *flow.Header, *delta.View) ([]byte, error)
	ComputeBlock(
		ctx context.Context,
		block *entity.ExecutableBlock,
		view *delta.View,
	) (*execution.ComputationResult, error)
	GetAccount(addr flow.Address, header *flow.Header, view *delta.View) (*flow.Account, error)
}

var DefaultScriptLogThreshold = 1 * time.Second

// Manager manages computation and execution
type Manager struct {
	log                zerolog.Logger
	me                 module.Local
	protoState         protocol.State
	vm                 VirtualMachine
	vmCtx              fvm.Context
	blockComputer      computer.BlockComputer
	scriptLogThreshold time.Duration
}

func New(
	logger zerolog.Logger,
	metrics module.ExecutionMetrics,
	tracer module.Tracer,
	me module.Local,
	protoState protocol.State,
	vm VirtualMachine,
	vmCtx fvm.Context,
	scriptLogThreshold time.Duration,
) (*Manager, error) {
	log := logger.With().Str("engine", "computation").Logger()

	blockComputer, err := computer.NewBlockComputer(
		vm,
		vmCtx,
		metrics,
		tracer,
		log.With().Str("component", "block_computer").Logger(),
	)

	if err != nil {
		return nil, fmt.Errorf("cannot create block computer: %w", err)
	}

	e := Manager{
		log:                log,
		me:                 me,
		protoState:         protoState,
		vm:                 vm,
		vmCtx:              vmCtx,
		blockComputer:      blockComputer,
		scriptLogThreshold: scriptLogThreshold,
	}

	return &e, nil
}

func (e *Manager) ExecuteScript(code []byte, arguments [][]byte, blockHeader *flow.Header, view *delta.View) ([]byte, error) {
	blockCtx := fvm.NewContextFromParent(e.vmCtx, fvm.WithBlockHeader(blockHeader))

	script := fvm.Script(code).WithArguments(arguments...)

	err := func() (err error) {

		start := time.Now()

		defer func() {

			prepareLog := func() *zerolog.Event {

				args := make([]string, 0, len(arguments))
				for _, a := range arguments {
					args = append(args, hex.EncodeToString(a))
				}
				return e.log.Error().
					Hex("script_hex", code).
					Str("args", strings.Join(args[:], ","))
			}

			elapsed := time.Since(start)

			if r := recover(); r != nil {
				prepareLog().
					Interface("recovered", r).
					Msg("script execution caused runtime panic")

				err = fmt.Errorf("cadence runtime error: %s", r)
				return
			}
			if elapsed >= e.scriptLogThreshold {
				prepareLog().
					Dur("duration", elapsed).
					Msg("script execution exceeded threshold")
			}
		}()

		return e.vm.Run(blockCtx, script, view)
	}()

	if err != nil {
		return nil, fmt.Errorf("failed to execute script (internal error): %w", err)
	}

	if script.Err != nil {
		return nil, fmt.Errorf("failed to execute script at block (%s): %s", blockHeader.ID(), script.Err.Error())
	}

	encodedValue, err := jsoncdc.Encode(script.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to encode runtime value: %w", err)
	}

	return encodedValue, nil
}

func (e *Manager) ComputeBlock(
	ctx context.Context,
	block *entity.ExecutableBlock,
	view *delta.View,
) (*execution.ComputationResult, error) {

	e.log.Debug().
		Hex("block_id", logging.Entity(block.Block)).
		Msg("received complete block")

	result, err := e.blockComputer.ExecuteBlock(ctx, block, view)
	if err != nil {
		e.log.Error().
			Hex("block_id", logging.Entity(block.Block)).
			Msg("failed to compute block result")

		return nil, fmt.Errorf("failed to execute block: %w", err)
	}

	e.log.Debug().
		Hex("block_id", logging.Entity(result.ExecutableBlock.Block)).
		Msg("computed block result")

	return result, nil
}

func (e *Manager) GetAccount(address flow.Address, blockHeader *flow.Header, view *delta.View) (*flow.Account, error) {
	blockCtx := fvm.NewContextFromParent(e.vmCtx, fvm.WithBlockHeader(blockHeader))

	account, err := e.vm.GetAccount(blockCtx, address, view)
	if err != nil {
		return nil, fmt.Errorf("failed to get account at block (%s): %w", blockHeader.ID(), err)
	}

	return account, nil
}
