package models

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

type Client struct {
	// Fabric 网络信息
	ConfigPath string
	OrgName    string
	OrgAdmin   string
	OrgUser    string

	//sdk 客户端
	SDK *fabsdk.FabricSDK
	Rc  *resmgmt.Client
	Cc  *channel.Client

	// 链码信息
	ChannelID     string //通道名
	ChaincodeID   string //链码ID或者名称
	ChaincodePath string //链码路径
	GoPath        string // GOPATH路径
}
