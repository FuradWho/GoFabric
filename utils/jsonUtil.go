package utils

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/ledger/rwset"
	"github.com/hyperledger/fabric-protos-go/ledger/rwset/kvrwset"
	"github.com/hyperledger/fabric-protos-go/msp"
	"github.com/hyperledger/fabric-protos-go/peer"
	"gofabric/models"
	"log"
	"time"
)

func GetEnvelopeFromBlock(data []byte) (*common.Envelope, error) {

	env := &common.Envelope{}

	if err := proto.Unmarshal(data, env); err != nil {
		log.Printf("Failed to block unmarshal :%s \n", err)
	}

	return env, nil
}

func GetTxFromEnvelopeDeep(rawEnvelope *common.Envelope) (*models.Transaction, error) {

	rawPayload := &common.Payload{}
	err := proto.Unmarshal(rawEnvelope.Payload, rawPayload)
	if err != nil {
		log.Printf("Failed to unmarshal Envelope's payload :%s \n", err)
	}

	channelHeader := &common.ChannelHeader{}
	err = proto.Unmarshal(rawPayload.Header.ChannelHeader, channelHeader)
	if err != nil {
		log.Printf("Failed to unmarshal Payload's ChannelHeader :%s \n", err)
	}

	signatureHeader := &common.SignatureHeader{}
	err = proto.Unmarshal(rawPayload.Header.SignatureHeader, signatureHeader)
	if err != nil {
		log.Printf("Failed to unmarshal Payload's signatureHeader :%s \n", err)
	}

	creator := &msp.SerializedIdentity{}
	err = proto.Unmarshal(signatureHeader.Creator, creator)
	if err != nil {
		log.Printf("Failed to unmarshal Payload's signatureHeader :%s \n", err)
	}

	// from creator get IdBytes (MSP)

	//uEnc := base64.URLEncoding.EncodeToString(creator.IdBytes)
	//
	//certText , err := base64.URLEncoding.DecodeString(uEnc)
	//if err != nil {
	//	log.Printf("Failed to unmarshal  certText :%s \n", err)
	//}
	//end , _ := pem.Decode(certText)
	//if end != nil{
	//	log.Printf("Failed to unmarshal  end :%s \n", end)
	//}

	//cert , err := x509.ParseCertificate(end.Bytes)
	//if err != nil {
	//	log.Printf("Failed to unmarshal  certText :%s \n", err)
	//}

	transaction := &peer.Transaction{}
	err = proto.Unmarshal(rawPayload.Data, transaction)
	if err != nil {
		log.Printf("Failed to unmarshal rawPayload's Data :%s \n", err)
	}

	var transactionActionList []*models.TransactionAction

	for i := range transaction.Actions {
		transactionAction, err := GetTxActionFromTxDeep(transaction.Actions[i])
		if err != nil {
			log.Printf("Failed to unmarshal transactionAction :%s \n", err)
		}
		transactionAction.TxId = channelHeader.TxId
		transactionAction.Type = string(channelHeader.Type)
		transactionAction.Timestamp = time.Unix(channelHeader.Timestamp.Seconds, 0).Format("2006-01-02 15:04:05")
		transactionAction.ChannelId = channelHeader.ChannelId
		transactionActionList = append(transactionActionList, transactionAction)
	}

	transactionInfo := models.Transaction{TransactionActionList: transactionActionList}
	return &transactionInfo, nil
}

func GetTxActionFromTxDeep(txAction *peer.TransactionAction) (*models.TransactionAction, error) {

	chaincodeActionPayload := &peer.ChaincodeActionPayload{}
	err := proto.Unmarshal(txAction.Payload, chaincodeActionPayload)
	if err != nil {
		log.Printf("Failed to unmarshal chaincodeActionPayload :%s \n", err)
	}

	proposalResponsePayload := &peer.ProposalResponsePayload{}
	chaincodeAction := &peer.ChaincodeAction{}
	chaincodeID := ""
	var nsReadWriteSetList []*rwset.NsReadWriteSet
	var readWriteSetList []*kvrwset.KVRWSet
	var readSetList []string
	var writeSetList []string

	if chaincodeActionPayload.GetAction() != nil {
		err = proto.Unmarshal(chaincodeActionPayload.Action.ProposalResponsePayload, proposalResponsePayload)
		if err != nil {
			log.Printf("Failed to unmarshal proposalResponsePayload :%s \n", err)
		}

		err = proto.Unmarshal(proposalResponsePayload.Extension, chaincodeAction)
		if err != nil {
			log.Printf("Failed to unmarshal chaincodeAction :%s \n", err)
		}

		if chaincodeAction.ChaincodeId == nil {
			chaincodeID = "none"
		} else {
			chaincodeID = chaincodeAction.ChaincodeId.Name

			txReadWriteSet := &rwset.TxReadWriteSet{}
			err = proto.Unmarshal(chaincodeAction.Results, txReadWriteSet)
			if err != nil {
				log.Printf("Failed to unmarshal txReadWriteSet :%s \n", err)
			}

			for i := range txReadWriteSet.NsRwset {
				readWriteSet := &kvrwset.KVRWSet{}

				err = proto.Unmarshal(txReadWriteSet.NsRwset[i].Rwset, readWriteSet)
				if err != nil {
					log.Printf("Failed to unmarshal readWriteSet :%s \n", err)
				}

				for i := range readWriteSet.Reads {
					readSetJsonStr, err := json.Marshal(readWriteSet.Reads[i])
					if err != nil {
						log.Printf("Failed to unmarshal readSetJsonStr :%s \n", err)
					}
					readSetList = append(readSetList, string(readSetJsonStr))

				}
				for i := range readWriteSet.Writes {

					writeSetItem := map[string]interface{}{
						"key":      readWriteSet.Writes[i].GetKey(),
						"Value":    string(readWriteSet.Writes[i].GetValue()),
						"IsDelete": readWriteSet.Writes[i].GetIsDelete(),
					}

					writeSetJsonStr, err := json.Marshal(writeSetItem)
					if err != nil {
						log.Printf("Failed to unmarshal writeSetJsonStr :%s \n", err)
					}
					writeSetList = append(writeSetList, string(writeSetJsonStr))
				}

				readWriteSetList = append(readWriteSetList, readWriteSet)
				nsReadWriteSetList = append(nsReadWriteSetList, txReadWriteSet.NsRwset[i])
			}
		}
	} else {
		chaincodeID = "none"
	}

	var endorsements []string

	if chaincodeActionPayload.Action.GetEndorsements() != nil {

		for i := range chaincodeActionPayload.Action.GetEndorsements() {

			endorser := &msp.SerializedIdentity{}
			err = proto.Unmarshal(chaincodeActionPayload.Action.Endorsements[i].Endorser, endorser)
			if err != nil {
				log.Printf("Failed to unmarshal endorser :%s \n", err)
			}
			endorsements = append(endorsements, endorser.Mspid)
		}
	}

	txActionInfo := models.TransactionAction{

		Endorsements: endorsements,
		ChaincodeId:  chaincodeID,
		ReadSetList:  readSetList,
		WriteSetList: writeSetList,
	}

	return &txActionInfo, nil

}
