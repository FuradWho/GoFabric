package controllers

import (
	"github.com/kataras/iris/v12"
	"gofabric/services"

	"github.com/iris-contrib/swagger/v12"
	"github.com/iris-contrib/swagger/v12/swaggerFiles"

	_ "gofabric/docs"
)

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
func StartIris() {
	app := iris.New()
	app.Use(Cors)

	config := &swagger.Config{
		URL:         "http://33p67e8007.qicp.vip/swagger/doc.json",
		DeepLinking: true,
	}

	app.Get("/swagger/{any:path}", swagger.CustomWrapHandler(config, swaggerFiles.Handler))

	testApi := app.Party("/")
	{
		testApi.Get("/test", services.Test)
		testApi.Get("/LifeCycleChaincodeTest", services.LifeCycleChaincodeTest)
		testApi.Get("/AuthenticateUser", services.AuthenticateUser)

	}

	// users API operate
	usersApi := app.Party("/user")
	{
		usersApi.Post("/CreateUser", services.CreateUser)
	}

	channelApi := app.Party("/channel")
	{
		channelApi.Post("/CreateChannel", services.CreateChannel)
		channelApi.Post("/JoinChannel", services.JoinChannel)
		channelApi.Get("/GetOrgTargetPeers", services.GetOrgTargetPeers)
		channelApi.Get("/GetNetworkConfig", services.GetNetworkConfig)

	}

	ccApi := app.Party("/cc")
	{
		ccApi.Post("/CreateCC", services.CreateCC)
		ccApi.Post("/InstallCC", services.InstallCC)
		ccApi.Post("/QueryInstalled", services.QueryInstalled)
		ccApi.Post("/ApproveCC", services.ApproveCC)
		ccApi.Post("/QueryApprovedCC", services.QueryApprovedCC)
		ccApi.Post("/CheckCCCommitReadiness", services.CheckCCCommitReadiness)
		ccApi.Post("/RequestInstallCCByOther", services.RequestInstallCCByOther)
		ccApi.Post("/RequestApproveCCByOther", services.RequestApproveCCByOther)
		ccApi.Post("/CommitCC", services.CommitCC)
	}

	app.Listen(":9099")

}

// Cors Resolve the CORS
func Cors(ctx iris.Context) {

	ctx.Header("Access-Control-Allow-Origin", "*")
	if ctx.Request().Method == "OPTIONS" {
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")
		ctx.StatusCode(204)
		return
	}
	ctx.Next()
}
