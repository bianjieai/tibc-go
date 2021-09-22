package tendermint

import (
	"github.com/bianjieai/tibc-go/modules/tibc/light-clients/07-tendermint/types"
)

// Name returns the TIBC client name
func Name() string {
	return types.SubModuleName
}
