package explore

import (
	"crypto/sha256"
	"encoding/asn1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config/lookup"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	log "github.com/sirupsen/logrus"
	com "gofabric/common"
	localConfig "gofabric/configs"
	"gofabric/models"
	"gofabric/utils"
	"strings"
	"time"
)

type ExploreClient struct {
	com.Foo
}

// QueryInstalledCC Query installed chaincode
func (e *ExploreClient) QueryInstalledCC() ([]*models.Chaincode, error) {

	configBackend, err := e.Option.MainSDK.Config()
	if err != nil {
		log.Printf("Failed to get mainSDK config:%s \n", err)
		return nil, err
	}

	targets, err := orgTargetPeers([]string{e.Option.OrgName}, configBackend)
	if err != nil {
		log.Printf("Failed to get targets:%s \n", err)
		return nil, err
	}
	peer := targets[0]

	var chaincodeInfos []*models.Chaincode

	installedCC, err := e.Option.ResMgmtClient.LifecycleQueryInstalledCC(resmgmt.WithTargetEndpoints(peer), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
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
func (e *ExploreClient) QueryLedgerInfo() (*fab.BlockchainInfoResponse, error) {

	log.Println("Query ledger info")

	ledgerInfo, err := e.Option.LedgerClient.QueryInfo()
	if err != nil {
		log.Printf("Failed to query blockChain info: %s\n ", err)
		return nil, err
	}

	return ledgerInfo, nil

}

// QueryLastesBlocksInfo  Query last 5 Blocks info
func (e *ExploreClient) QueryLastesBlocksInfo() ([]*models.Block, error) {

	ledgerInfo, err := e.Option.LedgerClient.QueryInfo()
	if err != nil {
		log.Printf("Failed to query last 5 Blocks info:%s \n", err)
		return nil, err
	}

	var lastesBlockList []*models.Block
	lastesBlockNum := ledgerInfo.BCI.Height - 1

	for i := lastesBlockNum; i > 0 && i > (lastesBlockNum-5); i-- {
		block, err := e.QueryBlockByBlockNum(int64(i))
		if err != nil {
			log.Printf("Failed to Query last 5 Blocks info:%s \n", err)
			return nil, err
		}
		lastesBlockList = append(lastesBlockList, block)
	}

	return lastesBlockList, nil
}

func (e *ExploreClient) QueryAllBlocksInfo() ([]*models.Block, error) {
	ledgerInfo, err := e.Option.LedgerClient.QueryInfo()
	if err != nil {
		log.Printf("Failed to query last 5 Blocks info:%s \n", err)
		return nil, err
	}

	var blockList []*models.Block
	blockNum := ledgerInfo.BCI.Height - 1

	for i := blockNum; i >= 3; i-- {
		block, err := e.QueryBlockByBlockNum(int64(i))
		if err != nil {
			log.Printf("Failed to Query last 5 Blocks info:%s \n", err)
			return nil, err
		}
		blockList = append(blockList, block)
	}

	return blockList, nil
}

// QueryBlockInfoByHash  Query one block by blockHash
func (e *ExploreClient) QueryBlockInfoByHash(blockHash []byte) (*models.Block, error) {
	rawBlockInfo, err := e.Option.LedgerClient.QueryBlockByHash(blockHash)
	if err != nil {
		log.Printf("Failed to query block by blockHash:%s \n", err)
		return nil, err
	}
	block, err := e.QueryBlockByBlockNum(int64(rawBlockInfo.GetHeader().Number))
	if err != nil {
		log.Printf("Failed to query block by blockHash QueryBlockByBlockNum:%s \n", err)
		return nil, err
	}

	return block, nil
}

// QueryBlockMainInfo Query the main config of channel
func (e *ExploreClient) QueryBlockMainInfo() (*models.BlockMainInfo, error) {

	ledgerInfo, err := e.Option.LedgerClient.QueryInfo()
	if err != nil {
		log.Printf("Failed to Query the main config of channel:%s \n", err)
		return nil, err
	}

	blockNum := ledgerInfo.BCI.Height - 1

	var txNum uint64
	for i := blockNum; i >= 7; i-- {
		rawBlock, err := e.Option.LedgerClient.QueryBlock(uint64(i))
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
func (e *ExploreClient) QueryBlockByBlockNum(num int64) (*models.Block, error) {

	rawBlock, err := e.Option.LedgerClient.QueryBlock(uint64(num))
	if err != nil {
		log.Printf("Failed to query Block info by block's number : %s \n", err)
		return nil, err
	}

	// parse the block body

	var txList []*models.Transaction

	for i := range rawBlock.Data.Data {
		rawEnvelope, err := utils.GetEnvelopeFromBlock(rawBlock.Data.Data[i])
		if err != nil {
			log.Printf("Failed to GetEnvelopeFromBlock: %s \n", err)
			return nil, err
		}

		transaction, err := utils.GetTxFromEnvelopeDeep(rawEnvelope)
		if err != nil {
			log.Printf("Failed to GetTxFromEnvelopeDeep: %s \n", err)
			return nil, err
		}

		for i := range transaction.TransactionActionList {
			transaction.TransactionActionList[i].BlockNum = rawBlock.Header.Number
		}

		txList = append(txList, transaction)
	}

	blockHash := e.GetBlockHash(rawBlock.Header)
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

func (e *ExploreClient) GetBlockHash(blockHeader *common.BlockHeader) []byte {

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

func (e *ExploreClient) QueryTxByTxId(txId string) (*models.Transaction, error) {
	rawTx, err := e.Option.LedgerClient.QueryTransaction(fab.TransactionID(txId))
	if err != nil {
		log.Printf("Failed to QueryTxByTxId rawTx: %s \n", err)
		return nil, err
	}

	transaction, err := utils.GetTxFromEnvelopeDeep(rawTx.TransactionEnvelope)
	if err != nil {
		log.Printf("Failed to QueryTxByTxId transaction: %s \n", err)
		return nil, err
	}

	block, err := e.Option.LedgerClient.QueryBlockByTxID(fab.TransactionID(txId))
	if err != nil {
		log.Printf("Failed to QueryTxByTxId block: %s \n", err)
		return nil, err
	}

	for i := range transaction.TransactionActionList {
		transaction.TransactionActionList[i].BlockNum = block.Header.Number
	}

	return transaction, nil
}

func (e *ExploreClient) QueryTxByTxIdJsonStr(txId string) (string, error) {
	transaction, err := e.QueryTxByTxId(txId)
	if err != nil {
		log.Printf("Failed to QueryTxByTxIdJsonStr transaction: %s \n", err)
		return "", err
	}

	jsonStr, err := json.Marshal(transaction)
	return string(jsonStr), err
}

// OperateLedgerTest Test for operate ledger
func (e *ExploreClient) OperateLedgerTest() {

	bci, err := e.Option.LedgerClient.QueryInfo()
	if err != nil {
		fmt.Printf("failed to query for blockchain info: %s\n", err)
	}

	if bci != nil {
		fmt.Println("Retrieved ledger info")
	}
	fmt.Println(bci.BCI.Height)
	fmt.Println(bci.BCI.PreviousBlockHash)

	rawBlock, err := e.Option.LedgerClient.QueryBlock(uint64(4))
	fmt.Println(rawBlock.GetHeader().GetDataHash())

}

// QueryChannels Query channels
func (e *ExploreClient) QueryChannels() []string {

	configBackend, err := e.Option.MainSDK.Config()
	if err != nil {
		log.Printf("Failed to get config backend from SDK: %s", err)
	}

	targets, err := orgTargetPeers([]string{localConfig.Org1}, configBackend)
	if err != nil {

		log.Printf("Creating peers failed: %s", err)
	}

	channelQueryResponse, err := e.Option.ResMgmtClient.QueryChannels(
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
func (e *ExploreClient) QueryChannelInfo() (*models.ChannelInfo, error) {

	ledgerInfo, err := e.Option.LedgerClient.QueryInfo()
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

func (e *ExploreClient) WalletTest() {
	wallet, err := gateway.NewFileSystemWallet("walletTest")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(wallet)
	fmt.Println(wallet.Exists("User"))

}

func (e *ExploreClient) QueryInfoByChaincode(uuid string) (*models.ChaincodeInfo, error) {

	var args [][]byte
	args = append(args, []byte(uuid))

	req := channel.Request{
		ChaincodeID: localConfig.ChaincodeID,
		Fcn:         "get",
		Args:        args,
	}

	resp, err := e.Option.ChannelClient.Query(req)
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

func (e *ExploreClient) InvokeInfoByChaincode(data string) (*models.ChaincodeInfo, error) {

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

	resp, err := e.Option.ChannelClient.Execute(req)
	if err != nil {
		log.Printf("Failed to invokeInfoByChaincode:%s\n", err)
		return nil, err
	}

	transaction, err := e.QueryTxByTxId(string(resp.TransactionID))
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
