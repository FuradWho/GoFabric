package explore

import (
	"encoding/hex"
	"github.com/kataras/iris/v12/context"
	log "github.com/sirupsen/logrus"
	"gofabric/pkg/explore"
	"strconv"
)

type ExploreService struct {
	exploreClient *explore.ExploreClient
}

func (e *ExploreService) InitExploreService(client *explore.ExploreClient) {
	e.exploreClient = client
}

func (e *ExploreService) GetLastesBlocksInfo(context context.Context) {
	blocks, err := e.exploreClient.QueryLastesBlocksInfo()
	if err != nil {
		log.Errorln(err)
	}
	context.JSON(blocks)
}

func (e *ExploreService) QueryAllBlocksInfo(context context.Context) {
	blocks, err := e.exploreClient.QueryAllBlocksInfo()
	if err != nil {
		log.Errorln(err)
	}
	context.JSON(blocks)
}

func (e *ExploreService) QueryTxByTxId(context context.Context) {

	txId := context.URLParam("txId")
	if txId == "" {
		context.JSON("fail")
	} else {
		transactions, err := e.exploreClient.QueryTxByTxId(txId)
		if err != nil {
			log.Errorln(err)
		}

		context.JSON(transactions)
	}
}

func (e *ExploreService) QueryTxByTxIdJsonStr(context context.Context) {

	txId := context.URLParam("txId")
	if txId == "" {
		context.JSON("fail")
	} else {
		transactions, err := e.exploreClient.QueryTxByTxId(txId)
		if err != nil {
			log.Errorln(err)
		}

		context.JSON(transactions)
	}
}

func (e *ExploreService) QueryBlockByBlockNum(context context.Context) {
	blockNum := context.URLParam("blockNum")
	if blockNum == "" {
		context.JSON("fail")
	} else {

		num, _ := strconv.ParseInt(blockNum, 10, 64)
		transactions, err := e.exploreClient.QueryBlockByBlockNum(num)
		if err != nil {
			log.Errorln(err)
		}

		context.JSON(transactions)
	}
}

func (e *ExploreService) QueryBlockInfoByHash(context context.Context) {
	blockHash := context.URLParam("blockHash")
	if blockHash == "" {
		context.JSON("fail")
	} else {
		byteBlockHash, err := hex.DecodeString(blockHash)
		if err != nil {
			log.Errorln(err)
		}
		blockInfo, err := e.exploreClient.QueryBlockInfoByHash(byteBlockHash)
		if err != nil {
			log.Errorln(err)
		}

		context.JSON(blockInfo)
	}
}

func (e *ExploreService) QueryBlockMainInfo(context context.Context) {
	blocks, err := e.exploreClient.QueryBlockMainInfo()
	if err != nil {
		log.Errorln(err)
	}
	context.JSON(blocks)
}

func (e *ExploreService) QueryInstalledCC(context context.Context) {

	chaincodeInfo, err := e.exploreClient.QueryInstalledCC()
	if err != nil {
		log.Errorln(err)
	}
	context.JSON(chaincodeInfo)
}

func (e *ExploreService) QueryChannelInfo(context context.Context) {
	channelInfo, err := e.exploreClient.QueryChannelInfo()
	if err != nil {
		log.Errorln(err)
	}
	context.JSON(channelInfo)
}

func (e *ExploreService) InvokeInfoByChaincode(context context.Context) {

	data := context.PostValue("data")
	if data == "" {
		context.JSON("fail")
	} else {
		chaincodeInfo, err := e.exploreClient.InvokeInfoByChaincode(data)
		if err != nil {
			log.Errorln(err)
		}
		context.JSON(chaincodeInfo)
	}

}

func (e *ExploreService) QueryInfoByChaincode(context context.Context) {

	uuid := context.URLParam("uuid")

	if uuid == "" {
		context.JSON("fail")
	} else {
		chaincodeInfo, err := e.exploreClient.QueryInfoByChaincode(uuid)
		if err != nil {
			log.Errorln(err)
		}
		context.JSON(chaincodeInfo)
	}

}
