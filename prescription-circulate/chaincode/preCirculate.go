package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"prescription-circulate/chaincode/smartcontract"
	//"github.com/hyperledger/fabric-samples/asset-transfer-basic/chaincode-go/chaincode"
)

func main() {
	preHashChaincode, err := contractapi.NewChaincode(&smartcontract.SmartContract{})
	if err != nil {
		log.Panicf("Error creating pescription-circulate chaincode: %v", err)
	}

	if err := preHashChaincode.Start(); err != nil {
		log.Panicf("Error starting prescription-circulate chaincode: %v", err)
	}
}
