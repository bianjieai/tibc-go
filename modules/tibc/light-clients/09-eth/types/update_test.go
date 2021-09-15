package types_test

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
	tibcethtypes "github.com/bianjieai/tibc-go/modules/tibc/light-clients/09-eth/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

var (
	ethurl     = "https://mainnet.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161"
	EthConType = "{\"@type\":\"/tibc.lightclients.eth.v1.ConsensusState\","
	EthStaType = "{\"@type\":\"/tibc.lightclients.eth.v1.ClientState\","
	chainName  = "eth"
	epoch      = uint64(200)
)

func (suite *ETHTestSuite) TestCheckHeaderAndUpdateState() {
	rc := NewRestClient()
	height, err := GetBlockHeight(rc, ethurl)

	suite.NoError(err)
	genesisHeight := height - height%epoch - 2*epoch
	header, err := GetNodeHeader(rc, ethurl, genesisHeight)
	suite.NoError(err)

	number := clienttypes.NewHeight(0, header.Number.Uint64())

	clientState := exported.ClientState(&tibcethtypes.ClientState{
		Header:          header.ToHeader(),
		ChainId:         1,
		ContractAddress: []byte("0x00"),
		TrustingPeriod:  200,
	})

	consensusState := exported.ConsensusState(&tibcethtypes.ConsensusState{
		Timestamp: header.Time,
		Number:    number,
		Root:      header.Root[:],
		Header:    header.ToHeader(),
	})

	suite.app.TIBCKeeper.ClientKeeper.SetClientConsensusState(suite.ctx, chainName, number, consensusState)

	store := suite.app.TIBCKeeper.ClientKeeper.ClientStore(suite.ctx, chainName)
	marshalInterface, err := suite.app.AppCodec().MarshalInterface(consensusState)
	suite.NoError(err)
	store.Set(tibcethtypes.ConsensusStateIndexKey(header.Hash()), marshalInterface)
	// change  uint64(1.5*float64(epoch)) to 50 for test time
	for i := uint64(1); i <= uint64(50); i++ {

		updateHeader, err := GetNodeHeader(rc, ethurl, genesisHeight+i)

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

func GetBlockHeight(rc *RestClient, url string) (height uint64, err error) {

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

func GetNodeHeader(restClient *RestClient, url string, height uint64) (*tibcethtypes.EthHeader, error) {
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
	return &tibcethtypes.EthHeader{
		ParentHash:  header.ParentHash,
		UncleHash:   header.UncleHash,
		Coinbase:    header.Coinbase,
		Root:        header.Root,
		TxHash:      header.TxHash,
		ReceiptHash: header.ReceiptHash,
		Bloom:       header.Bloom,
		Difficulty:  header.Difficulty,
		Number:      header.Number,
		GasLimit:    header.GasLimit,
		GasUsed:     header.GasUsed,
		Time:        header.Time,
		Extra:       header.Extra,
		MixDigest:   header.MixDigest,
		Nonce:       header.Nonce,
		BaseFee:     header.BaseFee,
	}, nil
}

func Test_getjson(test *testing.T) {
	rc := NewRestClient()
	height, err := GetBlockHeight(rc, ethurl)
	if err != nil {
		fmt.Println(err)
		return
	}
	height = height - 258
	header, err := GetNodeHeader(rc, ethurl, height)
	number := clienttypes.NewHeight(0, header.Number.Uint64())
	clientState := exported.ClientState(&tibcethtypes.ClientState{
		Header:          header.ToHeader(),
		ChainId:         56,
		ContractAddress: []byte("0x00"),
		TrustingPeriod:  200,
		TimeDelay:       0,
		BlockDelay:      7,
	})

	consensusState := exported.ConsensusState(&tibcethtypes.ConsensusState{
		Timestamp: header.Time,
		Number:    number,
		Root:      header.Root[:],
		Header:    header.ToHeader(),
	})
	b0, err := json.Marshal(clientState)
	if err != nil {
		return
	}
	b1, err := json.Marshal(consensusState)
	if err != nil {
		return
	}
	b0 = []byte(EthStaType + string(b0)[1:])
	clientStateName := "eth_client_state.json"
	err = ioutil.WriteFile(clientStateName, b0, os.ModeAppend)
	if err != nil {
		return
	}
	b1 = []byte(EthConType + string(b1)[1:])
	clientConsensusStateName := "eth_consensus_state.json"
	err = ioutil.WriteFile(clientConsensusStateName, b1, os.ModeAppend)
	if err != nil {
		return
	}
}
func TestVerifyHeader(test *testing.T) {
	cachedir, err := ioutil.TempDir("", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer os.RemoveAll(cachedir)
	rc := NewRestClient()
	height, err := GetBlockHeight(rc, ethurl)
	if err != nil {
		fmt.Println(err)
		return
	}
	height = height - 60

	nodeHeader, err := GetNodeHeader(rc, ethurl, 13177652)
	config := tibcethtypes.Config{
		CacheDir:     cachedir,
		CachesOnDisk: 1,
	}
	ethash := tibcethtypes.New(config, nil, false)
	defer ethash.Close()
	if err := ethash.VerifySeal(nodeHeader.ToHeader().ToVerifyHeader(), false); err != nil {
		fmt.Println(err)
		return
	}
}
