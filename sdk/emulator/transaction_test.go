package emulator_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dapperlabs/flow-go/crypto"
	"github.com/dapperlabs/flow-go/language/runtime"
	"github.com/dapperlabs/flow-go/model/flow"
	encodingValues "github.com/dapperlabs/flow-go/sdk/abi/encoding/values"
	"github.com/dapperlabs/flow-go/sdk/abi/types"
	"github.com/dapperlabs/flow-go/sdk/abi/values"
	"github.com/dapperlabs/flow-go/sdk/emulator"
	"github.com/dapperlabs/flow-go/sdk/keys"
)

func TestSubmitTransaction(t *testing.T) {
	b, err := emulator.NewBlockchain()
	require.NoError(t, err)

	addTwoScript, _ := deployAndGenerateAddTwoScript(t, b)

	tx1 := flow.Transaction{TransactionBody: flow.TransactionBody{
		Script:             []byte(addTwoScript),
		ReferenceBlockHash: nil,
		Nonce:              getNonce(),
		ComputeLimit:       10,
		PayerAccount:       b.RootAccountAddress(),
		ScriptAccounts:     []flow.Address{b.RootAccountAddress()},
	}}

	sig, err := keys.SignTransaction(tx1, b.RootKey())
	assert.NoError(t, err)

	tx1.AddSignature(b.RootAccountAddress(), sig)

	// Submit tx1
	err = b.AddTransaction(tx1)
	assert.NoError(t, err)

	// Execute tx1
	result, err := b.ExecuteNextTransaction()
	assert.NoError(t, err)
	assert.True(t, result.Succeeded())

	// tx1 status becomes TransactionFinalized
	tx2, err := b.GetTransaction(tx1.Hash())
	assert.NoError(t, err)
	assert.Equal(t, flow.TransactionFinalized, tx2.Status)
}

// TODO: Add test case for missing ReferenceBlockHash
func TestSubmitInvalidTransaction(t *testing.T) {
	b, err := emulator.NewBlockchain()
	require.NoError(t, err)

	addTwoScript, _ := deployAndGenerateAddTwoScript(t, b)

	t.Run("EmptyTransaction", func(t *testing.T) {
		// Create empty transaction (no required fields)
		tx := flow.Transaction{}

		sig, err := keys.SignTransaction(tx, b.RootKey())
		assert.NoError(t, err)

		tx.AddSignature(b.RootAccountAddress(), sig)

		// Submit tx
		err = b.AddTransaction(tx)
		assert.IsType(t, err, &emulator.ErrInvalidTransaction{})
	})

	t.Run("MissingScript", func(t *testing.T) {
		// Create transaction with no Script field
		tx := flow.Transaction{TransactionBody: flow.TransactionBody{
			ReferenceBlockHash: nil,
			Nonce:              getNonce(),
			ComputeLimit:       10,
			PayerAccount:       b.RootAccountAddress(),
		}}

		sig, err := keys.SignTransaction(tx, b.RootKey())
		assert.NoError(t, err)

		tx.AddSignature(b.RootAccountAddress(), sig)

		// Submit tx
		err = b.AddTransaction(tx)
		assert.IsType(t, err, &emulator.ErrInvalidTransaction{})
	})

	t.Run("MissingNonce", func(t *testing.T) {
		// Create transaction with no Nonce field
		tx := flow.Transaction{TransactionBody: flow.TransactionBody{
			Script:             []byte(addTwoScript),
			ReferenceBlockHash: nil,
			ComputeLimit:       10,
			PayerAccount:       b.RootAccountAddress(),
		}}

		sig, err := keys.SignTransaction(tx, b.RootKey())
		assert.NoError(t, err)

		tx.AddSignature(b.RootAccountAddress(), sig)

		// Submit tx
		err = b.AddTransaction(tx)
		assert.IsType(t, err, &emulator.ErrInvalidTransaction{})
	})

	t.Run("MissingComputeLimit", func(t *testing.T) {
		// Create transaction with no ComputeLimit field
		tx := flow.Transaction{TransactionBody: flow.TransactionBody{
			Script:             []byte(addTwoScript),
			ReferenceBlockHash: nil,
			Nonce:              getNonce(),
			PayerAccount:       b.RootAccountAddress(),
		}}

		sig, err := keys.SignTransaction(tx, b.RootKey())
		assert.NoError(t, err)

		tx.AddSignature(b.RootAccountAddress(), sig)

		// Submit tx
		err = b.AddTransaction(tx)
		assert.IsType(t, err, &emulator.ErrInvalidTransaction{})
	})

	t.Run("MissingPayerAccount", func(t *testing.T) {
		// Create transaction with no PayerAccount field
		tx := flow.Transaction{TransactionBody: flow.TransactionBody{
			Script:             []byte(addTwoScript),
			ReferenceBlockHash: nil,
			Nonce:              getNonce(),
			ComputeLimit:       10,
		}}

		sig, err := keys.SignTransaction(tx, b.RootKey())
		assert.NoError(t, err)

		tx.AddSignature(b.RootAccountAddress(), sig)

		// Submit tx
		err = b.AddTransaction(tx)
		assert.IsType(t, err, &emulator.ErrInvalidTransaction{})
	})
}

func TestSubmitDuplicateTransaction(t *testing.T) {
	b, err := emulator.NewBlockchain()
	require.NoError(t, err)

	addTwoScript, _ := deployAndGenerateAddTwoScript(t, b)

	accountAddress := b.RootAccountAddress()

	tx := flow.Transaction{TransactionBody: flow.TransactionBody{
		Script:             []byte(addTwoScript),
		ReferenceBlockHash: nil,
		Nonce:              getNonce(),
		ComputeLimit:       10,
		PayerAccount:       accountAddress,
		ScriptAccounts:     []flow.Address{accountAddress},
	}}

	sig, err := keys.SignTransaction(tx, b.RootKey())
	assert.NoError(t, err)

	tx.AddSignature(accountAddress, sig)

	// Submit tx
	err = b.AddTransaction(tx)
	assert.NoError(t, err)

	result, err := b.ExecuteNextTransaction()
	assert.NoError(t, err)
	assert.True(t, result.Succeeded())

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// Submit same tx again (errors)
	err = b.AddTransaction(tx)
	assert.IsType(t, err, &emulator.ErrDuplicateTransaction{})
}

func TestSubmitTransactionReverted(t *testing.T) {
	b, err := emulator.NewBlockchain()
	require.NoError(t, err)

	tx1 := flow.Transaction{TransactionBody: flow.TransactionBody{
		Script:             []byte("invalid script"),
		ReferenceBlockHash: nil,
		Nonce:              getNonce(),
		ComputeLimit:       10,
		PayerAccount:       b.RootAccountAddress(),
		ScriptAccounts:     []flow.Address{b.RootAccountAddress()},
	}}

	sig, err := keys.SignTransaction(tx1, b.RootKey())
	assert.NoError(t, err)

	tx1.AddSignature(b.RootAccountAddress(), sig)

	// Submit invalid tx1
	err = b.AddTransaction(tx1)
	assert.NoError(t, err)

	result, err := b.ExecuteNextTransaction()
	assert.NoError(t, err)
	assert.True(t, result.Reverted())

	// tx1 status becomes TransactionReverted
	tx, err := b.GetTransaction(tx1.Hash())
	assert.NoError(t, err)
	assert.Equal(t, flow.TransactionReverted, tx.Status)
}

func TestSubmitTransactionScriptAccounts(t *testing.T) {
	b, err := emulator.NewBlockchain()
	require.NoError(t, err)

	privateKeyA := b.RootKey()

	privateKeyB, _ := keys.GeneratePrivateKey(keys.ECDSA_P256_SHA3_256,
		[]byte("elephant ears space cowboy octopus rodeo potato cannon pineapple"))
	publicKeyB := privateKeyB.PublicKey(keys.PublicKeyWeightThreshold)

	accountAddressA := b.RootAccountAddress()
	accountAddressB, err := b.CreateAccount([]flow.AccountPublicKey{publicKeyB}, nil, getNonce())
	assert.NoError(t, err)

	t.Run("TooManyAccountsForScript", func(t *testing.T) {
		// script only supports one account
		script := []byte(`
		  transaction {
		    prepare(signer: Account) {}
		  }
		`)

		// create transaction with two accounts
		tx := flow.Transaction{TransactionBody: flow.TransactionBody{
			Script:             script,
			ReferenceBlockHash: nil,
			Nonce:              getNonce(),
			ComputeLimit:       10,
			PayerAccount:       accountAddressA,
			ScriptAccounts:     []flow.Address{accountAddressA, accountAddressB},
		}}

		sigA, err := keys.SignTransaction(tx, privateKeyA)
		assert.NoError(t, err)

		sigB, err := keys.SignTransaction(tx, privateKeyB)
		assert.NoError(t, err)

		tx.AddSignature(accountAddressA, sigA)
		tx.AddSignature(accountAddressB, sigB)

		err = b.AddTransaction(tx)
		assert.NoError(t, err)

		result, err := b.ExecuteNextTransaction()
		assert.NoError(t, err)
		assert.True(t, result.Reverted())

		_, err = b.CommitBlock()
		assert.NoError(t, err)
	})

	t.Run("NotEnoughAccountsForScript", func(t *testing.T) {
		// script requires two accounts
		script := []byte(`
		  transaction {
		    prepare(signerA: Account, signerB: Account) {}
		  }
		`)

		// create transaction with two accounts
		tx := flow.Transaction{TransactionBody: flow.TransactionBody{
			Script:             script,
			ReferenceBlockHash: nil,
			Nonce:              getNonce(),
			ComputeLimit:       10,
			PayerAccount:       accountAddressA,
			ScriptAccounts:     []flow.Address{accountAddressA},
		}}

		sig, err := keys.SignTransaction(tx, privateKeyA)
		assert.NoError(t, err)

		tx.AddSignature(accountAddressA, sig)

		err = b.AddTransaction(tx)
		assert.NoError(t, err)

		result, err := b.ExecuteNextTransaction()
		assert.NoError(t, err)
		assert.True(t, result.Reverted())

		_, err = b.CommitBlock()
		assert.NoError(t, err)
	})
}

func TestSubmitTransactionPayerSignature(t *testing.T) {
	t.Run("MissingPayerSignature", func(t *testing.T) {
		b, err := emulator.NewBlockchain()
		require.NoError(t, err)

		addTwoScript, _ := deployAndGenerateAddTwoScript(t, b)

		addressA := flow.HexToAddress("0000000000000000000000000000000000000002")

		tx := flow.Transaction{TransactionBody: flow.TransactionBody{
			Script:             []byte(addTwoScript),
			ReferenceBlockHash: nil,
			Nonce:              getNonce(),
			ComputeLimit:       10,
			PayerAccount:       addressA,
			ScriptAccounts:     []flow.Address{b.RootAccountAddress()},
		}}

		sig, err := keys.SignTransaction(tx, b.RootKey())
		assert.NoError(t, err)

		tx.AddSignature(b.RootAccountAddress(), sig)

		err = b.AddTransaction(tx)
		assert.IsType(t, err, &emulator.ErrMissingSignature{})
	})

	t.Run("InvalidAccount", func(t *testing.T) {
		b, err := emulator.NewBlockchain()
		require.NoError(t, err)

		invalidAddress := flow.HexToAddress("0000000000000000000000000000000000000002")

		tx := flow.Transaction{TransactionBody: flow.TransactionBody{
			Script:             []byte(`transaction { execute { } }`),
			ReferenceBlockHash: nil,
			Nonce:              getNonce(),
			ComputeLimit:       10,
			PayerAccount:       b.RootAccountAddress(),
			ScriptAccounts:     []flow.Address{b.RootAccountAddress()},
		}}

		sig, err := keys.SignTransaction(tx, b.RootKey())
		assert.NoError(t, err)

		tx.AddSignature(invalidAddress, sig)

		err = b.AddTransaction(tx)
		assert.IsType(t, err, &emulator.ErrInvalidSignatureAccount{})
	})

	t.Run("InvalidKeyPair", func(t *testing.T) {
		b, err := emulator.NewBlockchain()
		require.NoError(t, err)

		addTwoScript, _ := deployAndGenerateAddTwoScript(t, b)

		// use key-pair that does not exist on root account
		invalidKey, _ := keys.GeneratePrivateKey(keys.ECDSA_P256_SHA3_256,
			[]byte("invalid key elephant ears space cowboy octopus rodeo potato cannon"))

		tx := flow.Transaction{TransactionBody: flow.TransactionBody{
			Script:             []byte(addTwoScript),
			ReferenceBlockHash: nil,
			Nonce:              getNonce(),
			ComputeLimit:       10,
			PayerAccount:       b.RootAccountAddress(),
			ScriptAccounts:     []flow.Address{b.RootAccountAddress()},
		}}

		sig, err := keys.SignTransaction(tx, invalidKey)
		assert.NoError(t, err)

		tx.AddSignature(b.RootAccountAddress(), sig)

		err = b.AddTransaction(tx)
		assert.IsType(t, err, &emulator.ErrInvalidSignaturePublicKey{})
	})

	t.Run("KeyWeights", func(t *testing.T) {
		b, err := emulator.NewBlockchain()
		require.NoError(t, err)

		privateKeyA, _ := keys.GeneratePrivateKey(keys.ECDSA_P256_SHA3_256,
			[]byte("elephant ears space cowboy octopus rodeo potato cannon pineapple"))
		publicKeyA := privateKeyA.PublicKey(keys.PublicKeyWeightThreshold / 2)

		privateKeyB, _ := keys.GeneratePrivateKey(keys.ECDSA_P256_SHA3_256,
			[]byte("space cowboy elephant ears octopus rodeo potato cannon pineapple"))
		publicKeyB := privateKeyB.PublicKey(keys.PublicKeyWeightThreshold / 2)

		accountAddressA, err := b.CreateAccount([]flow.AccountPublicKey{publicKeyA, publicKeyB}, nil, getNonce())
		assert.NoError(t, err)

		script := []byte(`
		  transaction {
		    prepare(signer: Account) {}
		  }
		`)

		tx := flow.Transaction{TransactionBody: flow.TransactionBody{
			Script:             script,
			ReferenceBlockHash: nil,
			Nonce:              getNonce(),
			ComputeLimit:       10,
			PayerAccount:       accountAddressA,
			ScriptAccounts:     []flow.Address{accountAddressA},
		}}

		t.Run("InsufficientKeyWeight", func(t *testing.T) {
			sig, err := keys.SignTransaction(tx, privateKeyB)
			assert.NoError(t, err)

			tx.AddSignature(accountAddressA, sig)

			err = b.AddTransaction(tx)
			assert.IsType(t, err, &emulator.ErrMissingSignature{})
		})

		t.Run("SufficientKeyWeight", func(t *testing.T) {
			sigA, err := keys.SignTransaction(tx, privateKeyA)
			assert.NoError(t, err)

			sigB, err := keys.SignTransaction(tx, privateKeyB)
			assert.NoError(t, err)

			tx.AddSignature(accountAddressA, sigA)
			tx.AddSignature(accountAddressA, sigB)

			err = b.AddTransaction(tx)
			assert.NoError(t, err)

			result, err := b.ExecuteNextTransaction()
			assert.NoError(t, err)
			assert.True(t, result.Succeeded())
		})
	})
}

func TestSubmitTransactionScriptSignatures(t *testing.T) {
	t.Run("MissingScriptSignature", func(t *testing.T) {
		b, err := emulator.NewBlockchain()
		require.NoError(t, err)

		addressA := flow.HexToAddress("0000000000000000000000000000000000000002")

		addTwoScript, _ := deployAndGenerateAddTwoScript(t, b)

		tx := flow.Transaction{TransactionBody: flow.TransactionBody{
			Script:             []byte(addTwoScript),
			ReferenceBlockHash: nil,
			Nonce:              getNonce(),
			ComputeLimit:       10,
			PayerAccount:       b.RootAccountAddress(),
			ScriptAccounts:     []flow.Address{addressA},
		}}

		sig, err := keys.SignTransaction(tx, b.RootKey())
		assert.NoError(t, err)

		tx.AddSignature(b.RootAccountAddress(), sig)

		err = b.AddTransaction(tx)
		assert.IsType(t, err, &emulator.ErrMissingSignature{})
	})

	t.Run("MultipleAccounts", func(t *testing.T) {
		b, err := emulator.NewBlockchain()
		require.NoError(t, err)

		privateKeyA := b.RootKey()

		privateKeyB, _ := keys.GeneratePrivateKey(keys.ECDSA_P256_SHA3_256,
			[]byte("elephant ears space cowboy octopus rodeo potato cannon pineapple"))
		publicKeyB := privateKeyB.PublicKey(keys.PublicKeyWeightThreshold)

		accountAddressA := b.RootAccountAddress()
		accountAddressB, err := b.CreateAccount([]flow.AccountPublicKey{publicKeyB}, nil, getNonce())
		assert.NoError(t, err)

		multipleAccountScript := []byte(`
		  transaction {
		    prepare(signerA: Account, signerB: Account) {
		      log(signerA.address)
			  log(signerB.address)
		    }
		  }
		`)

		tx := flow.Transaction{TransactionBody: flow.TransactionBody{
			Script:             multipleAccountScript,
			ReferenceBlockHash: nil,
			Nonce:              getNonce(),
			ComputeLimit:       10,
			PayerAccount:       accountAddressA,
			ScriptAccounts:     []flow.Address{accountAddressA, accountAddressB},
		}}

		sigA, err := keys.SignTransaction(tx, privateKeyA)
		assert.NoError(t, err)

		sigB, err := keys.SignTransaction(tx, privateKeyB)
		assert.NoError(t, err)

		tx.AddSignature(accountAddressA, sigA)
		tx.AddSignature(accountAddressB, sigB)

		err = b.AddTransaction(tx)
		assert.NoError(t, err)

		result, err := b.ExecuteNextTransaction()
		assert.NoError(t, err)
		assert.True(t, result.Succeeded())

		assert.Contains(t, result.Logs, fmt.Sprintf("%x", accountAddressA.Bytes()))
		assert.Contains(t, result.Logs, fmt.Sprintf("%x", accountAddressB.Bytes()))
	})
}

func TestGetTransaction(t *testing.T) {
	b, err := emulator.NewBlockchain()
	require.NoError(t, err)

	myEventType := types.Event{
		Fields: []types.Field{
			{
				Identifier: "x",
				Type:       types.Int{},
			},
		},
	}

	eventsScript := `
		pub event MyEvent(x: Int)

		transaction {
		  execute {
		    emit MyEvent(x: 1)
		  }
		}
	`

	tx := flow.Transaction{TransactionBody: flow.TransactionBody{
		Script:       []byte(eventsScript),
		Nonce:        getNonce(),
		ComputeLimit: 10,
		PayerAccount: b.RootAccountAddress(),
	}}

	sig, err := keys.SignTransaction(tx, b.RootKey())
	assert.NoError(t, err)

	tx.AddSignature(b.RootAccountAddress(), sig)

	err = b.AddTransaction(tx)
	assert.NoError(t, err)

	result, err := b.ExecuteNextTransaction()
	assert.NoError(t, err)
	assert.True(t, result.Succeeded())

	t.Run("InvalidHash", func(t *testing.T) {
		_, err := b.GetTransaction(crypto.Hash{})
		assert.IsType(t, &emulator.ErrTransactionNotFound{}, err)
	})

	t.Run("ValidHash", func(t *testing.T) {
		resTx, err := b.GetTransaction(tx.Hash())
		require.NoError(t, err)

		assert.Equal(t, resTx.Status, flow.TransactionFinalized)

		require.Len(t, resTx.Events, 1)
		actualEvent := resTx.Events[0]

		eventValue, err := encodingValues.Decode(myEventType, actualEvent.Payload)
		require.NoError(t, err)

		decodedEvent := eventValue.(values.Event)

		location := runtime.TransactionLocation(tx.Hash())
		eventType := fmt.Sprintf("%s.MyEvent", location.ID())

		assert.Equal(t, tx.Hash(), actualEvent.TxHash)
		assert.Equal(t, eventType, actualEvent.Type)
		assert.Equal(t, uint(0), actualEvent.Index)
		assert.Equal(t, values.NewInt(1), decodedEvent.Fields[0])
	})
}
