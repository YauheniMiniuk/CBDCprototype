package chaincode

import (
	"encoding/json"
	"fmt"
	"time"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Define key names for options
var MINTER = "Org2MSP"

// Define objectType names for prefix
const allowancePrefix = "allowance"

// TeaContract provides functions for managing a car
type TeaContract struct {
	contractapi.Contract
}
// Car describes basic details of what makes up a car
type Tea struct {
	Name   string `json:"name"`
	Price  float32 `json:"price"`
	Amount float32 `json:"amount"`
	Owner  string `json:"owner"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"key"`
	Record *Tea
}
type QueryHistory struct {
	TxID string `json:"txID"`
	TimeStamp time.Time `json:"timestamp"`
	IsDelete bool `json:"isDelete"`
	Record *Tea
}

// event provides an organized struct for emitting events
type event struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Record	*Tea    `json:"value"`
}

func (s *TeaContract) Mint(ctx contractapi.TransactionContextInterface, name string, price float32, amount float32, recipient string)(string, error){

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
		Amount: amount,
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
func (s *TeaContract) QueryAllTokens(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
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

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		tea := new(Tea)
		_ = json.Unmarshal(queryResponse.Value, tea)

		queryResult := QueryResult{Key: queryResponse.Key, Record: tea}
		results = append(results, queryResult)
	}

	return results, nil
}

// QueryAllCars returns all cars found in world state
func (s *TeaContract) QueryTokensByClientID(ctx contractapi.TransactionContextInterface, clientID string) ([]QueryResult, error) {
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

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		token := new(Tea)
		_ = json.Unmarshal(queryResponse.Value, token)

		if token.Owner == clientID{
			queryResult := QueryResult{Key: queryResponse.Key, Record: token}
			results = append(results, queryResult)
		}
	}

	return results, nil
}

func (s *TeaContract) QueryClientTokens(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
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

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		token := new(Tea)
		_ = json.Unmarshal(queryResponse.Value, token)

		if token.Owner == clientID{
			queryResult := QueryResult{Key: queryResponse.Key, Record: token}
			results = append(results, queryResult)
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

func (s *TeaContract) Sub(ctx contractapi.TransactionContextInterface, tokenId string, amount float32) string {
	// Check minter authorization - this sample assumes Org1 is the central banker with privilege to burn new tokens
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "Не удается получить MSPID"
	}

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

	if token.Amount < amount {
		token.Amount = 0
	} else {
		token.Amount -= amount
	}

	if token.Amount == 0{
		ctx.GetStub().DelState(tokenId)
	}

	return fmt.Sprintf("Количество токена было уменьшено на %s единиц", amount)
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
func (s *TeaContract) GetHistoryForKey(ctx contractapi.TransactionContextInterface, key string) ([]QueryHistory, error) {
	iterator, err := ctx.GetStub().GetHistoryForKey(key)
	if err != nil {
		return nil, fmt.Errorf("Не удалось получить историю транзакций токена %s: %v", key, err)
	}
	defer iterator.Close()

	var result []QueryHistory
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

		tea := QueryHistory{
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




// ??????????????????????????????????????????? 
func (s *TeaContract) Approve(ctx contractapi.TransactionContextInterface, spender string, tokenId string) error {
	// Get ID of submitting client identity
	owner, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client id: %v", err)
	}

	// Create allowanceKey
	allowanceKey, err := ctx.GetStub().CreateCompositeKey(allowancePrefix, []string{owner, spender})
	if err != nil {
		return fmt.Errorf("failed to create the composite key for prefix %s: %v", allowancePrefix, err)
	}

	// Update the state of the smart contract by adding the allowanceKey and value
	err = ctx.GetStub().PutState(allowanceKey, []byte(tokenId))
	if err != nil {
		return fmt.Errorf("failed to update state of smart contract for key %s: %v", allowanceKey, err)
	}
	token, err := s.QueryToken(ctx, tokenId)
	if err != nil {
		return fmt.Errorf("Не удалось найти указанный токен")
	}
	// Emit the Approval event
	approvalEvent := event{owner, spender, token}
	approvalEventJSON, err := json.Marshal(approvalEvent)
	if err != nil {
		return fmt.Errorf("failed to obtain JSON encoding: %v", err)
	}
	err = ctx.GetStub().SetEvent("Approval", approvalEventJSON)
	if err != nil {
		return fmt.Errorf("failed to set event: %v", err)
	}

	// log.Printf("client %s approved a withdrawal allowance of %d for spender %s", owner, tokenId, spender)

	return nil
}

// Allowance returns the amount still available for the spender to withdraw from the owner
func (s *TeaContract) Allowance(ctx contractapi.TransactionContextInterface, owner string, spender string) (int, error) {
	// Create allowanceKey
	allowanceKey, err := ctx.GetStub().CreateCompositeKey(allowancePrefix, []string{owner, spender})
	if err != nil {
		return 0, fmt.Errorf("failed to create the composite key for prefix %s: %v", allowancePrefix, err)
	}

	// Read the allowance amount from the world state
	allowanceBytes, err := ctx.GetStub().GetState(allowanceKey)
	if err != nil {
		return 0, fmt.Errorf("failed to read allowance for %s from world state: %v", allowanceKey, err)
	}

	var allowance int

	// If no current allowance, set allowance to 0
	if allowanceBytes == nil {
		allowance = '0'
	} else {
		allowance, err = strconv.Atoi(string(allowanceBytes)) // Error handling not needed since Itoa() was used when setting the totalSupply, guaranteeing it was an integer.
	}

	// log.Printf("The allowance left for spender %s to withdraw from owner %s: %d", spender, owner, allowance)

	return allowance, nil
}

// TransferFrom transfers the value amount from the "from" address to the "to" address
// This function triggers a Transfer event
func (s *TeaContract) TransferFrom(ctx contractapi.TransactionContextInterface, from string, to string, tokenId string) error {
	// Get ID of submitting client identity
	spender, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client id: %v", err)
	}

	// Create allowanceKey
	allowanceKey, err := ctx.GetStub().CreateCompositeKey(allowancePrefix, []string{from, spender})
	if err != nil {
		return fmt.Errorf("failed to create the composite key for prefix %s: %v", allowancePrefix, err)
	}

	// Retrieve the allowance of the spender
	currentAllowanceBytes, err := ctx.GetStub().GetState(allowanceKey)
	if err != nil {
		return fmt.Errorf("failed to retrieve the allowance for %s from world state: %v", allowanceKey, err)
	}

	var currentAllowance int
	currentAllowance, _ = strconv.Atoi(string(currentAllowanceBytes)) // Error handling not needed since Itoa() was used when setting the totalSupply, guaranteeing it was an integer.

	// Check if transferred value is less than allowance
	if currentAllowance <= 0 {
		return fmt.Errorf("spender does not have enough allowance for transfer")
	}

	// Get token
	token, err := s.QueryToken(ctx, tokenId)
	if err != nil{
		return fmt.Errorf("Не удалось найти указанный токен")
	}
	// Initiate the transfer
	token.Owner = to


	err = ctx.GetStub().PutState(allowanceKey, []byte(strconv.Itoa(0)))
	if err != nil {
		return err
	}

	// Emit the Transfer event
	transferEvent := event{from, to, token}
	transferEventJSON, err := json.Marshal(transferEvent)
	if err != nil {
		return fmt.Errorf("failed to obtain JSON encoding: %v", err)
	}
	err = ctx.GetStub().SetEvent("Transfer", transferEventJSON)
	if err != nil {
		return fmt.Errorf("failed to set event: %v", err)
	}

	// log.Printf("spender %s allowance updated from %d to %d", spender, currentAllowance, 0)

	return nil
}