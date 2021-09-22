package bsc

import (
	"github.com/bianjieai/tibc-go/modules/tibc/light-clients/08-bsc/types"
)

// Name returns the TIBC client name
func Name() string {
	return types.SubModuleName
}
