package models

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config/lookup"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	lcpackager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/lifecycle"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/policydsl"
	"github.com/sirupsen/logrus"
	"os/exec"
	"time"

	//	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"os"
	"strings"
)

const (
	Admin = "User3"
	peer1 = "peer1-org1"
	peer2 = "peer2-org1"
)

var log = logrus.New()

type FabricClient struct {
	ConnectionFile []byte
	Orgs []string
	ChannelId string
	GoPath string

	userName string
	userOrg string

	resmgmtClients map[string]*resmgmt.Client
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
		log.Error("Failed to setup main sdk ")
		return err
	}
	f.sdk = sdk

	resmgmtClients := make(map[string]*resmgmt.Client)
	for _, v := range f.Orgs {
		resmgmtClient , err := resmgmt.New(sdk.Context(fabsdk.WithUser(Admin),fabsdk.WithOrg(v)))
		if err != nil {
			log.Errorf("Failed to create channel management client : %s \n",err)
			return err
		}
		resmgmtClients[v] = resmgmtClient
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

func (f *FabricClient) CreateUser(userName string,secret string,userType string,orgName string,caName string) (priFile string, pubFile string, err error) {

	ctx := f.sdk.Context()
	mspClient, err := mspclient.New(ctx,mspclient.WithOrg(orgName))
	if err != nil {
		log.Errorf("Failed to create msp client: %s\n", err)
		return "","",err
	}

	caInfo, err := mspClient.GetCAInfo()
	if err != nil {
		log.Errorf("Failed to get CA Info :%s \n",err)
		return"","",err
	}

	if caInfo.CAName != caName {
		log.Errorln("Not match ca ")
		return "","",errors.New("Not match ca")
	}
	/*
	affiliations, err := mspClient.GetAllIdentities()
	if err != nil {
		log.Printf("%s \n",err)
	}

	for _ , info := range affiliations{
		fmt.Println(info.ID)
		fmt.Println(info.Type)
		fmt.Println(info.Attributes)
		fmt.Println("----------------------")
	}
	*/

	id, err := mspClient.GetSigningIdentity(userName)
	if err == nil {
		log.Infof("user exists: %s\n", userName)
		priFile, pubFile = f.GetKeyFile(id)
		return priFile,pubFile,errors.New("user exists")
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
		CAName: caName,
		Secret: secret,
		Attributes: attributes,
		Affiliation: orgName,
	}

	_, err = mspClient.Register(req)
	if err != nil && !strings.Contains(err.Error(), "is already registered") {
		log.Errorf("register %s [%s]\n", userName, err)
		return "","",err
	}

	err = mspClient.Enroll(userName,mspclient.WithSecret(secret))
	if err != nil {
		log.Errorf("Failed to enroll user: %s\n", err)
		return "","",err
	}

	id, err = mspClient.GetSigningIdentity(userName)
	if err != nil {
		log.Errorf("Failed to get signing identity : %s \n",err)
		return "","",err
	}
	priFile, pubFile = f.GetKeyFile(id)
	log.Infof("register %s successfully \n", userName)

	return priFile,pubFile,nil
}



func (f *FabricClient) GetKeyFile(id msp.SigningIdentity) (string,string){

	priFile := hex.EncodeToString(id.PrivateKey().SKI())+"_sk"
	pubFile := id.Identifier().ID + "@" +id.Identifier().MSPID + "-cert.pem"

	return priFile,pubFile
}

func (f *FabricClient) osCmd (channelId string) (string , error) {

	channelTx := "/usr/local/hyper/test2/configtx/channel-artifacts/"

	cmdStr := "cd /usr/local/hyper/test2" +
		"&&export PATH=${PATH}/../bin:${PWD}:$PATH" +
		"&&export FABRIC_CFG_PATH=${PWD}" +
		"&&configtxgen -profile TwoOrgsChannel -outputCreateChannelTx "+
		channelTx + "xxx_channel.tx -channelID xxx_channel"

	cmdStr = strings.ReplaceAll(cmdStr,"xxx_channel",channelId)
	cmd := exec.Command("/bin/bash","-c",cmdStr)
	stdout, _ := cmd.StdoutPipe()
	defer stdout.Close()

	if err := cmd.Start(); err != nil {
		log.Errorf("cmd.Start: %s \n",err)
		return "", err
	}

	time.Sleep(time.Duration(5)*time.Second)

	err := os.Chmod("/usr/local/hyper/test2/configtx/channel-artifacts/"+channelId+".tx", 0777)
	if err != nil {
		log.Errorf("Failed to os chmod file : %s \n",err)
		return "", err
	}

	log.Infof("cmd create channel tx file on cmd :%s \n",cmd.Args)
	channelTx = channelTx + channelId + ".tx"
	return  channelTx, nil
}

func (f *FabricClient) CreateChannel(org ,userName , channelId string) (string , error) {

	channelTx, err := f.osCmd(channelId)
	if err != nil {
		return "", err
	}

	mspClient , err := mspclient.New(f.sdk.Context(),mspclient.WithOrg(org))
	if err != nil {
		log.Errorf("Failed to new mspClient : %s \n",err)
		return "", err
	}

	adminidentity, err := mspClient.GetSigningIdentity(userName)
	if err != nil {
		log.Errorf("Failed to get signIdentity : %s \n",err)
		return "" , err
	}

	resmgmtClient , err := resmgmt.New(f.sdk.Context(fabsdk.WithUser(userName),fabsdk.WithOrg(org)))
	if err != nil {
		log.Errorf("Failed to create channel management client : %s \n",err)
		return "", err
	}

	req := resmgmt.SaveChannelRequest{
		ChannelID: channelId,
		ChannelConfigPath: channelTx,
		SigningIdentities: []msp.SigningIdentity{adminidentity},
	}

	resp , err := resmgmtClient.SaveChannel(req,resmgmt.WithOrdererEndpoint("orderer1-org0"))
	if err != nil {
		log.Errorf("Failed to save channel : %s \n",err)
		return "", err
	}

	return string(resp.TransactionID), nil
}

func (f *FabricClient) UpdateChannel(anchorsTx []string) error {

	for org, c := range f.resmgmtClients {

		mspClient, err := mspclient.New(f.sdk.Context(), mspclient.WithOrg(org))
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
			ChannelConfigPath: anchorsTx[0],
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

func (f *FabricClient) JoinChannel(channelId , user, org string) error {

	resmgmtClient , err := resmgmt.New(f.sdk.Context(fabsdk.WithUser(user),fabsdk.WithOrg(org)))
	if err != nil {
		log.Errorf("Failed to create channel management client : %s \n",err)
		return err
	}

	err = resmgmtClient.JoinChannel(channelId)
	if err != nil && !strings.Contains(err.Error(), "LedgerID already exists") {
		log.Errorf("Org peers failed to JoinChannel: %s \n", err)
		return err
	}

	log.Infof("%s join channel \n", org)

	return nil
}

func (f *FabricClient) CreateCC(chaincodeId, chaincodePath, version, org , userName, channelId string ) (string, error) {

	ccPkg, err := packager.NewCCPackage(chaincodePath, f.GoPath)
	if err != nil {
		log.Errorf("Failed to CreateCC :%s \n",err)
		return "", err
	}

	installCCReq := resmgmt.InstallCCRequest{
		Name: chaincodeId,
		Path: chaincodePath,
		Version: version,
		Package: ccPkg,
	}

	resmgmtClient , err := resmgmt.New(f.sdk.Context(fabsdk.WithUser(userName),fabsdk.WithOrg(org)))
	if err != nil {
		log.Errorf("Failed to create channel management client : %s \n",err)
		return "", err
	}

	_, err = resmgmtClient.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Errorf("Failed to create chaincode :%s \n",err)
		return "", err
	}

	ccPolicy := policydsl.SignedByAnyMember([]string{"org1MSP","org2MSP"})

	resp, err := resmgmtClient.InstantiateCC( channelId,
		resmgmt.InstantiateCCRequest{
			Name: chaincodeId,
			Path: chaincodePath,
			Version: version,
			Policy: ccPolicy,
		},
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
	)
	if err != nil {
		log.Errorf("Failed to instantiate chaincode :%s \n",err)
		return "", err
	}

	log.Infoln(resp.TransactionID)
	return string(resp.TransactionID), nil
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

	//	configBackend , _ := f.sdk.Config()
	//	targets, err := orgTargetPeers([]string{"org1"},configBackend )
	//	if err != nil {
	//		log.Printf("Failed to get targets:%s \n",err)
	//	}
	//	peer := targets[0]
	//
	//	installedCC, err := f.resmgmtClients[0].LifecycleQueryInstalledCC(resmgmt.WithTargetEndpoints(peer))
	//	if err != nil {
	//		return err
	//	}
	//for _, cc := range installedCC {
	//	fmt.Println(cc.Label)
	//	fmt.Println(cc.PackageID)
	//}

	//}

	return nil
}


func (f *FabricClient) GetOrgsTargetPeers(orgs []string) ([]string, error) {

	configBackend , _ := f.sdk.Config()
	networkConfig := fab.NetworkConfig{}

	err := lookup.New(configBackend).UnmarshalKey("organizations", &networkConfig.Organizations)
	if err != nil {
		log.Errorf("Failed to unmarsha org :%s \n",err)
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

func (f *FabricClient) GetOrgTargetPeers(org string) ([]string, error) {

	configBackend , _ := f.sdk.Config()
	networkConfig := fab.NetworkConfig{}

	err := lookup.New(configBackend).UnmarshalKey("organizations", &networkConfig.Organizations)
	if err != nil {
		log.Errorf("Failed to unmarsha org :%s \n",err)
		return nil, err
	}

	var peers []string
	orgConfig, ok := networkConfig.Organizations[strings.ToLower(org)]
	if !ok {
		log.Errorf("%s dont exeits",org)
		return nil, errors.New("org dont in list")
	}
	peers = append(peers, orgConfig.Peers...)

	return peers, nil
}

//func (f *FabricClient) PackageCC() (string, []byte) {
//	desc := &lcpackager.Descriptor{
//		Path:  "/usr/local/soft/fabric-test5/chaincode/newchaincode/test",
//		Type:  pb.ChaincodeSpec_GOLANG,
//		Label: "newcc0",
//	}
//	ccPkg, err := lcpackager.NewCCPackage(desc)
//	if err != nil {
//		log.Errorf()("Failed to package chaincode : %s \n",err)
//	}
//	return desc.Label, ccPkg
//}


func (f *FabricClient) InstallCC(label string, ccPkg []byte , org string) {
	installCCReq := resmgmt.LifecycleInstallCCRequest{
		Label:   label,
		Package: ccPkg,
	}

	packageID := lcpackager.ComputePackageID(installCCReq.Label, installCCReq.Package)

	resp, err := f.resmgmtClients[org].LifecycleInstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Errorf("Failed to install chaincode : %s \n",err)
	}
	fmt.Println(packageID)
	fmt.Println(resp[0].PackageID)
}

func (f *FabricClient) GetInstalledCCPackage(packageID string , org string) {
	_, err := f.resmgmtClients[org].LifecycleGetInstalledCCPackage(packageID,resmgmt.WithTargetEndpoints(peer1), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Errorf("Failed to GetInstalledCCPackage chaincode : %s \n",err)
	}
	//fmt.Println(resp)
}

func (f *FabricClient) QueryInstalled(label string, packageID string , org string) {
	resp, err := f.resmgmtClients[org].LifecycleQueryInstalledCC(resmgmt.WithTargetEndpoints(peer2), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Errorf("Failed to QueryInstalled chaincode : %s \n",err)
	}
	fmt.Println(resp[0].PackageID)
	fmt.Println(resp[0].Label)
}

func  (f *FabricClient) ApproveCC( packageID string , org string) {
	//queryApprovedCCReq := resmgmt.LifecycleQueryApprovedCCRequest{
	//	Name:       "newcc0",
	//	Sequence: 1,
	//}
	//resp, err := f.resmgmtClients[0].LifecycleQueryApprovedCC("mychannel", queryApprovedCCReq, resmgmt.WithTargetEndpoints(peer1), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(resp.PackageID)
	//
	//
	ccPolicy := policydsl.SignedByAnyMember([]string{"org1MSP"})

	approveCCReq := resmgmt.LifecycleApproveCCRequest{
		Name:              "newcc0",
		Version:           "0",
		PackageID:         packageID,
		Sequence:          1,
		EndorsementPlugin: "escc",
		ValidationPlugin:  "vscc",
		SignaturePolicy:   ccPolicy,
		InitRequired:      true,
	}
	fmt.Println("!!!")

	txnID, err := f.resmgmtClients[org].LifecycleApproveCC("mychannel", approveCCReq, resmgmt.WithTargetEndpoints(peer1), resmgmt.WithOrdererEndpoint("orderer1-org0"))
	fmt.Println("???")
	if err != nil {
		log.Errorf("Failed to ApproveCC chaincode : %s \n",err)
	}

	fmt.Println(txnID)

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




func (f *FabricClient) QueryLedger() error {

	ledger, err := ledger.New(f.sdk.ChannelContext(f.ChannelId, fabsdk.WithUser("User3"), fabsdk.WithOrg("org1")))
	if err != nil {
		return err
	}

	bci, err := ledger.QueryInfo()
	if err != nil {
		return err
	}
	fmt.Println(bci.BCI.Height)

	return nil
}

