package services

import (
	"archive/zip"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/kataras/iris/v12/context"
	"github.com/sirupsen/logrus"
	"gofabric/models"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

const (
	channelId = "mychannel"
	//connectConfigDir = "connect-config/channel-connection.yaml"
	connectConfigDir = "connect-config/orgcpp-config.yaml"
	chaincodePath    = "/usr/local/soft/fabric-test5/chaincode/newchaincode/test"
	Admin            = "Admin"
)

var fabricClient *models.FabricClient
var log = logrus.New()
var orgs = []string{"org1", "org2"}

func NewFabricClient() {

	connectConfig, _ := ioutil.ReadFile(connectConfigDir)
	fabricClient = models.NewFabricClient(connectConfig, channelId, orgs)
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
		Secret:   context.PostValueTrim("secret"),
		UserType: context.PostValueTrim("user_type"),
		OrgName:  context.PostValueTrim("org_name"),
		CaName:   context.PostValueTrim("ca_name"),
	}

	fmt.Println(user.CaName)

	priFile, pubFile, err := fabricClient.CreateUser(user.UserName, user.Secret, user.UserType, user.OrgName, user.CaName)
	if err != nil {
		if priFile != "" && pubFile != "" {

			//context.JSON(models.FailedData(err.Error(),models.UserData{
			//	PriFile: priFile,
			//	PubFile: pubFile,
			//}))

			priFileDir := "/tmp/channel-msp/keystore/" + priFile
			pubFileDir := "/tmp/channel-store/" + pubFile
			fileName := "/home/fabric/ideaProject/GoFabric/cafiles/" + user.UserName + ".zip"

			err := zipFiles(fileName, []string{priFileDir, pubFileDir})
			if err != nil {
				return
			}

			context.Header("Content-Type", "application/zip")
			err = context.SendFile(fileName, "cafiles.zip")
			if err != nil {
				log.Errorln(err)
			}

		} else {
			context.JSON(models.FailedMsg(err.Error()))
		}
		return
	}

	context.JSON(models.SuccessData(models.UserData{
		PriFile: priFile,
		PubFile: pubFile,
	}))
	return

}

func CreateChannel(context context.Context) {

	path := context.Path()
	log.Infoln(path)

	info := models.CreateChannelInfo{
		ChannelId: context.PostValueTrim("channel_id"),
		UserName:  context.PostValueTrim("user_name"),
		Org:       context.PostValueTrim("org"),
	}

	_, err := fabricClient.GetOrgTargetPeers(info.Org)
	if err != nil {
		context.JSON(models.FailedMsg(err.Error()))
		return
	}

	txId, err := fabricClient.CreateChannel(info.Org, info.UserName, info.ChannelId)
	if err != nil {
		context.JSON(models.FailedMsg("Failed to create channel"))
		return
	}

	context.JSON(models.SuccessData(map[string]string{
		"txId": txId,
	}))

}

func JoinChannel(context context.Context) {

	path := context.Path()
	log.Infoln(path)

	info := models.CreateChannelInfo{
		ChannelId: context.PostValueTrim("channel_id"),
		UserName:  context.PostValueTrim("user_name"),
		Org:       context.PostValueTrim("org"),
	}

	log.Infof("join channel info : %+v \n", info)

	err := fabricClient.JoinChannel(info.ChannelId, info.UserName, info.Org)
	if err != nil {
		context.JSON(models.FailedMsg("Failed to join channel"))
		return
	}

	context.JSON(models.SuccessMsg("success to join channel"))

}

func CreateCC(context context.Context) {

	path := context.Path()
	log.Infoln(path)

	/*
		ChaincodeId string `json:"chaincode_id"`
		ChaincodePath string `json:"chaincode_path"`
		Version string `json:"version"`
		Org string `json:"org"`
		UserName string `json:"user_name"`
		ChannelId string `json:"channel_id"`
	*/
	info := models.CcInfo{
		ChannelId:     context.PostValueTrim("channel_id"),
		UserName:      context.PostValueTrim("user_name"),
		Org:           context.PostValueTrim("org"),
		Version:       context.PostValueTrim("version"),
		ChaincodeId:   context.PostValueTrim("chaincode_id"),
		ChaincodePath: chaincodePath,
	}

	log.Infof("create chaincode info : %+v \n", info)

	// chaincodeId, chaincodePath, version, org , userName, channelId string
	txId, err := fabricClient.CreateCC(info.ChaincodeId, info.ChaincodePath, info.Version, info.Org, info.UserName, info.ChannelId)
	if err != nil {
		context.JSON(models.FailedMsg("Failed to create chaincode"))
		return
	}

	context.JSON(models.SuccessData(map[string]string{
		"txId": txId,
	}))

}

func InstallCC(ctx context.Context) {

	path := ctx.Path()
	log.Infoln(path)

	info := models.CcInfo{
		UserName:      ctx.PostValueTrim("user_name"),
		Org:           ctx.PostValueTrim("org"),
		ChaincodeId:   ctx.PostValueTrim("chaincode_id"),
		ChaincodePath: chaincodePath,
		Peer:          ctx.PostValueTrim("peer"),
	}

	log.Infof("InstallCC info : %+v \n", info)

	installed, err := fabricClient.QueryInstalled(info.UserName, info.Org, info.Peer)
	if err != nil {
		ctx.JSON(models.FailedMsg("Failed to QueryInstalled chaincode"))
		return
	}

	flag := false
	for _, chaincode := range installed {
		if info.ChaincodeId != chaincode.Label {
			continue
		} else {
			flag = true
		}
	}
	if flag {
		ctx.JSON(models.FailedMsg("The chaincode has installed "))
		return
	}

	txId, err := fabricClient.InstallCC(info.ChaincodeId, info.ChaincodePath, info.Org, info.UserName, info.Peer)
	if err != nil {
		ctx.JSON(models.FailedMsg("Failed to Install chaincode"))
		return
	}
	log.Infof("txId : %s \n", txId)
	ctx.JSON(models.SuccessData(map[string]string{
		"txId": txId,
	}))

}

func QueryInstalled(ctx context.Context) {

	path := ctx.Path()
	log.Infoln(path)

	info := models.CcInfo{
		UserName: ctx.PostValueTrim("user_name"),
		Org:      ctx.PostValueTrim("org"),
		Peer:     ctx.PostValueTrim("peer"),
	}

	log.Infof("QueryInstalled info : %+v \n", info)

	installed, err := fabricClient.QueryInstalled(info.UserName, info.Org, info.Peer)
	if err != nil {
		ctx.JSON(models.FailedMsg("Failed to QueryInstalled chaincode"))
		return
	}
	ctx.JSON(models.SuccessData(map[string][]resmgmt.LifecycleInstalledCC{
		"chaincodes": installed,
	}))

}

func ApproveCC(ctx context.Context) {

	path := ctx.Path()
	log.Infoln(path)

	info := models.CcInfo{
		PackageId:   ctx.PostValueTrim("package_id"),
		UserName:    ctx.PostValueTrim("user_name"),
		Org:         ctx.PostValueTrim("org"),
		Peer:        ctx.PostValueTrim("peer"),
		ChaincodeId: ctx.PostValueTrim("chaincode_id"),
		Version:     ctx.PostValueTrim("version"),
		ChannelId:   ctx.PostValueTrim("channel_id"),
		Orderer:     ctx.PostValueTrim("orderer"),
		Sequence:    ctx.PostValueTrim("sequence"),
	}

	log.Infof("ApproveCC info : %+v \n", info)

	installed, err := fabricClient.QueryInstalled(info.UserName, info.Org, info.Peer)
	if err != nil {
		ctx.JSON(models.FailedMsg("Failed to QueryInstalled chaincode"))
		return
	}

	flag := false
	for _, chaincode := range installed {
		if info.ChaincodeId != chaincode.Label {
			continue
		} else {
			flag = true
		}
	}
	if !flag {
		ctx.JSON(models.FailedMsg("The chaincode has not installed "))
		return
	}
	sequence, _ := strconv.Atoi(info.Sequence)

	txnID, err := fabricClient.ApproveCC(info.PackageId, info.ChaincodeId, info.Version, info.ChannelId, info.UserName, info.Org, info.Peer, info.Orderer, sequence)
	if err != nil {
		ctx.JSON(models.FailedMsg("Failed to approve the chaincode "))
		return
	}

	ctx.JSON(models.SuccessData(map[string]fab.TransactionID{
		"txnID": txnID,
	}))

}

func QueryApprovedCC(ctx context.Context) {
	path := ctx.Path()
	log.Infoln(path)

	info := models.CcInfo{
		UserName:    ctx.PostValueTrim("user_name"),
		Org:         ctx.PostValueTrim("org"),
		Peer:        ctx.PostValueTrim("peer"),
		ChaincodeId: ctx.PostValueTrim("chaincode_id"),
		ChannelId:   ctx.PostValueTrim("channel_id"),
		Sequence:    ctx.PostValueTrim("sequence"),
	}

	log.Infof("QueryApprovedCC info : %+v \n", info)

	installed, err := fabricClient.QueryInstalled(info.UserName, info.Org, info.Peer)
	if err != nil {
		ctx.JSON(models.FailedMsg("Failed to QueryInstalled chaincode"))
		return
	}

	flag := false
	for _, chaincode := range installed {
		if info.ChaincodeId != chaincode.Label {
			continue
		} else {
			flag = true
		}
	}
	if !flag {
		ctx.JSON(models.FailedMsg("The chaincode has not installed "))
		return
	}

	sequence, _ := strconv.Atoi(info.Sequence)

	packageId, err := fabricClient.QueryApprovedCC(info.ChaincodeId, info.UserName, info.Org, info.ChannelId, info.Peer, sequence)
	if err != nil {
		ctx.JSON(models.FailedMsg("Failed to QueryApprovedCC the chaincode "))
		return
	}

	ctx.JSON(models.SuccessData(map[string]string{
		"packageId": packageId,
	}))

}

func CheckCCCommitReadiness(ctx context.Context) {

	path := ctx.Path()
	log.Infoln(path)

	info := models.CcInfo{
		UserName:    ctx.PostValueTrim("user_name"),
		Org:         ctx.PostValueTrim("org"),
		Peer:        ctx.PostValueTrim("peer"),
		ChaincodeId: ctx.PostValueTrim("chaincode_id"),
		ChannelId:   ctx.PostValueTrim("channel_id"),
		Sequence:    ctx.PostValueTrim("sequence"),
		Version:     ctx.PostValueTrim("version"),
	}

	log.Infof("CheckCCCommitReadiness info : %+v \n", info)

	sequence, _ := strconv.Atoi(info.Sequence)

	installed, err := fabricClient.QueryInstalled(info.UserName, info.Org, info.Peer)
	if err != nil {
		ctx.JSON(models.FailedMsg("Failed to QueryInstalled chaincode"))
		return
	}

	flag := false
	for _, chaincode := range installed {
		if info.ChaincodeId != chaincode.Label {
			continue
		} else {
			flag = true
		}
	}
	if !flag {
		ctx.JSON(models.FailedMsg("The chaincode has not installed "))
		return
	}

	// func (f *FabricClient) CheckCCCommitReadiness(ccID, version, user, org, channelId, peer string, sequence int) (map[string]bool, error) {

	readiness, err := fabricClient.CheckCCCommitReadiness(info.ChaincodeId, info.Version, info.UserName, info.Org, info.ChannelId, info.Peer, sequence)
	if err != nil {
		ctx.JSON(models.FailedMsg("Failed to CheckCCCommitReadiness the chaincode "))
		return
	}

	ctx.JSON(models.SuccessData(readiness))
}

func RequestInstallCCByOther(ctx context.Context) {

	path := ctx.Path()
	log.Infoln(path)

	info := models.CcInfo{
		UserName:      Admin,
		Org:           ctx.PostValueTrim("org"),
		ChaincodeId:   ctx.PostValueTrim("chaincode_id"),
		ChaincodePath: chaincodePath,
		Peer:          ctx.PostValueTrim("peer"),
	}

	log.Infof("RequestInstallCCByOther info : %+v \n", info)

	installed, err := fabricClient.QueryInstalled(info.UserName, info.Org, info.Peer)
	if err != nil {
		ctx.JSON(models.FailedMsg("Failed to QueryInstalled chaincode"))
		return
	}

	flag := false
	for _, chaincode := range installed {
		if info.ChaincodeId != chaincode.Label {
			continue
		} else {
			flag = true
		}
	}
	if flag {
		ctx.JSON(models.FailedMsg("The chaincode has installed "))
		return
	}

	txId, err := fabricClient.InstallCC(info.ChaincodeId, info.ChaincodePath, info.Org, info.UserName, info.Peer)
	if err != nil {
		ctx.JSON(models.FailedMsg("Failed to RequestInstallCCByOther "))
		return
	}
	log.Infof("txId : %s \n", txId)
	ctx.JSON(models.SuccessData(map[string]string{
		"txId": txId,
	}))

}

func RequestApproveCCByOther(ctx context.Context) {

	path := ctx.Path()
	log.Infoln(path)

	info := models.CcInfo{
		PackageId:   ctx.PostValueTrim("package_id"),
		UserName:    Admin,
		Org:         ctx.PostValueTrim("org"),
		Peer:        ctx.PostValueTrim("peer"),
		ChaincodeId: ctx.PostValueTrim("chaincode_id"),
		Version:     ctx.PostValueTrim("version"),
		ChannelId:   ctx.PostValueTrim("channel_id"),
		Orderer:     ctx.PostValueTrim("orderer"),
		Sequence:    ctx.PostValueTrim("sequence"),
	}

	log.Infof("RequestApproveCCByOther info : %+v \n", info)

	installed, err := fabricClient.QueryInstalled(info.UserName, info.Org, info.Peer)
	if err != nil {
		ctx.JSON(models.FailedMsg("Failed to QueryInstalled chaincode"))
		return
	}

	flag := false
	for _, chaincode := range installed {
		if info.ChaincodeId != chaincode.Label {
			continue
		} else {
			flag = true
		}
	}
	if !flag {
		ctx.JSON(models.FailedMsg("The chaincode has not installed "))
		return
	}
	sequence, _ := strconv.Atoi(info.Sequence)

	txnID, err := fabricClient.ApproveCC(info.PackageId, info.ChaincodeId, info.Version, info.ChannelId, info.UserName, info.Org, info.Peer, info.Orderer, sequence)
	if err != nil {
		ctx.JSON(models.FailedMsg("Failed to RequestApproveCCByOther the chaincode "))
		return
	}

	ctx.JSON(models.SuccessData(map[string]fab.TransactionID{
		"txnID": txnID,
	}))

}

func CommitCC(ctx context.Context){
	path := ctx.Path()
	log.Infoln(path)

	info := models.CcInfo{
		UserName:    ctx.PostValueTrim("user_name"),
		Org:         ctx.PostValueTrim("org"),
		ChaincodeId: ctx.PostValueTrim("chaincode_id"),
		Version:     ctx.PostValueTrim("version"),
		Peer:        ctx.PostValueTrim("peer"),
		ChannelId:   ctx.PostValueTrim("channel_id"),
		Orderer:     ctx.PostValueTrim("orderer"),
		Sequence:    ctx.PostValueTrim("sequence"),
	}

	log.Infof("RequestApproveCCByOther info : %+v \n", info)

	installed, err := fabricClient.QueryInstalled(info.UserName, info.Org, info.Peer)
	if err != nil {
		ctx.JSON(models.FailedMsg("Failed to QueryInstalled chaincode"))
		return
	}

	flag := false
	for _, chaincode := range installed {
		if info.ChaincodeId != chaincode.Label {
			continue
		} else {
			flag = true
		}
	}
	if !flag {
		ctx.JSON(models.FailedMsg("The chaincode has installed "))
		return
	}

	sequence, _ := strconv.Atoi(info.Sequence)

	txId , err := fabricClient.CommitCC(info.ChaincodeId, info.UserName, info.Org, info.ChannelId,info.Orderer,info.Version,sequence)
	if err != nil {
		ctx.JSON(models.FailedMsg("Failed to CommitCC "))
		return
	}
	ctx.JSON(models.SuccessData(map[string]string{
		"txId": string(txId),
	}))
}


func GetOrgTargetPeers(ctx context.Context) {

	path := ctx.Path()
	log.Infoln(path)

	info := models.CcInfo{
		Org: ctx.URLParam("org"),
	}

	log.Infof("GetOrgTargetPeers info : %+v \n", info)

	peers, err := fabricClient.GetOrgTargetPeers(info.Org)
	if err != nil {
		ctx.JSON(models.FailedMsg("Failed to GetOrgTargetPeers"))
		return
	}

	ctx.JSON(models.SuccessData(map[string][]string{
		"peers": peers,
	}))

}

func GetNetworkConfig(ctx context.Context) {

	path := ctx.Path()
	log.Infoln(path)

	config, err := fabricClient.GetNetworkConfig()
	if err != nil {
		ctx.JSON(models.FailedMsg("Failed to GetOrgTargetPeers"))
		return
	}

	ctx.JSON(models.SuccessData(config))

}

func LifeCycleChaincodeTest(ctx context.Context) {

	path := ctx.Path()
	log.Infoln(path)

	// chaincodeId, chaincodePath, org , user string

	//txId, err := fabricClient.InstallCC("Test3","/usr/local/soft/fabric-test5/chaincode/newchaincode/test","org1","Admin")
	//if err != nil {
	//	return
	//}
	//log.Infoln(txId)
	//
	//tx2Id, err := fabricClient.InstallCC("Test3","/usr/local/soft/fabric-test5/chaincode/newchaincode/test","org2","Admin")
	//if err != nil {
	//	return
	//}
	//log.Infoln(tx2Id)

	// Test0:5d6f5940712a57ee77265c718ec9f25c9683f286d7450338f3e47e1a46fcf52d

	// Test0:5d6f5940712a57ee77265c718ec9f25c9683f286d7450338f3e47e1a46fcf52d
	//  Test1:f9785b613f60c15c518fdab380e42c05938112b211fa632b75797f5fe4680855

	//Test2:792f96243801760b2dfcbae9b5a505aedcde14a63e8f6dcea01125f6ec0ce4a0

	// Test3:c11fd6513a390b097694f72dc0a089e27bf633481ae37e4ce9b06cdea3bc5b80

	//err := fabricClient.ApproveCC("Test3:c11fd6513a390b097694f72dc0a089e27bf633481ae37e4ce9b06cdea3bc5b80", "org1", "Test3", "0", "mychannel", "Admin","peer0.org1.example.com")
	//if err != nil {
	//	return
	//}
	//
	//err = fabricClient.ApproveCC("Test3:c11fd6513a390b097694f72dc0a089e27bf633481ae37e4ce9b06cdea3bc5b80", "org2", "Test3", "0", "mychannel", "Admin","peer0.org2.example.com")
	//if err != nil {
	//	return
	//}
	//
	//err = fabricClient.QueryInstalled("Admin", "org1")
	//if err != nil {
	//	return
	//}

	//err = fabricClient.GetInstalledCCPackage("Admin", "Test0:5d6f5940712a57ee77265c718ec9f25c9683f286d7450338f3e47e1a46fcf52d", "org1")
	//if err != nil {
	//	return
	//}
	//

	//time.Sleep(time.Duration(10)*time.Second)
	//err :=  fabricClient.CheckCCCommitReadiness("Test3", "Admin", "org1", "mychannel","peer0.org1.example.com")
	//if err != nil {
	//	return
	//}
	//
	//time.Sleep(time.Duration(5)*time.Second)
	//
	//err =  fabricClient.CheckCCCommitReadiness("Test3", "Admin", "org2", "mychannel","peer0.org2.example.com")
	//if err != nil {
	//	return
	//}

	//err := fabricClient.QueryApprovedCC("Test1", "Admin", "org1", "mychannel")
	//if err != nil {
	//	return
	//}
	//err = fabricClient.CommitCC("Test3", "Admin", "org2", "mychannel", "peer0.org2.example.com")
	//if err != nil {
	//	return
	//}

	//
	//err = fabricClient.CommitCC("Test3", "Admin", "org1", "mychannel", "peer0.org1.example.com")
	//if err != nil {
	//	return
	//}

	fabricClient.QueryCommittedCC("Test3", "Admin", "org1", "mychannel", "peer0.org1.example.com")

}

func zipFiles(filename string, files []string) error {
	newZipFile, err := os.Create(filename)
	if err != nil {
		log.Errorln(err)
		return err
	}

	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	for _, file := range files {
		fileToZip, err := os.Open(file)
		if err != nil {
			log.Errorln(err)
			return err
		}
		defer fileToZip.Close()

		info, err := fileToZip.Stat()
		if err != nil {
			log.Errorln(err)
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			log.Errorln(err)
			return err
		}

		header.Name = fileToZip.Name()
		header.Method = zip.Deflate

		w, err := zipWriter.CreateHeader(header)
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
