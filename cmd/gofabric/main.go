package main

import (
	"gofabric/common"
	localConfig "gofabric/configs"
	"gofabric/controllers"
	"gofabric/services"
	_ "gofabric/third_party/logger"
)

func main() {

	common.InitChainExploreService(localConfig.OrgGoConfig, localConfig.OrgGo, localConfig.Admin, localConfig.User)
	services.NewFabricClient()
	controllers.StartIris()

}
