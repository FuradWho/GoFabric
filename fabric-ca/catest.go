package fabric_ca

import (
	"encoding/hex"
	"fmt"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"log"
	"os"
	"strings"
)

const (
	Admin = "admin-org2"
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



func InitCaClient()  {
	sdk, err := fabsdk.New(config.FromFile("./connect-config/User10-config.yaml"))
	if err != nil {
		fmt.Println(err)
	}

	ctx := sdk.Context()

	client, err := mspclient.New(ctx)
	if err != nil {
		fmt.Println(err)
	}

	//resp, err := client.GetCAInfo()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(resp.CAName)
	//
	//affiliations, err := client.GetAllIdentities()
	//if err != nil {
	//	return
	//}

	//req := &msp.RegistrationRequest{
	//	Name: "User8",
	//	Type: "client",
	//	CAName: "",
	//	Secret: "123456",
	//}
	////
	//register, err := client.Register(req)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(register)
	//
	//err = client.Enroll(register)
	//if err != nil {
	//	fmt.Println(err)
	//}

	info, err := client.GetCAInfo()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(info.CAName)
	fmt.Println(info.Version)



	//identity, err := client.CreateIdentity(&msp.IdentityRequest{ID: "123", Affiliation: "org2",
	//	Attributes: []msp.Attribute{{Name: "attName1", Value: "attValue1"}}})
	//if err != nil {
	//	fmt.Printf("Create identity return error %s\n", err)
	//	return
	//}
	//fmt.Printf("identity '%s' created\n", identity.ID)


	//identity, err := client.GetIdentity("User4")
	//if err != nil {
	//	fmt.Printf("Get Identity : %s \n",err)
	//}
	//fmt.Println(identity.ID)
	//

	//req := &msp.RegistrationRequest{
	//	Name: "User10",
	//	Type: "client",
	//	CAName: "",
	//	Secret: "123456",
	//}
	////
	//register, err := client.Register(req)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//fmt.Println(register)
	//
	//err = client.Enroll("User10",msp.WithSecret("123456"))
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	signingIdentity, err := client.GetSigningIdentity("admin-org2")
	if err != nil {
		fmt.Printf("GetSigningIdentity : %s \n",err)
	}
	fmt.Println(signingIdentity.PrivateKey().SKI())

	fmt.Println(GetKeyFile(signingIdentity))

	//identity, err := client.GetSigningIdentity("User5")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//key := identity.Identifier().ID
	//fmt.Println(key)

	

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

	affiliations, err := mspClient.GetAllIdentities()
	if err != nil {
		return nil
	}

	for _ , info := range affiliations{
		fmt.Println(info.ID)
		fmt.Println(info.Type)
		fmt.Println(info.Attributes)
		fmt.Println(info.CAName)

		fmt.Println("----------------------")
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

	txId , err := f.resmgmtClients[0].SaveChannel(req,f.retry)
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
	ccPkg, err := gopackager.NewCCPackage(chaincodePath, f.GoPath)
	if err != nil {
		log.Printf("Org peers failed to InstallChaincode: %s \n", err)
		return err
	}

	req := resmgmt.InstallCCRequest{
		Name:    chaincodeId,
		Path:    chaincodePath,
		Version: version,
		Package: ccPkg,
	}

	for _, c := range f.resmgmtClients {
		res, err := c.InstallCC(req, f.retry)
		if err != nil {
			log.Printf("Org peers failed to InstallChaincode: %s \n", err)
			return err
		}
		log.Printf("%s \n", res)
	}

	return nil
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
