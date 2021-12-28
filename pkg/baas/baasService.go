package baas

import (
	"archive/zip"
	"github.com/go-playground/validator/v10"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/kataras/iris/v12/context"
	log "github.com/sirupsen/logrus"
	"gofabric/models"
	"io"
	"os"
	"strconv"
	"sync"
)

var validate = validator.New()

type BaasService struct {
	baasClient *BaasClient
}

func (b *BaasService) InitBaasService(baas *BaasClient) {
	b.baasClient = baas
}

// Test
// @Summary 网络连通性测试
// @Description Test-网络连通性测试
// @Tags 测试接口
// @Produce json
// @Success 200 {json} json "{ "code": 200, "msg": "connection success" }"
// @Failure 400 {json} json "{ " " }"
// @Router /Test [get]
func (b *BaasService) Test(ctx context.Context) {
	ctx.JSON(models.SuccessMsg("connection success"))
}

// CreateUser
// @Summary　创建用户
// @Description CreateUser-创建用户
// @Tags 用户通道操作
// @Accept mpfd
// @Produce json
// @Param UserName formData string true "所在组织的用户名"
// @Param Secret formData string true "所在组织的用户密码"
// @Param UserType formData string true "用户类型"
// @Param OrgName formData string true "安装链码所在的组织"
// @Param CaName formData string true "所申请组织的CA服务器"
// @Success 200 {json} json "{ "code": 200, "msg": "success","data": { "PriFile": priFile, "PubFile" : pubFile }}"
// @Failure 400 {json} json "{ "code": 400, "msg": "Field to CreateUser","data": "" }"
// @Router /user/CreateUser [post]
func (b *BaasService) CreateUser(context context.Context) {

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

	err := validate.Struct(user)
	if err != nil {
		log.Errorln(err)
		context.JSON(models.FailedData("Field does not match", err))
		return
	}

	log.Infof("CreateUser info : %+v \n", user)

	priFile, pubFile, err := b.baasClient.CreateUser(user.UserName, user.Secret, user.UserType, user.OrgName, user.CaName)
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

// CreateChannel
// @Summary　创建通道
// @Description CreateChannel-创建通道
// @Tags 用户通道操作
// @Accept mpfd
// @Produce json
// @Param UserName formData string true "所在组织的用户名"
// @Param OrgName formData string true "安装链码所在的组织"
// @Param ChannelId formData string true "创建通道名称"
// @Param Orderer formData string true "排序节点"
// @Success 200 {json} json "{ "code": 200, "msg": "success","data": {"txId":""} }"
// @Failure 400 {json} json "{ "code": 400, "msg": "Failed to create channel","data": "" }"
// @Router /channel/CreateChannel [post]
func (b *BaasService) CreateChannel(context context.Context) {

	path := context.Path()
	log.Infoln(path)

	info := models.CreateChannelInfo{
		ChannelId: context.PostValueTrim("channel_id"),
		UserName:  context.PostValueTrim("user_name"),
		Org:       context.PostValueTrim("org"),
		Orderer:   context.PostValueTrim("orderer"),
	}

	err := validate.Struct(info)
	if err != nil {
		log.Errorln(err)
		context.JSON(models.FailedData("Field does not match", err))
		return
	}

	log.Infof("CreateChannel info : %+v \n", info)

	_, err = b.baasClient.GetOrgTargetPeers(info.Org)
	if err != nil {
		context.JSON(models.FailedMsg(err.Error()))
		return
	}

	txId, err := b.baasClient.CreateChannel(info.Org, info.UserName, info.ChannelId, info.Orderer)
	if err != nil {
		context.JSON(models.FailedMsg("Failed to create channel"))
		return
	}

	context.JSON(models.SuccessData(map[string]string{
		"txId": txId,
	}))

}

// JoinChannel
// @Summary　加入当前通道
// @Description JoinChannel-加入当前通道
// @Tags 用户通道操作
// @Accept mpfd
// @Produce json
// @Param UserName formData string true "所在组织的用户名"
// @Param OrgName formData string true "安装链码所在的组织"
// @Param ChannelId formData string true "创建通道名称"
// @Success 200 {json} json "{ "code": 200, "msg": "success","data": {"success to join channel"} }"
// @Failure 400 {json} json "{ "code": 400, "msg": "Failed to join channel","data": "" }"
// @Router /channel/JoinChannel [post]
func (b *BaasService) JoinChannel(context context.Context) {

	path := context.Path()
	log.Infoln(path)

	info := models.JoinChannelInfo{
		ChannelId: context.PostValueTrim("channel_id"),
		UserName:  context.PostValueTrim("user_name"),
		Org:       context.PostValueTrim("org"),
	}

	err := validate.Struct(info)
	if err != nil {
		log.Errorln(err)
		context.JSON(models.FailedData("Field does not match", err))
		return
	}

	log.Infof("JoinChannel info : %+v \n", info)

	err = b.baasClient.JoinChannel(info.ChannelId, info.UserName, info.Org)
	if err != nil {
		context.JSON(models.FailedMsg("Failed to join channel"))
		return
	}

	context.JSON(models.SuccessMsg("success to join channel"))

}

func (b *BaasService) CreateCC(context context.Context) {

	path := context.Path()
	log.Infoln(path)

	info := models.CcInfo{
		ChannelId:     context.PostValueTrim("channel_id"),
		UserName:      context.PostValueTrim("user_name"),
		Org:           context.PostValueTrim("org"),
		Version:       context.PostValueTrim("version"),
		ChaincodeId:   context.PostValueTrim("chaincode_id"),
		ChaincodePath: b.baasClient.Option.ChaincodePath,
	}
	log.Infof("create chaincode info : %+v \n", info)

	// chaincodeId, chaincodePath, version, org , userName, channelId string
	txId, err := b.baasClient.CreateCC(info.ChaincodeId, info.ChaincodePath, info.Version, info.Org, info.UserName, info.ChannelId)
	if err != nil {
		context.JSON(models.FailedMsg("Failed to create chaincode"))
		return
	}

	context.JSON(models.SuccessData(map[string]string{
		"txId": txId,
	}))

}

// InstallCC
// @Summary 安装链码
// @Description InstallCC-安装链码
// @Tags 链码操作
// @Accept mpfd
// @Produce json
// @Param UserName formData string true "所在组织的用户名"
// @Param Org formData string true "安装链码所在的组织"
// @Param ChaincodeId formData string true "链码ID(名称)"
// @Param Peer formData string true "安装链码所在的节点"
// @Success 200 {json} json "{ "code": 200, "msg": "success","data": {"txId":""} }"
// @Failure 400 {json} json "{ "code": 400, "msg": "The chaincode has installed ","data": "" }"
// @Router /cc/InstallCC [post]
func (b *BaasService) InstallCC(ctx context.Context) {

	path := ctx.Path()
	log.Infoln(path)

	info := models.InstallCCInfo{
		UserName:      ctx.PostValueTrim("user_name"),
		Org:           ctx.PostValueTrim("org"),
		ChaincodeId:   ctx.PostValueTrim("chaincode_id"),
		ChaincodePath: b.baasClient.Option.ChaincodePath,
		Peer:          ctx.PostValueTrim("peer"),
	}

	err := validate.Struct(info)
	if err != nil {
		log.Errorln(err)
		ctx.JSON(models.FailedData("Field does not match", err))
		return
	}

	log.Infof("InstallCC info : %+v \n", info)

	installed, err := b.baasClient.QueryInstalled(info.UserName, info.Org, info.Peer)
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
	var lck sync.Mutex
	lck.Lock()

	log.Infoln("locking ....")

	txId, err := b.baasClient.InstallCC(info.ChaincodeId, info.ChaincodePath, info.Org, info.UserName, info.Peer)
	if err != nil {
		ctx.JSON(models.FailedMsg("Failed to Install chaincode"))
		return
	}

	defer lck.Unlock()
	log.Infoln("Unlock ....")

	log.Infof("txId : %s \n", txId)
	ctx.JSON(models.SuccessData(map[string]string{
		"txId": txId,
	}))

}

// QueryInstalled
// @Summary 请求当前节点安装的链码
// @Description QueryInstalled-请求当前节点安装的链码
// @Tags 链码操作
// @Accept mpfd
// @Produce json
// @Param UserName formData string true "所在组织的用户名"
// @Param Org formData string true "链码所在的组织"
// @Param Peer formData string true "链码所在的节点"
// @Success 200 {json} json "{ "code": 200, "msg": "success","data": {"chaincodes":[]} }"
// @Failure 400 {json} json "{ "code": 400, "msg": "Failed to QueryInstalled chaincode","data": "" }"
// @Router /cc/QueryInstalled [post]
func (b *BaasService) QueryInstalled(ctx context.Context) {

	path := ctx.Path()
	log.Infoln(path)

	info := models.QueryInstalledInfo{
		UserName: ctx.PostValueTrim("user_name"),
		Org:      ctx.PostValueTrim("org"),
		Peer:     ctx.PostValueTrim("peer"),
	}

	err := validate.Struct(info)
	if err != nil {
		log.Errorln(err)
		ctx.JSON(models.FailedData("Field does not match", err))
		return
	}

	log.Infof("QueryInstalled info : %+v \n", info)

	installed, err := b.baasClient.QueryInstalled(info.UserName, info.Org, info.Peer)
	if err != nil {
		ctx.JSON(models.FailedMsg("Failed to QueryInstalled chaincode"))
		return
	}
	ctx.JSON(models.SuccessData(map[string][]resmgmt.LifecycleInstalledCC{
		"chaincodes": installed,
	}))

}

// ApproveCC
// @Summary 请求当前组织批准链码
// @Description ApproveCC-请求当前组织批准链码
// @Tags 链码操作
// @Accept mpfd
// @Produce json
// @Param UserName formData string true "所在组织的用户名"
// @Param Org formData string true "链码所在的组织"
// @Param Peer formData string true "链码所在的节点"
// @Param ChannelId formData string true "链码所在的通道"
// @Param Orderer formData string true "请求排序节点"
// @Param ChaincodeId formData string true "链码ID(名称)"
// @Param PackageId formData string true "链码包ID"
// @Param Version formData string true "链码版本"
// @Param Sequence formData string true "链码更新次数"
// @Success 200 {json} json "{ "code": 200, "msg": "success","data": {"txId":""} }"
// @Failure 400 {json} json "{ "code": 400, "msg": "Failed to approve the chaincode ","data": "" }"
// @Router /cc/ApproveCC [post]
func (b *BaasService) ApproveCC(ctx context.Context) {

	path := ctx.Path()
	log.Infoln(path)

	info := models.ApproveCCInfo{
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

	err := validate.Struct(info)
	if err != nil {
		log.Errorln(err)
		ctx.JSON(models.FailedData("Field does not match", err))
		return
	}

	log.Infof("ApproveCC info : %+v \n", info)

	installed, err := b.baasClient.QueryInstalled(info.UserName, info.Org, info.Peer)
	if err != nil {
		log.Errorln(err)
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
	var lck sync.Mutex
	lck.Lock()

	txnID, err := b.baasClient.ApproveCC(info.PackageId, info.ChaincodeId, info.Version, info.ChannelId, info.UserName, info.Org, info.Peer, info.Orderer, sequence)
	if err != nil {
		ctx.JSON(models.FailedMsg("Failed to approve the chaincode "))
		return
	}
	lck.Unlock()

	ctx.JSON(models.SuccessData(map[string]fab.TransactionID{
		"txId": txnID,
	}))

}

// QueryApprovedCC
// @Summary 查询是否通过当前组织批准链码
// @Description QueryApprovedCC-查询是否通过当前组织批准链码
// @Tags 链码操作
// @Accept mpfd
// @Produce json
// @Param UserName formData string true "所在组织的用户名"
// @Param Org formData string true "链码所在的组织"
// @Param Peer formData string true "链码所在的节点"
// @Param ChannelId formData string true "链码所在的通道"
// @Param ChaincodeId formData string true "链码ID(名称)"
// @Param Sequence formData string true "链码更新次数"
// @Success 200 {json} json "{ "code": 200, "msg": "success","data": {"packageId":""} }"
// @Failure 400 {json} json "{ "code": 400, "msg": "Failed to QueryApprovedCC the chaincode ","data": "" }"
// @Router /cc/QueryApprovedCC [post]
func (b *BaasService) QueryApprovedCC(ctx context.Context) {
	path := ctx.Path()
	log.Infoln(path)

	info := models.QueryApprovedCCInfo{
		UserName:    ctx.PostValueTrim("user_name"),
		Org:         ctx.PostValueTrim("org"),
		Peer:        ctx.PostValueTrim("peer"),
		ChaincodeId: ctx.PostValueTrim("chaincode_id"),
		ChannelId:   ctx.PostValueTrim("channel_id"),
		Sequence:    ctx.PostValueTrim("sequence"),
	}

	err := validate.Struct(info)
	if err != nil {
		log.Errorln(err)
		ctx.JSON(models.FailedData("Field does not match", err))
		return
	}

	log.Infof("QueryApprovedCC info : %+v \n", info)

	installed, err := b.baasClient.QueryInstalled(info.UserName, info.Org, info.Peer)
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

	packageId, err := b.baasClient.QueryApprovedCC(info.ChaincodeId, info.UserName, info.Org, info.ChannelId, info.Peer, sequence)
	if err != nil {
		ctx.JSON(models.FailedMsg("Failed to QueryApprovedCC the chaincode "))
		return
	}

	ctx.JSON(models.SuccessData(map[string]string{
		"packageId": packageId,
	}))

}

func (b *BaasService) CheckCCCommitReadiness(ctx context.Context) {

	path := ctx.Path()
	log.Infoln(path)

	info := models.CheckCCCommitReadinessInfo{
		UserName:    ctx.PostValueTrim("user_name"),
		Org:         ctx.PostValueTrim("org"),
		Peer:        ctx.PostValueTrim("peer"),
		ChaincodeId: ctx.PostValueTrim("chaincode_id"),
		ChannelId:   ctx.PostValueTrim("channel_id"),
		Sequence:    ctx.PostValueTrim("sequence"),
		Version:     ctx.PostValueTrim("version"),
	}

	err := validate.Struct(info)
	if err != nil {
		log.Errorln(err)
		ctx.JSON(models.FailedData("Field does not match", err))
		return
	}

	log.Infof("CheckCCCommitReadiness info : %+v \n", info)

	sequence, _ := strconv.Atoi(info.Sequence)

	installed, err := b.baasClient.QueryInstalled(info.UserName, info.Org, info.Peer)
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

	// func (f *b.baasClient) CheckCCCommitReadiness(ccID, version, user, org, channelId, peer string, sequence int) (map[string]bool, error) {

	var lck sync.Mutex
	lck.Lock()

	readiness, err := b.baasClient.CheckCCCommitReadiness(info.ChaincodeId, info.Version, info.UserName, info.Org, info.ChannelId, info.Peer, sequence)
	if err != nil {
		ctx.JSON(models.FailedMsg("Failed to CheckCCCommitReadiness the chaincode "))
		return
	}
	lck.Unlock()

	ctx.JSON(models.SuccessData(readiness))
}

// RequestInstallCCByOther
// @Summary 请求其他组织安装链码
// @Description RequestInstallCCByOther-请求其他组织安装链码
// @Tags 链码操作
// @Accept mpfd
// @Produce json
// @Param Org formData string true "链码所在的组织"
// @Param Peer formData string true "链码所在的节点"
// @Param ChaincodeId formData string true "链码ID(名称)"
// @Success 200 {json} json "{ "code": 200, "msg": "success","data": {"txId":""} }"
// @Failure 400 {json} json "{ "code": 400, "msg": "Failed to RequestInstallCCByOther ","data": "" }"
// @Router /cc/RequestInstallCCByOther [post]
func (b *BaasService) RequestInstallCCByOther(ctx context.Context) {

	path := ctx.Path()
	log.Infoln(path)

	info := models.RequestInstallCCByOtherInfo{
		UserName:      b.baasClient.Option.OrgUser,
		Org:           ctx.PostValueTrim("org"),
		ChaincodeId:   ctx.PostValueTrim("chaincode_id"),
		ChaincodePath: b.baasClient.Option.ChaincodePath,
		Peer:          ctx.PostValueTrim("peer"),
	}

	err := validate.Struct(info)
	if err != nil {
		log.Errorln(err)
		ctx.JSON(models.FailedData("Field does not match", err))
		return
	}

	log.Infof("RequestInstallCCByOther info : %+v \n", info)

	installed, err := b.baasClient.QueryInstalled(info.UserName, info.Org, info.Peer)
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

	var lck sync.Mutex
	lck.Lock()

	txId, err := b.baasClient.InstallCC(info.ChaincodeId, info.ChaincodePath, info.Org, info.UserName, info.Peer)
	if err != nil {
		ctx.JSON(models.FailedMsg("Failed to RequestInstallCCByOther "))
		return
	}
	log.Infof("txId : %s \n", txId)
	lck.Unlock()
	ctx.JSON(models.SuccessData(map[string]string{
		"txId": txId,
	}))

}

// RequestApproveCCByOther
// @Summary 请求其他组织批准链码
// @Description RequestApproveCCByOther-请求其他组织批准链码
// @Tags 链码操作
// @Accept mpfd
// @Produce json
// @Param Org formData string true "链码所在的组织"
// @Param Peer formData string true "链码所在的节点"
// @Param ChannelId formData string true "链码所在的通道"
// @Param Orderer formData string true "请求排序节点"
// @Param ChaincodeId formData string true "链码ID(名称)"
// @Param PackageId formData string true "链码包ID"
// @Param Version formData string true "链码版本"
// @Param Sequence formData string true "链码更新次数"
// @Success 200 {json} json "{ "code": 200, "msg": "success","data": {"txId":""} }"
// @Failure 400 {json} json "{ "code": 400, "msg": "Failed to RequestApproveCCByOther the chaincode ","data": "" }"
// @Router /cc/RequestApproveCCByOther [post]
func (b *BaasService) RequestApproveCCByOther(ctx context.Context) {

	path := ctx.Path()
	log.Infoln(path)

	info := models.RequestApproveCCByOtherInfo{
		PackageId:   ctx.PostValueTrim("package_id"),
		UserName:    b.baasClient.Option.UserName,
		Org:         ctx.PostValueTrim("org"),
		Peer:        ctx.PostValueTrim("peer"),
		ChaincodeId: ctx.PostValueTrim("chaincode_id"),
		Version:     ctx.PostValueTrim("version"),
		ChannelId:   ctx.PostValueTrim("channel_id"),
		Orderer:     ctx.PostValueTrim("orderer"),
		Sequence:    ctx.PostValueTrim("sequence"),
	}

	err := validate.Struct(info)
	if err != nil {
		log.Errorln(err)
		ctx.JSON(models.FailedData("Field does not match", err))
		return
	}

	log.Infof("RequestApproveCCByOther info : %+v \n", info)

	installed, err := b.baasClient.QueryInstalled(info.UserName, info.Org, info.Peer)
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
	var lck sync.Mutex
	lck.Lock()
	txnID, err := b.baasClient.ApproveCC(info.PackageId, info.ChaincodeId, info.Version, info.ChannelId, info.UserName, info.Org, info.Peer, info.Orderer, sequence)
	if err != nil {
		ctx.JSON(models.FailedMsg("Failed to RequestApproveCCByOther the chaincode "))
		return
	}
	lck.Unlock()

	ctx.JSON(models.SuccessData(map[string]fab.TransactionID{
		"txnID": txnID,
	}))

}

// CommitCC
// @Summary 提交已通过MSP认证的链码
// @Description CommitCC-提交已通过MSP认证的链码
// @Tags 链码操作
// @Accept mpfd
// @Produce json
// @Param user_name formData string true "所在组织的用户名"
// @Param org formData string true "链码所在的组织"
// @Param peer formData string true "链码所在的节点"
// @Param channel_id formData string true "链码所在的通道"
// @Param orderer formData string true "请求排序节点"
// @Param chaincode_id formData string true "链码ID(名称)"
// @Param version formData string true "链码版本"
// @Param sequence formData string true "链码更新次数"
// @Success 200 {json} json "{ "code": 200, "msg": "success","data": {"txId":""} }"
// @Failure 400 {json} json "{ "code": 400, "msg": "Failed to CommitCC ","data": "" }"
// @Router /cc/CommitCC [post]
func (b *BaasService) CommitCC(ctx context.Context) {
	path := ctx.Path()
	log.Infoln(path)

	info := models.CommitCCInfo{
		UserName:    ctx.PostValueTrim("user_name"),
		Org:         ctx.PostValueTrim("org"),
		ChaincodeId: ctx.PostValueTrim("chaincode_id"),
		Version:     ctx.PostValueTrim("version"),
		Peer:        ctx.PostValueTrim("peer"),
		ChannelId:   ctx.PostValueTrim("channel_id"),
		Orderer:     ctx.PostValueTrim("orderer"),
		Sequence:    ctx.PostValueTrim("sequence"),
	}

	err := validate.Struct(info)
	if err != nil {
		log.Errorln(err)
		ctx.JSON(models.FailedData("Field does not match", err))
		return
	}

	log.Infof("RequestApproveCCByOther info : %+v \n", info)

	installed, err := b.baasClient.QueryInstalled(info.UserName, info.Org, info.Peer)
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

	var lck sync.Mutex
	lck.Lock()

	txId, err := b.baasClient.CommitCC(info.ChaincodeId, info.UserName, info.Org, info.ChannelId, info.Orderer, info.Version, sequence)
	if err != nil {
		ctx.JSON(models.FailedMsg("Failed to CommitCC "))
		return
	}

	lck.Unlock()

	ctx.JSON(models.SuccessData(map[string]string{
		"txId": string(txId),
	}))
}

// GetOrgTargetPeers
// @Summary 获取组织节点信息
// @Description GetOrgTargetPeers-获取组织节点信息
// @Tags 网络通道操作
// @Produce json
// @Param org path string true "所在组织名"
// @Success 200 {json} json "{ "code": 200, "msg": "success","data": {"peers":[]} }"
// @Failure 400 {json} json "{ "code": 400, "msg": "Failed to GetOrgTargetPeers ","data": "" }"
// @Router /channel/GetOrgTargetPeers [get]
func (b *BaasService) GetOrgTargetPeers(ctx context.Context) {

	path := ctx.Path()
	log.Infoln(path)

	info := models.CcInfo{
		Org: ctx.URLParam("org"),
	}

	log.Infof("GetOrgTargetPeers info : %+v \n", info)

	peers, err := b.baasClient.GetOrgTargetPeers(info.Org)
	if err != nil {
		ctx.JSON(models.FailedMsg("Failed to GetOrgTargetPeers"))
		return
	}

	ctx.JSON(models.SuccessData(map[string][]string{
		"peers": peers,
	}))

}

// GetNetworkConfig
// @Summary 获取网络信息
// @Description GetNetworkConfig-获取"网络通道组织节点"信息
// @Tags 网络通道操作
// @Produce json
// @Success 200 {json} json "{ "code": 200, "msg": "success","data": {} }"
// @Failure 400 {json} json "{ "code": 400, "msg": "Failed to GetNetworkConfig ","data": "" }"
// @Router /channel/GetNetworkConfig [get]
func (b *BaasService) GetNetworkConfig(ctx context.Context) {

	path := ctx.Path()
	log.Infoln(path)

	config, err := b.baasClient.GetNetworkConfig()
	if err != nil {
		ctx.JSON(models.FailedMsg("Failed to GetOrgTargetPeers"))
		return
	}

	ctx.JSON(models.SuccessData(config))

}

// LifeCycleChaincodeTest
// @Summary 操作测试接口
// @Description LifeCycleChaincodeTest-操作测试接口
// @Tags 测试接口
// @Produce json
// @Success 200 {json} json "{ "code": 200, "msg": "connection success" }"
// @Failure 400 {json} json "{ " " }"
// @Router /LifeCycleChaincodeTest [get]
func (b *BaasService) LifeCycleChaincodeTest(ctx context.Context) {

	path := ctx.Path()
	log.Infoln(path)

	// chaincodeId, chaincodePath, org , user string

	//txId, err := b.baasClient.InstallCC("Test3","/usr/local/soft/fabric-test5/chaincode/newchaincode/test","org1","Admin")
	//if err != nil {
	//	return
	//}
	//log.Infoln(txId)
	//
	//tx2Id, err := b.baasClient.InstallCC("Test3","/usr/local/soft/fabric-test5/chaincode/newchaincode/test","org2","Admin")
	//if err != nil {
	//	return
	//}
	//log.Infoln(tx2Id)

	// Test0:5d6f5940712a57ee77265c718ec9f25c9683f286d7450338f3e47e1a46fcf52d

	// Test0:5d6f5940712a57ee77265c718ec9f25c9683f286d7450338f3e47e1a46fcf52d
	//  Test1:f9785b613f60c15c518fdab380e42c05938112b211fa632b75797f5fe4680855

	//Test2:792f96243801760b2dfcbae9b5a505aedcde14a63e8f6dcea01125f6ec0ce4a0

	// Test3:c11fd6513a390b097694f72dc0a089e27bf633481ae37e4ce9b06cdea3bc5b80

	//err := b.baasClient.ApproveCC("Test3:c11fd6513a390b097694f72dc0a089e27bf633481ae37e4ce9b06cdea3bc5b80", "org1", "Test3", "0", "mychannel", "Admin","peer0.org1.example.com")
	//if err != nil {
	//	return
	//}
	//
	//err = b.baasClient.ApproveCC("Test3:c11fd6513a390b097694f72dc0a089e27bf633481ae37e4ce9b06cdea3bc5b80", "org2", "Test3", "0", "mychannel", "Admin","peer0.org2.example.com")
	//if err != nil {
	//	return
	//}
	//
	//err = b.baasClient.QueryInstalled("Admin", "org1")
	//if err != nil {
	//	return
	//}

	//err = b.baasClient.GetInstalledCCPackage("Admin", "Test0:5d6f5940712a57ee77265c718ec9f25c9683f286d7450338f3e47e1a46fcf52d", "org1")
	//if err != nil {
	//	return
	//}
	//

	//time.Sleep(time.Duration(10)*time.Second)
	//err :=  b.baasClient.CheckCCCommitReadiness("Test3", "Admin", "org1", "mychannel","peer0.org1.example.com")
	//if err != nil {
	//	return
	//}
	//
	//time.Sleep(time.Duration(5)*time.Second)
	//
	//err =  b.baasClient.CheckCCCommitReadiness("Test3", "Admin", "org2", "mychannel","peer0.org2.example.com")
	//if err != nil {
	//	return
	//}

	//err := b.baasClient.QueryApprovedCC("Test1", "Admin", "org1", "mychannel")
	//if err != nil {
	//	return
	//}
	//err = b.baasClient.CommitCC("Test3", "Admin", "org2", "mychannel", "peer0.org2.example.com")
	//if err != nil {
	//	return
	//}

	//
	//err = b.baasClient.CommitCC("Test3", "Admin", "org1", "mychannel", "peer0.org1.example.com")
	//if err != nil {
	//	return
	//}

	//b.baasClient.QueryCommittedCC("Test3", "Admin", "org1", "mychannel", "peer0.org1.example.com")

	order, err := b.baasClient.QueryConfigBlockFromOrder("Admin", "Ordererorg", "channel", "orderer.example.com")
	if err != nil {
		return
	}

	ctx.JSON(models.SuccessData(order))
}

func QueryConfigBlock(ctx context.Context) {

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

func (b *BaasService) AuthenticateUser(ctx context.Context) {

	path := ctx.Path()
	log.Infoln(path)

	b.baasClient.AuthenticateUser("org1")

}
