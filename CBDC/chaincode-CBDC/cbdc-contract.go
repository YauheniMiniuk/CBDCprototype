/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// CbdcContract contract for managing CRUD for Cbdc
type CbdcContract struct {
	contractapi.Contract
}

// CbdcExists returns true when asset with given ID exists in world state
func (c *CbdcContract) CbdcExists(ctx contractapi.TransactionContextInterface, cbdcID string) (bool, error) {
	data, err := ctx.GetStub().GetState(cbdcID)

	if err != nil {
		return false, err
	}

	return data != nil, nil
}

// CreateCbdc creates a new instance of Cbdc
func (c *CbdcContract) CreateCbdc(ctx contractapi.TransactionContextInterface, cbdcID string, value string) error {
	exists, err := c.CbdcExists(ctx, cbdcID)
	if err != nil {
		return fmt.Errorf("Could not read from world state. %s", err)
	} else if exists {
		return fmt.Errorf("The asset %s already exists", cbdcID)
	}

	cbdc := new(Cbdc)
	cbdc.Value = value

	bytes, _ := json.Marshal(cbdc)

	return ctx.GetStub().PutState(cbdcID, bytes)
}

// ReadCbdc retrieves an instance of Cbdc from the world state
func (c *CbdcContract) ReadCbdc(ctx contractapi.TransactionContextInterface, cbdcID string) (*Cbdc, error) {
	exists, err := c.CbdcExists(ctx, cbdcID)
	if err != nil {
		return nil, fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return nil, fmt.Errorf("The asset %s does not exist", cbdcID)
	}

	bytes, _ := ctx.GetStub().GetState(cbdcID)

	cbdc := new(Cbdc)

	err = json.Unmarshal(bytes, cbdc)

	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal world state data to type Cbdc")
	}

	return cbdc, nil
}

// UpdateCbdc retrieves an instance of Cbdc from the world state and updates its value
func (c *CbdcContract) UpdateCbdc(ctx contractapi.TransactionContextInterface, cbdcID string, newValue string) error {
	exists, err := c.CbdcExists(ctx, cbdcID)
	if err != nil {
		return fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return fmt.Errorf("The asset %s does not exist", cbdcID)
	}

	cbdc := new(Cbdc)
	cbdc.Value = newValue

	bytes, _ := json.Marshal(cbdc)

	return ctx.GetStub().PutState(cbdcID, bytes)
}

// DeleteCbdc deletes an instance of Cbdc from the world state
func (c *CbdcContract) DeleteCbdc(ctx contractapi.TransactionContextInterface, cbdcID string) error {
	exists, err := c.CbdcExists(ctx, cbdcID)
	if err != nil {
		return fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return fmt.Errorf("The asset %s does not exist", cbdcID)
	}

	return ctx.GetStub().DelState(cbdcID)
}
