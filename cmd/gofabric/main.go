package main

import (
	log "github.com/sirupsen/logrus"
	"gofabric/common"
	"gofabric/controllers"
	"gofabric/pkg/baas"
	"gofabric/pkg/explore"
	_ "gofabric/third_party/logger"
)

func main() {

	newFabricOption, err := common.NewFabricOption(func(option *common.Option) {
	})
	if err != nil {
		log.Errorln(err)
	}

	baasClient := baas.BaasClient{Foo: newFabricOption}
	baasService := new(baas.BaasService)
	baasService.InitBaasService(&baasClient)

	exploreClient := explore.ExploreClient{Foo: newFabricOption}
	exploreService := new(explore.ExploreService)
	exploreService.InitExploreService(&exploreClient)

	irisClient := new(controllers.IrisClient)
	irisClient.StartIris(baasService, exploreService)

}
