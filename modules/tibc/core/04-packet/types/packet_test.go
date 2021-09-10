package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	"github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
)

func TestCommitPacket(t *testing.T) {
	packet := types.NewPacket(validPacketData, 1, sourceChain, destChain, relayChain, port)

	registry := codectypes.NewInterfaceRegistry()
	clienttypes.RegisterInterfaces(registry)
	types.RegisterInterfaces(registry)

	commitment := types.CommitPacket(&packet)
	require.NotNil(t, commitment)
}

func TestPacketValidateBasic(t *testing.T) {
	testCases := []struct {
		packet  types.Packet
		expPass bool
		errMsg  string
	}{
		{types.NewPacket(validPacketData, 1, sourceChain, destChain, relayChain, port), true, ""},
		{types.NewPacket(validPacketData, 0, sourceChain, destChain, relayChain, port), false, "invalid sequence"},
		// {types.NewPacket(validPacketData, 1, invalidPort, destChain, relayChain, port), false, "invalid source port"},
		// {types.NewPacket(validPacketData, 1, sourceChain, destChain, relayChain, invalidPort), false, "invalid port"},
		{types.NewPacket(unknownPacketData, 1, sourceChain, destChain, relayChain, port), true, ""},
	}

	for i, tc := range testCases {
		err := tc.packet.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "Case %d failed: %s", i, tc.errMsg)
		} else {
			require.Error(t, err, "Invalid Case %d passed: %s", i, tc.errMsg)
		}
	}
}
