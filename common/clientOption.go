package common

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	contextApi "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	log "github.com/sirupsen/logrus"
	localConfig "gofabric/configs"
	"io/ioutil"
	"os"
)

type Foo struct {
	Option Option
}
type Option struct {

	// Fabric 网络信息
	ConfigPath     string
	OrgName        string
	OrgUser        string
	ConnectionFile []byte
	Orgs           []string
	UserName       string
	UserOrg        string

	// 链码信息
	ChannelID     string //通道名
	ChaincodeID   string //链码ID或者名称
	ChaincodePath string //链码路径
	GoPath        string // GOPATH路径

	Ctx            contextApi.ClientProvider
	MainSDK        fabsdk.FabricSDK
	LedgerClient   ledger.Client
	ChannelClient  channel.Client
	ResMgmtClient  resmgmt.Client
	ResMgmtClients map[string]resmgmt.Client
	Retry          resmgmt.RequestOption
}

type ModOption func(option *Option)

func NewFabricOption(modOption ModOption) (Foo, error) {
	log.Infoln("init the option.")

	// load the local config
	option := Option{
		ConfigPath:    localConfig.ConnectConfigDir,
		OrgName:       localConfig.Org1,
		OrgUser:       localConfig.Admin,
		ChannelID:     localConfig.ChannelID,
		ChaincodeID:   localConfig.ChaincodeID,
		ChaincodePath: localConfig.ChaincodePath,
		GoPath:        os.Getenv("GOPATH"),
		Orgs:          []string{"org1", "org2"},
	}

	// read the sdk config
	connectConfig, _ := ioutil.ReadFile(option.ConfigPath)
	option.ConnectionFile = connectConfig

	// init the fab sdk
	sdk, err := fabsdk.New(config.FromRaw(option.ConnectionFile, "yaml"))
	if err != nil {
		return Foo{}, err
	}
	option.MainSDK = *sdk

	option.Ctx = option.MainSDK.Context()

	// init the resmgmt sdk for orgs
	resMgmtClients := make(map[string]resmgmt.Client)
	for _, v := range option.Orgs {
		resMgmtClient, err := resmgmt.New(sdk.Context(fabsdk.WithUser(option.OrgUser), fabsdk.WithOrg(v)))
		if err != nil {
			return Foo{}, err
		}
		resMgmtClients[v] = *resMgmtClient

		// rand one org for the resMgmtClient to be as the explore client
		option.ResMgmtClient = *resMgmtClient
	}
	option.ResMgmtClients = resMgmtClients

	// init the channel client sdk
	channelContext := option.MainSDK.ChannelContext(option.ChannelID, fabsdk.WithUser(option.OrgUser))
	// init the ledger client
	ledgerClient, err := ledger.New(channelContext)
	if err != nil {
		return Foo{}, err
	}
	option.LedgerClient = *ledgerClient

	// init the channel client
	channelClient, err := channel.New(channelContext)
	if err != nil {
		return Foo{}, err
	}
	option.ChannelClient = *channelClient
	// init the retry
	option.Retry = resmgmt.WithRetry(retry.DefaultResMgmtOpts)

	modOption(&option)

	return Foo{
		Option: option,
	}, nil
}
