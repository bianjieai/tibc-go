package tibctesting

import (
	"bytes"
	"fmt"

	packettypes "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
)

// Path contains two endpoints representing two chains connected over TIBC
type Path struct {
	EndpointA *Endpoint
	EndpointB *Endpoint
}

// NewPath constructs an endpoint for each chain using the default values
// for the endpoints. Each endpoint is updated to have a pointer to the
// counterparty endpoint.
func NewPath(chainA, chainB *TestChain) *Path {
	endpointA := NewDefaultEndpoint(chainA)
	endpointB := NewDefaultEndpoint(chainB)

	endpointA.Counterparty = endpointB
	endpointB.Counterparty = endpointA

	return &Path{
		EndpointA: endpointA,
		EndpointB: endpointB,
	}
}

// RelayPacket attempts to relay the packet first on EndpointA and then on EndpointB
// if EndpointA does not contain a packet commitment for that packet. An error is returned
// if a relay step fails or the packet commitment does not exist on either endpoint.
func (path *Path) RelayPacket(packet packettypes.Packet, ack []byte) error {
	pc := path.EndpointA.Chain.App.TIBCKeeper.PacketKeeper.GetPacketCommitment(
		path.EndpointA.Chain.GetContext(), packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence(),
	)
	if bytes.Equal(pc, packettypes.CommitPacket(packet)) {

		// packet found, relay from A to B
		_ = path.EndpointB.UpdateClient()

		if err := path.EndpointB.RecvPacket(packet); err != nil {
			return err
		}

		if err := path.EndpointA.AcknowledgePacket(packet, ack); err != nil {
			return err
		}
		return nil

	}

	pc = path.EndpointB.Chain.App.TIBCKeeper.PacketKeeper.GetPacketCommitment(
		path.EndpointB.Chain.GetContext(), packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence(),
	)
	if bytes.Equal(pc, packettypes.CommitPacket(packet)) {
		// packet found, relay B to A
		_ = path.EndpointA.UpdateClient()

		if err := path.EndpointA.RecvPacket(packet); err != nil {
			return err
		}
		if err := path.EndpointB.AcknowledgePacket(packet, ack); err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("packet commitment does not exist on either endpoint for provided packet")
}
