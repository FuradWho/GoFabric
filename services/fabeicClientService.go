package services

import (
	"archive/zip"
	"fmt"
	"github.com/kataras/iris/v12/context"
	"github.com/sirupsen/logrus"
	"gofabric/models"
	"io"
	"io/ioutil"
	"os"
)

const (
	channelId = "mychannel"
	connectConfigDir = "connect-config/channel-connection.yaml"
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

			//context.JSON(models.FailedData(err.Error(),models.UserData{
			//	PriFile: priFile,
			//	PubFile: pubFile,
			//}))

			priFileDir := "/tmp/channel-msp/keystore/"+priFile
			pubFileDir := "/tmp/channel-store/"+pubFile
			fileName := "/home/fabric/ideaProject/GoFabric/cafiles/"+user.UserName + ".zip"

			err := zipFiles(fileName, []string{priFileDir, pubFileDir})
			if err != nil {
				return
			}

			context.Header("Content-Type","application/zip")
			err = context.SendFile(fileName, "cafiles.zip")
			if err != nil {
				log.Errorln(err)
			}

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
	}

	_ , err := fabricClient.GetOrgTargetPeers(info.Org)
	if err != nil {
		context.JSON(models.FailedMsg(err.Error()))
		return
	}

	txId , err := fabricClient.CreateChannel(info.Org,info.UserName,info.ChannelId)
	if err != nil {
		context.JSON(models.FailedMsg("Failed to create channel"))
		return
	}

	context.JSON(models.SuccessData(map[string]string{
		"txId":txId,
	}))

}

func JoinChannel(context context.Context)  {

	path := context.Path()
	log.Infoln(path)

	info := models.CreateChannelInfo{
		ChannelId: context.PostValueTrim("channel_id"),
		UserName: context.PostValueTrim("user_name"),
		Org: context.PostValueTrim("org"),
	}

	log.Infof("userName : %s \n",info.UserName)
	log.Infof("org : %s \n ", info.Org)
	log.Infof( "channelId : %s \n",info.ChannelId)

	err := fabricClient.JoinChannel(info.ChannelId,info.UserName,info.Org)
	if err != nil {
		context.JSON(models.FailedMsg("Failed to join channel"))
		return
	}

	context.JSON(models.SuccessMsg("success to join channel"))

}

func CreateCC(context context.Context)  {

	path := context.Path()
	log.Infoln(path)

	// chaincodeId, chaincodePath, version, org , userName, channelId string

}

func zipFiles(filename string , files []string) error  {
	newZipFile, err := os.Create(filename)
	if err != nil {
		log.Errorln(err)
		return err
	}

	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	for _,file := range files{
		fileToZip, err := os.Open(file)
		if err != nil {
			log.Errorln(err)
			return err
		}
		defer fileToZip.Close()

		info ,err := fileToZip.Stat()
		if err != nil {
			log.Errorln(err)
			return err
		}

		header , err := zip.FileInfoHeader(info)
		if err != nil {
			log.Errorln(err)
			return err
		}

		header.Name = fileToZip.Name()
		header.Method = zip.Deflate

		w , err := zipWriter.CreateHeader(header)
		if err != nil {
			log.Errorln(err)
			return err
		}

		_, err = io.Copy(w, fileToZip)
		if err != nil {
			log.Errorln(err)
			return err
		}
	}

	return nil
}
