package main

import (
	"gofabric/common"
)

func main() {

	common.NewFabricClient()

	// common.CreateChannel()

	common.InstallChaincode()
}