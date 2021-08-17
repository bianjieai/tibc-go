package exported

// PacketI defines the standard interface for TIBC clean packets
type PacketI interface {
	GetSequence() uint64
	GetPort() string
	GetSourceChain() string
	GetDestChain() string
	GetRelayChain() string
	GetData() []byte
	ValidateBasic() error
}

// CleanPacketI defines the standard interface for TIBC clean packets
type CleanPacketI interface {
	GetSequence() uint64
	GetSourceChain() string
	GetDestChain() string
	GetRelayChain() string
	ValidateBasic() error
}
