package types_test

import (
	"encoding/json"
	"io/ioutil"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
	bsctypes "github.com/bianjieai/tibc-go/modules/tibc/light-clients/08-bsc/types"
)

var (
	chainName = "bsc"
	epoch     = uint64(200)
)

func (suite *BSCTestSuite) TestCheckHeaderAndUpdateState() {
	var genesisState GenesisState
	genesisStateBz, _ := ioutil.ReadFile("testdata/genesis_state.json")
	err := json.Unmarshal(genesisStateBz, &genesisState)
	suite.NoError(err)

	header := genesisState.GenesisHeader
	genesisValidatorHeader := genesisState.GenesisValidatorHeader

	validators, err := bsctypes.ParseValidators(header.Extra)
	suite.NoError(err)

	genesisValidators, err := bsctypes.ParseValidators(genesisValidatorHeader.Extra)
	suite.NoError(err)

	number := clienttypes.NewHeight(0, header.Number.Uint64())

	clientState := exported.ClientState(&bsctypes.ClientState{
		Header:          header.ToHeader(),
		ChainId:         56,
		Epoch:           epoch,
		BlockInteval:    3,
		Validators:      genesisValidators,
		ContractAddress: []byte("0x00"),
		TrustingPeriod:  200,
	})

	consensusState := exported.ConsensusState(&bsctypes.ConsensusState{
		Timestamp: header.Time,
		Number:    number,
		Root:      header.Root[:],
	})

	suite.app.TIBCKeeper.ClientKeeper.SetClientConsensusState(suite.ctx, chainName, number, consensusState)

	bsctypes.SetPendingValidators(suite.app.TIBCKeeper.ClientKeeper.ClientStore(suite.ctx, chainName), suite.app.AppCodec(), validators)

	var updateHeaders []*bsctypes.BscHeader
	updateHeadersBz, _ := ioutil.ReadFile("testdata/update_headers.json")
	err = json.Unmarshal(updateHeadersBz, &updateHeaders)
	suite.NoError(err)
	suite.Equal(int(1.5*float64(epoch)), len(updateHeaders))
	for i, updateHeader := range updateHeaders {
		updateHeaders = append(updateHeaders, updateHeader)

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

		recentSigners, err := bsctypes.GetRecentSigners(suite.app.TIBCKeeper.ClientKeeper.ClientStore(suite.ctx, chainName))
		suite.NoError(err)

		validatorCount := len(clientState.(*bsctypes.ClientState).Validators)
		if i+1 <= validatorCount/2+1 {
			suite.Equal(i+1, len(recentSigners))
		} else {
			suite.Equal(validatorCount/2+1, len(recentSigners))
		}
		suite.Equal(updateHeader.Number.Uint64(), clientState.GetLatestHeight().GetRevisionHeight())
	}
}

type GenesisState struct {
	GenesisHeader          *bsctypes.BscHeader `json:"genesis_header"`
	GenesisValidatorHeader *bsctypes.BscHeader `json:"genesis_validator_header"`
}
