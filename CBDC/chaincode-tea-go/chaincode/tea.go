package chaincode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Define key names for options
var MINTER = "Org2MSP"

// TeaContract provides functions for managing a car
type TeaContract struct {
	contractapi.Contract
}
// Car describes basic details of what makes up a car
type Tea struct {
	Name   string `json:"name"`
	Price  float32 `json:"price"`
	// Amount float32 `json:"amount"`
	Owner  string `json:"owner"`
}

// QueryResult structure used for handling result of query
type QueryTea struct {
	Key    string `json:"key"`
	Record *Tea
}
type TeaHistory struct {
	TxID string `json:"txID"`
	TimeStamp time.Time `json:"timestamp"`
	IsDelete bool `json:"isDelete"`
	Record *Tea
}

// event provides an organized struct for emitting events
type TeaEvent struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Record	*Tea    `json:"value"`
}

func (s *TeaContract) MintToken(ctx contractapi.TransactionContextInterface, name string, price float32, amount float32, recipient string)(string, error){

	// Check minter authorization - this sample assumes Org1 is the central banker with privilege to mint new tokens
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("Не удается получить MSPID")
	}
	if clientMSPID != MINTER {
		return "", fmt.Errorf("Клиент не авторизован для выпуска токенов!")
	}
	
	// Get ID of submitting client identity
	minter, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("Не удается получить идентификатор клиента")
	}

	if price < 0 {
		return "", fmt.Errorf("Цена токена не может быть отрицательной!")
	}
	if amount <= 0 {
		return "", fmt.Errorf("Размер токена не может меньше либо равным нулю")
	}

	tea := Tea{
		Name: name,
		Price: price,
		Owner: minter,
	}

	// Mint token
	tokenAsBytes, _ := json.Marshal(tea)

	id := ctx.GetStub().GetTxID()
	ctx.GetStub().PutState(id, tokenAsBytes)
	// if result == nil {
	// 	return "", fmt.Errorf("Не удалось выпустить токен")
	// }
	// Transfer token to the owner
	if recipient == ""{
		return "Не найден получатель. Токен был успешно выпущен и помещен в кошелек администратора.", nil
	}
	s.Transfer(ctx, id, recipient)

	// Return
	return "Токен был успешно выпущен", nil
}

// QueryTea returns the tea stored in the world state with given id
func (s *TeaContract) QueryToken(ctx contractapi.TransactionContextInterface, tokenId string) (*Tea, error) {
	tokenAsBytes, err := ctx.GetStub().GetState(tokenId)

	if err != nil{
		return nil, fmt.Errorf("Не удалось прочитать world state")
	}

	if tokenAsBytes == nil {
		return nil, fmt.Errorf("%s не существует", tokenAsBytes)
	}

	tea := new(Tea)
	_ = json.Unmarshal(tokenAsBytes, tea)
	
	return tea, nil
}

// QueryAllCars returns all cars found in world state
func (s *TeaContract) QueryAllTokens(ctx contractapi.TransactionContextInterface) ([]QueryTea, error) {
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("Не удается получить MSPID")
	}
	if clientMSPID != MINTER {
		return nil, fmt.Errorf("Клиент не авторизован для просмотра токенов!")
	}
	
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryTea{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		tea := new(Tea)
		_ = json.Unmarshal(queryResponse.Value, tea)

		queryTea := QueryTea{Key: queryResponse.Key, Record: tea}
		results = append(results, queryTea)
	}

	return results, nil
}

// QueryAllCars returns all cars found in world state
func (s *TeaContract) QueryTokensByClientID(ctx contractapi.TransactionContextInterface, clientID string) ([]QueryTea, error) {
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("Не удается получить MSPID")
	}
	if clientMSPID != MINTER {
		return nil, fmt.Errorf("Клиент не авторизован для просмотра токенов!")
	}
	
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryTea{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		token := new(Tea)
		_ = json.Unmarshal(queryResponse.Value, token)

		if token.Owner == clientID{
			queryTea := QueryTea{Key: queryResponse.Key, Record: token}
			results = append(results, queryTea)
		}
	}

	return results, nil
}

func (s *TeaContract) QueryClientTokens(ctx contractapi.TransactionContextInterface) ([]QueryTea, error) {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return nil, fmt.Errorf("Не удается получить ID")
	}
	
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryTea{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		token := new(Tea)
		_ = json.Unmarshal(queryResponse.Value, token)

		if token.Owner == clientID{
			queryTea := QueryTea{Key: queryResponse.Key, Record: token}
			results = append(results, queryTea)
		}
	}

	return results, nil
}

func (s *TeaContract) Transfer(ctx contractapi.TransactionContextInterface, tokenId string, recipientId string) string {
	// Get ID of submitting client identity
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "Не удалось получить ID владельца"
	}
	_, err = ctx.GetStub().GetState(recipientId)
	if err != nil {
		return "Не удалось получить ID получателя"
	}

	token, err := s.QueryToken(ctx, tokenId)
	if err != nil{
		return "Не удалось найти токен"
	}
	if token.Owner != clientID{
		return "Вы не можете отправить выбранный токен"
	}

	token.Owner = recipientId

	tokenAsBytes, _ := json.Marshal(token)
	ctx.GetStub().PutState(tokenId, tokenAsBytes)

	return "Токен был успешно передан"
}

func (s *TeaContract) Burn(ctx contractapi.TransactionContextInterface, tokenId string) (string){
	
	// Check minter authorization - this sample assumes Org1 is the central banker with privilege to burn new tokens
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "Не удается получить MSPID"
	}
	// if clientMSPID != MINTER {
	// 	return "Клиент не может сжигать токены"
	// }

	// Get ID of submitting client identity
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "failed to get client id"
	}
	
	token, err := s.QueryToken(ctx, tokenId)
	if err != nil{
		return "Не удалось найти указанный токен"
	}

	if token.Owner != clientID  || clientMSPID != MINTER{
		return "Вы не можете удалить токен"
	}

	ctx.GetStub().DelState(tokenId)

	return "Токен был успешно удалён"

}

func (s *TeaContract) ClientAccountID(ctx contractapi.TransactionContextInterface) (string, error) {
	// Get ID of submitting client identity
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("Не удалось получить ID: %v", err)
	}

	return clientID, nil
}

// GetHistoryForKey returns all transactions for a given key (token)
func (s *TeaContract) GetHistoryForKey(ctx contractapi.TransactionContextInterface, key string) ([]TeaHistory, error) {
	iterator, err := ctx.GetStub().GetHistoryForKey(key)
	if err != nil {
		return nil, fmt.Errorf("Не удалось получить историю транзакций токена %s: %v", key, err)
	}
	defer iterator.Close()

	var result []TeaHistory
	for iterator.HasNext() {
		response, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("Не удалось получить следующий блок транзкций: %v", err)
		}
		txID := response.TxId
		timestamp := time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos))
		isDelete := response.IsDelete
		valueBytes := response.Value

		value := new(Tea)
	_ = json.Unmarshal(valueBytes, &value)

		tea := TeaHistory{
			TxID: txID,
			TimeStamp: timestamp,
			IsDelete: isDelete,
			Record: value,
		}
		// teaBytes, _ := json.Marshal(tea)
		result = append(result, tea)
		// result = append(result, fmt.Sprintf("TxID: %s Timestamp: %s IsDelete: %t Value: %s", txID, timestamp, isDelete, value))
	}

	return result, nil
}