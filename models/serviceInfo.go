package models

type User struct {
	UserName string `json:"user_name" validate:"required,gt=0"`
	Secret   string `json:"secret" validate:"required,gt=0"`
	UserType string `json:"user_type" validate:"required,gt=0"`
	OrgName  string `json:"org_name" validate:"required,gt=0"`
	CaName   string `json:"ca_name" validate:"required,gt=0"`
}

type ResponseBean struct {
	Code int         `json:"code" validate:"required,gt=0"`
	Msg  string      `json:"msg" validate:"required,gt=0"`
	Data interface{} `json:"data" validate:"required,gt=0"`
}

type UserData struct {
	PriFile string `json:"pri_file" validate:"required,gt=0"`
	PubFile string `json:"pub_file" validate:"required,gt=0"`
}

type CreateChannelInfo struct {
	Org       string `json:"org" validate:"required,gt=0"`
	UserName  string `json:"user_name" validate:"required,gt=0"`
	ChannelId string `json:"channel_id" validate:"required,gt=0"`
	Orderer   string `json:"orderer" validate:"required,gt=0"`
}

type JoinChannelInfo struct {
	Org       string `json:"org" validate:"required,gt=0"`
	UserName  string `json:"user_name" validate:"required,gt=0"`
	ChannelId string `json:"channel_id" validate:"required,gt=0"`
}

type CcInfo struct {
	PackageId     string `json:"package_id"`
	ChaincodeId   string `json:"chaincode_id"`
	ChaincodePath string `json:"chaincode_path"`
	Version       string `json:"version"`
	Org           string `json:"org"`
	UserName      string `json:"user_name"`
	ChannelId     string `json:"channel_id"`
	Peer          string `json:"peer"`
	Orderer       string `json:"orderer"`
	Sequence      string `json:"sequence"`
}

type QueryApprovedCCInfo struct {
	ChaincodeId string `json:"chaincode_id" validate:"required,gt=0"`
	Org         string `json:"org" validate:"required,gt=0"`
	UserName    string `json:"user_name" validate:"required,gt=0"`
	ChannelId   string `json:"channel_id" validate:"required,gt=0"`
	Peer        string `json:"peer" validate:"required,gt=0"`
	Sequence    string `json:"sequence" validate:"required,gt=0"`
}

type ApproveCCInfo struct {
	PackageId   string `json:"package_id" validate:"required,gt=0"`
	ChaincodeId string `json:"chaincode_id" validate:"required,gt=0"`
	//ChaincodePath string `json:"chaincode_path" validate:"required,gt=0"`
	Version   string `json:"version" validate:"required,gt=0"`
	Org       string `json:"org" validate:"required,gt=0"`
	UserName  string `json:"user_name" validate:"required,gt=0"`
	ChannelId string `json:"channel_id" validate:"required,gt=0"`
	Peer      string `json:"peer" validate:"required,gt=0"`
	Orderer   string `json:"orderer" validate:"required,gt=0"`
	Sequence  string `json:"sequence" validate:"required,gt=0"`
}

type CheckCCCommitReadinessInfo struct {
	// PackageId   string `json:"package_id" validate:"required,gt=0"`
	ChaincodeId string `json:"chaincode_id" validate:"required,gt=0"`
	// ChaincodePath string `json:"chaincode_path" validate:"required,gt=0"`
	Version   string `json:"version" validate:"required,gt=0"`
	Org       string `json:"org" validate:"required,gt=0"`
	UserName  string `json:"user_name" validate:"required,gt=0"`
	ChannelId string `json:"channel_id" validate:"required,gt=0"`
	Peer      string `json:"peer" validate:"required,gt=0"`
	// Orderer   string `json:"orderer" validate:"required,gt=0"`
	Sequence string `json:"sequence" validate:"required,gt=0"`
}

type CommitCCInfo struct {
	// PackageId   string `json:"package_id" validate:"required,gt=0"`
	ChaincodeId string `json:"chaincode_id" validate:"required,gt=0"`
	// ChaincodePath string `json:"chaincode_path" validate:"required,gt=0"`
	Version   string `json:"version" validate:"required,gt=0"`
	Org       string `json:"org" validate:"required,gt=0"`
	UserName  string `json:"user_name" validate:"required,gt=0"`
	ChannelId string `json:"channel_id" validate:"required,gt=0"`
	Peer      string `json:"peer" validate:"required,gt=0"`
	Orderer   string `json:"orderer" validate:"required,gt=0"`
	Sequence  string `json:"sequence" validate:"required,gt=0"`
}

type RequestApproveCCByOtherInfo struct {
	PackageId   string `json:"package_id" validate:"required,gt=0"`
	ChaincodeId string `json:"chaincode_id" validate:"required,gt=0"`
	// ChaincodePath string `json:"chaincode_path" validate:"required,gt=0"`
	Version   string `json:"version" validate:"required,gt=0"`
	Org       string `json:"org" validate:"required,gt=0"`
	UserName  string `json:"user_name" validate:"required,gt=0"`
	ChannelId string `json:"channel_id" validate:"required,gt=0"`
	Peer      string `json:"peer" validate:"required,gt=0"`
	Orderer   string `json:"orderer" validate:"required,gt=0"`
	Sequence  string `json:"sequence" validate:"required,gt=0"`
}

type RequestInstallCCByOtherInfo struct {

	//PackageId   string `json:"package_id" validate:"required,gt=0"`
	ChaincodeId   string `json:"chaincode_id" validate:"required,gt=0"`
	ChaincodePath string `json:"chaincode_path" validate:"required,gt=0"`
	//Version   string `json:"version" validate:"required,gt=0"`
	Org      string `json:"org" validate:"required,gt=0"`
	UserName string `json:"user_name" validate:"required,gt=0"`
	//ChannelId string `json:"channel_id" validate:"required,gt=0"`
	Peer string `json:"peer" validate:"required,gt=0"`
	//Orderer   string `json:"orderer" validate:"required,gt=0"`
	//Sequence  string `json:"sequence" validate:"required,gt=0"`
}

type QueryInstalledInfo struct {

	//PackageId   string `json:"package_id" validate:"required,gt=0"`
	//ChaincodeId   string `json:"chaincode_id" validate:"required,gt=0"`
	//ChaincodePath string `json:"chaincode_path" validate:"required,gt=0"`
	//Version   string `json:"version" validate:"required,gt=0"`
	Org      string `json:"org" validate:"required,gt=0"`
	UserName string `json:"user_name" validate:"required,gt=0"`
	//ChannelId string `json:"channel_id" validate:"required,gt=0"`
	Peer string `json:"peer" validate:"required,gt=0"`
	//Orderer   string `json:"orderer" validate:"required,gt=0"`
	//Sequence  string `json:"sequence" validate:"required,gt=0"`
}

type InstallCCInfo struct {

	//PackageId   string `json:"package_id" validate:"required,gt=0"`
	ChaincodeId   string `json:"chaincode_id" validate:"required,gt=0"`
	ChaincodePath string `json:"chaincode_path" validate:"required,gt=0"`
	//Version   string `json:"version" validate:"required,gt=0"`
	Org      string `json:"org" validate:"required,gt=0"`
	UserName string `json:"user_name" validate:"required,gt=0"`
	//ChannelId string `json:"channel_id" validate:"required,gt=0"`
	Peer string `json:"peer" validate:"required,gt=0"`
	//Orderer   string `json:"orderer" validate:"required,gt=0"`
	//Sequence  string `json:"sequence" validate:"required,gt=0"`
}

func SuccessData(data interface{}) *ResponseBean {
	return &ResponseBean{
		200,
		"success",
		data,
	}
}

func SuccessMsg(msg string) *ResponseBean {
	return &ResponseBean{
		200,
		msg,
		"",
	}
}

func FailedMsg(errMsg string) *ResponseBean {
	return &ResponseBean{
		400,
		errMsg,
		"",
	}
}

func FailedData(errMsg string, data interface{}) *ResponseBean {
	return &ResponseBean{
		400,
		errMsg,
		data,
	}
}
