package services

import (
	"fmt"
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


func CreateUser(context context.Context) {

	path := context.Path()
	log.Infoln(path)

	/*
	//
	var user models.User
	if err := context.ReadJSON(&user);err != nil{
		log.Errorf("failed to read user info to json : %s \n",err)
		context.JSON( models.FailedMsg("failed to create user"))
		return
	}
	 */

	// request type json form
	user := models.User{
		UserName: context.PostValueTrim("user_name"),
		Secret: context.PostValueTrim("secret"),
		UserType: context.PostValueTrim("user_type") ,
		OrgName: context.PostValueTrim("org_name"),
		CaName: context.PostValueTrim("ca_name"),
	}

	fmt.Println(user.CaName)

	priFile , pubFile , err := fabricClient.CreateUser(user.UserName,user.Secret,user.UserType,user.OrgName,user.CaName)
	if err != nil {
		if priFile != "" && pubFile != ""{
			context.JSON(models.FailedData(err.Error(),models.UserData{
				PriFile: priFile,
				PubFile: pubFile,
			}))
		}else{
			context.JSON( models.FailedMsg(err.Error()))
		}
		return
	}

	context.JSON(models.SuccessData(models.UserData{
		PriFile: priFile,
		PubFile: pubFile,
	}))
	return

}

func CreateChannel(context context.Context)  {

	path := context.Path()
	log.Infoln(path)

	info := models.CreateChannelInfo{
		ChannelId: context.PostValueTrim("channel_id"),
		UserName: context.PostValueTrim("user_name"),
		Org: context.PostValueTrim("org"),
		ChannelTx: channelTx,
	}

	_ , err := fabricClient.GetOrgTargetPeers(info.Org)
	if err != nil {
		context.JSON(models.FailedMsg(err.Error()))
		return
	}

	txId , err := fabricClient.CreateChannel(info.ChannelTx,info.Org,info.UserName,info.ChannelId)
	if err != nil {
		context.JSON(models.FailedMsg("Failed to create channel"))
		return
	}

	context.JSON(models.SuccessData(map[string]string{
		"txId":txId,
	}))

}

//func JoinChannel(context context.Context)  {
//
//	path := context.Path()
//	log.Infoln(path)
//
//	err := fabricClient.JoinChannel()
//	if err != nil {
//		log.Panicf("Failed to JoinChannel: %s \n",err)
//	}
//
//}
