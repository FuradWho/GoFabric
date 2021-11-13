package models

import (
	"encoding/hex"
	"fmt"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config/lookup"
		//	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"log"
	"os"
	"strings"

	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	lcpackager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/lifecycle"
)

const (
	Admin = "User2"
)

type FabricClient struct {
	ConnectionFile []byte
	Orgs []string
	ChannelId string
	GoPath string

	userName string
	userOrg string

	resmgmtClients []*resmgmt.Client
	sdk *fabsdk.FabricSDK
	retry resmgmt.RequestOption
}

func NewFabricClient(connectionFile []byte, channelId string, orgs []string) *FabricClient {
	fabric := &FabricClient{
		ConnectionFile: connectionFile,
		ChannelId:      channelId,
		Orgs:           orgs,
		GoPath:         os.Getenv("GOPATH"),
	}
	return fabric
}


func (f *FabricClient) Setup() error {


	sdk, err := fabsdk.New(config.FromRaw(f.ConnectionFile,"yaml"))
	if err != nil {
		log.Printf("Failed to Create SDK : %s \n",err)
		return err
	}
	f.sdk = sdk
	resmgmtClients := make([]*resmgmt.Client,0)
	for _, v := range f.Orgs {
		resmgmtClient , err := resmgmt.New(sdk.Context(fabsdk.WithUser(Admin),fabsdk.WithOrg(v)))
		if err != nil {
			log.Printf("Failed to create channel management client : %s \n",err)
		}
		resmgmtClients = append(resmgmtClients,resmgmtClient)
	}

	f.resmgmtClients = resmgmtClients
	f.retry = resmgmt.WithRetry(retry.DefaultResMgmtOpts)
	return nil
}

func (f *FabricClient) Close()  {
	if f.sdk != nil {
		f.sdk.Close()
	}
}

func (f *FabricClient) SetUser(userName , userOrg string)  {
	f.userName = userName
	f.userOrg = userOrg
}


func GetKeyFile(id msp.SigningIdentity) (string,string){

	priFile := hex.EncodeToString(id.PrivateKey().SKI())+"_sk"
	pubFile := id.Identifier().ID + "@" +id.Identifier().MSPID + "-cert.pem"
	return priFile,pubFile
}


func (f *FabricClient) CreateChannel(channelTx string) error {

	mspClient , err := mspclient.New(f.sdk.Context(),mspclient.WithOrg(f.Orgs[0]))
	if err != nil {
		log.Printf("Failed to new mspClient : %s \n",err)
		return err
	}

	adminidentity, err := mspClient.GetSigningIdentity(Admin)
	if err != nil {
		log.Printf("Failed to get signIdentity : %s \n",err)
		return err
	}

	req := resmgmt.SaveChannelRequest{
		ChannelID: f.ChannelId,
		ChannelConfigPath: channelTx,
		SigningIdentities: []msp.SigningIdentity{adminidentity},
	}

	txId , err := f.resmgmtClients[0].SaveChannel(req)
	if err != nil {
		log.Printf("Failed to save channel : %s \n",err)
		return err
	}

	fmt.Printf("txId : %s \n",txId)
	return nil
}


func (f *FabricClient) UpdateChannel(anchorsTx []string) error {

	for i, c := range f.resmgmtClients {

		mspClient, err := mspclient.New(f.sdk.Context(), mspclient.WithOrg(f.Orgs[i]))
		if err != nil {
			log.Printf("Failed to UpdateChannel : %s \n",err)
			return err
		}
		adminIdentity, err := mspClient.GetSigningIdentity(Admin)
		if err != nil {
			log.Printf("Failed to UpdateChannel : %s \n",err)
			return err
		}
		req := resmgmt.SaveChannelRequest{
			ChannelID:         f.ChannelId,
			ChannelConfigPath: anchorsTx[i],
			SigningIdentities: []msp.SigningIdentity{adminIdentity},
		}
		txId, err := c.SaveChannel(req, f.retry)
		if err != nil {
			log.Printf("Failed to UpdateChannel : %s \n",err)
			return err
		}
		log.Printf("Failed to UpdateChannel : %s \n",txId)
	}

	return nil
}

func (f *FabricClient) JoinChannel() error {

	for i, c := range f.resmgmtClients {
		err := c.JoinChannel(f.ChannelId, f.retry)
		if err != nil && !strings.Contains(err.Error(), "LedgerID already exists") {
			log.Printf("Org peers failed to JoinChannel: %s \n", err)
			return err
		}
		log.Printf("%s join channel", f.Orgs[i])

	}
	return nil

}

func (f *FabricClient) InstallChaincode(chaincodeId, chaincodePath, version string) error {
	//ccPkg, err := gopackager.NewCCPackage(chaincodePath, f.GoPath)
	//if err != nil {
	//	log.Printf("Org peers failed to InstallChaincode: %s \n", err)
	//	return err
	//}

	//req := resmgmt.InstallCCRequest{
	//	Name:    chaincodeId,
	//	Path:    chaincodePath,
	//	Version: version,
	//	Package: ccPkg,
	//}

	//for _, c := range f.resmgmtClients {
		//res, err := c.InstallCC(req, f.retry)
		//if err != nil {
		//	log.Printf("Org peers failed to InstallChaincode: %s \n", err)
		//	return err
		//}
		//log.Printf("%s \n", res)
		configBackend , _ := f.sdk.Config()
		targets, err := orgTargetPeers([]string{"org1"},configBackend )
		if err != nil {
			log.Printf("Failed to get targets:%s \n",err)
		}
		peer := targets[0]

		installedCC, err := f.resmgmtClients[0].QueryInstalledChaincodes(resmgmt.WithTargetEndpoints(peer))
		if err != nil {
			return err
		}
		fmt.Printf("%v", installedCC)

	//}

	return nil
}


func orgTargetPeers(orgs []string, configBackend ...core.ConfigBackend) ([]string, error) {

	networkConfig := fab.NetworkConfig{}

	err := lookup.New(configBackend...).UnmarshalKey("organizations", &networkConfig.Organizations)
	if err != nil {
		return nil, err
	}

	var peers []string
	for _, org := range orgs {
		orgConfig, ok := networkConfig.Organizations[strings.ToLower(org)]
		if !ok {
			continue
		}
		peers = append(peers, orgConfig.Peers...)
	}
	return peers, nil
}





//func (f *FabricClient) InstantiateChaincode(chaincodeId, chaincodePath, version string, policy string, args [][]byte) (string, error) {
//
//	//"OR ('Org1MSP.member','Org2MSP.member')"
//	ccPolicy, err := cauthdsl.FromString(policy)
//	if err != nil {
//		return "", err
//	}
//	resp, err := f.resmgmtClients[0].InstantiateCC(
//		f.ChannelId,
//		resmgmt.InstantiateCCRequest{
//			Name:    chaincodeId,
//			Path:    chaincodePath,
//			Version: version,
//			Args:    args,
//			Policy:  ccPolicy,
//		},
//		f.retry,
//	)
//
//	return string(resp.TransactionID), nil
//}




