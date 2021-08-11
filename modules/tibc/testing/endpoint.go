package ibctesting

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	//	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

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
// set in the configuration structs when executing IBC messages.
type Endpoint struct {
	Chain        *TestChain
	Counterparty *Endpoint
	ChainName    string
	ClientConfig ClientConfig
}

// NewEndpoint constructs a new endpoint without the counterparty.
// CONTRACT: the counterparty endpoint must be set by the caller.
func NewEndpoint(
	chain *TestChain, clientConfig ClientConfig,
) *Endpoint {
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
		ChainName:    chain.ChainID,
		ClientConfig: NewTendermintConfig(),
	}
}

// QueryProof queries proof associated with this endpoint using the lastest client state
// height on the counterparty chain.
func (endpoint *Endpoint) QueryProof(key []byte) ([]byte, clienttypes.Height) {
	// obtain the counterparty client representing the chain associated with the endpoint
	clientState := endpoint.Counterparty.Chain.GetClientState(endpoint.Counterparty.ChainName)

	// query proof on the counterparty using the latest height of the IBC client
	return endpoint.QueryProofAtHeight(key, clientState.GetLatestHeight().GetRevisionHeight())
}

// QueryProofAtHeight queries proof associated with this endpoint using the proof height
// providied
func (endpoint *Endpoint) QueryProofAtHeight(key []byte, height uint64) ([]byte, clienttypes.Height) {
	// query proof on the counterparty using the latest height of the IBC client
	return endpoint.Chain.QueryProofAtHeight(key, int64(height))
}

// CreateClient creates an IBC client on the endpoint. It will update the
// clientID for the endpoint if the message is successfully executed.
// NOTE: a solo machine client will be created with an empty diversifier.
func (endpoint *Endpoint) CreateClient() (err error) {
	// ensure counterparty has committed state
	endpoint.Chain.Coordinator.CommitBlock(endpoint.Counterparty.Chain)

	// ensure the chain has the latest time
	endpoint.Chain.Coordinator.UpdateTimeForChain(endpoint.Chain)

	var (
		clientState    exported.ClientState
		consensusState exported.ConsensusState
	)

	switch endpoint.ClientConfig.GetClientType() {
	case exported.Tendermint:
		tmConfig, ok := endpoint.ClientConfig.(*TendermintConfig)
		require.True(endpoint.Chain.t, ok)

		height := endpoint.Counterparty.Chain.LastHeader.GetHeight().(clienttypes.Height)
		clientState = ibctmtypes.NewClientState(
			endpoint.Counterparty.Chain.ChainID, tmConfig.TrustLevel, tmConfig.TrustingPeriod, tmConfig.UnbondingPeriod, tmConfig.MaxClockDrift,
			height, commitmenttypes.GetSDKSpecs(), Prefix, 0,
		)
		consensusState = endpoint.Counterparty.Chain.LastHeader.ConsensusState()

	default:
		err = fmt.Errorf("client type %s is not supported", endpoint.ClientConfig.GetClientType())
	}

	if err != nil {
		return err
	}
	ctx := endpoint.Chain.GetContext()

	endpoint.Chain.App.IBCKeeper.Packetkeeper.SetNextSequenceSend(ctx, endpoint.ChainName, endpoint.Counterparty.ChainName, 1)

	relayers := []string{endpoint.Chain.SenderAccount.GetAddress().String()}
	endpoint.Chain.App.IBCKeeper.ClientKeeper.SetChainName(ctx, endpoint.ChainName)
	endpoint.Chain.App.IBCKeeper.ClientKeeper.RegisterRelayers(endpoint.Chain.GetContext(), endpoint.Counterparty.ChainName, relayers)
	err = endpoint.Chain.App.IBCKeeper.ClientKeeper.CreateClient(
		endpoint.Chain.GetContext(),
		endpoint.Counterparty.ChainName,
		clientState, consensusState,
	)
	require.NoError(endpoint.Chain.t, err)

	endpoint.Chain.App.EndBlock(abci.RequestEndBlock{})
	endpoint.Chain.App.Commit()

	endpoint.Chain.NextBlock()
	endpoint.Chain.Coordinator.IncrementTime()
	return nil
}

// UpdateClient updates the IBC client associated with the endpoint.
func (endpoint *Endpoint) UpdateClient() (err error) {
	// ensure counterparty has committed state
	endpoint.Chain.Coordinator.CommitBlock(endpoint.Counterparty.Chain)

	var (
		header exported.Header
	)

	switch endpoint.ClientConfig.GetClientType() {
	case exported.Tendermint:
		header, err = endpoint.Chain.ConstructUpdateTMClientHeader(endpoint.Counterparty.Chain, endpoint.Counterparty.ChainName)

	default:
		err = fmt.Errorf("client type %s is not supported", endpoint.ClientConfig.GetClientType())
	}

	if err != nil {
		return err
	}

	msg, err := clienttypes.NewMsgUpdateClient(
		endpoint.Counterparty.ChainName, header,
		endpoint.Chain.SenderAccount.GetAddress(),
	)
	require.NoError(endpoint.Chain.t, err)

	return endpoint.Chain.sendMsgs(msg)

}

// SendPacket sends a packet through the channel keeper using the associated endpoint
// The counterparty client is updated so proofs can be sent to the counterparty chain.
func (endpoint *Endpoint) SendPacket(packet exported.PacketI) error {
	// no need to send message, acting as a module
	err := endpoint.Chain.App.IBCKeeper.Packetkeeper.SendPacket(endpoint.Chain.GetContext(), packet)
	if err != nil {
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
	err := endpoint.Chain.App.IBCKeeper.Packetkeeper.WriteAcknowledgement(endpoint.Chain.GetContext(), packet, acknowledgement)
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
func (endpoint *Endpoint) CleanPacket(packet packettypes.Packet) error {
	// commit changes since no message was sent
	endpoint.Chain.Coordinator.CommitBlock(endpoint.Chain)
	cleanMsg := packettypes.NewMsgCleanPacket(packet, endpoint.Chain.SenderAccount.GetAddress())

	return endpoint.Chain.sendMsgs(cleanMsg)
}

// AcknowledgePacket sends a MsgAcknowledgement to the channel associated with the endpoint.
func (endpoint *Endpoint) RecvCleanPacket(packet packettypes.Packet) error {
	// get proof of acknowledgement on counterparty
	packetKey := host.PacketCleanKey(packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())
	proof, proofHeight := endpoint.Counterparty.QueryProof(packetKey)

	recvCleanMsg := packettypes.NewMsgRecvCleanPacket(packet, proof, proofHeight, endpoint.Chain.SenderAccount.GetAddress())

	return endpoint.Chain.sendMsgs(recvCleanMsg)
}

func (endpoint *Endpoint) ClientStore() sdk.KVStore {
	return endpoint.Chain.App.IBCKeeper.ClientKeeper.ClientStore(endpoint.Chain.GetContext(), endpoint.Counterparty.ChainName)
}

// GetClientState retrieves the Client State for this endpoint. The
// client state is expected to exist otherwise testing will fail.
func (endpoint *Endpoint) GetClientState() exported.ClientState {
	return endpoint.Chain.GetClientState(endpoint.Counterparty.ChainName)
}

// SetClientState sets the client state for this endpoint.
func (endpoint *Endpoint) SetClientState(clientState exported.ClientState) {
	endpoint.Chain.App.IBCKeeper.ClientKeeper.SetClientState(endpoint.Chain.GetContext(), endpoint.Counterparty.ChainName, clientState)
}

// GetConsensusState retrieves the Consensus State for this endpoint at the provided height.
// The consensus state is expected to exist otherwise testing will fail.
func (endpoint *Endpoint) GetConsensusState(height exported.Height) exported.ConsensusState {
	consensusState, found := endpoint.Chain.GetConsensusState(endpoint.Counterparty.ChainName, height)
	require.True(endpoint.Chain.t, found)

	return consensusState
}

// SetConsensusState sets the consensus state for this endpoint.
func (endpoint *Endpoint) SetConsensusState(consensusState exported.ConsensusState, height exported.Height) {
	endpoint.Chain.App.IBCKeeper.ClientKeeper.SetClientConsensusState(endpoint.Chain.GetContext(), endpoint.Counterparty.ChainName, height, consensusState)
}

// QueryClientStateProof performs and abci query for a client stat associated
// with this endpoint and returns the ClientState along with the proof.
func (endpoint *Endpoint) QueryClientStateProof() (exported.ClientState, []byte) {
	// retrieve client state to provide proof for
	clientState := endpoint.GetClientState()

	clientKey := host.FullClientStateKey(endpoint.Counterparty.ChainName)
	proofClient, _ := endpoint.QueryProof(clientKey)

	return clientState, proofClient
}
