/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

var num = 0
var MINTER = "Org2MSP"

// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}

// Car describes basic details of what makes up a car
type Tea struct {
	Name   string `json:"name"`
	Price  string `json:"price"`
	Amount string `json:"amount"`
	Owner  string `json:"owner"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"key"`
	Record *Tea
}

// // InitLedger adds a base set of cars to the ledger
// func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
// 	cars := []Car{
// 		Car{Make: "Toyota", Model: "Prius", Colour: "blue", Owner: "Tomoko"},
// 		Car{Make: "Ford", Model: "Mustang", Colour: "red", Owner: "Brad"},
// 		Car{Make: "Hyundai", Model: "Tucson", Colour: "green", Owner: "Jin Soo"},
// 		Car{Make: "Volkswagen", Model: "Passat", Colour: "yellow", Owner: "Max"},
// 		Car{Make: "Tesla", Model: "S", Colour: "black", Owner: "Adriana"},
// 		Car{Make: "Peugeot", Model: "205", Colour: "purple", Owner: "Michel"},
// 		Car{Make: "Chery", Model: "S22L", Colour: "white", Owner: "Aarav"},
// 		Car{Make: "Fiat", Model: "Punto", Colour: "violet", Owner: "Pari"},
// 		Car{Make: "Tata", Model: "Nano", Colour: "indigo", Owner: "Valeria"},
// 		Car{Make: "Holden", Model: "Barina", Colour: "brown", Owner: "Shotaro"},
// 	}

// 	for i, car := range cars {
// 		carAsBytes, _ := json.Marshal(car)
// 		err := ctx.GetStub().PutState("CAR"+strconv.Itoa(i), carAsBytes)

// 		if err != nil {
// 			return fmt.Errorf("Failed to put to world state. %s", err.Error())
// 		}
// 	}

// 	return nil
// }

func (s *SmartContract) Mint(ctx contractapi.TransactionContextInterface, name string, price float32, amount float32, recipient string)(string){

	// Check minter authorization - this sample assumes Org1 is the central banker with privilege to mint new tokens
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("Не удается получить MSPID: %v", err)
	}
	if clientMSPID != MINTER {
		return fmt.Errorf("Клиент не авторизован для выпуска токенов!")
	}
	
	// Get ID of submitting client identity
	minter, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("Не удается получить идентификатор клиента: %v", err)
	}
	ownerIdentity, err := ctx.GetClientIdentity().GetID()

	if price < 0 {
		return fmt.Errorf("Цена токена не может быть отрицательной!")
	}
	if amount <= 0 {
		return fmt.Errorf("Размер токена не может меньше либо равным нулю")
	}

	tea := Tea{
		Name: name,
		Price: price,
		Amount: amount,
		Owner: minter,
	}

	// Mint token
	tokenAsBytes, _ := json.Marshal(tea)

	result := ctx.GetStub().PutState(i, tokenAsBytes)
	if !result {
		return fmt.Errorf("Не удалось выпустить токен")
	}
	// Transfer token to the owner

	// Return
	return fmt.Printf("Токен был успешно выпущен")
}

func (s *SmartContract) Burn(ctx contractapi.TransactionContextInterface, tokenId string) (string){
	
	// Check minter authorization - this sample assumes Org1 is the central banker with privilege to burn new tokens
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("Не удается получить MSPID: %v", err)
	}
	if clientMSPID != MINTER {
		return fmt.Errorf("Клиент не может сжигать токены")
	}

	// Get ID of submitting client identity
	minter, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client id: %v", err)
	}
	
	token, err := QueryToken(tokenId)
	if err != nil{
		return fmt.Errorf("Не удалось найти указанный токен")
	}


	return fmt.Printf("Токен был успешно сожжён")

}


// QueryTea returns the car stored in the world state with given id
func (s *SmartContract) QueryToken(ctx contractapi.TransactionContextInterface, tokenId string) (*Tea, error) {
	tokenAsBytes, err := ctx.GetStub().GetState(tokenId)

	if err != nil{
		return nil, fmt.Errorf("Не удалось прочитать world state. %s", err.Error())
	}

	if tokenAsBytes == nil {
		return nil, fmt.Errorf("%s не существует", tokenAsBytes)
	}

	tea := new(Tea)
	_ = json.Unmarshal(tokenAsBytes, tea)

	return tea, nil
}



// // QueryAllCars returns all cars found in world state
// func (s *SmartContract) QueryAllCars(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
// 	startKey := ""
// 	endKey := ""

// 	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resultsIterator.Close()

// 	results := []QueryResult{}

// 	for resultsIterator.HasNext() {
// 		queryResponse, err := resultsIterator.Next()

// 		if err != nil {
// 			return nil, err
// 		}

// 		tea := new(Tea)
// 		_ = json.Unmarshal(queryResponse.Value, tea)

// 		queryResult := QueryResult{Key: queryResponse.Key, Record: tea}
// 		results = append(results, queryResult)
// 	}

// 	return results, nil
// }

// // ChangeCarOwner updates the owner field of car with given id in world state
// func (s *SmartContract) ChangeCarOwner(ctx contractapi.TransactionContextInterface, carNumber string, newOwner string) error {
// 	car, err := s.QueryCar(ctx, carNumber)

// 	if err != nil {
// 		return err
// 	}

// 	car.Owner = newOwner

// 	carAsBytes, _ := json.Marshal(car)

// 	return ctx.GetStub().PutState(carNumber, carAsBytes)
// }
