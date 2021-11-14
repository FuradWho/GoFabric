package services

import (
	"github.com/kataras/iris/v12/context"
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
		return
	}
}


func CreateUser(context *context.Context) interface{}  {

	//path := context.Path()
	//log.Infoln(path)
	//
	//var user models.User
	//if err := context.ReadJSON(&user);err != nil{
	//	log.Errorf("failed to read user info to json : %s \n",err)
	//}

	priFile , pubFile , err := fabricClient.CreateUser(user.UserName,user.Secret,user.UserType,user.OrgName,user.CaName)
	if err != nil {
		return models.FailedMsg("failed to create user")
	}

	return models.SuccessData([]string{priFile,pubFile})

}