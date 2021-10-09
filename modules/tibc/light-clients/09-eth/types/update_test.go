package types_test

import (
	"encoding/json"
	"io/ioutil"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
	tibcethtypes "github.com/bianjieai/tibc-go/modules/tibc/light-clients/09-eth/types"
)

var chainName = "eth"

func (suite *ETHTestSuite) TestCheckHeaderAndUpdateState() {
	var updateHeaders []*tibcethtypes.EthHeader
	updateHeadersBz, _ := ioutil.ReadFile("testdata/update_headers.json")
	err := json.Unmarshal(updateHeadersBz, &updateHeaders)
	suite.NoError(err)
	suite.GreaterOrEqual(len(updateHeaders), 1)

	header := updateHeaders[0]

	number := clienttypes.NewHeight(0, header.Number.Uint64())

	clientState := exported.ClientState(&tibcethtypes.ClientState{
		Header:          header.ToHeader(),
		ChainId:         1,
		ContractAddress: []byte("0x00"),
		TrustingPeriod:  200000,
		TimeDelay:       0,
		BlockDelay:      1,
	})

	consensusState := exported.ConsensusState(&tibcethtypes.ConsensusState{
		Timestamp: header.Time,
		Number:    number,
		Root:      header.Root[:],
	})

	suite.app.TIBCKeeper.ClientKeeper.SetClientConsensusState(suite.ctx, chainName, number, consensusState)
	protoHeader := header.ToHeader()
	store := suite.app.TIBCKeeper.ClientKeeper.ClientStore(suite.ctx, chainName)
	headerBytes, err := suite.app.AppCodec().MarshalInterface(&protoHeader)
	suite.NoError(err)

	tibcethtypes.SetEthHeaderIndex(store, protoHeader, headerBytes)
	tibcethtypes.SetEthConsensusRoot(store, protoHeader.Height.RevisionHeight, protoHeader.ToEthHeader().Root, header.Hash())

	for _, updateHeader := range updateHeaders[1:5] {
		protoHeader := updateHeader.ToHeader()
		suite.NoError(err)

		clientState, consensusState, err = clientState.CheckHeaderAndUpdateState(
			suite.ctx,
			suite.app.AppCodec(),
			suite.app.TIBCKeeper.ClientKeeper.ClientStore(suite.ctx, chainName), // pass in chainName prefixed clientStore
			&protoHeader,
		)

		suite.NoError(err)

		number.RevisionHeight = protoHeader.Height.RevisionHeight
		suite.app.TIBCKeeper.ClientKeeper.SetClientConsensusState(suite.ctx, chainName, number, consensusState)

		suite.Equal(updateHeader.Number.Uint64(), clientState.GetLatestHeight().GetRevisionHeight())
	}
}
