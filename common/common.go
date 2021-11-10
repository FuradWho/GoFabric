package common

import (
	"fmt"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"gofabric/models"
	"io/ioutil"
	"log"
)



const (
	channelId = "mychannel"
	channelTx = "/usr/local/hyper/test2/configtx/channel-artifacts/mychannel.tx"
	connectConfigDir = "connect-config/channel-connection.yaml"
	chaincodeId = "mycc"
	chaincodePath = "newchaincode/test"
	ccVersion = "0"

)

var fabricClient *models.FabricClient

var orgs = []string{"org1","org2"}

func CreateUser(userName string,password string,userType string,orgName string)  {

	sdk, err := fabsdk.New(config.FromFile("connect-config/channel-connection.yaml"))
	if err != nil {
		fmt.Println(err)
	}

	ctx := sdk.Context()
	client, err := mspclient.New(ctx)
	if err != nil {
		fmt.Println(err)
	}

	info, err := client.GetCAInfo()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(info.CAName)
	fmt.Println(info.Version)

	affiliations, err := client.GetAllIdentities()
	if err != nil {
		log.Printf("%s \n",err)
	}

	for _ , info := range affiliations{
		fmt.Println(info.ID)
		fmt.Println(info.Type)
		fmt.Println(info.Attributes)
		fmt.Println("----------------------")
	}

	a1 := mspclient.Attribute{
		Name: "hf.Registrar.Roles",
		Value:"client,orderer,peer,user",
	}
	a2 := mspclient.Attribute{
		Name: "hf.Registrar.DelegateRoles",
		Value:"client,orderer,peer,user",
	}
	a3 := mspclient.Attribute{
		Name: "hf.Registrar.Attributes",
		Value:"*",
	}
	a4 := mspclient.Attribute{
		Name: "hf.GenCRL",
		Value:"true",
	}
	a5 := mspclient.Attribute{
		Name: "hf.Revoker",
		Value:"true",
	}
	a6 := mspclient.Attribute{
		Name: "hf.AffiliationMgr",
		Value:"true",
	}
	a7 := mspclient.Attribute{
		Name: "hf.IntermediateCA",
		Value:"true",
	}

	var attributes []mspclient.Attribute
	attributes = append(attributes,a1,a2,a3,a4,a5,a6,a7)

	req := &mspclient.RegistrationRequest{
		Name: userName,
		Type: userType,
		CAName: "ca-org2",
		Secret: password,
		Attributes: attributes,
		Affiliation: orgName,
	}
	register, err := client.Register(req)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(register)

	err = client.Enroll(userName,mspclient.WithSecret(password))
	if err != nil {
		fmt.Println(err)
	}

	signingIdentity, err := client.GetSigningIdentity(userName)
	if err != nil {
		fmt.Printf("GetSigningIdentity : %s \n",err)
	}
	fmt.Println(signingIdentity.Identifier().ID)
}

func NewFabricClient()  {
	
	connectConfig, _ := ioutil.ReadFile(connectConfigDir)

	fabricClient =  models.NewFabricClient(connectConfig,channelId,orgs)
	//defer fabricClient.Close()

	err := fabricClient.Setup()
	if err != nil {
		log.Panicf("Failed to setup fabricclient: %s \n",err)
	}
	
}

func CreateChannel()  {
	err := fabricClient.CreateChannel(channelTx)
	if err != nil {
		log.Panicf("Failed to CreateChannel: %s \n",err)
	}
	err = fabricClient.JoinChannel()
	if err != nil {
		log.Panicf("Failed to JoinChannel: %s \n",err)
	}
}
func InstallChaincode()  {

	err := fabricClient.InstallChaincode(chaincodeId,chaincodePath,ccVersion)
	if err != nil {
		log.Panicf("Failed to InstallChaincode: %s \n",err)
	}

}






/*
	// ordererDomain  := "orderer.example.com"
	orgs :=[]string{"org2"}
	channelId := "mychannel"
	connectConfig,_ := ioutil.ReadFile("./connect-config/client-network.yaml")

	// chaincodeId := "mycc"
	// chaincodePath := "/usr/local/hyper/fabric-ca/chaincode/newchaincode"

	fabric := fabric_ca.NewFabricClient(connectConfig, channelId ,orgs)
	defer fabric.Close()
	fabric.Setup()
	//创建channel
	//fabric.CreateChannel(channelTx)
	//加入channel
	//fabric.JoinChannel()
*/


//sdkClient , err := fabsdk.New(config.FromFile("connect-config/channel-connection.yaml"))
//if err != nil {
//	log.Panicf("Failed to create a sdkClient :%s \n",err)
//
//}
//resourceProvider := sdkClient.Context(fabsdk.WithUser("User2"),fabsdk.WithOrg("org2"))
//
//resourceClient , err := resmgmt.New(resourceProvider)
//if err != nil {
//	log.Panicf("Failed to create a resourceClient : %s \n",err)
//}


//mspClient , err := mspclient.New(sdkClient.Context(),mspclient.WithOrg("org1"))
//if err != nil {
//	log.Printf("Failed to new mspClient : %s \n",err)
//}

//adminidentity, err := mspClient.GetSigningIdentity("User2")
//if err != nil {
//	log.Printf("Failed to get signIdentity : %s \n",err)
//}

//channelTx := "/usr/local/hyper/test2/configtx/channel-artifacts/mychannel.tx"
// channelId := "mychannel"
//
//req := resmgmt.SaveChannelRequest{
//	ChannelID: channelId,
//	ChannelConfigPath: channelTx,
//	SigningIdentities: []msp.SigningIdentity{adminidentity},
//}
//
//txId , err := resourceClient.SaveChannel(req)
//if err != nil {
//	log.Printf("Failed to save channel : %s \n",err)
//}
//
//fmt.Println(txId)

//err = resourceClient.JoinChannel(channelId)
//if err != nil && !strings.Contains(err.Error(), "LedgerID already exists") {
//	log.Printf("Org peers failed to JoinChannel: %s \n", err)
//}



//p, err := peer.New()
//if err != nil {
//	return
//}

