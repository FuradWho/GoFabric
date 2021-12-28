package models

import "github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"

type Chaincode struct {
	Label      string                           `json:"label"`
	PackageID  string                           `json:"packageID"`
	References map[string][]resmgmt.CCReference `json:"references"`
}

type ChannelInfo struct {
	Name      string `json:"name"`
	Blocks    int    `json:"blocks"`
	Txs       int    `json:"txs"`
	Timestamp string `json:"timestamp"`
}
