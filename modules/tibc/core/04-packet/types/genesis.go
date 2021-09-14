package types

import (
	"errors"
	"fmt"

	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
)

// NewPacketState creates a new PacketState instance.
func NewPacketState(sourceChain, destinationChain string, seq uint64, data []byte) PacketState {
	return PacketState{
		SourceChain:      sourceChain,
		DestinationChain: destinationChain,
		Sequence:         seq,
		Data:             data,
	}
}

// Validate performs basic validation of fields returning an error upon any
// failure.
func (pa PacketState) Validate() error {
	if pa.Data == nil {
		return errors.New("data bytes cannot be nil")
	}
	return validateGenFields(pa.SourceChain, pa.DestinationChain, pa.Sequence)
}

// NewPacketSequence creates a new PacketSequences instance.
func NewPacketSequence(sourceChain, destinationChain string, seq uint64) PacketSequence {
	return PacketSequence{
		SourceChain:      sourceChain,
		DestinationChain: destinationChain,
		Sequence:         seq,
	}
}

// Validate performs basic validation of fields returning an error upon any
// failure.
func (ps PacketSequence) Validate() error {
	return validateGenFields(ps.SourceChain, ps.DestinationChain, ps.Sequence)
}

// NewGenesisState creates a GenesisState instance.
func NewGenesisState(
	acks, commitments, receipts []PacketState,
	sendSeqs, recvSeqs, ackSeqs []PacketSequence,
) GenesisState {
	return GenesisState{
		Acknowledgements: acks,
		Commitments:      commitments,
		Receipts:         receipts,
		SendSequences:    sendSeqs,
		RecvSequences:    recvSeqs,
		AckSequences:     ackSeqs,
	}
}

// DefaultGenesisState returns the tibc packet submodule's default genesis state.
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Acknowledgements: []PacketState{},
		Receipts:         []PacketState{},
		Commitments:      []PacketState{},
		SendSequences:    []PacketSequence{},
		RecvSequences:    []PacketSequence{},
		AckSequences:     []PacketSequence{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	for i, ack := range gs.Acknowledgements {
		if err := ack.Validate(); err != nil {
			return fmt.Errorf("invalid acknowledgement %v ack index %d: %w", ack, i, err)
		}
		if len(ack.Data) == 0 {
			return fmt.Errorf("invalid acknowledgement %v ack index %d: data bytes cannot be empty", ack, i)
		}
	}

	for i, receipt := range gs.Receipts {
		if err := receipt.Validate(); err != nil {
			return fmt.Errorf("invalid acknowledgement %v ack index %d: %w", receipt, i, err)
		}
	}

	for i, commitment := range gs.Commitments {
		if err := commitment.Validate(); err != nil {
			return fmt.Errorf("invalid commitment %v index %d: %w", commitment, i, err)
		}
		if len(commitment.Data) == 0 {
			return fmt.Errorf("invalid acknowledgement %v ack index %d: data bytes cannot be empty", commitment, i)
		}
	}

	for i, ss := range gs.SendSequences {
		if err := ss.Validate(); err != nil {
			return fmt.Errorf("invalid send sequence %v index %d: %w", ss, i, err)
		}
	}

	for i, rs := range gs.RecvSequences {
		if err := rs.Validate(); err != nil {
			return fmt.Errorf("invalid receive sequence %v index %d: %w", rs, i, err)
		}
	}

	for i, as := range gs.AckSequences {
		if err := as.Validate(); err != nil {
			return fmt.Errorf("invalid acknowledgement sequence %v index %d: %w", as, i, err)
		}
	}

	return nil
}

func validateGenFields(sourceChain, destChain string, sequence uint64) error {
	if err := host.SourceChainValidator(sourceChain); err != nil {
		return fmt.Errorf("invalid port Id: %w", err)
	}
	if err := host.DestChainValidator(destChain); err != nil {
		return fmt.Errorf("invalid channel Id: %w", err)
	}
	if sequence == 0 {
		return errors.New("sequence cannot be 0")
	}
	return nil
}
