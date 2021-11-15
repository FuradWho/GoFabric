package models

type User struct {
	UserName string `json:"user_name,omitempty"`
	Secret string `json:"secret,omitempty"`
	UserType string `json:"user_type,omitempty"`
	OrgName string `json:"org_name,omitempty"`
	CaName string `json:"ca_name,omitempty"`
}

type ResponseBean struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

type UserData struct {
	PriFile string `json:"pri_file"`
	PubFile string `json:"pub_file"`
}

type CreateChannelInfo struct {

	ChannelTx string `json:"channel_tx"`
	Org string `json:"org"`
	UserName string `json:"user_name"`
	ChannelId string `json:"channel_id"`

}


func SuccessData(data interface{}) *ResponseBean{
	return &ResponseBean{
		200,
		"success",
		data,
	}
}

func SuccessMsg(msg string) *ResponseBean  {
	return &ResponseBean{
		200,
		msg,
		"",
	}
}

func FailedMsg(errMsg string) *ResponseBean  {
	return &ResponseBean{
		400,
		errMsg,
		"",
	}
}

func FailedData(errMsg string,data interface{}) *ResponseBean  {
	return &ResponseBean{
		400,
		errMsg,
		data,
	}
}


