/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-contract-api-go/metadata"
)

func main() {
	cbdcContract := new(CbdcContract)
	cbdcContract.Info.Version = "0.0.1"
	cbdcContract.Info.Description = "My Smart Contract"
	cbdcContract.Info.License = new(metadata.LicenseMetadata)
	cbdcContract.Info.License.Name = "Apache-2.0"
	cbdcContract.Info.Contact = new(metadata.ContactMetadata)
	cbdcContract.Info.Contact.Name = "John Doe"

	chaincode, err := contractapi.NewChaincode(cbdcContract)
	chaincode.Info.Title = "chaincode-CBDC chaincode"
	chaincode.Info.Version = "0.0.1"

	if err != nil {
		panic("Could not create chaincode from CbdcContract." + err.Error())
	}

	err = chaincode.Start()

	if err != nil {
		panic("Failed to start chaincode. " + err.Error())
	}
}
