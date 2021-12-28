package models

type Block struct {
	Number          uint64         `json:"number"`          //区块号
	PreviousHash    string         `json:"previousHash"`    //前区块Hash
	DataHash        string         `json:"dataHash"`        //交易体Hash
	BlockHash       string         `json:"blockHash"`       //区块Hash
	TxNum           int            `json:"txNum"`           //区块内交易个数
	TransactionList []*Transaction `json:"transactionList"` //交易列表
	CreateTime      string         `json:"createTime"`      //区块生成时间
}

type BlockHeader struct {
	Number       int8   `json:"number"`       //区块号
	PreviousHash []byte `json:"previousHash"` //前区块Hash
	DataHash     []byte `json:"dataHash"`     //交易体Hash
}

type Transaction struct {
	TransactionActionList []*TransactionAction `json:"transactionActionList"` //交易列表
}

type TransactionAction struct {
	TxId         string   `json:"txId"`         //交易ID
	BlockNum     uint64   `json:"blockNum"`     //区块号
	Type         string   `json:"type"`         //交易类型
	Timestamp    string   `json:"timestamp"`    //交易创建时间
	ChannelId    string   `json:"channelId"`    //通道ID
	Endorsements []string `json:"endorsements"` //背书组织ID列表
	ChaincodeId  string   `json:"chaincodeId"`  //链代码名称
	ReadSetList  []string `json:"readSetList"`  //读集
	WriteSetList []string `json:"writeSetList"` //写集
}

type BlockMainInfo struct {
	BlockNum       uint64 `json:"blockNum"`       // the number of blocks
	TransactionNum uint64 `json:"transactionNum"` // the number of transactions
	ChaincodeNum   uint64 `json:"chaincodeNum"`   // the number of chain-codes
	NodeNum        uint64 `json:"nodeNum"`        // the number of nodes
}
