/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const getStateError = "world state get error"

type MockStub struct {
	shim.ChaincodeStubInterface
	mock.Mock
}

func (ms *MockStub) GetState(key string) ([]byte, error) {
	args := ms.Called(key)

	return args.Get(0).([]byte), args.Error(1)
}

func (ms *MockStub) PutState(key string, value []byte) error {
	args := ms.Called(key, value)

	return args.Error(0)
}

func (ms *MockStub) DelState(key string) error {
	args := ms.Called(key)

	return args.Error(0)
}

type MockContext struct {
	contractapi.TransactionContextInterface
	mock.Mock
}

func (mc *MockContext) GetStub() shim.ChaincodeStubInterface {
	args := mc.Called()

	return args.Get(0).(*MockStub)
}

func configureStub() (*MockContext, *MockStub) {
	var nilBytes []byte

	testCbdc := new(Cbdc)
	testCbdc.Value = "set value"
	cbdcBytes, _ := json.Marshal(testCbdc)

	ms := new(MockStub)
	ms.On("GetState", "statebad").Return(nilBytes, errors.New(getStateError))
	ms.On("GetState", "missingkey").Return(nilBytes, nil)
	ms.On("GetState", "existingkey").Return([]byte("some value"), nil)
	ms.On("GetState", "cbdckey").Return(cbdcBytes, nil)
	ms.On("PutState", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)
	ms.On("DelState", mock.AnythingOfType("string")).Return(nil)

	mc := new(MockContext)
	mc.On("GetStub").Return(ms)

	return mc, ms
}

func TestCbdcExists(t *testing.T) {
	var exists bool
	var err error

	ctx, _ := configureStub()
	c := new(CbdcContract)

	exists, err = c.CbdcExists(ctx, "statebad")
	assert.EqualError(t, err, getStateError)
	assert.False(t, exists, "should return false on error")

	exists, err = c.CbdcExists(ctx, "missingkey")
	assert.Nil(t, err, "should not return error when can read from world state but no value for key")
	assert.False(t, exists, "should return false when no value for key in world state")

	exists, err = c.CbdcExists(ctx, "existingkey")
	assert.Nil(t, err, "should not return error when can read from world state and value exists for key")
	assert.True(t, exists, "should return true when value for key in world state")
}

func TestCreateCbdc(t *testing.T) {
	var err error

	ctx, stub := configureStub()
	c := new(CbdcContract)

	err = c.CreateCbdc(ctx, "statebad", "some value")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors")

	err = c.CreateCbdc(ctx, "existingkey", "some value")
	assert.EqualError(t, err, "The asset existingkey already exists", "should error when exists returns true")

	err = c.CreateCbdc(ctx, "missingkey", "some value")
	stub.AssertCalled(t, "PutState", "missingkey", []byte("{\"value\":\"some value\"}"))
}

func TestReadCbdc(t *testing.T) {
	var cbdc *Cbdc
	var err error

	ctx, _ := configureStub()
	c := new(CbdcContract)

	cbdc, err = c.ReadCbdc(ctx, "statebad")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors when reading")
	assert.Nil(t, cbdc, "should not return Cbdc when exists errors when reading")

	cbdc, err = c.ReadCbdc(ctx, "missingkey")
	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when reading")
	assert.Nil(t, cbdc, "should not return Cbdc when key does not exist in world state when reading")

	cbdc, err = c.ReadCbdc(ctx, "existingkey")
	assert.EqualError(t, err, "Could not unmarshal world state data to type Cbdc", "should error when data in key is not Cbdc")
	assert.Nil(t, cbdc, "should not return Cbdc when data in key is not of type Cbdc")

	cbdc, err = c.ReadCbdc(ctx, "cbdckey")
	expectedCbdc := new(Cbdc)
	expectedCbdc.Value = "set value"
	assert.Nil(t, err, "should not return error when Cbdc exists in world state when reading")
	assert.Equal(t, expectedCbdc, cbdc, "should return deserialized Cbdc from world state")
}

func TestUpdateCbdc(t *testing.T) {
	var err error

	ctx, stub := configureStub()
	c := new(CbdcContract)

	err = c.UpdateCbdc(ctx, "statebad", "new value")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors when updating")

	err = c.UpdateCbdc(ctx, "missingkey", "new value")
	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when updating")

	err = c.UpdateCbdc(ctx, "cbdckey", "new value")
	expectedCbdc := new(Cbdc)
	expectedCbdc.Value = "new value"
	expectedCbdcBytes, _ := json.Marshal(expectedCbdc)
	assert.Nil(t, err, "should not return error when Cbdc exists in world state when updating")
	stub.AssertCalled(t, "PutState", "cbdckey", expectedCbdcBytes)
}

func TestDeleteCbdc(t *testing.T) {
	var err error

	ctx, stub := configureStub()
	c := new(CbdcContract)

	err = c.DeleteCbdc(ctx, "statebad")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors")

	err = c.DeleteCbdc(ctx, "missingkey")
	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when deleting")

	err = c.DeleteCbdc(ctx, "cbdckey")
	assert.Nil(t, err, "should not return error when Cbdc exists in world state when deleting")
	stub.AssertCalled(t, "DelState", "cbdckey")
}
