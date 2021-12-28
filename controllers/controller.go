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
	baasService    *baas.BaasService
	exploreService *explore.ExploreService
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
func (i *IrisClient) StartIris(baasService *baas.BaasService, exploreService *explore.ExploreService) {

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
		baasTestApi.Get("/LifeCycleChaincodeTest", i.baasService.LifeCycleChaincodeTest)
		baasTestApi.Get("/AuthenticateUser", i.baasService.AuthenticateUser)

	}

	// users API operate
	baasUsersApi := app.Party("/baasUser")
	{
		baasUsersApi.Post("/baasCreateUser", i.baasService.CreateUser)
	}

	baasChannelApi := app.Party("/baasChannel")
	{
		baasChannelApi.Post("/baasCreateChannel", i.baasService.CreateChannel)
		baasChannelApi.Post("/baasJoinChannel", i.baasService.JoinChannel)
		baasChannelApi.Get("/baasGetOrgTargetPeers", i.baasService.GetOrgTargetPeers)
		baasChannelApi.Get("/baasGetNetworkConfig", i.baasService.GetNetworkConfig)

	}

	exploreChannelApi := app.Party("/exploreChannel")
	{
		exploreChannelApi.Get("/exploreQueryChannelInfo", i.exploreService.QueryChannelInfo)
	}

	baasCcApi := app.Party("/baasCc")
	{
		baasCcApi.Post("/baasCreateCC", i.baasService.CreateCC)
		baasCcApi.Post("/baasInstallCC", i.baasService.InstallCC)
		baasCcApi.Post("/baasQueryInstalled", i.baasService.QueryInstalled)
		baasCcApi.Post("/baasApproveCC", i.baasService.ApproveCC)
		baasCcApi.Post("/baasQueryApprovedCC", i.baasService.QueryApprovedCC)
		baasCcApi.Post("/baasCheckCCCommitReadiness", i.baasService.CheckCCCommitReadiness)
		baasCcApi.Post("/baasRequestInstallCCByOther", i.baasService.RequestInstallCCByOther)
		baasCcApi.Post("/baasRequestApproveCCByOther", i.baasService.RequestApproveCCByOther)
		baasCcApi.Post("/baasCommitCC", i.baasService.CommitCC)

	}

	exploreCcApi := app.Party("/exploreCc")
	{
		exploreCcApi.Get("/exploreQueryInstalledCC", i.exploreService.QueryInstalledCC)
		exploreCcApi.Post("/exploreInvokeInfoByChaincode", i.exploreService.InvokeInfoByChaincode)
		exploreCcApi.Get("/exploreQueryInfoByChaincode", i.exploreService.QueryInfoByChaincode)
	}

	exploreBlocksApi := app.Party("/exploreBlocks")
	{
		exploreBlocksApi.Get("/exploreQueryLastesBlocksInfo", i.exploreService.GetLastesBlocksInfo)
		exploreBlocksApi.Get("/exploreQueryBlockByBlockNum", i.exploreService.QueryBlockByBlockNum)
		exploreBlocksApi.Get("/exploreQueryAllBlocksInfo", i.exploreService.QueryAllBlocksInfo)
		exploreBlocksApi.Get("/exploreQueryBlockInfoByHash", i.exploreService.QueryBlockInfoByHash)
		exploreBlocksApi.Get("/exploreQueryBlockMainInfo", i.exploreService.QueryBlockMainInfo)

	}

	exploreTxsApi := app.Party("/exploreTxs")
	{
		exploreTxsApi.Get("/exploreQueryTxByTxId", i.exploreService.QueryTxByTxId)
		exploreTxsApi.Get("/exploreQueryTxByTxIdJsonStr", i.exploreService.QueryTxByTxIdJsonStr)
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
