/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

var USER_NAME = "admin"

func main() {

	log.Println("============ application-golang starts ============")

	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environment variable: %v", err)
	}

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	if !wallet.Exists(USER_NAME) {
		err = populateWallet(wallet)
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}

	// ccpPath := filepath.Join(
	// 	"..",
	// 	"..",
	// 	"test-network",
	// 	"organizations",
	// 	"peerOrganizations",
	// 	"org1.example.com",
	// 	"connection-org1.yaml",
	// )
	ccpPath := "D:\\Programming\\Hyperledger\\Vagrant\\data\\Server\\CBDC\\test-network\\organizations\\peerOrganizations\\org2.example.com\\connection-org2.yaml"
	cleanFile := gateway.WithConfig(config.FromFile(filepath.Clean((ccpPath))))
	println(cleanFile)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, USER_NAME),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	channelName := "mychannel"
	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
		channelName = cname
	}

	log.Println("--> Connecting to channel", channelName)
	network, err := gw.GetNetwork(channelName)
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}

	chaincodeName := "cbdc"
	if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
		chaincodeName = ccname
	}

	log.Println("--> Using chaincode", chaincodeName)
	contract := network.GetContract(chaincodeName)

	// log.Println("--> Submit Transaction: Initialize, function creates the initial set of assets on the ledger")
	// result, err := contract.SubmitTransaction("Initialize", "CBR", "CBR", "2")
	// if err != nil {
	// 	log.Fatalf("Failed to Submit transaction: %v", err)
	// }
	// log.Println(string(result))

	log.Println("--> Submit Transaction: Mint, function creates the initial set of assets on the ledger")
	result, err := contract.SubmitTransaction("Mint", "5000")
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	log.Println(string(result))

	// log.Println("--> Evaluate Transaction: ClientAccountID")
	// result, err := contract.SubmitTransaction("ClientAccountID")
	// if err != nil {
	// 	log.Fatalf("Failed to evaluate transaction: %v", err)
	// }
	// log.Println(string(result))

	// log.Println("--> Evaluate Transaction: ClientAccountBalance")
	// result, err = contract.SubmitTransaction("ClientAccountBalance")
	// if err != nil {
	// 	log.Fatalf("Failed to evaluate transaction: %v", err)
	// }
	// log.Println(string(result))

	// log.Println("--> Submit Transaction: CreateAsset, creates new asset with ID, color, owner, size, and appraisedValue arguments")
	// result, err = contract.SubmitTransaction("CreateAsset", "asset13", "yellow", "5", "Tom", "1300")
	// if err != nil {
	// 	log.Fatalf("Failed to Submit transaction: %v", err)
	// }
	// log.Println(string(result))

	// log.Println("--> Evaluate Transaction: ReadAsset, function returns an asset with a given assetID")
	// result, err = contract.EvaluateTransaction("ReadAsset", "asset13")
	// if err != nil {
	// 	log.Fatalf("Failed to evaluate transaction: %v\n", err)
	// }
	// log.Println(string(result))

	// log.Println("--> Evaluate Transaction: AssetExists, function returns 'true' if an asset with given assetID exist")
	// result, err = contract.EvaluateTransaction("AssetExists", "asset1")
	// if err != nil {
	// 	log.Fatalf("Failed to evaluate transaction: %v\n", err)
	// }
	// log.Println(string(result))

	// log.Println("--> Submit Transaction: TransferAsset asset1, transfer to new owner of Tom")
	// _, err = contract.SubmitTransaction("TransferAsset", "asset1", "Tom")
	// if err != nil {
	// 	log.Fatalf("Failed to Submit transaction: %v", err)
	// }

	// log.Println("--> Evaluate Transaction: ReadAsset, function returns 'asset1' attributes")
	// result, err = contract.EvaluateTransaction("ReadAsset", "asset1")
	// if err != nil {
	// 	log.Fatalf("Failed to evaluate transaction: %v", err)
	// }
	// log.Println(string(result))
	log.Println("============ application-golang ends ============")
}

func populateWallet(wallet *gateway.Wallet) error {
	log.Println("============ Populating wallet ============")
	// credPath := filepath.Join(
	// 	"..",
	// 	"..",
	// 	"test-network",
	// 	"organizations",
	// 	"peerOrganizations",
	// 	"org1.example.com",
	// 	"users",
	// 	"User1@org1.example.com",
	// 	"msp",
	// )
	credPath := "D:/Programming/Hyperledger/Vagrant/data/Server/CBDC/test-network/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp"

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err := os.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := os.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := os.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org2MSP", string(cert), string(key))

	return wallet.Put(USER_NAME, identity)
}
