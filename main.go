package main

import (
	"github.com/sirupsen/logrus"
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

	// common.NewFabricClient()
	//common.CreateChannel()
	// common.InstallChaincode()
	//label ,ccpkg := common.PackageCC()
	//fmt.Println(label)
	//common.InstallCC(label,ccpkg)
	// common.GetInstalledCCPackage("newcc0:e09a3747940e77d2bc96d7555672eb42f37a6d3af9a2e69b53454a133248e074")
	//common.CreateUser("User2","User2","admin","org1")
	//common.QueryInstalled("newcc0","newcc0:e09a3747940e77d2bc96d7555672eb42f37a6d3af9a2e69b53454a133248e074")
	//common.ApproveCC("newcc0:e09a3747940e77d2bc96d7555672eb42f37a6d3af9a2e69b53454a133248e074")
	// common.CreateCC()
	// common.QueryLedger()

	services.NewFabricClient()
	controllers.StartIris()

}

// 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
