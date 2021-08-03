package localhost

import (
	"github.com/bianjieai/tibc-go/modules/tibc/light-clients/09-localhost/types"
)

// Name returns the IBC client name
func Name() string {
	return types.SubModuleName
}
