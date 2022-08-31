package tibctesting_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	tmtypes "github.com/tendermint/tendermint/types"

	tibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
	"github.com/bianjieai/tibc-go/modules/tibc/testing/mock"
)

func TestCreateSortedSignerArray(t *testing.T) {
	privVal1 := mock.NewPV()
	pubKey1, err := privVal1.GetPubKey()
	require.NoError(t, err)

	privVal2 := mock.NewPV()
	pubKey2, err := privVal2.GetPubKey()
	require.NoError(t, err)

	validator1 := tmtypes.NewValidator(pubKey1, 1)
	validator2 := tmtypes.NewValidator(pubKey2, 2)

	expected := []tmtypes.PrivValidator{privVal2, privVal1}

	actual := tibctesting.CreateSortedSignerArray(privVal1, privVal2, validator1, validator2)
	require.Equal(t, expected, actual)

	// swap order
	actual = tibctesting.CreateSortedSignerArray(privVal2, privVal1, validator2, validator1)
	require.Equal(t, expected, actual)

	// smaller address
	validator1.Address = []byte{1}
	validator2.Address = []byte{2}
	validator2.VotingPower = 1

	expected = []tmtypes.PrivValidator{privVal1, privVal2}

	actual = tibctesting.CreateSortedSignerArray(privVal1, privVal2, validator1, validator2)
	require.Equal(t, expected, actual)

	// swap order
	actual = tibctesting.CreateSortedSignerArray(privVal2, privVal1, validator2, validator1)
	require.Equal(t, expected, actual)
}
