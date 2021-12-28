package common

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config/lookup"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/sirupsen/logrus"
	"gofabric/models"
	"os"
	"strings"
	"time"

	localConfig "gofabric/configs"
	util "gofabric/utils"

	"crypto/sha256"
	"encoding/asn1"
	"github.com/gofrs/uuid"
)

var log = logrus.New()

var mainSDK *fabsdk.FabricSDK
var ledgerClient *ledger.Client
var resMgmtClient *resmgmt.Client
var client *models.Client
var channelClient *channel.Client

// InitChainExploreService Init ChainExplore Client (ledger client)
func InitChainExploreService(cfg, org, admin, user string) {

	log.Println("Initialize the client")

	client = &models.Client{
		ConfigPath: cfg,
		OrgName:    org,
		OrgAdmin:   admin,
		OrgUser:    user,

		ChaincodeID:   localConfig.ChaincodeID,
		ChaincodePath: localConfig.ChaincodePath,
		GoPath:        os.Getenv("GOPATH"),
		ChannelID:     localConfig.ChannelID,
	}

	var err error
	// create fabsdk SDk
	mainSDK, err = fabsdk.New(config.FromFile(client.ConfigPath))
	if err != nil {
		log.Panicf("Failed to create an new SDK:%s\n", err)
	}

	client.SDK = mainSDK
	log.Println("Success to create an new SDK")
	// get channel context
	userChannelContext := mainSDK.ChannelContext(client.ChannelID, fabsdk.WithUser(client.OrgUser))
	// ledger client
	ledgerClient, err = ledger.New(userChannelContext)
	if err != nil {
		log.Printf("Failed to create an new ledgerClient:%s\n", err)
	}
	log.Println("Success to create an new ledgerClient ")

	adminContext := mainSDK.Context(fabsdk.WithUser(client.OrgAdmin), fabsdk.WithOrg(client.OrgName))

	resMgmtClient, err = resmgmt.New(adminContext)
	if err != nil {
		log.Printf("Failed to create an new orgResMgmt:%s\n", err)
	}

	channelClient, err = channel.New(userChannelContext)
	if err != nil {
		log.Printf("Failed to create an new channelClient:%s\n", err)
	}

	//defer client.SDK.Close()
	//ccfig,err := ledgerClient.QueryConfig()
	//fmt.Println(ccfig)
}

// QueryInstalledCC Query installed chaincode
func QueryInstalledCC() ([]*models.Chaincode, error) {

	configBackend, err := mainSDK.Config()
	if err != nil {
		log.Printf("Failed to get mainSDK config:%s \n", err)
		return nil, err
	}

	targets, err := orgTargetPeers([]string{client.OrgName}, configBackend)
	if err != nil {
		log.Printf("Failed to get targets:%s \n", err)
		return nil, err
	}
	peer := targets[0]

	var chaincodeInfos []*models.Chaincode

	installedCC, err := resMgmtClient.LifecycleQueryInstalledCC(resmgmt.WithTargetEndpoints(peer), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	for _, cc := range installedCC {

		chaincodeInfo := &models.Chaincode{}
		chaincodeInfo.Label = cc.Label
		chaincodeInfo.PackageID = cc.PackageID
		chaincodeInfo.References = cc.References

		chaincodeInfos = append(chaincodeInfos, chaincodeInfo)

	}

	return chaincodeInfos, nil
}

func orgTargetPeers(orgs []string, configBackend ...core.ConfigBackend) ([]string, error) {

	networkConfig := fab.NetworkConfig{}

	err := lookup.New(configBackend...).UnmarshalKey("organizations", &networkConfig.Organizations)
	if err != nil {
		return nil, err
	}

	var peers []string
	for _, org := range orgs {
		orgConfig, ok := networkConfig.Organizations[strings.ToLower(org)]
		if !ok {
			continue
		}
		peers = append(peers, orgConfig.Peers...)
	}
	return peers, nil
}

// QueryLedgerInfo Query ledger info
func QueryLedgerInfo() (*fab.BlockchainInfoResponse, error) {

	log.Println("Query ledger info")

	ledgerInfo, err := ledgerClient.QueryInfo()
	if err != nil {
		log.Printf("Failed to query blockChain info: %s\n ", err)
		return nil, err
	}

	return ledgerInfo, nil

}

// QueryLastesBlocksInfo  Query last 5 Blocks info
func QueryLastesBlocksInfo() ([]*models.Block, error) {

	ledgerInfo, err := ledgerClient.QueryInfo()
	if err != nil {
		log.Printf("Failed to query last 5 Blocks info:%s \n", err)
		return nil, err
	}

	var lastesBlockList []*models.Block
	lastesBlockNum := ledgerInfo.BCI.Height - 1

	for i := lastesBlockNum; i > 0 && i > (lastesBlockNum-5); i-- {
		block, err := QueryBlockByBlockNum(int64(i))
		if err != nil {
			log.Printf("Failed to Query last 5 Blocks info:%s \n", err)
			return nil, err
		}
		lastesBlockList = append(lastesBlockList, block)
	}

	return lastesBlockList, nil
}

func QueryAllBlocksInfo() ([]*models.Block, error) {
	ledgerInfo, err := ledgerClient.QueryInfo()
	if err != nil {
		log.Printf("Failed to query last 5 Blocks info:%s \n", err)
		return nil, err
	}

	var blockList []*models.Block
	blockNum := ledgerInfo.BCI.Height - 1

	for i := blockNum; i >= 3; i-- {
		block, err := QueryBlockByBlockNum(int64(i))
		if err != nil {
			log.Printf("Failed to Query last 5 Blocks info:%s \n", err)
			return nil, err
		}
		blockList = append(blockList, block)
	}

	return blockList, nil
}

// QueryBlockInfoByHash  Query one block by blockHash
func QueryBlockInfoByHash(blockHash []byte) (*models.Block, error) {
	rawBlockInfo, err := ledgerClient.QueryBlockByHash(blockHash)
	if err != nil {
		log.Printf("Failed to query block by blockHash:%s \n", err)
		return nil, err
	}
	block, err := QueryBlockByBlockNum(int64(rawBlockInfo.GetHeader().Number))
	if err != nil {
		log.Printf("Failed to query block by blockHash QueryBlockByBlockNum:%s \n", err)
		return nil, err
	}

	return block, nil
}

// QueryBlockMainInfo Query the main config of channel
func QueryBlockMainInfo() (*models.BlockMainInfo, error) {

	ledgerInfo, err := ledgerClient.QueryInfo()
	if err != nil {
		log.Printf("Failed to Query the main config of channel:%s \n", err)
		return nil, err
	}

	blockNum := ledgerInfo.BCI.Height - 1

	var txNum uint64
	for i := blockNum; i >= 7; i-- {
		rawBlock, err := ledgerClient.QueryBlock(uint64(i))
		if err != nil {
			log.Printf("Failed to Query the main config of channel:%s \n", err)
			return nil, err
		}

		txNum = txNum + uint64(len(rawBlock.Data.Data))
	}

	mainBlockInfo := &models.BlockMainInfo{
		BlockNum:       ledgerInfo.BCI.Height,
		TransactionNum: txNum,
		ChaincodeNum:   2,
		NodeNum:        3,
	}

	return mainBlockInfo, err
}

// QueryBlockByBlockNum Query Block info by block's number
func QueryBlockByBlockNum(num int64) (*models.Block, error) {

	rawBlock, err := ledgerClient.QueryBlock(uint64(num))
	if err != nil {
		log.Printf("Failed to query Block info by block's number : %s \n", err)
		return nil, err
	}

	// parse the block body

	var txList []*models.Transaction

	for i := range rawBlock.Data.Data {
		rawEnvelope, err := util.GetEnvelopeFromBlock(rawBlock.Data.Data[i])
		if err != nil {
			log.Printf("Failed to GetEnvelopeFromBlock: %s \n", err)
			return nil, err
		}

		transaction, err := util.GetTxFromEnvelopeDeep(rawEnvelope)
		if err != nil {
			log.Printf("Failed to GetTxFromEnvelopeDeep: %s \n", err)
			return nil, err
		}

		for i := range transaction.TransactionActionList {
			transaction.TransactionActionList[i].BlockNum = rawBlock.Header.Number
		}

		txList = append(txList, transaction)
	}

	blockHash := GetBlockHash(rawBlock.Header)
	block := models.Block{

		Number:          rawBlock.Header.Number,
		PreviousHash:    hex.EncodeToString(rawBlock.Header.PreviousHash),
		DataHash:        hex.EncodeToString(rawBlock.Header.DataHash),
		BlockHash:       hex.EncodeToString(blockHash),
		TxNum:           len(rawBlock.Data.Data),
		TransactionList: txList,
		CreateTime:      txList[0].TransactionActionList[0].Timestamp,
	}

	return &block, nil
}

func GetBlockHash(blockHeader *common.BlockHeader) []byte {

	rawBlockHeader := models.BlockHeader{
		Number:       int8(blockHeader.Number),
		PreviousHash: blockHeader.PreviousHash,
		DataHash:     blockHeader.DataHash,
	}

	data, err := asn1.Marshal(rawBlockHeader)
	if err != nil {
		log.Printf("Failed to GetBlockHash : %s \n", err)
	}

	h := sha256.New()
	h.Write(data)
	byteHash := h.Sum(nil)
	return byteHash
}

func QueryTxByTxId(txId string) (*models.Transaction, error) {
	rawTx, err := ledgerClient.QueryTransaction(fab.TransactionID(txId))
	if err != nil {
		log.Printf("Failed to QueryTxByTxId rawTx: %s \n", err)
		return nil, err
	}

	transaction, err := util.GetTxFromEnvelopeDeep(rawTx.TransactionEnvelope)
	if err != nil {
		log.Printf("Failed to QueryTxByTxId transaction: %s \n", err)
		return nil, err
	}

	block, err := ledgerClient.QueryBlockByTxID(fab.TransactionID(txId))
	if err != nil {
		log.Printf("Failed to QueryTxByTxId block: %s \n", err)
		return nil, err
	}

	for i := range transaction.TransactionActionList {
		transaction.TransactionActionList[i].BlockNum = block.Header.Number
	}

	return transaction, nil
}

func QueryTxByTxIdJsonStr(txId string) (string, error) {
	transaction, err := QueryTxByTxId(txId)
	if err != nil {
		log.Printf("Failed to QueryTxByTxIdJsonStr transaction: %s \n", err)
		return "", err
	}

	jsonStr, err := json.Marshal(transaction)
	return string(jsonStr), err
}

// OperateLedgerTest Test for operate ledger
func OperateLedgerTest() {

	bci, err := ledgerClient.QueryInfo()
	if err != nil {
		fmt.Printf("failed to query for blockchain info: %s\n", err)
	}

	if bci != nil {
		fmt.Println("Retrieved ledger info")
	}
	fmt.Println(bci.BCI.Height)
	fmt.Println(bci.BCI.PreviousBlockHash)

	rawBlock, err := ledgerClient.QueryBlock(uint64(4))
	fmt.Println(rawBlock.GetHeader().GetDataHash())

}

// QueryChannels Query channels
func QueryChannels() []string {

	configBackend, err := mainSDK.Config()
	if err != nil {
		log.Printf("Failed to get config backend from SDK: %s", err)
	}

	targets, err := orgTargetPeers([]string{localConfig.OrgGo}, configBackend)
	if err != nil {

		log.Printf("Creating peers failed: %s", err)
	}

	channelQueryResponse, err := resMgmtClient.QueryChannels(
		resmgmt.WithTargetEndpoints(targets[0]),
		resmgmt.WithRetry(retry.DefaultResMgmtOpts))

	if err != nil {
		log.Printf("QueryChannels return error: %s", err)
	}

	var channels []string

	for _, channel := range channelQueryResponse.Channels {
		channels = append(channels, channel.ChannelId)
	}
	return channels
}

// QueryChannelInfo Query channel info
func QueryChannelInfo() (*models.ChannelInfo, error) {

	ledgerInfo, err := ledgerClient.QueryInfo()
	if err != nil {
		log.Printf("Failed to Query the main config of channel:%s \n", err)
		return nil, err
	}

	cInfo := &models.ChannelInfo{
		Name:      localConfig.ChannelID,
		Blocks:    int(ledgerInfo.BCI.Height),
		Txs:       int(ledgerInfo.BCI.Height) - 3,
		Timestamp: "2021-10-13T07:41:30.000Z",
	}

	return cInfo, nil

}

func WalletTest() {
	wallet, err := gateway.NewFileSystemWallet("walletTest")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(wallet)
	fmt.Println(wallet.Exists("User"))

}

func QueryInfoByChaincode(uuid string) (*models.ChaincodeInfo, error) {

	var args [][]byte
	args = append(args, []byte(uuid))

	req := channel.Request{
		ChaincodeID: localConfig.ChaincodeID,
		Fcn:         "get",
		Args:        args,
	}

	resp, err := channelClient.Query(req)
	if err != nil {
		log.Printf("Failed to QueryInfoByChaincode:%s\n", err)
		return nil, err
	}

	chaincodeInfo := &models.ChaincodeInfo{
		Uuid:    uuid,
		TxId:    string(resp.TransactionID),
		Time:    time.Now(),
		Payload: string(resp.Payload),
	}

	log.Printf("query chaincode tx : %s", resp.TransactionID)
	log.Printf("result : %v", string(resp.Payload))

	return chaincodeInfo, nil

}

func InvokeInfoByChaincode(data string) (*models.ChaincodeInfo, error) {

	v4, err := uuid.NewV4()
	key := []byte(v4.String())

	var args [][]byte
	args = append(args, key)
	args = append(args, []byte(data))

	log.Printf("uuid : %v", string(key))

	req := channel.Request{
		ChaincodeID: localConfig.ChaincodeID,
		Fcn:         "set",
		Args:        args,
	}

	resp, err := channelClient.Execute(req)
	if err != nil {
		log.Printf("Failed to invokeInfoByChaincode:%s\n", err)
		return nil, err
	}

	transaction, err := QueryTxByTxId(string(resp.TransactionID))
	if err != nil {
		return nil, err
	}

	chaincodeInfo := &models.ChaincodeInfo{
		Uuid:            v4.String(),
		TxId:            string(resp.TransactionID),
		Time:            time.Now(),
		Payload:         string(resp.Payload),
		TransactionInfo: *transaction,
	}

	log.Printf("query chaincode tx : %s", resp.TransactionID)
	log.Printf("result : %v", string(resp.Payload))

	return chaincodeInfo, nil
}
