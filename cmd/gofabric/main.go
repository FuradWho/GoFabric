package main

import (
	log "github.com/sirupsen/logrus"
	"gofabric/common"
	"gofabric/controllers"
	"gofabric/pkg/baas"
	"gofabric/pkg/explore"
	baas2 "gofabric/service/baas"
	explore2 "gofabric/service/explore"
	_ "gofabric/third_party/logger"
)

func main() {

	newFabricOption, err := common.NewFabricOption(func(option *common.Option) {
	})
	if err != nil {
		log.Errorln(err)
	}

	baasClient := baas.BaasClient{Foo: newFabricOption}
	baasService := new(baas2.BaasService)
	baasService.InitBaasService(baasClient)

	exploreClient := explore.ExploreClient{Foo: newFabricOption}
	exploreService := new(explore2.ExploreService)
	exploreService.InitExploreService(exploreClient)

	irisClient := new(controllers.IrisClient)
	irisClient.StartIris(*baasService, *exploreService)

}
