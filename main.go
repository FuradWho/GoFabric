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


	// test for configtxlator

	//get, err := http.Post("http://127.0.0.1:7059/protolator/decode/common.Config")
	//if err != nil {
	//	return
	//}
	//log.Infof("%+v \n",get)

	services.NewFabricClient()
	controllers.StartIris()

}

