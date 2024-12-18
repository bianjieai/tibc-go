package tibctesting

import (
	"fmt"

	"github.com/stretchr/testify/require"

	storetypes "cosmossdk.io/store/types"
	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	packettypes "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	commitmenttypes "github.com/bianjieai/tibc-go/modules/tibc/core/23-commitment/types"
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
	ibctmtypes "github.com/bianjieai/tibc-go/modules/tibc/light-clients/07-tendermint/types"
)

// Endpoint is a which represents a channel endpoint and its associated
// client and connections. It contains client, connection, and channel
// configuration parameters. Endpoint functions will utilize the parameters
// set in the configuration structs when executing TIBC messages.
type Endpoint struct {
	Chain        *TestChain
	Counterparty *Endpoint
	ClientConfig ClientConfig
}

// NewEndpoint constructs a new endpoint without the counterparty.
// CONTRACT: the counterparty endpoint must be set by the caller.
func NewEndpoint(chain *TestChain, clientConfig ClientConfig) *Endpoint {
	return &Endpoint{
		Chain:        chain,
		ClientConfig: clientConfig,
	}
}

// NewDefaultEndpoint constructs a new endpoint using default values.
// CONTRACT: the counterparty endpoitn must be set by the caller.
func NewDefaultEndpoint(chain *TestChain) *Endpoint {
	return &Endpoint{
		Chain:        chain,
		ClientConfig: NewTendermintConfig(),
	}
}

// QueryProof queries proof associated with this endpoint using the lastest client state
// height on the counterparty chain.
func (endpoint *Endpoint) QueryProof(key []byte) ([]byte, clienttypes.Height) {
	// obtain the counterparty client representing the chain associated with the endpoint
	clientState := endpoint.Counterparty.Chain.GetClientState(endpoint.Chain.ChainName)
	// query proof on the counterparty using the latest height of the TIBC client
	return endpoint.QueryProofAtHeight(key, clientState.GetLatestHeight().GetRevisionHeight())
}

// QueryProofAtHeight queries proof associated with this endpoint using the proof height
// providied
func (endpoint *Endpoint) QueryProofAtHeight(key []byte, height uint64) ([]byte, clienttypes.Height) {
	// query proof on the counterparty using the latest height of the TIBC client
	return endpoint.Chain.QueryProofAtHeight(key, int64(height))
}

// CreateClient creates an TIBC client on the endpoint. It will update the
// chainName for the endpoint if the message is successfully executed.
// NOTE: a solo machine client will be created with an empty diversifier.
func (endpoint *Endpoint) CreateClient() error {
	// ensure counterparty has committed state
	endpoint.Counterparty.Chain.NextBlock()

	// ensure the chain has the latest time
	endpoint.Chain.Coordinator.UpdateTimeForChain(endpoint.Chain)

	if endpoint.ClientConfig.GetClientType() != exported.Tendermint {
		return fmt.Errorf("client type %s is not supported", endpoint.ClientConfig.GetClientType())
	}

	tmConfig, ok := endpoint.ClientConfig.(*TendermintConfig)
	require.True(endpoint.Chain.T, ok)

	height := endpoint.Counterparty.Chain.LastHeader.GetHeight().(clienttypes.Height)
	clientState := ibctmtypes.NewClientState(
		endpoint.Counterparty.Chain.ChainID, tmConfig.TrustLevel,
		tmConfig.TrustingPeriod, tmConfig.UnbondingPeriod, tmConfig.MaxClockDrift,
		height, commitmenttypes.GetSDKSpecs(), Prefix, 0,
	)
	consensusState := endpoint.Counterparty.Chain.LastHeader.ConsensusState()

	ctx := endpoint.Chain.GetContext()

	// set self chain name
	endpoint.Chain.App.TIBCKeeper.ClientKeeper.SetChainName(ctx, endpoint.Chain.ChainName)

	// set send sequence
	endpoint.Chain.App.TIBCKeeper.PacketKeeper.SetNextSequenceSend(ctx, endpoint.Chain.ChainName, endpoint.Counterparty.Chain.ChainName, 1)

	// set relayers
	relayers := []string{endpoint.Chain.SenderAccount.GetAddress().String()}
	endpoint.Chain.App.TIBCKeeper.ClientKeeper.RegisterRelayers(ctx, endpoint.Counterparty.Chain.ChainName, relayers)

	// create counterparty chain light client
	err := endpoint.Chain.App.TIBCKeeper.ClientKeeper.CreateClient(
		ctx,
		endpoint.Counterparty.Chain.ChainName,
		clientState,
		consensusState,
	)
	require.NoError(endpoint.Chain.T, err)
	// NextBlock calls app.Commit()
	endpoint.Chain.NextBlock()
	// increment sequence for successful transaction execution
	endpoint.Chain.Coordinator.IncrementTime()
	return nil
}

// UpdateClient updates the TIBC client associated with the endpoint.
func (endpoint *Endpoint) UpdateClient() error {
	// ensure counterparty has committed state
	endpoint.Chain.Coordinator.CommitBlock(endpoint.Counterparty.Chain)

	if endpoint.ClientConfig.GetClientType() != exported.Tendermint {
		return fmt.Errorf("client type %s is not supported", endpoint.ClientConfig.GetClientType())
	}

	header, err := endpoint.Chain.ConstructUpdateTMClientHeader(endpoint.Counterparty.Chain, endpoint.Counterparty.Chain.ChainName)
	if err != nil {
		return err
	}

	msg, err := clienttypes.NewMsgUpdateClient(
		endpoint.Counterparty.Chain.ChainName, header,
		endpoint.Chain.SenderAccount.GetAddress(),
	)
	require.NoError(endpoint.Chain.T, err)

	return endpoint.Chain.sendMsgs(msg)
}

// SendPacket sends a packet through the channel keeper using the associated endpoint
// The counterparty client is updated so proofs can be sent to the counterparty chain.
func (endpoint *Endpoint) SendPacket(packet exported.PacketI) error {
	// no need to send message, acting as a module
	if err := endpoint.Chain.App.TIBCKeeper.PacketKeeper.SendPacket(endpoint.Chain.GetContext(), packet); err != nil {
		return err
	}

	// commit changes since no message was sent
	endpoint.Chain.Coordinator.CommitBlock(endpoint.Chain)

	return endpoint.Counterparty.UpdateClient()
}

// RecvPacket receives a packet on the associated endpoint.
// The counterparty client is updated.
func (endpoint *Endpoint) RecvPacket(packet packettypes.Packet) error {
	// get proof of packet commitment on source
	packetKey := host.PacketCommitmentKey(packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())
	proof, proofHeight := endpoint.Counterparty.Chain.QueryProof(packetKey)

	recvMsg := packettypes.NewMsgRecvPacket(packet, proof, proofHeight, endpoint.Chain.SenderAccount.GetAddress())

	// receive on counterparty and update source client
	if err := endpoint.Chain.sendMsgs(recvMsg); err != nil {
		return err
	}

	return endpoint.Counterparty.UpdateClient()
}

// WriteAcknowledgement writes an acknowledgement on the channel associated with the endpoint.
// The counterparty client is updated.
func (endpoint *Endpoint) WriteAcknowledgement(acknowledgement []byte, packet exported.PacketI) error {
	// no need to send message, acting as a handler
	err := endpoint.Chain.App.TIBCKeeper.PacketKeeper.WriteAcknowledgement(endpoint.Chain.GetContext(), packet, acknowledgement)
	if err != nil {
		return err
	}

	// commit changes since no message was sent
	endpoint.Chain.Coordinator.CommitBlock(endpoint.Chain)

	return endpoint.Counterparty.UpdateClient()
}

// AcknowledgePacket sends a MsgAcknowledgement to the channel associated with the endpoint.
func (endpoint *Endpoint) AcknowledgePacket(packet packettypes.Packet, ack []byte) error {
	// get proof of acknowledgement on counterparty
	packetKey := host.PacketAcknowledgementKey(packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())
	proof, proofHeight := endpoint.Counterparty.QueryProof(packetKey)

	ackMsg := packettypes.NewMsgAcknowledgement(packet, ack, proof, proofHeight, endpoint.Chain.SenderAccount.GetAddress())

	return endpoint.Chain.sendMsgs(ackMsg)
}

// AcknowledgePacket sends a MsgAcknowledgement to the channel associated with the endpoint.
func (endpoint *Endpoint) CleanPacket(cleanPacket packettypes.CleanPacket) error {
	// commit changes since no message was sent
	endpoint.Chain.Coordinator.CommitBlock(endpoint.Chain)
	cleanMsg := packettypes.NewMsgCleanPacket(cleanPacket, endpoint.Chain.SenderAccount.GetAddress())
	endpoint.Chain.sendMsgs(cleanMsg)

	// commit changes since no message was sent
	endpoint.Chain.Coordinator.CommitBlock(endpoint.Chain)

	return endpoint.Counterparty.UpdateClient()
}

// AcknowledgePacket sends a MsgAcknowledgement to the channel associated with the endpoint.
func (endpoint *Endpoint) RecvCleanPacket(cleanPacket packettypes.CleanPacket) error {
	// get proof of acknowledgement on counterparty
	packetKey := host.CleanPacketCommitmentKey(cleanPacket.GetSourceChain(), cleanPacket.GetDestChain())
	proof, proofHeight := endpoint.Counterparty.QueryProof(packetKey)

	recvCleanMsg := packettypes.NewMsgRecvCleanPacket(cleanPacket, proof, proofHeight, endpoint.Chain.SenderAccount.GetAddress())

	return endpoint.Chain.sendMsgs(recvCleanMsg)
}

func (endpoint *Endpoint) ClientStore() storetypes.KVStore {
	return endpoint.Chain.App.TIBCKeeper.ClientKeeper.ClientStore(endpoint.Chain.GetContext(), endpoint.Counterparty.Chain.ChainName)
}

// GetClientState retrieves the Client State for this endpoint. The
// client state is expected to exist otherwise testing will fail.
func (endpoint *Endpoint) GetClientState() exported.ClientState {
	return endpoint.Chain.GetClientState(endpoint.Counterparty.Chain.ChainName)
}

// SetClientState sets the client state for this endpoint.
func (endpoint *Endpoint) SetClientState(clientState exported.ClientState) {
	endpoint.Chain.App.TIBCKeeper.ClientKeeper.SetClientState(endpoint.Chain.GetContext(), endpoint.Counterparty.Chain.ChainName, clientState)
}

// GetConsensusState retrieves the Consensus State for this endpoint at the provided height.
// The consensus state is expected to exist otherwise testing will fail.
func (endpoint *Endpoint) GetConsensusState(height exported.Height) exported.ConsensusState {
	consensusState, found := endpoint.Chain.GetConsensusState(endpoint.Counterparty.Chain.ChainName, height)
	require.True(endpoint.Chain.T, found)

	return consensusState
}

// SetConsensusState sets the consensus state for this endpoint.
func (endpoint *Endpoint) SetConsensusState(consensusState exported.ConsensusState, height exported.Height) {
	endpoint.Chain.App.TIBCKeeper.ClientKeeper.SetClientConsensusState(endpoint.Chain.GetContext(), endpoint.Counterparty.Chain.ChainName, height, consensusState)
}

// QueryClientStateProof performs and abci query for a client stat associated
// with this endpoint and returns the ClientState along with the proof.
func (endpoint *Endpoint) QueryClientStateProof() (exported.ClientState, []byte) {
	// retrieve client state to provide proof for
	clientState := endpoint.GetClientState()

	clientKey := host.FullClientStateKey(endpoint.Counterparty.Chain.ChainName)
	proofClient, _ := endpoint.QueryProof(clientKey)

	return clientState, proofClient
}
