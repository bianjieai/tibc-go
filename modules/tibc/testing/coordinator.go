package ibctesting

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	packettypes "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

var (
	ChainIDPrefix   = "testchain"
	globalStartTime = time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	TimeIncrement   = time.Second * 5
)

// Coordinator is a testing struct which contains N TestChain's. It handles keeping all chains
// in sync with regards to time.
type Coordinator struct {
	t *testing.T

	Chains map[string]*TestChain
}

// NewCoordinator initializes Coordinator with N TestChain's
func NewCoordinator(t *testing.T, n int) *Coordinator {
	chains := make(map[string]*TestChain)

	for i := 0; i < n; i++ {
		chainID := GetChainID(i)
		chains[chainID] = NewTestChain(t, chainID)
	}
	return &Coordinator{
		t:      t,
		Chains: chains,
	}
}

// Setup constructs a TM client, connection, and channel on both chains provided. It will
// fail if any error occurs. The clientID's, TestConnections, and TestChannels are returned
// for both chains. The channels created are connected to the ibc-transfer application.
func (coord *Coordinator) Setup(
	chainA, chainB *TestChain,
) (string, string) {
	clientA, clientB := coord.SetupClientConnections(chainA, chainB, exported.Tendermint)

	return clientA, clientB
}

// SetupClients is a helper function to create clients on both chains. It assumes the
// caller does not anticipate any errors.
func (coord *Coordinator) SetupClients(
	chainA, chainB *TestChain,
	clientType string,
) (string, string) {

	clientA, err := coord.CreateClient(chainA, chainB, clientType)
	require.NoError(coord.t, err)
	chainA.App.IBCKeeper.ClientKeeper.RegisterRelayers(chainA.GetContext(),clientA,[]string{chainA.SenderAccount.GetAddress().String()})

	clientB, err := coord.CreateClient(chainB, chainA, clientType)
	require.NoError(coord.t, err)
	chainB.App.IBCKeeper.ClientKeeper.RegisterRelayers(chainB.GetContext(),clientA,[]string{chainB.SenderAccount.GetAddress().String()})
	return clientA, clientB
}

// SetupClientConnections is a helper function to create clients and the appropriate
// connections on both the source and counterparty chain. It assumes the caller does not
// anticipate any errors.
func (coord *Coordinator) SetupClientConnections(
	chainA, chainB *TestChain,
	clientType string,
) (string, string) {

	clientA, clientB := coord.SetupClients(chainA, chainB, clientType)

	return clientA, clientB
}

// CreateClient creates a counterparty client on the source chain and returns the clientID.
func (coord *Coordinator) CreateClient(
	source, counterparty *TestChain,
	clientType string,
) (clientID string, err error) {
	coord.CommitBlock(source, counterparty)

	clientID = source.NewClientID(clientType)

	switch clientType {
	case exported.Tendermint:
		err = source.CreateTMClient(counterparty, clientID)

	default:
		err = fmt.Errorf("client type %s is not supported", clientType)
	}

	if err != nil {
		return "", err
	}

	coord.IncrementTime()

	return clientID, nil
}

// UpdateClient updates a counterparty client on the source chain.
func (coord *Coordinator) UpdateClient(
	source, counterparty *TestChain,
	clientID string,
	clientType string,
) (err error) {
	coord.CommitBlock(source, counterparty)

	switch clientType {
	case exported.Tendermint:
		err = source.UpdateTMClient(counterparty, clientID)

	default:
		err = fmt.Errorf("client type %s is not supported", clientType)
	}

	if err != nil {
		return err
	}

	coord.IncrementTime()

	return nil
}

// SendPacket sends a packet through the channel keeper on the source chain and updates the
// counterparty client for the source chain.
func (coord *Coordinator) SendPacket(
	source, counterparty *TestChain,
	packet exported.PacketI,
	counterpartyClientID string,
) error {
	if err := source.SendPacket(packet); err != nil {
		return err
	}
	coord.IncrementTime()

	// update source client on counterparty connection
	return coord.UpdateClient(
		counterparty, source,
		counterpartyClientID, exported.Tendermint,
	)
}

// RecvPacket receives a channel packet on the counterparty chain and updates
// the client on the source chain representing the counterparty.
func (coord *Coordinator) RecvPacket(
	source, counterparty *TestChain,
	sourceClient string,
	packet packettypes.Packet,
) error {
	// get proof of packet commitment on source
	packetKey := host.PacketCommitmentKey(packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())
	proof, proofHeight := source.QueryProof(packetKey)

	// Increment time and commit block so that 5 second delay period passes between send and receive
	coord.IncrementTime()
	coord.CommitBlock(source, counterparty)

	recvMsg := packettypes.NewMsgRecvPacket(packet, proof, proofHeight, counterparty.SenderAccount.GetAddress())

	// receive on counterparty and update source client
	return coord.SendMsgs(counterparty, source, sourceClient, []sdk.Msg{recvMsg})
}

// WriteAcknowledgement writes an acknowledgement to the channel keeper on the source chain and updates the
// counterparty client for the source chain.
func (coord *Coordinator) WriteAcknowledgement(
	source, counterparty *TestChain,
	packet exported.PacketI,
	counterpartyClientID string,
) error {
	if err := source.WriteAcknowledgement(packet); err != nil {
		return err
	}
	coord.IncrementTime()

	// update source client on counterparty connection
	return coord.UpdateClient(
		counterparty, source,
		counterpartyClientID, exported.Tendermint,
	)
}

// AcknowledgePacket acknowledges on the source chain the packet received on
// the counterparty chain and updates the client on the counterparty representing
// the source chain.
// TODO: add a query for the acknowledgement by events
// - https://github.com/cosmos/cosmos-sdk/issues/6509
func (coord *Coordinator) AcknowledgePacket(
	source, counterparty *TestChain,
	counterpartyClient string,
	packet packettypes.Packet, ack []byte,
) error {
	// get proof of acknowledgement on counterparty
	packetKey := host.PacketAcknowledgementKey(packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())
	proof, proofHeight := counterparty.QueryProof(packetKey)

	// Increment time and commit block so that 5 second delay period passes between send and receive
	coord.IncrementTime()
	coord.CommitBlock(source, counterparty)

	ackMsg := packettypes.NewMsgAcknowledgement(packet, ack, proof, proofHeight, source.SenderAccount.GetAddress())
	return coord.SendMsgs(source, counterparty, counterpartyClient, []sdk.Msg{ackMsg})
}

// RelayPacket receives a channel packet on counterparty, queries the ack
// and acknowledges the packet on source. The clients are updated as needed.
func (coord *Coordinator) RelayPacket(
	source, counterparty *TestChain,
	sourceClient, counterpartyClient string,
	packet packettypes.Packet, ack []byte,
) error {
	// Increment time and commit block so that 5 second delay period passes between send and receive
	coord.IncrementTime()
	coord.CommitBlock(counterparty)

	if err := coord.RecvPacket(source, counterparty, sourceClient, packet); err != nil {
		return err
	}

	// Increment time and commit block so that 5 second delay period passes between send and receive
	coord.IncrementTime()
	coord.CommitBlock(source)

	return coord.AcknowledgePacket(source, counterparty, counterpartyClient, packet, ack)
}

// IncrementTime iterates through all the TestChain's and increments their current header time
// by 5 seconds.
//
// CONTRACT: this function must be called after every commit on any TestChain.
func (coord *Coordinator) IncrementTime() {
	for _, chain := range coord.Chains {
		chain.CurrentHeader.Time = chain.CurrentHeader.Time.Add(TimeIncrement)
		chain.App.BeginBlock(abci.RequestBeginBlock{Header: chain.CurrentHeader})
	}
}

// IncrementTimeBy iterates through all the TestChain's and increments their current header time
// by specified time.
func (coord *Coordinator) IncrementTimeBy(increment time.Duration) {
	for _, chain := range coord.Chains {
		chain.CurrentHeader.Time = chain.CurrentHeader.Time.Add(increment)
		chain.App.BeginBlock(abci.RequestBeginBlock{Header: chain.CurrentHeader})
	}
}

// SendMsg delivers a single provided message to the chain. The counterparty
// client is update with the new source consensus state.
func (coord *Coordinator) SendMsg(source, counterparty *TestChain, counterpartyClientID string, msg sdk.Msg) error {
	return coord.SendMsgs(source, counterparty, counterpartyClientID, []sdk.Msg{msg})
}

// SendMsgs delivers the provided messages to the chain. The counterparty
// client is updated with the new source consensus state.
func (coord *Coordinator) SendMsgs(source, counterparty *TestChain, counterpartyClientID string, msgs []sdk.Msg) error {
	if err := source.sendMsgs(msgs...); err != nil {
		return err
	}

	coord.IncrementTime()

	// update source client on counterparty connection
	return coord.UpdateClient(
		counterparty, source,
		counterpartyClientID, exported.Tendermint,
	)
}

// GetChain returns the TestChain using the given chainID and returns an error if it does
// not exist.
func (coord *Coordinator) GetChain(chainID string) *TestChain {
	chain, found := coord.Chains[chainID]
	require.True(coord.t, found, fmt.Sprintf("%s chain does not exist", chainID))
	return chain
}

// GetChainID returns the chainID used for the provided index.
func GetChainID(index int) string {
	return ChainIDPrefix + strconv.Itoa(index)
}

// CommitBlock commits a block on the provided indexes and then increments the global time.
//
// CONTRACT: the passed in list of indexes must not contain duplicates
func (coord *Coordinator) CommitBlock(chains ...*TestChain) {
	for _, chain := range chains {
		chain.App.Commit()
		chain.NextBlock()
	}
	coord.IncrementTime()
}

// CommitNBlocks commits n blocks to state and updates the block height by 1 for each commit.
func (coord *Coordinator) CommitNBlocks(chain *TestChain, n uint64) {
	for i := uint64(0); i < n; i++ {
		chain.App.BeginBlock(abci.RequestBeginBlock{Header: chain.CurrentHeader})
		chain.App.Commit()
		chain.NextBlock()
		coord.IncrementTime()
	}
}
