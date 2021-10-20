package types

import (
	"fmt"
	"sort"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

var (
	_ codectypes.UnpackInterfacesMessage = IdentifiedClientState{}
	_ codectypes.UnpackInterfacesMessage = ClientsConsensusStates{}
	_ codectypes.UnpackInterfacesMessage = ClientConsensusStates{}
	_ codectypes.UnpackInterfacesMessage = GenesisState{}
)

var (
	_              sort.Interface           = ClientsConsensusStates{}
	_              exported.GenesisMetadata = GenesisMetadata{}
	defaultGenesis                          = GenesisState{
		Clients:          []IdentifiedClientState{},
		ClientsConsensus: ClientsConsensusStates{},
		NativeChainName:  "tibc-test",
	}
)

// ClientsConsensusStates defines a slice of ClientConsensusStates that supports the sort interface
type ClientsConsensusStates []ClientConsensusStates

// Len implements sort.Interface
func (ccs ClientsConsensusStates) Len() int { return len(ccs) }

// Less implements sort.Interface
func (ccs ClientsConsensusStates) Less(i, j int) bool { return ccs[i].ChainName < ccs[j].ChainName }

// Swap implements sort.Interface
func (ccs ClientsConsensusStates) Swap(i, j int) { ccs[i], ccs[j] = ccs[j], ccs[i] }

// Sort is a helper function to sort the set of ClientsConsensusStates in place
func (ccs ClientsConsensusStates) Sort() ClientsConsensusStates {
	sort.Sort(ccs)
	return ccs
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (ccs ClientsConsensusStates) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	for _, clientConsensus := range ccs {
		if err := clientConsensus.UnpackInterfaces(unpacker); err != nil {
			return err
		}
	}
	return nil
}

// NewClientConsensusStates creates a new ClientConsensusStates instance.
func NewClientConsensusStates(chainName string, consensusStates []ConsensusStateWithHeight) ClientConsensusStates {
	return ClientConsensusStates{
		ChainName:       chainName,
		ConsensusStates: consensusStates,
	}
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (ccs ClientConsensusStates) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	for _, consStateWithHeight := range ccs.ConsensusStates {
		if err := consStateWithHeight.UnpackInterfaces(unpacker); err != nil {
			return err
		}
	}
	return nil
}

// NewGenesisState creates a GenesisState instance.
func NewGenesisState(
	clients []IdentifiedClientState,
	clientsConsensus ClientsConsensusStates,
	clientsMetadata []IdentifiedGenesisMetadata,
	nativeChainName string,
) GenesisState {
	return GenesisState{
		Clients:          clients,
		ClientsConsensus: clientsConsensus,
		ClientsMetadata:  clientsMetadata,
		NativeChainName:  nativeChainName,
	}
}

// DefaultGenesisState returns the tibc client submodule's default genesis state.
func DefaultGenesisState() GenesisState {
	return defaultGenesis
}

// DefaultGenesisState returns the tibc client submodule's default genesis state.
func SetDefaultGenesisState(genesis GenesisState) {
	defaultGenesis = genesis
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (gs GenesisState) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	for _, client := range gs.Clients {
		if err := client.UnpackInterfaces(unpacker); err != nil {
			return err
		}
	}

	return gs.ClientsConsensus.UnpackInterfaces(unpacker)
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	validClients := make(map[string]string)

	for i, client := range gs.Clients {
		if err := host.ClientIdentifierValidator(client.ChainName); err != nil {
			return fmt.Errorf("invalid client consensus state identifier %s index %d: %w", client.ChainName, i, err)
		}

		clientState, ok := client.ClientState.GetCachedValue().(exported.ClientState)
		if !ok {
			return fmt.Errorf("invalid client state with ID %s", client.ChainName)
		}

		if err := clientState.Validate(); err != nil {
			return fmt.Errorf("invalid client %v index %d: %w", client, i, err)
		}
		// add chain name to validClients map
		validClients[client.ChainName] = clientState.ClientType()
	}

	for _, cc := range gs.ClientsConsensus {
		// check that consensus state is for a client in the genesis clients list
		clientType, ok := validClients[cc.ChainName]
		if !ok {
			return fmt.Errorf("consensus state in genesis has a chain name %s that does not map to a genesis client", cc.ChainName)
		}

		for i, consensusState := range cc.ConsensusStates {
			if consensusState.Height.IsZero() {
				return fmt.Errorf("consensus state height cannot be zero")
			}

			cs, ok := consensusState.ConsensusState.GetCachedValue().(exported.ConsensusState)
			if !ok {
				return fmt.Errorf("invalid consensus state with client ID %s at height %s", cc.ChainName, consensusState.Height)
			}

			if err := cs.ValidateBasic(); err != nil {
				return fmt.Errorf("invalid client consensus state %v chainName %s index %d: %w", cs, cc.ChainName, i, err)
			}

			// ensure consensus state type matches client state type
			if clientType != cs.ClientType() {
				return fmt.Errorf("consensus state client type %s does not equal client state client type %s", cs.ClientType(), clientType)
			}

		}
	}

	for _, clientMetadata := range gs.ClientsMetadata {
		// check that metadata is for a client in the genesis clients list
		_, ok := validClients[clientMetadata.ChainName]
		if !ok {
			return fmt.Errorf("metadata in genesis has a chain name %s that does not map to a genesis client", clientMetadata.ChainName)
		}

		for i, gm := range clientMetadata.Metadata {
			if err := gm.Validate(); err != nil {
				return fmt.Errorf("invalid client metadata %v chainName %s index %d: %w", gm, clientMetadata.ChainName, i, err)
			}

		}

	}

	if err := host.ClientIdentifierValidator(gs.NativeChainName); err != nil {
		return err
	}
	return nil
}

// NewGenesisMetadata is a constructor for GenesisMetadata
func NewGenesisMetadata(key, val []byte) GenesisMetadata {
	return GenesisMetadata{
		Key:   key,
		Value: val,
	}
}

// GetKey returns the key of metadata. Implements exported.GenesisMetadata interface.
func (gm GenesisMetadata) GetKey() []byte {
	return gm.Key
}

// GetValue returns the value of metadata. Implements exported.GenesisMetadata interface.
func (gm GenesisMetadata) GetValue() []byte {
	return gm.Value
}

// Validate ensures key and value of metadata are not empty
func (gm GenesisMetadata) Validate() error {
	if len(gm.Key) == 0 {
		return fmt.Errorf("genesis metadata key cannot be empty")
	}
	if len(gm.Value) == 0 {
		return fmt.Errorf("genesis metadata value cannot be empty")
	}
	return nil
}

// NewIdentifiedGenesisMetadata takes in a client ID and list of genesis metadata for that client
// and constructs a new IdentifiedGenesisMetadata.
func NewIdentifiedGenesisMetadata(chainName string, gms []GenesisMetadata) IdentifiedGenesisMetadata {
	return IdentifiedGenesisMetadata{
		ChainName: chainName,
		Metadata:  gms,
	}
}
