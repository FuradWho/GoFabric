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

	services.NewFabricClient()
	controllers.StartIris()

}
