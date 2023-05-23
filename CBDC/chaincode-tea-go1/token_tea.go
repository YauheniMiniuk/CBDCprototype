package main

import (
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)


func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		log.Panicf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting fabcar chaincode: %s", err.Error())
	}
}