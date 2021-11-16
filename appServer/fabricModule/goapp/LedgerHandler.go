package main

import (
	"encoding/json"
	"fmt"
	"github.com/MasterWigu/Thesis/appServer/AssetRepresentation"
)

func getPermissions(user *User) (string, error) {
	network := user.network

	contract := network.GetContract("userCheck")

	result, err := contract.EvaluateTransaction("GetPerms")
	if err != nil {
		return "", fmt.Errorf("failed to Submit transaction: %v", err)
	}

	return string(result), nil
}



/* Asset Management */

func GetAssetTypes(user *User) (*AssetRepresentation.TypeTracker, error) {
	network := user.network

	contract := network.GetContract("inventoryMgt")

	result, err := contract.EvaluateTransaction("GetAssetTypes")
	if err != nil {
		return nil, fmt.Errorf("failed to Submit transaction: %v", err)
	}

	typeTracker := AssetRepresentation.TypeTracker{}

	err = json.Unmarshal(result, &typeTracker)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal type tracker: %v", err)
	}

	return &typeTracker, nil
}

func AddAssetType(user *User, newAssetType string) (*AssetRepresentation.TypeTracker, error) {
	network := user.network

	contract := network.GetContract("inventoryMgt")

	result, err := contract.SubmitTransaction("AddAssetType", newAssetType)
	if err != nil {
		return nil, fmt.Errorf("failed to Submit transaction: %v", err)
	}

	typeTracker := AssetRepresentation.TypeTracker{}

	err = json.Unmarshal(result, &typeTracker)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal type tracker: %v", err)
	}

	return &typeTracker, nil
}

func GetAssetsByType(user *User, assetType string) (*AssetRepresentation.AssetList, error) {
	network := user.network

	contract := network.GetContract("inventoryMgt")

	result, err := contract.SubmitTransaction("GetAssetByType", assetType)
	if err != nil {
		return nil, fmt.Errorf("failed to Submit transaction: %v", err)
	}

	assetList := AssetRepresentation.AssetList{}

	err = json.Unmarshal(result, &assetList)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal asset: %v", err)
	}

	return &assetList, nil
}

func GetAsset(user *User, assetId string) (*AssetRepresentation.Asset, error) {
	network := user.network

	contract := network.GetContract("inventoryMgt")

	result, err := contract.SubmitTransaction("GetAsset", assetId)
	if err != nil {
		return nil, fmt.Errorf("failed to Submit transaction: %v", err)
	}

	asset := AssetRepresentation.Asset{}

	err = json.Unmarshal(result, &asset)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal asset: %v", err)
	}

	return &asset, nil
}

func RegisterAsset(user *User, asset *AssetRepresentation.Asset) (*AssetRepresentation.Asset, error) {
	network := user.network

	contract := network.GetContract("inventoryMgt")

	assetJson, err := json.Marshal(asset)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal asset")
	}

	result, err := contract.SubmitTransaction("RegisterAsset", string(assetJson))
	if err != nil {
		return nil, fmt.Errorf("failed to Submit transaction: %v", err)
	}

	newAsset := AssetRepresentation.Asset{}

	err = json.Unmarshal(result, &newAsset)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal asset: %v", err)
	}

	return &newAsset, nil
}

func RegisterAssetPlan(user *User, asset *AssetRepresentation.Asset) (*AssetRepresentation.Asset, error) {
	network := user.network

	contract := network.GetContract("inventoryMgt")

	assetJson, err := json.Marshal(asset)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal asset")
	}

	result, err := contract.EvaluateTransaction("RegisterAsset", string(assetJson))
	if err != nil {
		return nil, fmt.Errorf("failed to Submit transaction: %v", err)
	}

	newAsset := AssetRepresentation.Asset{}

	err = json.Unmarshal(result, &newAsset)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal asset: %v", err)
	}

	return &newAsset, nil
}

func ModifyAsset(user *User, asset *AssetRepresentation.Asset) (*AssetRepresentation.Asset, error) {
	network := user.network

	contract := network.GetContract("inventoryMgt")

	assetJson, err := json.Marshal(asset)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal asset")
	}

	result, err := contract.SubmitTransaction("ModifyAsset", string(assetJson))
	if err != nil {
		return nil, fmt.Errorf("failed to Submit transaction: %v", err)
	}

	newAsset := AssetRepresentation.Asset{}

	err = json.Unmarshal(result, &newAsset)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal asset: %v", err)
	}

	return &newAsset, nil
}

func ModifyAssetPlan(user *User, asset *AssetRepresentation.Asset) (*AssetRepresentation.Asset, error) {
	network := user.network

	contract := network.GetContract("inventoryMgt")

	assetJson, err := json.Marshal(asset)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal asset")
	}

	result, err := contract.EvaluateTransaction("ModifyAsset", string(assetJson))
	if err != nil {
		return nil, fmt.Errorf("failed to Submit transaction: %v", err)
	}

	newAsset := AssetRepresentation.Asset{}

	err = json.Unmarshal(result, &newAsset)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal asset: %v", err)
	}

	return &newAsset, nil
}

func ConfirmAsset(user *User, assetId string) (*AssetRepresentation.Asset, error) {
	network := user.network

	contract := network.GetContract("inventoryMgt")

	result, err := contract.SubmitTransaction("ConfirmAsset", assetId)
	if err != nil {
		return nil, fmt.Errorf("failed to Submit transaction: %v", err)
	}

	asset := AssetRepresentation.Asset{}

	err = json.Unmarshal(result, &asset)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal asset: %v", err)
	}

	return &asset, nil
}

func RemoveAsset(user *User, assetId string) error {
	network := user.network

	contract := network.GetContract("inventoryMgt")

	_, err := contract.SubmitTransaction("RemoveAsset", assetId)
	if err != nil {
		return fmt.Errorf("failed to Submit transaction: %v", err)
	}

	return nil
}

func ListDependencies(user *User, assetId string) (*AssetRepresentation.DependencyList, error) {
	network := user.network

	contract := network.GetContract("inventoryMgt")

	result, err := contract.SubmitTransaction("ListDependencies", assetId)
	if err != nil {
		return nil, fmt.Errorf("failed to Submit transaction: %v", err)
	}

	depList := AssetRepresentation.DependencyList{}

	err = json.Unmarshal(result, &depList)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dependency list: %v", err)
	}

	return &depList, nil
}

func ListDependants(user *User, assetId string) (*AssetRepresentation.DependantList, error) {
	network := user.network

	contract := network.GetContract("inventoryMgt")

	result, err := contract.SubmitTransaction("ListDependants", assetId)
	if err != nil {
		return nil, fmt.Errorf("failed to Submit transaction: %v", err)
	}

	depList := AssetRepresentation.DependantList{}

	err = json.Unmarshal(result, &depList)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dependant list: %v", err)
	}

	return &depList, nil
}

func AddDependency(user *User, assetId string, dependencyId string, originId string) error {
	network := user.network

	contract := network.GetContract("inventoryMgt")

	_, err := contract.SubmitTransaction("AddDependency", assetId, dependencyId, originId)
	if err != nil {
		return fmt.Errorf("failed to Submit transaction: %v", err)
	}
	return nil
}

func RemoveDependency(user *User, assetId string, dependencyId string) error {
	network := user.network

	contract := network.GetContract("inventoryMgt")

	_, err := contract.SubmitTransaction("RemoveDependency", assetId, dependencyId)
	if err != nil {
		return fmt.Errorf("failed to Submit transaction: %v", err)
	}
	return nil
}



/* Applied Tools */

func CreateAppliedTool(user *User, appliedTool *AssetRepresentation.AppliedTool) (*AssetRepresentation.AppliedTool, error) {
	network := user.network

	contract := network.GetContract("inventoryMgt")

	appliedToolJson, err := json.Marshal(appliedTool)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal asset")
	}

	result, err := contract.SubmitTransaction("CreateAppliedTool", string(appliedToolJson))
	if err != nil {
		return nil, fmt.Errorf("failed to Submit transaction: %v", err)
	}

	newAppliedTool := AssetRepresentation.AppliedTool{}

	err = json.Unmarshal(result, &newAppliedTool)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal applied tool: %v", err)
	}

	return &newAppliedTool, nil
}

func CreateAppliedToolPlan(user *User, appliedTool *AssetRepresentation.AppliedTool) (*AssetRepresentation.AppliedTool, error) {
	network := user.network

	contract := network.GetContract("inventoryMgt")

	appliedToolJson, err := json.Marshal(appliedTool)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal asset")
	}

	result, err := contract.EvaluateTransaction("CreateAppliedTool", string(appliedToolJson))
	if err != nil {
		return nil, fmt.Errorf("failed to Submit transaction: %v", err)
	}

	newAppliedTool := AssetRepresentation.AppliedTool{}

	err = json.Unmarshal(result, &newAppliedTool)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal applied tool: %v", err)
	}

	return &newAppliedTool, nil
}

func FinishAppliedTool(user *User, appliedTool *AssetRepresentation.AppliedTool, appliedTo []string) (*AssetRepresentation.AppliedTool, error) {
	network := user.network

	contract := network.GetContract("inventoryMgt")

	appliedToolJson, err := json.Marshal(appliedTool)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal asset")
	}

	appliedToString := ""
	first := true
	for _, s := range appliedTo {
		if !first {
			appliedToString += ","
			first = false
		}
		appliedToString += s
	}

	result, err := contract.SubmitTransaction("FinishAppliedTool", string(appliedToolJson), appliedToString)
	if err != nil {
		return nil, fmt.Errorf("failed to Submit transaction: %v", err)
	}

	newAppliedTool := AssetRepresentation.AppliedTool{}

	err = json.Unmarshal(result, &newAppliedTool)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal applied tool: %v", err)
	}

	return &newAppliedTool, nil
}


func RevertAppliedTool(user *User, appliedTool *AssetRepresentation.AppliedTool) (*AssetRepresentation.AppliedTool, error) {
	network := user.network

	contract := network.GetContract("inventoryMgt")

	appliedToolJson, err := json.Marshal(appliedTool)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal asset")
	}

	result, err := contract.SubmitTransaction("RevertAppliedTool", string(appliedToolJson))
	if err != nil {
		return nil, fmt.Errorf("failed to Submit transaction: %v", err)
	}

	newAppliedTool := AssetRepresentation.AppliedTool{}

	err = json.Unmarshal(result, &newAppliedTool)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal applied tool: %v", err)
	}

	return &newAppliedTool, nil
}


func GetAppliedTool(user *User, appliedToolId string) (*AssetRepresentation.AppliedTool, error) {
	network := user.network

	contract := network.GetContract("inventoryMgt")

	result, err := contract.SubmitTransaction("GetAppliedTool", appliedToolId)
	if err != nil {
		return nil, fmt.Errorf("failed to Submit transaction: %v", err)
	}

	newAppliedTool := AssetRepresentation.AppliedTool{}

	err = json.Unmarshal(result, &newAppliedTool)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal applied tool: %v", err)
	}

	return &newAppliedTool, nil
}







