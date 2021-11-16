package chaincode

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) GetPerms(ctx contractapi.TransactionContextInterface) (string, error) {
	role, found, err := ctx.GetClientIdentity().GetAttributeValue("hf.Type")
	//cert, err := ctx.GetClientIdentity().GetX509Certificate()
	if err != nil {
		return "1", err
	}
	/*mgr := attrmgr.New()
	attrs, err := mgr.GetAttributesFromCert(cert)
	if err != nil {
		return "2", err
	}


	attrsJson, err := json.Marshal(attrs)
	*/
	if !found {
		return "None", nil
	}
	/*if err != nil {
		return "3", err
	}*/
	return role, nil
}
