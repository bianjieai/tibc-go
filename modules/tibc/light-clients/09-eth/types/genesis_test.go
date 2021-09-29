package types_test

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/common"

	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
	tibcethtypes "github.com/bianjieai/tibc-go/modules/tibc/light-clients/09-eth/types"
)

func (suite *ETHTestSuite) TestExportMetadata() {
	var updateHeaders []*tibcethtypes.EthHeader
	updateHeadersBz, _ := ioutil.ReadFile("testdata/update_headers.json")
	err := json.Unmarshal(updateHeadersBz, &updateHeaders)
	suite.NoError(err)
	suite.GreaterOrEqual(len(updateHeaders), 1)
	header := updateHeaders[0]
	clientState := exported.ClientState(&tibcethtypes.ClientState{
		Header:          header.ToHeader(),
		ChainId:         1,
		ContractAddress: []byte("0x00"),
		TrustingPeriod:  200000,
		TimeDelay:       0,
		BlockDelay:      1,
	})
	suite.app.TIBCKeeper.ClientKeeper.SetClientState(suite.ctx, chainName, clientState)
	gm := clientState.ExportMetadata(suite.app.TIBCKeeper.ClientKeeper.ClientStore(suite.ctx, chainName))
	suite.Require().Nil(gm, "client with no metadata returned non-nil exported metadata")

	protoHeader := header.ToHeader()
	store := suite.app.TIBCKeeper.ClientKeeper.ClientStore(suite.ctx, chainName)
	headerBytes, err := suite.app.AppCodec().MarshalInterface(&protoHeader)
	suite.NoError(err)

	tibcethtypes.SetEthHeaderIndex(store, protoHeader, headerBytes)
	tibcethtypes.SetEthConsensusRoot(store, protoHeader.Height.RevisionHeight, protoHeader.ToEthHeader().Root, header.Hash())

	gm = clientState.ExportMetadata(suite.app.TIBCKeeper.ClientKeeper.ClientStore(suite.ctx, chainName))
	suite.Require().NotNil(gm, "client with metadata returned nil exported metadata")
	suite.Require().Len(gm, 2, "exported metadata has unexpected length")

	suite.Require().Equal(tibcethtypes.EthHeaderIndexKey(protoHeader.Hash(), protoHeader.Height.RevisionHeight), gm[0].GetKey(), "metadata has unexpected key")
	suite.Require().Equal(headerBytes, gm[0].GetValue(), "metadata has unexpected value")

	suite.Require().Equal(tibcethtypes.EthRootMainKey(common.BytesToHash(protoHeader.Root), protoHeader.Height.RevisionHeight), gm[1].GetKey(), "metadata has unexpected key")
	suite.Require().Equal(gm[0].GetKey(), gm[1].GetValue(), "metadata has unexpected value")
}
