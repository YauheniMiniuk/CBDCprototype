/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-samples/token-erc-20/chaincode-go/chaincode"
)

func main() {
	chaincode, err := contractapi.NewChaincode(&chaincode.TeaContract{}, &chaincode.Erc20Contract{})
	// chaincode, err := contractapi.NewChaincode(&chaincode.TeaContract{})

	if err != nil {
		log.Panicf("Error creating token-tea chaincode: %v", err)
	}
	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting token-tea chaincode: %v", err)
	}

	// erc20Chaincode, err := contractapi.NewChaincode(&chaincode.Erc20Contract{})
	// if err != nil {
	// 	log.Panicf("Error creating token-erc-20 chaincode: %v", err)
	// }
	// if err := erc20Chaincode.Start(); err != nil {
	// 	log.Panicf("Error starting token-erc-20 chaincode: %v", err)
	// }
}
