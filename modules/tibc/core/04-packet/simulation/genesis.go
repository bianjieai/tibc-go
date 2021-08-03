package simulation

import (
	"math/rand"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
)

// GenpacketGenesis returns the default packet genesis state.
func GenpacketGenesis(_ *rand.Rand, _ []simtypes.Account) types.GenesisState {
	return types.DefaultGenesisState()
}
