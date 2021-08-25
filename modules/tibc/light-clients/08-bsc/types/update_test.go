package types_test

import (
	"crypto/tls"
	"encoding/json"
	fmt "fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	bsctypes "github.com/bianjieai/tibc-go/modules/tibc/light-clients/08-bsc/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

var (
	url = "https://bsc-dataseed1.binance.org"
)

func (suite *BSCTestSuite) TestCheckHeaderAndUpdateState() {
	rc := NewRestClient()
	height, err := getBlockHeight(rc)

	suite.NoError(err)

	header, err := GetNodeHeader(rc, height)
	suite.NoError(err)

	str, _ := json.Marshal(header)
	println(string(str))

	validators, err := bsctypes.ParseValidators(header.Extra)
	suite.NoError(err)

	clientState := bsctypes.ClientState{
		Header:          header.ToHeader(),
		ChainId:         1,
		Epoch:           200,
		BlockInteval:    3,
		Validators:      validators,
		ContractAddress: []byte("0x00"),
		TrustingPeriod:  200,
	}

	consensusState := bsctypes.ConsensusState{
		Timestamp: header.Time,
		Number:    clienttypes.NewHeight(0, header.Number.Uint64()),
		Root:      header.Root[:],
	}
	suite.app.TIBCKeeper.ClientKeeper.SetClientConsensusState(suite.ctx, clientState, consensusState)
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
		return nil, fmt.Errorf("GetNodeHeight, unmarshal resp err: %s", rsp.Error.Message)
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
