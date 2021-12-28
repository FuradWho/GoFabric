package controllers

import (
	"github.com/kataras/iris/v12"
	"gofabric/service/baas"
	"gofabric/service/explore"

	"github.com/iris-contrib/swagger/v12"
	"github.com/iris-contrib/swagger/v12/swaggerFiles"

	_ "gofabric/docs"
)

type IrisClient struct {
	baasService    baas.BaasService
	exploreService explore.ExploreService
}

// @title GO Fabric 对于Fabric网络的操作
// @version 1.0
// @description go sdk for Fabric

// @contact.name FuradWho
// @contact.email liu1337543811@gmail.com

// @license.name Fabric 2.3.3
// @license.url https://hyperledger-fabric.readthedocs.io/zh_CN/release-2.2/who_we_are.html

// StartIris
// @host localhost:9099
// @BasePath /
func (i *IrisClient) StartIris(baasService baas.BaasService, exploreService explore.ExploreService) {

	i.baasService = baasService
	i.exploreService = exploreService

	app := iris.New()
	app.Use(cors)

	config := &swagger.Config{
		URL:         "http://33p67e8007.qicp.vip/swagger/doc.json",
		DeepLinking: true,
	}

	app.Get("/swagger/{any:path}", swagger.CustomWrapHandler(config, swaggerFiles.Handler))

	baasTestApi := app.Party("/")
	{
		baasTestApi.Get("/bassTest", i.baasService.Test)
		baasTestApi.Get("/lifeCycleChaincodeTest", i.baasService.LifeCycleChaincodeTest)
		baasTestApi.Get("/authenticateUser", i.baasService.AuthenticateUser)

	}

	// users API operate
	baasUsersApi := app.Party("/baasUser")
	{
		baasUsersApi.Post("/createUser", i.baasService.CreateUser)
	}

	baasChannelApi := app.Party("/baasChannel")
	{
		baasChannelApi.Post("/createChannel", i.baasService.CreateChannel)
		baasChannelApi.Post("/joinChannel", i.baasService.JoinChannel)
		baasChannelApi.Get("/getOrgTargetPeers", i.baasService.GetOrgTargetPeers)
		baasChannelApi.Get("/getNetworkConfig", i.baasService.GetNetworkConfig)

	}

	exploreChannelApi := app.Party("/exploreChannel")
	{
		exploreChannelApi.Get("/queryChannelInfo", i.exploreService.QueryChannelInfo)
	}

	baasCcApi := app.Party("/baasCc")
	{
		baasCcApi.Post("/createCC", i.baasService.CreateCC)
		baasCcApi.Post("/installCC", i.baasService.InstallCC)
		baasCcApi.Post("/queryInstalled", i.baasService.QueryInstalled)
		baasCcApi.Post("/approveCC", i.baasService.ApproveCC)
		baasCcApi.Post("/queryApprovedCC", i.baasService.QueryApprovedCC)
		baasCcApi.Post("/checkCCCommitReadiness", i.baasService.CheckCCCommitReadiness)
		baasCcApi.Post("/requestInstallCCByOther", i.baasService.RequestInstallCCByOther)
		baasCcApi.Post("/requestApproveCCByOther", i.baasService.RequestApproveCCByOther)
		baasCcApi.Post("/commitCC", i.baasService.CommitCC)

	}

	exploreCcApi := app.Party("/exploreCc")
	{
		exploreCcApi.Get("/queryInstalledCC", i.exploreService.QueryInstalledCC)
		exploreCcApi.Post("/invokeInfoByChaincode", i.exploreService.InvokeInfoByChaincode)
		exploreCcApi.Get("/queryInfoByChaincode", i.exploreService.QueryInfoByChaincode)
	}

	exploreBlocksApi := app.Party("/exploreBlocks")
	{
		exploreBlocksApi.Get("/queryLastesBlocksInfo", i.exploreService.GetLastesBlocksInfo)
		exploreBlocksApi.Get("/queryBlockByBlockNum", i.exploreService.QueryBlockByBlockNum)
		exploreBlocksApi.Get("/queryAllBlocksInfo", i.exploreService.QueryAllBlocksInfo)
		exploreBlocksApi.Get("/queryBlockInfoByHash", i.exploreService.QueryBlockInfoByHash)
		exploreBlocksApi.Get("/queryBlockMainInfo", i.exploreService.QueryBlockMainInfo)

	}

	exploreTxsApi := app.Party("/exploreTxs")
	{
		exploreTxsApi.Get("/queryTxByTxId", i.exploreService.QueryTxByTxId)
		exploreTxsApi.Get("/queryTxByTxIdJsonStr", i.exploreService.QueryTxByTxIdJsonStr)
	}

	app.Listen(":9099")

}

// Cors Resolve the CORS
func cors(ctx iris.Context) {

	ctx.Header("Access-Control-Allow-Origin", "*")
	if ctx.Request().Method == "OPTIONS" {
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")
		ctx.StatusCode(204)
		return
	}
	ctx.Next()
}
