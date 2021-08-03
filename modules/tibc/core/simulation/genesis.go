package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"

	clientsims "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/simulation"
	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	packetsims "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/simulation"
	packettypes "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	"github.com/bianjieai/tibc-go/modules/tibc/core/types"
)

// Simulation parameter constants
const (
	clientGenesis = "client_genesis"
	packetGenesis = "packet_genesis"
)

// RandomizedGenState generates a random GenesisState for evidence
func RandomizedGenState(simState *module.SimulationState) {
	var (
		clientGenesisState clienttypes.GenesisState
		packetGenesisState packettypes.GenesisState
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, clientGenesis, &clientGenesisState, simState.Rand,
		func(r *rand.Rand) { clientGenesisState = clientsims.GenClientGenesis(r, simState.Accounts) },
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, packetGenesis, &packetGenesisState, simState.Rand,
		func(r *rand.Rand) { packetGenesisState = packetsims.GenpacketGenesis(r, simState.Accounts) },
	)

	ibcGenesis := types.GenesisState{
		ClientGenesis: clientGenesisState,
		PacketGenesis: packetGenesisState,
	}

	bz, err := json.MarshalIndent(&ibcGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated %s parameters:\n%s\n", host.ModuleName, bz)
	simState.GenState[host.ModuleName] = simState.Cdc.MustMarshalJSON(&ibcGenesis)
}
