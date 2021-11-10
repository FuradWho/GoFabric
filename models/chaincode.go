package models

import "time"

type ChaincodeInfo struct {

	Uuid string `json:"uuid"`
	TxId string `json:"txId"`
	Time time.Time `json:"time"`
	Payload string `json:"payload"`
}