package controllers

import (
	"github.com/kataras/iris/v12"
	"gofabric/pkg/baas"
	"gofabric/pkg/explore"

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

	testApi := app.Party("/")
	{
		testApi.Get("/test", i.baasService.Test)
		testApi.Get("/LifeCycleChaincodeTest", i.baasService.LifeCycleChaincodeTest)
		testApi.Get("/AuthenticateUser", i.baasService.AuthenticateUser)

	}

	// users API operate
	usersApi := app.Party("/user")
	{
		usersApi.Post("/CreateUser", i.baasService.CreateUser)
	}

	channelApi := app.Party("/channel")
	{
		channelApi.Post("/CreateChannel", i.baasService.CreateChannel)
		channelApi.Post("/JoinChannel", i.baasService.JoinChannel)
		channelApi.Get("/GetOrgTargetPeers", i.baasService.GetOrgTargetPeers)
		channelApi.Get("/GetNetworkConfig", i.baasService.GetNetworkConfig)
		channelApi.Get("/QueryChannelInfo", i.exploreService.QueryChannelInfo)

	}

	ccApi := app.Party("/cc")
	{
		ccApi.Post("/CreateCC", i.baasService.CreateCC)
		ccApi.Post("/InstallCC", i.baasService.InstallCC)
		ccApi.Post("/QueryInstalled", i.baasService.QueryInstalled)
		ccApi.Post("/ApproveCC", i.baasService.ApproveCC)
		ccApi.Post("/QueryApprovedCC", i.baasService.QueryApprovedCC)
		ccApi.Post("/CheckCCCommitReadiness", i.baasService.CheckCCCommitReadiness)
		ccApi.Post("/RequestInstallCCByOther", i.baasService.RequestInstallCCByOther)
		ccApi.Post("/RequestApproveCCByOther", i.baasService.RequestApproveCCByOther)
		ccApi.Post("/CommitCC", i.baasService.CommitCC)
		ccApi.Get("/QueryInstalledCC", i.exploreService.QueryInstalledCC)
		ccApi.Post("/InvokeInfoByChaincode", i.exploreService.InvokeInfoByChaincode)
		ccApi.Get("/QueryInfoByChaincode", i.exploreService.QueryInfoByChaincode)
	}

	blocksApi := app.Party("/blocks")
	{
		blocksApi.Get("/QueryLastesBlocksInfo", i.exploreService.GetLastesBlocksInfo)
		blocksApi.Get("/QueryBlockByBlockNum", i.exploreService.QueryBlockByBlockNum)
		blocksApi.Get("/QueryAllBlocksInfo", i.exploreService.QueryAllBlocksInfo)
		blocksApi.Get("/QueryBlockInfoByHash", i.exploreService.QueryBlockInfoByHash)
		blocksApi.Get("/QueryBlockMainInfo", i.exploreService.QueryBlockMainInfo)

	}

	txsApi := app.Party("/txs")
	{
		txsApi.Get("/QueryTxByTxId", i.exploreService.QueryTxByTxId)
		txsApi.Get("/QueryTxByTxIdJsonStr", i.exploreService.QueryTxByTxIdJsonStr)
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
