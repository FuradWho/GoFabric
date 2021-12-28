package main

import (
	"github.com/sirupsen/logrus"
	"gofabric/common"
	localConfig "gofabric/configs"
	"gofabric/controllers"
	"gofabric/services"
	"os"
	"time"
)

var log = logrus.New()

func main() {
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		TimestampFormat: time.Now().String(),
	})
	log.Out = os.Stdout
	common.InitChainExploreService(localConfig.OrgGoConfig, localConfig.OrgGo, localConfig.Admin, localConfig.User)
	services.NewFabricClient()
	controllers.StartIris()

}
