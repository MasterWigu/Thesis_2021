/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"github.com/MasterWigu/Thesis/chaincode/inventory-management/chaincode"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	assetChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating inventory management chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting inventory management chaincode: %v", err)
	}
}
