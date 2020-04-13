package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dapperlabs/flow/protobuf/go/flow/access"
	"github.com/dapperlabs/testingdock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dapperlabs/flow-go/integration/testnet"
	"github.com/dapperlabs/flow-go/model/flow"
	"github.com/dapperlabs/flow-go/utils/unittest"
)

func TestAccess(t *testing.T) {

	var (
		colNode1 = testnet.NewNodeConfig(flow.RoleCollection, testnet.WithLogLevel("error"))
		colNode2 = testnet.NewNodeConfig(flow.RoleCollection, testnet.WithLogLevel("error"))
		conNode1 = testnet.NewNodeConfig(flow.RoleConsensus, testnet.WithLogLevel("error"))
		conNode2 = testnet.NewNodeConfig(flow.RoleConsensus, testnet.WithLogLevel("error"))
		conNode3 = testnet.NewNodeConfig(flow.RoleConsensus, testnet.WithLogLevel("error"))
		exeNode  = testnet.NewNodeConfig(flow.RoleExecution, testnet.WithLogLevel("error"))
		verNode  = testnet.NewNodeConfig(flow.RoleVerification, testnet.WithLogLevel("error"))
		accNode  = testnet.NewNodeConfig(flow.RoleAccess, testnet.WithLogLevel("debug"))
	)

	nodes := []testnet.NodeConfig{colNode1, colNode2, conNode1, conNode2, conNode3, exeNode, verNode, accNode}
	conf := testnet.NetworkConfig{Nodes: nodes}

	testingdock.Verbose = true

	ctx := context.Background()

	net, err := testnet.PrepareFlowNetwork(t, "access", conf)
	require.Nil(t, err)

	net.Start(ctx)
	defer net.Cleanup()

	accessContainer, ok := net.ContainerByID(accNode.Identifier)
	assert.True(t, ok)

	port, ok := accessContainer.Ports[testnet.AccessNodeAPIPort]
	assert.True(t, ok)

	client, err := testnet.NewClient(fmt.Sprintf(":%s", port))
	assert.Nil(t, err)

	rpcClient := *client.GetClient()

	t.Run("get genesis block", func(t *testing.T) {
		expectedHeight := uint64(0)
		req := access.GetBlockByHeightRequest{Height: expectedHeight}
		resp, err := rpcClient.GetBlockByHeight(ctx, &req)
		assert.Nil(t, err)
		t.Log("block: ", resp.GetBlock())
		block := resp.GetBlock()
		assert.NotNil(t, block)
		assert.Equal(t, expectedHeight, block.GetHeight())
		assert.NotNil(t, block.GetId())
	})

	t.Run("get genesis block header", func(t *testing.T) {
		expectedHeight := uint64(0)
		req := access.GetBlockHeaderByHeightRequest{Height: expectedHeight}
		resp, err := rpcClient.GetBlockHeaderByHeight(ctx, &req)
		assert.Nil(t, err)
		t.Log("block: ", resp.GetBlock())
		block := resp.GetBlock()
		assert.NotNil(t, block)
		assert.Equal(t, expectedHeight, block.GetHeight())
		assert.NotNil(t, block.GetId())
	})

	tx := unittest.TransactionBodyFixture()
	tx, err = client.SignTransaction(tx)
	assert.Nil(t, err)

	t.Run("send transaction", func(t *testing.T) {

		t.Log("sending transaction: ", tx.ID())

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		err = client.SendTransaction(ctx, tx)
		assert.Nil(t, err)
	})

	//t.Run("get transaction", func(t *testing.T) {
	//
	//	t.Log("get transaction: ", tx.ID())
	//
	//	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	//	defer cancel()
	//
	//	txID := tx.ID()
	//	req := access.GetTransactionRequest{
	//		Id: txID[:],
	//	}
	//
	//	assert.Eventually(t, func() bool {
	//		resp, err := rpcClient.GetTransaction(ctx, &req)
	//		if err != nil {
	//			return false
	//		}
	//		t := resp.GetTransaction()
	//		if t == nil {
	//			return false
	//		}
	//
	//		actualTx, err := convert.MessageToTransaction(t)
	//
	//		return (tx == *actualTx)
	//	}, time.Minute, time.Second)
	//
	//	assert.Nil(t, err)
	//})

	//
	//t.Run("send transaction", func(t *testing.T) {
	//	rpcClient := *client.GetClient()
	//
	//	req := access.GetBlockByHeightRequest{Height: 0}
	//	resp, err := rpcClient.GetBlockByHeight(ctx, &req)
	//	assert.Nil(t, err)
	//	t.Log("Block: ", resp.GetBlock())
	//
	//	time.Sleep(2 * time.Minute)
	//	txID := tx.ID()
	//	trreq := access.GetTransactionRequest{
	//		Id: txID[:],
	//	}
	//
	//	trresp, err := rpcClient.GetTransaction(ctx, &trreq)
	//	assert.Nil(t, err)
	//	fmt.Println(trresp.GetTransaction())
	//	assert.NotNil(t, trresp.GetTransaction())
	//	assert.Equal(t, entities.TransactionStatus_STATUS_SEALED, trresp.GetTransaction().GetStatus())
	//
	//	// wait for consensus to complete
	//	//TODO we should listen for collection guarantees instead, but this is blocked
	//	// ref: https://github.com/dapperlabs/flow-go/issues/3021
	//	//time.Sleep(10 * time.Second)
	//
	//	// TODO stop then start containers
	//	err = net.StopContainers()
	//	assert.Nil(t, err)
	//	//
	//	//identities := net.Identities()
	//	//
	//	//chainID := protocol.ChainIDForCluster(identities.Filter(filter.HasRole(flow.RoleCollection)))
	//	//
	//	//// get database for COL1
	//	//db, err := colContainer1.DB()
	//	//require.Nil(t, err)
	//	//
	//	//state, err := clusterstate.NewState(db, chainID)
	//	//assert.Nil(t, err)
	//	//
	//	//// the transaction should be included in exactly one collection
	//	//head, err := state.Final().Head()
	//	//assert.Nil(t, err)
	//	//
	//	//foundTx := false
	//	//for head.Height > 0 {
	//	//	collection, err := state.AtBlockID(head.ID()).Collection()
	//	//	assert.Nil(t, err)
	//	//
	//	//	head, err = state.AtBlockID(head.ParentID).Head()
	//	//	assert.Nil(t, err)
	//	//
	//	//	if collection.Len() == 0 {
	//	//		continue
	//	//	}
	//	//
	//	//	for _, txID := range collection.Transactions {
	//	//		assert.Equal(t, tx.ID(), txID, "found unexpected transaction")
	//	//		if txID == tx.ID() {
	//	//			assert.False(t, foundTx, "found duplicate transaction")
	//	//			foundTx = true
	//	//		}
	//	//	}
	//	//}
	//	//
	//	//assert.True(t, foundTx)
	//})
	//fmt.Println(client)

	//time.Sleep(10 * time.Second)
	err = net.StopContainers()

}
