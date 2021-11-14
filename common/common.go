package common

import (
	"github.com/sirupsen/logrus"
	"gofabric/models"
	"io/ioutil"
)

const (
	channelId = "mychannel"
	channelTx = "/usr/local/hyper/test2/configtx/channel-artifacts/mychannel.tx"
	connectConfigDir = "connect-config/channel-connection.yaml"
	chaincodeId = "mycc_0"
	chaincodePath = "newchaincode/test"
	ccVersion = "0"

)

var fabricClient *models.FabricClient
var log = logrus.New()
var orgs = []string{"org1","org2"}

func NewFabricClient()  {
	
	connectConfig, _ := ioutil.ReadFile(connectConfigDir)

	fabricClient =  models.NewFabricClient(connectConfig,channelId,orgs)
	//defer fabricClient.Close()

	err := fabricClient.Setup()
	if err != nil {
		log.Errorf("Failed to New Fabric Client: %s \n",err)
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

//func PackageCC() (string, []byte) {
//	return fabricClient.PackageCC()
//}

func InstallCC(label string, ccPkg []byte)  {
	fabricClient.InstallCC(label,ccPkg)
}

func GetInstalledCCPackage(label string)  {
	fabricClient.GetInstalledCCPackage(label)
}

func QueryInstalled(label string, packageID string)  {
	fabricClient.QueryInstalled(label,packageID)
}

func ApproveCC(label string)  {
	fabricClient.ApproveCC(label)
}

func CreateCC(){
	fabricClient.CreateCC(chaincodeId,chaincodePath,ccVersion)
}

func QueryLedger(){
	fabricClient.QueryLedger()
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

