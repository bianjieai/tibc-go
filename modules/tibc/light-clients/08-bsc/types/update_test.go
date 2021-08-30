package types_test

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
	bsctypes "github.com/bianjieai/tibc-go/modules/tibc/light-clients/08-bsc/types"
)

var (
	url       = "https://bsc-dataseed1.binance.org"
	chainName = "bsc"
	epoch     = uint64(200)
)

func (suite *BSCTestSuite) TestCheckHeaderAndUpdateState() {
	rc := NewRestClient()
	height, err := getBlockHeight(rc)

	suite.NoError(err)

	genesisHeight := height - height%epoch - 2*epoch

	header, err := GetNodeHeader(rc, genesisHeight)
	suite.NoError(err)

	validators, err := bsctypes.ParseValidators(header.Extra)
	suite.NoError(err)

	genesisValidatorHeader, err := GetNodeHeader(rc, genesisHeight-epoch)
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

	for i := uint64(1); i <= uint64(1.5*float64(epoch)); i++ {
		updateHeader, err := GetNodeHeader(rc, genesisHeight+i)

		// skip some connection error on getting header
		if err != nil {
			i--
			continue
		}

		protoHeader := updateHeader.ToHeader()
		suite.NoError(err)

		clientState, consensusState, err = clientState.CheckHeaderAndUpdateState(
			suite.ctx,
			suite.app.AppCodec(),
			suite.app.TIBCKeeper.ClientKeeper.ClientStore(suite.ctx, chainName), // pass in chainName prefixed clientStore
			&protoHeader,
		)

		suite.NoError(err)

		number.RevisionHeight = genesisHeight + i
		suite.app.TIBCKeeper.ClientKeeper.SetClientConsensusState(suite.ctx, chainName, number, consensusState)

		recentSigners, err := bsctypes.GetRecentSigners(suite.app.TIBCKeeper.ClientKeeper.ClientStore(suite.ctx, chainName))
		suite.NoError(err)

		validatorCount := len(clientState.(*bsctypes.ClientState).Validators)
		if i <= uint64(validatorCount)/2+1 {
			suite.Equal(i, uint64(len(recentSigners)))
		} else {
			suite.Equal(uint64(validatorCount)/2+1, uint64(len(recentSigners)))
		}
		suite.Equal(updateHeader.Number.Uint64(), clientState.GetLatestHeight().GetRevisionHeight())
	}
}

type RestClient struct {
	Addr       string
	restClient *http.Client
}

func NewRestClient() *RestClient {
	return &RestClient{
		restClient: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost:   5,
				DisableKeepAlives:     false,
				IdleConnTimeout:       time.Second * 300,
				ResponseHeaderTimeout: time.Second * 300,
				TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: time.Second * 300,
		},
	}
}

func (self *RestClient) SendRestRequest(addr string, data []byte) ([]byte, error) {
	resp, err := self.restClient.Post(addr, "application/json", strings.NewReader(string(data)))
	if err != nil {
		return nil, fmt.Errorf("http post request:%s error:%s", data, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read rest response body error:%s", err)
	}
	return body, nil
}

type heightReq struct {
	JsonRpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	Id      uint     `json:"id"`
}

type heightRsp struct {
	JsonRpc string     `json:"jsonrpc"`
	Result  string     `json:"result,omitempty"`
	Error   *jsonError `json:"error,omitempty"`
	Id      uint       `json:"id"`
}

type jsonError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type blockReq struct {
	JsonRpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      uint          `json:"id"`
}

type blockRsp struct {
	JsonRPC string           `json:"jsonrpc"`
	Result  *ethtypes.Header `json:"result,omitempty"`
	Error   *jsonError       `json:"error,omitempty"`
	Id      uint             `json:"id"`
}

func getBlockHeight(rc *RestClient) (height uint64, err error) {

	req := &heightReq{
		JsonRpc: "2.0",
		Method:  "eth_blockNumber",
		Params:  make([]string, 0),
		Id:      1,
	}
	reqData, err := json.Marshal(req)
	if err != nil {
		return 0, fmt.Errorf("GetNodeHeight: marshal req err: %s", err)
	}
	rspData, err := rc.SendRestRequest(url, reqData)
	if err != nil {
		return 0, fmt.Errorf("GetNodeHeight err: %s", err)
	}

	rsp := &heightRsp{}
	err = json.Unmarshal(rspData, rsp)
	if err != nil {
		return 0, fmt.Errorf("GetNodeHeight, unmarshal resp err: %s", err)
	}
	if rsp.Error != nil {
		return 0, fmt.Errorf("GetNodeHeight, unmarshal resp err: %s", rsp.Error.Message)
	}
	height, err = strconv.ParseUint(rsp.Result, 0, 64)
	if err != nil {
		return 0, fmt.Errorf("GetNodeHeight, parse resp height %s failed", rsp.Result)
	} else {
		return height, nil
	}
}

func GetNodeHeader(restClient *RestClient, height uint64) (*bsctypes.BscHeader, error) {
	params := []interface{}{fmt.Sprintf("0x%x", height), true}
	req := &blockReq{
		JsonRpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  params,
		Id:      1,
	}
	reqdata, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("GetNodeHeight: marshal req err: %s", err)
	}
	rspdata, err := restClient.SendRestRequest(url, reqdata)
	if err != nil {
		return nil, fmt.Errorf("GetNodeHeight err: %s", err)
	}
	rsp := &blockRsp{}
	err = json.Unmarshal(rspdata, rsp)
	if err != nil {
		return nil, fmt.Errorf("GetNodeHeight, unmarshal resp err: %s", err)
	}
	if rsp.Error != nil {
		return nil, fmt.Errorf("GetNodeHeight, return error: %s", rsp.Error.Message)
	}

	if rsp.Result == nil {
		return nil, errors.New("GetNodeHeight, no result")
	}

	header := rsp.Result
	return &bsctypes.BscHeader{
		ParentHash:  header.ParentHash,
		UncleHash:   header.UncleHash,
		Coinbase:    header.Coinbase,
		Root:        header.Root,
		TxHash:      header.TxHash,
		ReceiptHash: header.ReceiptHash,
		Bloom:       bsctypes.Bloom(header.Bloom),
		Difficulty:  header.Difficulty,
		Number:      header.Number,
		GasLimit:    header.GasLimit,
		GasUsed:     header.GasUsed,
		Time:        header.Time,
		Extra:       header.Extra,
		MixDigest:   header.MixDigest,
		Nonce:       bsctypes.BlockNonce(header.Nonce),
	}, nil
}
