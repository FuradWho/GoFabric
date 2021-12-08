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

	pb "github.com/hyperledger/fabric-protos-go/peer"
)

const (
	Admin = "Admin"
	peer1 = "peer0.org1.example.com"
	peer2 = "peer0.org2.example.com"
)

var log = logrus.New()

type FabricClient struct {
	ConnectionFile []byte
	Orgs           []string
	ChannelId      string
	GoPath         string

	userName string
	userOrg  string

	resmgmtClients map[string]*resmgmt.Client
	sdk            *fabsdk.FabricSDK
	retry          resmgmt.RequestOption
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

	sdk, err := fabsdk.New(config.FromRaw(f.ConnectionFile, "yaml"))
	if err != nil {
		log.Error("Failed to setup main sdk ")
		return err
	}
	f.sdk = sdk

	resmgmtClients := make(map[string]*resmgmt.Client)
	for _, v := range f.Orgs {
		resmgmtClient, err := resmgmt.New(sdk.Context(fabsdk.WithUser(Admin), fabsdk.WithOrg(v)))
		if err != nil {
			log.Errorf("Failed to create channel management client : %s \n", err)
			return err
		}
		resmgmtClients[v] = resmgmtClient
	}

	f.resmgmtClients = resmgmtClients
	f.retry = resmgmt.WithRetry(retry.DefaultResMgmtOpts)

	return nil
}

func (f *FabricClient) Close() {
	if f.sdk != nil {
		f.sdk.Close()
	}
}

func (f *FabricClient) SetUser(userName, userOrg string) {

	f.userName = userName
	f.userOrg = userOrg

}

func (f *FabricClient) CreateUser(userName string, secret string, userType string, orgName string, caName string) (priFile string, pubFile string, err error) {

	ctx := f.sdk.Context()
	mspClient, err := mspclient.New(ctx, mspclient.WithOrg(orgName))
	if err != nil {
		log.Errorf("Failed to create msp client: %s\n", err)
		return "", "", err
	}

	caInfo, err := mspClient.GetCAInfo()
	if err != nil {
		log.Errorf("Failed to get CA Info :%s \n", err)
		return "", "", err
	}

	if caInfo.CAName != caName {
		log.Errorln("Not match ca ")
		return "", "", errors.New("Not match ca")
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
		return priFile, pubFile, errors.New("user exists")
	}

	a1 := mspclient.Attribute{
		Name:  "hf.Registrar.Roles",
		Value: "client,orderer,peer,user",
	}
	a2 := mspclient.Attribute{
		Name:  "hf.Registrar.DelegateRoles",
		Value: "client,orderer,peer,user",
	}
	a3 := mspclient.Attribute{
		Name:  "hf.Registrar.Attributes",
		Value: "*",
	}
	a4 := mspclient.Attribute{
		Name:  "hf.GenCRL",
		Value: "true",
	}
	a5 := mspclient.Attribute{
		Name:  "hf.Revoker",
		Value: "true",
	}
	a6 := mspclient.Attribute{
		Name:  "hf.AffiliationMgr",
		Value: "true",
	}
	a7 := mspclient.Attribute{
		Name:  "hf.IntermediateCA",
		Value: "true",
	}

	var attributes []mspclient.Attribute
	attributes = append(attributes, a1, a2, a3, a4, a5, a6, a7)

	req := &mspclient.RegistrationRequest{
		Name:        userName,
		Type:        userType,
		CAName:      caName,
		Secret:      secret,
		Attributes:  attributes,
		Affiliation: orgName,
	}

	_, err = mspClient.Register(req)
	if err != nil && !strings.Contains(err.Error(), "is already registered") {
		log.Errorf("register %s [%s]\n", userName, err)
		return "", "", err
	}

	err = mspClient.Enroll(userName, mspclient.WithSecret(secret))
	if err != nil {
		log.Errorf("Failed to enroll user: %s\n", err)
		return "", "", err
	}

	id, err = mspClient.GetSigningIdentity(userName)
	if err != nil {
		log.Errorf("Failed to get signing identity : %s \n", err)
		return "", "", err
	}
	priFile, pubFile = f.GetKeyFile(id)
	log.Infof("register %s successfully \n", userName)

	return priFile, pubFile, nil
}

func (f *FabricClient) AuthenticateUser(orgName string) {

	ctx := f.sdk.Context()
	mspClient, err := mspclient.New(ctx, mspclient.WithOrg(orgName))
	if err != nil {
		log.Errorf("Failed to create msp client: %s\n", err)
		return
	}

	caInfo, err := mspClient.GetCAInfo()
	if err != nil {
		log.Errorf("Failed to get CA Info :%s \n", err)
		return
	}

	log.Infoln(caInfo.CAName)

	affiliations, err := mspClient.GetAllIdentities()
	if err != nil {
		log.Printf("%s \n", err)
	}

	for _, info := range affiliations {
		fmt.Println(info.ID)
		fmt.Println(info.Type)
		fmt.Println(info.Attributes)
		fmt.Println("----------------------")
	}

}

func (f *FabricClient) GetKeyFile(id msp.SigningIdentity) (string, string) {

	priFile := hex.EncodeToString(id.PrivateKey().SKI()) + "_sk"
	pubFile := id.Identifier().ID + "@" + id.Identifier().MSPID + "-cert.pem"

	return priFile, pubFile
}

func (f *FabricClient) osCmd(channelId string) (string, error) {

	channelTx := "/usr/local/hyper/test5/channel-artifacts/"

	cmdStr := "cd /usr/local/hyper/test5" +
		"&&export PATH=${PATH}/../bin:${PWD}:$PATH" +
		"&&export FABRIC_CFG_PATH=${PWD}" +
		"&&configtxgen -profile TwoOrgsChannels -outputCreateChannelTx " +
		channelTx + "xxx_channel.tx -channelID xxx_channel"

	cmdStr = strings.ReplaceAll(cmdStr, "xxx_channel", channelId)
	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	stdout, _ := cmd.StdoutPipe()
	defer stdout.Close()

	if err := cmd.Start(); err != nil {
		log.Errorf("cmd.Start: %s \n", err)
		return "", err
	}

	time.Sleep(time.Duration(5) * time.Second)

	err := os.Chmod("/usr/local/hyper/test5/channel-artifacts/"+channelId+".tx", 0777)
	if err != nil {
		log.Errorf("Failed to os chmod file : %s \n", err)
		return "", err
	}

	log.Infof("cmd create channel tx file on cmd :%s \n", cmd.Args)
	channelTx = channelTx + channelId + ".tx"
	return channelTx, nil
}

func (f *FabricClient) CreateChannel(org, userName, channelId, orderer string) (string, error) {

	channelTx, err := f.osCmd(channelId)
	if err != nil {
		return "", err
	}

	mspClient, err := mspclient.New(f.sdk.Context(), mspclient.WithOrg(org))
	if err != nil {
		log.Errorf("Failed to new mspClient : %s \n", err)
		return "", err
	}

	adminidentity, err := mspClient.GetSigningIdentity(userName)
	if err != nil {
		log.Errorf("Failed to get signIdentity : %s \n", err)
		return "", err
	}

	resmgmtClient, err := resmgmt.New(f.sdk.Context(fabsdk.WithUser(userName), fabsdk.WithOrg(org)))
	if err != nil {
		log.Errorf("Failed to create channel management client : %s \n", err)
		return "", err
	}

	req := resmgmt.SaveChannelRequest{
		ChannelID:         channelId,
		ChannelConfigPath: channelTx,
		SigningIdentities: []msp.SigningIdentity{adminidentity},
	}

	resp, err := resmgmtClient.SaveChannel(req, resmgmt.WithOrdererEndpoint(orderer))
	if err != nil {
		log.Errorf("Failed to save channel : %s \n", err)
		return "", err
	}

	log.Infof("channel resp : %+v \n", resp)

	return string(resp.TransactionID), nil
}

func (f *FabricClient) UpdateChannel(anchorsTx []string) error {

	for org, c := range f.resmgmtClients {

		mspClient, err := mspclient.New(f.sdk.Context(), mspclient.WithOrg(org))
		if err != nil {
			log.Printf("Failed to UpdateChannel : %s \n", err)
			return err
		}
		adminIdentity, err := mspClient.GetSigningIdentity(Admin)
		if err != nil {
			log.Printf("Failed to UpdateChannel : %s \n", err)
			return err
		}
		req := resmgmt.SaveChannelRequest{
			ChannelID:         f.ChannelId,
			ChannelConfigPath: anchorsTx[0],
			SigningIdentities: []msp.SigningIdentity{adminIdentity},
		}
		txId, err := c.SaveChannel(req, f.retry)
		if err != nil {
			log.Printf("Failed to UpdateChannel : %s \n", err)
			return err
		}
		log.Printf("Failed to UpdateChannel : %s \n", txId)
	}

	return nil
}

func (f *FabricClient) JoinChannel(channelId, user, org string) error {

	mspClient, err := mspclient.New(f.sdk.Context(), mspclient.WithOrg(org))
	if err != nil {
		log.Errorf("Failed to new mspClient : %s \n", err)
		return err
	}

	adminidentity, err := mspClient.GetSigningIdentity(user)
	if err != nil {
		log.Errorf("Failed to get signIdentity : %s \n", err)
		return err
	}

	log.Infoln(string(adminidentity.PrivateKey().SKI()))

	resmgmtClient, err := resmgmt.New(f.sdk.Context(fabsdk.WithUser(user), fabsdk.WithOrg(org)))
	if err != nil {
		log.Errorf("Failed to create channel management client : %s \n", err)
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

func (f *FabricClient) CreateCC(chaincodeId, chaincodePath, version, org, userName, channelId string) (string, error) {

	ccPkg, err := packager.NewCCPackage(chaincodePath, f.GoPath)
	if err != nil {
		log.Errorf("Failed to CreateCC :%s \n", err)
		return "", err
	}

	installCCReq := resmgmt.InstallCCRequest{
		Name:    chaincodeId,
		Path:    chaincodePath,
		Version: version,
		Package: ccPkg,
	}

	resmgmtClient, err := resmgmt.New(f.sdk.Context(fabsdk.WithUser(userName), fabsdk.WithOrg(org)))
	if err != nil {
		log.Errorf("Failed to create channel management client : %s \n", err)
		return "", err
	}

	_, err = resmgmtClient.InstallCC(installCCReq)
	if err != nil {
		log.Errorf("Failed to create chaincode :%s \n", err)
		return "", err
	}

	ccPolicy := policydsl.SignedByAnyMember([]string{"org1MSP"})

	initArgs := [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")}
	resp, err := resmgmtClient.InstantiateCC(channelId,
		resmgmt.InstantiateCCRequest{
			Name:    chaincodeId,
			Path:    chaincodePath,
			Version: version,
			Args:    initArgs,
			Policy:  ccPolicy,
		},
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
	)
	if err != nil {
		log.Errorf("Failed to instantiate chaincode :%s \n", err)
		return "", err
	}

	log.Infoln(resp.TransactionID)
	return string(resp.TransactionID), nil
}

func (f *FabricClient) GetInstalledCCPackage(user, packageID, org string) error {

	resmgmtClient, err := resmgmt.New(f.sdk.Context(fabsdk.WithUser(user), fabsdk.WithOrg(org)))
	if err != nil {
		log.Errorf("Failed to create channel management client : %s \n", err)
		return err
	}

	resp, err := resmgmtClient.LifecycleGetInstalledCCPackage(packageID, resmgmt.WithTargetEndpoints(peer1), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Errorf("Failed to GetInstalledCCPackage chaincode : %s \n", err)
		return err
	}
	log.Infoln(resp)
	return nil
}

func (f *FabricClient) QueryInstalled(user, org, peer string) ([]resmgmt.LifecycleInstalledCC, error) {

	resmgmtClient, err := resmgmt.New(f.sdk.Context(fabsdk.WithUser(user), fabsdk.WithOrg(org)))
	if err != nil {
		log.Errorf("Failed to create channel management client : %s \n", err)
		return nil, err
	}

	resp, err := resmgmtClient.LifecycleQueryInstalledCC(resmgmt.WithTargetEndpoints(peer), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Errorf("Failed to QueryInstalled chaincode : %s \n", err)
		return nil, err
	}
	return resp, nil
}

func (f *FabricClient) InstallCC(chaincodeId, chaincodePath, org, user, peer string) (string, error) {

	desc := &lcpackager.Descriptor{
		Path:  chaincodePath,
		Type:  pb.ChaincodeSpec_GOLANG,
		Label: chaincodeId,
	}

	ccPkg, err := lcpackager.NewCCPackage(desc)
	if err != nil {
		log.Errorf("Failed to NewCCPackage client : %s \n", err)
		return "", err
	}

	installCCReq := resmgmt.LifecycleInstallCCRequest{
		Label:   chaincodeId,
		Package: ccPkg,
	}

	packageID := lcpackager.ComputePackageID(installCCReq.Label, installCCReq.Package)

	resmgmtClient, err := resmgmt.New(f.sdk.Context(fabsdk.WithUser(user), fabsdk.WithOrg(org)))
	if err != nil {
		log.Errorf("Failed to create channel management client : %s \n", err)
		return "", err
	}

	resp, err := resmgmtClient.LifecycleInstallCC(installCCReq, resmgmt.WithTargetEndpoints(peer), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Errorf("Failed to install chaincode : %s \n", err)
		return "", err
	}
	log.Infof("resp[0].PackageID : %s  and packageID : %s  \n", resp[0].PackageID, packageID)
	return resp[0].PackageID, nil
}

func (f *FabricClient) ApproveCC(packageID, chaincodeId, version, channelId, user, org, peer, orderer string, sequence int) (fab.TransactionID, error) {

	ccPolicy := policydsl.SignedByAnyMember([]string{"Org1MSP", "Org2MSP"})

	approveCCReq := resmgmt.LifecycleApproveCCRequest{
		Name:              chaincodeId,
		Version:           version,
		PackageID:         packageID,
		Sequence:          int64(sequence),
		EndorsementPlugin: "escc",
		ValidationPlugin:  "vscc",
		SignaturePolicy:   ccPolicy,
		InitRequired:      true,
	}

	log.Infoln("!!!")

	resmgmtClient, err := resmgmt.New(f.sdk.Context(fabsdk.WithUser(user), fabsdk.WithOrg(org)))
	if err != nil {
		log.Errorf("Failed to create channel management client : %s \n", err)
		return "", err
	}

	txnID, err := resmgmtClient.LifecycleApproveCC(channelId, approveCCReq, resmgmt.WithOrdererEndpoint(orderer), resmgmt.WithTargetEndpoints(peer))

	log.Infoln("???")
	if err != nil {
		log.Errorf("Failed to ApproveCC chaincode : %s \n", err)
		return "", err
	}
	log.Infoln(txnID)
	return txnID, nil
}

func (f *FabricClient) QueryApprovedCC(ccID, user, org, channelId, peer string, sequence int) (string, error) {
	queryApprovedCCReq := resmgmt.LifecycleQueryApprovedCCRequest{
		Name:     ccID,
		Sequence: int64(sequence),
	}

	resmgmtClient, err := resmgmt.New(f.sdk.Context(fabsdk.WithUser(user), fabsdk.WithOrg(org)))
	if err != nil {
		log.Errorf("Failed to create channel management client : %s \n", err)
		return "", err
	}

	resp, err := resmgmtClient.LifecycleQueryApprovedCC(channelId, queryApprovedCCReq, resmgmt.WithTargetEndpoints(peer), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Errorf("Failed to LifecycleQueryApprovedCC : %s \n", err)
		return "", err
	}
	log.Infoln(resp.PackageID)
	return resp.PackageID, nil
}

func (f *FabricClient) CheckCCCommitReadiness(ccID, version, user, org, channelId, peer string, sequence int) (map[string]bool, error) {
	ccPolicy := policydsl.SignedByAnyMember([]string{"Org1MSP", "Org2MSP"})
	req := resmgmt.LifecycleCheckCCCommitReadinessRequest{
		Name:              ccID,
		Version:           version,
		EndorsementPlugin: "escc",
		ValidationPlugin:  "vscc",
		SignaturePolicy:   ccPolicy,
		Sequence:          int64(sequence),
		InitRequired:      true,
	}

	resmgmtClient, err := resmgmt.New(f.sdk.Context(fabsdk.WithUser(user), fabsdk.WithOrg(org)))
	if err != nil {
		log.Errorf("Failed to create channel management client : %s \n", err)
		return nil, err
	}

	resp, err := resmgmtClient.LifecycleCheckCCCommitReadiness(channelId, req, resmgmt.WithTargetEndpoints(peer), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Errorf("Failed to LifecycleCheckCCCommitReadiness : %s \n", err)
		return nil, err
	}
	log.Infof("%+v \n", resp.Approvals)
	return resp.Approvals, nil
}

func (f *FabricClient) CommitCC(ccID, user, org, channelId, orderer, version string, sequence int) (fab.TransactionID, error) {
	ccPolicy := policydsl.SignedByAnyMember([]string{"Org1MSP", "Org2MSP"})

	req := resmgmt.LifecycleCommitCCRequest{
		Name:              ccID,
		Version:           version,
		Sequence:          int64(sequence),
		EndorsementPlugin: "escc",
		ValidationPlugin:  "vscc",
		SignaturePolicy:   ccPolicy,
		InitRequired:      true,
	}

	resmgmtClient, err := resmgmt.New(f.sdk.Context(fabsdk.WithUser(user), fabsdk.WithOrg(org)))
	if err != nil {
		log.Errorf("Failed to create channel management client : %s \n", err)
		return "", err
	}

	txnID, err := resmgmtClient.LifecycleCommitCC(channelId, req, resmgmt.WithOrdererEndpoint(orderer))
	if err != nil {
		log.Errorf("Failed to LifecycleCommitCC : %s \n", err)
		return "", err
	}
	log.Infof("%+v \n", txnID)
	return txnID, nil
}

func (f *FabricClient) QueryCommittedCC(ccID, user, org, channelId, peer string) error {
	req := resmgmt.LifecycleQueryCommittedCCRequest{
		Name: ccID,
	}

	resmgmtClient, err := resmgmt.New(f.sdk.Context(fabsdk.WithUser(user), fabsdk.WithOrg(org)))
	if err != nil {
		log.Errorf("Failed to create channel management client : %s \n", err)
		return err
	}
	resp, err := resmgmtClient.LifecycleQueryCommittedCC(channelId, req, resmgmt.WithTargetEndpoints(peer), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Errorf("Failed to LifecycleQueryCommittedCC : %s \n", err)
		return err
	}
	log.Infoln(resp[0].Name)
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

	configBackend, _ := f.sdk.Config()
	networkConfig := fab.NetworkConfig{}

	err := lookup.New(configBackend).UnmarshalKey("organizations", &networkConfig.Organizations)
	if err != nil {
		log.Errorf("Failed to unmarsha org :%s \n", err)
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

	configBackend, _ := f.sdk.Config()
	networkConfig := fab.NetworkConfig{}

	err := lookup.New(configBackend).UnmarshalKey("organizations", &networkConfig.Organizations)
	if err != nil {
		log.Errorf("Failed to unmarsha org :%s \n", err)
		return nil, err
	}

	var peers []string
	orgConfig, ok := networkConfig.Organizations[strings.ToLower(org)]
	if !ok {
		log.Errorf("%s dont exeits", org)
		return nil, errors.New("org dont in list")
	}
	peers = append(peers, orgConfig.Peers...)

	return peers, nil
}

func (f *FabricClient) GetNetworkConfig() (fab.NetworkConfig, error) {

	configBackend, _ := f.sdk.Config()
	networkConfig := fab.NetworkConfig{}

	err := lookup.New(configBackend).UnmarshalKey("organizations", &networkConfig.Organizations)
	if err != nil {
		log.Errorf("Failed to unmarsha org :%s \n", err)
		return networkConfig, err
	}

	err = lookup.New(configBackend).UnmarshalKey("orderers", &networkConfig.Orderers)
	if err != nil {
		log.Errorf("Failed to unmarsha org :%s \n", err)
		return networkConfig, err
	}

	err = lookup.New(configBackend).UnmarshalKey("channels", &networkConfig.Channels)
	if err != nil {
		log.Errorf("Failed to unmarsha org :%s \n", err)
		return networkConfig, err
	}

	err = lookup.New(configBackend).UnmarshalKey("peers", &networkConfig.Peers)
	if err != nil {
		log.Errorf("Failed to unmarsha org :%s \n", err)
		return networkConfig, err
	}

	return networkConfig, nil
}

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

func (f *FabricClient) QueryConfigBlockFromOrder(user, org, channelId, orderer string) error {

	resmgmtClient, err := resmgmt.New(f.sdk.Context(fabsdk.WithUser(user), fabsdk.WithOrg(org)))
	if err != nil {
		log.Errorf("Failed to create channel management client : %s \n", err)
		return err
	}

	fromOrderer, err := resmgmtClient.QueryConfigFromOrderer(channelId, resmgmt.WithOrdererEndpoint(orderer))
	if err != nil {
		log.Errorf("Failed to QueryConfigFromOrderer : %s \n", err)
		return err
	}
	log.Infof("Config Block : %+v \n", fromOrderer)
	return nil

}

//func (f *FabricClient)  InitCC(ccID , user , org ,channelId ,peer string) {
//	//prepare channel client context using client context
//	clientChannelContext := f.sdk.ChannelContext(channelId, fabsdk.WithUser(user), fabsdk.WithOrg(org))
//	// Channel client is used to query and execute transactions (Org1 is default org)
//	client, err := channel.New(clientChannelContext)
//	if err != nil {
//		t.Fatalf("Failed to create new channel client: %s", err)
//	}
//
//	// init
//	_, err = client.Execute(channel.Request{ChaincodeID: ccID, Fcn: "init", Args: integration.ExampleCCInitArgsLc(), IsInit: true},
//		channel.WithRetry(retry.DefaultChannelOpts))
//	if err != nil {
//		t.Fatalf("Failed to init: %s", err)
//	}
//}
