package chaincode

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/MasterWigu/Thesis/appServer/AssetRepresentation"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"math/rand"
	"strings"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Init

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	typeTracker := AssetRepresentation.TypeTracker{
		AssetTypes: make([]string, 0),
	}

	typeTracker.AssetTypes = append(typeTracker.AssetTypes, "Server", "VM")

	newTypeTrackerJSON, err := json.Marshal(typeTracker)
	if err != nil {
		return fmt.Errorf("could not marshal counter: %v", err)
	}

	err = ctx.GetStub().PutState("TypeTracker", newTypeTrackerJSON)
	if err != nil {
		return fmt.Errorf("could not put typeTracker in world state: %v", err)
	}
	return nil
}


// Asset Types

func (s *SmartContract) GetAssetTypes(ctx contractapi.TransactionContextInterface) (*AssetRepresentation.TypeTracker, error) {
	typeTrackerJSON, err := ctx.GetStub().GetState("TypeTracker")
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	if typeTrackerJSON == nil {
		return nil, fmt.Errorf("could not find typeTracker in world state")
	}

	typeTracker := AssetRepresentation.TypeTracker{}

	err = json.Unmarshal(typeTrackerJSON, &typeTracker)
	if err != nil {
		return  nil, fmt.Errorf("failed to unmarshal typeTracker: %v", err)
	}

	return &typeTracker, nil
}

func (s *SmartContract) AddAssetType(ctx contractapi.TransactionContextInterface, assetType string) (*AssetRepresentation.TypeTracker, error) {
	// Only admins can add asset types

	isAdmin, err := s.CheckUserAdmin(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not verify if user is admin: %v", err)
	}
	if !isAdmin {
		return nil, fmt.Errorf("user is not admin, access denied")
	}


	typeTrackerJSON, err := ctx.GetStub().GetState("TypeTracker")
	if err != nil {
		return  nil, fmt.Errorf("failed to read from world state: %v", err)
	}


	typeTracker := AssetRepresentation.TypeTracker{}

	err = json.Unmarshal(typeTrackerJSON, &typeTracker)
	if err != nil {
		return  nil, fmt.Errorf("failed to unmarshal typeTracker: %v", err)
	}

	if s.Contains(typeTracker.AssetTypes, assetType) {
		return &typeTracker, nil
	}

	typeTracker.AssetTypes = append(typeTracker.AssetTypes, assetType)


	newTypeTrackerJSON, err := json.Marshal(typeTracker)
	if err != nil {
		return nil, fmt.Errorf("could not marshal counter: %v", err)
	}

	err = ctx.GetStub().PutState("TypeTracker", newTypeTrackerJSON)
	if err != nil {
		return nil, fmt.Errorf("could not put typeTracker in world state: %v", err)
	}

	return &typeTracker, nil
}

// Asset Management

func (s *SmartContract) GetAssetByType(ctx contractapi.TransactionContextInterface, assetType string) (*AssetRepresentation.AssetList, error) {
	isAdmin, err := s.CheckUserAdmin(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not verify if user is admin: %v", err)
	}

	clientID, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get client ID: %v", err)
	}


	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(assetType, make([]string, 0))
	if err != nil {
		return nil, fmt.Errorf("could not get world state: %v", err)
	}
	defer resultsIterator.Close()

	assetList := AssetRepresentation.AssetList{}
	assetList.Assets = make([]*AssetRepresentation.Asset,0)

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("could not get next state in iterator: %v", err)
		}

		var asset AssetRepresentation.Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal asset: %v", err)
		}

		// if client is admin, can see all assets
		if isAdmin || asset.Owner == clientID {
			assetList.Assets = append(assetList.Assets, &asset)
		}
	}


	return &assetList, nil
}

func (s *SmartContract) GetAsset(ctx contractapi.TransactionContextInterface, assetId string) (*AssetRepresentation.Asset, error) {
	isAdmin, err := s.CheckUserAdmin(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not verify if user is admin: %v", err)
	}

	clientID, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get client ID: %v", err)
	}

	asset, _, err := GetAssetIntern(s, ctx, assetId)
	if err != nil {
		return nil, fmt.Errorf("could not get asset: %v", err)
	}

	if !(isAdmin || asset.Owner == clientID) {
		return nil, fmt.Errorf("asset with id %v not found", assetId)
	}

	return asset, nil
}

func (s *SmartContract) RegisterAsset(ctx contractapi.TransactionContextInterface, asset *AssetRepresentation.Asset) (*AssetRepresentation.Asset, error) {
	clientID, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get client ID: %v", err)
	}

	// validate type
	trackerJson, err := ctx.GetStub().GetState("TypeTracker")
	if err != nil {
		return nil, fmt.Errorf("could not get world state: %v", err)
	}

	if trackerJson == nil {
		return nil, fmt.Errorf(" Type Tracker not found")
	}

	tracker := AssetRepresentation.TypeTracker{}

	err = json.Unmarshal(trackerJson, &tracker)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal asset: %v", err)
	}

	if !s.Contains(tracker.AssetTypes, asset.Type) {
		return nil, fmt.Errorf("specified asset type does not exist")
	}

	// generate random asset id
	newRandID := s.RandStringBytes(7)

	newID, err := ctx.GetStub().CreateCompositeKey(asset.Type, []string{newRandID})

	newIDReadable, err := s.CompToReadableKey(ctx, newID)
	if err != nil {
		return nil, fmt.Errorf("could not parse generated id for asset: %v", err)
	}

	// populate remaining fields
	asset.ID = newIDReadable
	asset.Owner = clientID
	asset.Implemented = false

	// doesn't make sense to creat an asset with dependants
	asset.Dependants = make([]*AssetRepresentation.DependantRelation, 0)

	// populate if empty
	if asset.IpAddrs == nil || len(asset.IpAddrs) == 0 {
		asset.IpAddrs = make([]string, 0)
	}

	// we have to check if any dependencies are created and implement them as dependants too
	// for that, we get all dependency assets, and add this new asset as dependant
	if asset.Dependencies == nil || len(asset.Dependencies) == 0 {
		asset.Dependencies = make([]*AssetRepresentation.DependencyRelation, 0)
	} else {
		for _, dependencyObj := range asset.Dependencies {
			dependencyAsset, depAssCompId, err := GetAssetIntern(s, ctx, dependencyObj.Dependency)
			if err != nil {
				return nil, fmt.Errorf("could not get asset for dependency creation: %v", err)
			}

			newDependant := AssetRepresentation.DependantRelation{
				Dependant: newIDReadable,
				OriginID:  "Dependant_Creation",
			}

			dependencyAsset.Dependants = append(dependencyAsset.Dependants, &newDependant)
			dependencyAssetJson, err := json.Marshal(dependencyAsset)
			if err != nil {
				return nil, fmt.Errorf("could not marshal asset for dependency creation: %v", err)
			}


			err = ctx.GetStub().PutState(depAssCompId, dependencyAssetJson)
			if err != nil {
				return nil, fmt.Errorf("could not update asset in ledger for dependency creation: %v", err)
			}

		}
	}


	if asset.AppliedTools == nil || len(asset.AppliedTools) == 0 {
		asset.AppliedTools = make([]string, 0)
	} else if len(asset.AppliedTools) > 1 {
		return nil, fmt.Errorf("too many applied tools for asset creation: %v", err)
	} else {
		_, _, err = GetAssetIntern(s, ctx, asset.AppliedTools[0])
		if err != nil {
			return nil, fmt.Errorf("unable to find applied tool asset: %v", err)
		}
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return nil, fmt.Errorf("could not marshal asset: %v", err)
	}


	err = ctx.GetStub().PutState(newID,assetJSON)
	if err != nil {
		return nil, fmt.Errorf("could not put asset in ledger: %v", err)
	}

	return asset, nil
}

func (s *SmartContract) ModifyAsset(ctx contractapi.TransactionContextInterface, newAsset *AssetRepresentation.Asset) (*AssetRepresentation.Asset, error) {
	isAdmin, err := s.CheckUserAdmin(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not verify if user is admin: %v", err)
	}

	clientID, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get client ID: %v", err)
	}

	oldAsset, assetCompId, err := GetAssetIntern(s, ctx, newAsset.ID)
	if err != nil {
		return nil, fmt.Errorf("could not get asset: %v", err)
	}

	if !(isAdmin || oldAsset.Owner == clientID) {
		return nil, fmt.Errorf("asset with id %v not found", newAsset.ID)
	}

	// TODO check modified fields, if every change is acceptable, treat new dependencies, do not allow new dependants

	// check modified fields
	// no need to check ID
	// type
	if oldAsset.Type != newAsset.Type {
		return nil, fmt.Errorf("type change not allowed")
	}

	// no need to check location

	// owner
	if oldAsset.Owner != newAsset.Owner {
		return nil, fmt.Errorf("owner change not allowed")
	}

	// no need to check ram or cpu cores

	// applied tools must be valid (equal but can have new ones)
	for i, appliedToolID := range newAsset.AppliedTools {
		if i < len(oldAsset.AppliedTools) {
			if appliedToolID != oldAsset.AppliedTools[i] {
				return nil, fmt.Errorf("applied tool removal not allowed")
			}
		}

		// new applied tools
		_, _, err = GetAssetIntern(s, ctx, appliedToolID)
		if err != nil {
			return nil, fmt.Errorf("cannot verify applied tool: %v", err)
		}
	}

	// no need to verify ip addrs

	// do not allow dependency changing (use correct method for that)
	if len(oldAsset.Dependencies) != len(newAsset.Dependencies) {
		return nil, fmt.Errorf("dependency changing not allowed")
	}
	for _, newDep := range newAsset.Dependencies {
		found := false
		for _, oldDep := range oldAsset.Dependencies {
			if newDep.OriginID == oldDep.OriginID && newDep.Dependency == oldDep.Dependency {
				found = true
			}
		}
		if !found {
			return nil, fmt.Errorf("dependency changing not allowed")
		}
	}

	// do not allow dependants changing (use correct method for that)
	if len(oldAsset.Dependants) != len(newAsset.Dependants) {
		return nil, fmt.Errorf("dependency changing not allowed")
	}
	for _, newDep := range newAsset.Dependants {
		found := false
		for _, oldDep := range oldAsset.Dependants {
			if newDep.OriginID == oldDep.OriginID && newDep.Dependant == oldDep.Dependant {
				found = true
			}
		}
		if !found {
			return nil, fmt.Errorf("dependency changing not allowed")
		}
	}

	// do not allow change of "implemented"
	if newAsset.Implemented != oldAsset.Implemented {
		return nil, fmt.Errorf("implementation status change not allowed")
	}


	assetJSON, err := json.Marshal(newAsset)
	if err != nil {
		return nil, fmt.Errorf("could not marshal asset: %v", err)
	}


	err = ctx.GetStub().PutState(assetCompId,assetJSON)
	if err != nil {
		return nil, fmt.Errorf("could not put asset in ledger: %v", err)
	}

	return newAsset, nil
}

func (s *SmartContract) ConfirmAsset(ctx contractapi.TransactionContextInterface, assetId string) (*AssetRepresentation.Asset, error) {
	isAdmin, err := s.CheckUserAdmin(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not verify if user is admin: %v", err)
	}

	clientID, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get client ID: %v", err)
	}

	asset, compositeId, err := GetAssetIntern(s, ctx, assetId)
	if err != nil {
		return nil, fmt.Errorf("could not get asset: %v", err)
	}

	if !(isAdmin || clientID == asset.Owner) {
		return nil, fmt.Errorf("access denied to change asset")
	}

	asset.Implemented = true

	assetUpJson, err := json.Marshal(asset)
	if err != nil {
		return nil, fmt.Errorf("could not marshal asset: %v", err)
	}

	err = ctx.GetStub().PutState(compositeId, assetUpJson)
	if err != nil {
		return nil, fmt.Errorf("could not update asset in ledger: %v", err)
	}
	return asset, nil
}

func (s *SmartContract) RemoveAsset(ctx contractapi.TransactionContextInterface, assetId string) error {
	isAdmin, err := s.CheckUserAdmin(ctx)
	if err != nil {
		return fmt.Errorf("could not verify if user is admin: %v", err)
	}

	clientID, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return fmt.Errorf("could not get client ID: %v", err)
	}

	asset, compositeId, err := GetAssetIntern(s, ctx, assetId)
	if err != nil {
		return fmt.Errorf("could not get asset: %v", err)
	}

	if !(isAdmin || clientID == asset.Owner) {
		return fmt.Errorf("access denied to change asset")
	}

	err = ctx.GetStub().DelState(compositeId)
	if err != nil {
		return fmt.Errorf("could not remove asset in ledger: %v", err)
	}
	return nil
}

func (s *SmartContract) ListDependencies(ctx contractapi.TransactionContextInterface, assetId string) (*AssetRepresentation.DependencyList, error) {
	isAdmin, err := s.CheckUserAdmin(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not verify if user is admin: %v", err)
	}

	clientID, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get client ID: %v", err)
	}

	asset, _, err := GetAssetIntern(s, ctx, assetId)
	if err != nil {
		return nil, fmt.Errorf("could not get asset: %v", err)
	}

	if !(isAdmin || asset.Owner == clientID) {
		return nil, fmt.Errorf("asset with id %v not found", assetId)
	}

	depList := AssetRepresentation.DependencyList{
		Dependencies: make([]*AssetRepresentation.DependencyRelation, 0),
	}
	for _, dependencyObj := range asset.Dependencies {
		depList.Dependencies = append(depList.Dependencies, dependencyObj)
	}

	return &depList, nil
}

func (s *SmartContract) ListDependants(ctx contractapi.TransactionContextInterface, assetId string) (*AssetRepresentation.DependantList, error) {
	isAdmin, err := s.CheckUserAdmin(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not verify if user is admin: %v", err)
	}

	clientID, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get client ID: %v", err)
	}

	asset, _, err := GetAssetIntern(s, ctx, assetId)
	if err != nil {
		return nil, fmt.Errorf("could not get asset: %v", err)
	}

	if !(isAdmin || asset.Owner == clientID) {
		return nil, fmt.Errorf("asset with id %v not found", assetId)
	}

	depList := AssetRepresentation.DependantList{
		Dependants: make([]*AssetRepresentation.DependantRelation, 0),
	}
	for _, dependantObj := range asset.Dependants {
		depList.Dependants = append(depList.Dependants, dependantObj)
	}

	return &depList, nil
}

func (s *SmartContract) AddDependency(ctx contractapi.TransactionContextInterface, assetId string, dependencyId string, originId string) error {
	// AssetId is the asset that depends on dependencyId (i.e. the asset is the VM, the dependency is the Host)

	isAdmin, err := s.CheckUserAdmin(ctx)
	if err != nil {
		return fmt.Errorf("could not verify if user is admin: %v", err)
	}

	clientID, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return fmt.Errorf("could not get client ID: %v", err)
	}


	asset, compositeId, err := GetAssetIntern(s, ctx, assetId)
	if err != nil {
		return fmt.Errorf("could not get asset: %v", err)
	}

	dependency, depCompId, err := GetAssetIntern(s, ctx, dependencyId)
	if err != nil {
		return fmt.Errorf("could not get asset: %v", err)
	}


	// we have to be owners of the assets
	if !(isAdmin || (asset.Owner == clientID && dependency.Owner == clientID)) {
		return fmt.Errorf("permission denied to change assets")
	}

	if originId != "" {
		asset.AppliedTools = append(asset.AppliedTools, originId)
	}

	// if dependency already exists (not an error)
	for _, dep := range asset.Dependencies {
		if dep.Dependency == dependencyId {
			return nil
		}
	}

	newDependency := AssetRepresentation.DependencyRelation{
		Dependency: dependencyId,
		OriginID:   originId,
	}

	newDependant := AssetRepresentation.DependantRelation{
		Dependant: assetId,
		OriginID:  originId,
	}

	asset.Dependencies = append(asset.Dependencies, &newDependency)
	dependency.Dependants = append(dependency.Dependants, &newDependant)

	assetJson, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("could not marshal asset: %v", err)
	}

	dependencyJson, err := json.Marshal(dependency)
	if err != nil {
		return fmt.Errorf("could not marshal asset: %v", err)
	}

	err = ctx.GetStub().PutState(compositeId, assetJson)
	if err != nil {
		return fmt.Errorf("could not update asset in ledger: %v", err)
	}

	err = ctx.GetStub().PutState(depCompId, dependencyJson)
	if err != nil {
		return fmt.Errorf("could not update asset in ledger: %v", err)
	}

	return nil
}

func (s *SmartContract) RemoveDependency(ctx contractapi.TransactionContextInterface, assetId string, dependencyId string) error {
	isAdmin, err := s.CheckUserAdmin(ctx)
	if err != nil {
		return fmt.Errorf("could not verify if user is admin: %v", err)
	}

	clientID, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return fmt.Errorf("could not get client ID: %v", err)
	}


	asset, compositeId, err := GetAssetIntern(s, ctx, assetId)
	if err != nil {
		return fmt.Errorf("could not get asset: %v", err)
	}

	dependency, depCompId, err := GetAssetIntern(s, ctx, dependencyId)
	if err != nil {
		return fmt.Errorf("could not get asset: %v", err)
	}

	// we have to be owners of the assets
	if !(isAdmin || (asset.Owner == clientID && dependency.Owner == clientID)) {
		return fmt.Errorf("permission denied to change assets")
	}

	// remove dependency
	index := -1
	for i, dep := range asset.Dependencies {
		if dep.Dependency == dependencyId {
			index = i
		}
	}
	if index != -1 {
		// remove the item at position "index"
		asset.Dependencies[index] = asset.Dependencies[len(asset.Dependencies)-1]
		asset.Dependencies = asset.Dependencies[:len(asset.Dependencies)-1]
	}

	// remove dependant
	index = -1
	for i, dep := range asset.Dependants {
		if dep.Dependant == assetId {
			index = i
		}
	}
	if index != -1 {
		// remove the item at position "index"
		asset.Dependants[index] = asset.Dependants[len(asset.Dependants)-1]
		asset.Dependants = asset.Dependants[:len(asset.Dependants)-1]
	}

	assetJson, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("could not marshal asset: %v", err)
	}

	dependencyJson, err := json.Marshal(dependency)
	if err != nil {
		return fmt.Errorf("could not marshal asset: %v", err)
	}

	err = ctx.GetStub().PutState(compositeId, assetJson)
	if err != nil {
		return fmt.Errorf("could not update asset in ledger: %v", err)
	}

	err = ctx.GetStub().PutState(depCompId, dependencyJson)
	if err != nil {
		return fmt.Errorf("could not update asset in ledger: %v", err)
	}

	return nil
}


/*func (s *SmartContract) ChangeOwnership(ctx contractapi.TransactionContextInterface, id string, owner string) error {
	// TODO
	return nil
}*/


// Applied tool management

func (s *SmartContract) CreateAppliedTool(ctx contractapi.TransactionContextInterface, appliedTool *AssetRepresentation.AppliedTool) (*AssetRepresentation.AppliedTool, error) {
	isAdmin, err := s.CheckUserAdmin(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not verify if user is admin: %v", err)
	}

	clientID, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get client ID: %v", err)
	}

	if appliedTool.ToolName == "" {
		return nil, fmt.Errorf("tool name cannot be empty")
	}

	newRandID := s.RandStringBytes(7)

	newID, err := ctx.GetStub().CreateCompositeKey(appliedTool.ToolName, []string{newRandID})

	newIDReadable, err := s.CompToReadableKey(ctx, newID)
	if err != nil {
		return nil, fmt.Errorf("could not parse generated id for asset: %v", err)
	}


	appliedTool.ID = newIDReadable

	if (appliedTool.FileName == "" && appliedTool.FileHash != "") || (appliedTool.FileName != "" && appliedTool.FileHash == "") {
		return nil, fmt.Errorf("file and hash mismatch")
	}

	if appliedTool.Finished == true {
		appliedTool.Finished = true
	} else {
		appliedTool.Finished = false
	}

	appliedTool.Reverted = ""

	for _, assetId := range appliedTool.AppliedTo {
		asset, assetCompId, err := GetAssetIntern(s, ctx, assetId)
		if err != nil {
			return nil, fmt.Errorf("could not get asset: %v", err)
		}

		if !(isAdmin || clientID == asset.Owner) {
			return nil, fmt.Errorf("access denied to change asset")
		}

		asset.AppliedTools = append(asset.AppliedTools, appliedTool.ID)

		assetUpJson, err := json.Marshal(asset)
		if err != nil {
			return nil, fmt.Errorf("could not marshal asset: %v", err)
		}

		err = ctx.GetStub().PutState(assetCompId, assetUpJson)
		if err != nil {
			return nil, fmt.Errorf("could not update asset in ledger: %v", err)
		}
	}

	appliedToolJSON, err := json.Marshal(appliedTool)
	if err != nil {
		return nil, fmt.Errorf("could not marshal applied tool: %v", err)
	}

	err = ctx.GetStub().PutState(newID,appliedToolJSON)
	if err != nil {
		return nil, fmt.Errorf("could not put asset in ledger: %v", err)
	}


	return appliedTool, nil
}

func (s *SmartContract) FinishAppliedTool(ctx contractapi.TransactionContextInterface, appliedTool *AssetRepresentation.AppliedTool, appliedToString string) (*AssetRepresentation.AppliedTool, error) {
	// maybe verify owner of applied tool?

	oldAppliedTool, appToolCompId, err := GetAppliedToolIntern(s, ctx, appliedTool.ID)
	if err != nil {
		return nil, fmt.Errorf("could not get applied tool: %v", err)
	}

	oldAppliedTool.Finished = true
	oldAppliedTool.FinalState = appliedTool.FinalState

	appliedToolJSON, err := json.Marshal(oldAppliedTool)
	if err != nil {
		return nil, fmt.Errorf("could not marshal applied tool: %v", err)
	}

	err = ctx.GetStub().PutState(appToolCompId,appliedToolJSON)
	if err != nil {
		return nil, fmt.Errorf("could not put asset in ledger: %v", err)
	}

	appliedTo := strings.Split(appliedToString, ",")
	//add tool to assets
	for _, assetId := range appliedTo {
		for _, depId := range appliedTool.AssociatedDependencies {
			err = s.AddDependency(ctx, assetId, depId, appliedTool.ID)
			if err != nil {
				return nil, fmt.Errorf("unable to add associated dependencies: %v", err)
			}
		}

	}

	return oldAppliedTool, nil
}

func (s *SmartContract) RevertAppliedTool(ctx contractapi.TransactionContextInterface, appliedTool *AssetRepresentation.AppliedTool) (*AssetRepresentation.AppliedTool, error) {
	// maybe verify owner of applied tool?

	oldAppliedTool, appToolCompId, err := GetAppliedToolIntern(s, ctx, appliedTool.ID)
	if err != nil {
		return nil, fmt.Errorf("could not get applied tool: %v", err)
	}

	oldAppliedTool.Reverted = appliedTool.Reverted

	appliedToolJSON, err := json.Marshal(oldAppliedTool)
	if err != nil {
		return nil, fmt.Errorf("could not marshal applied tool: %v", err)
	}

	err = ctx.GetStub().PutState(appToolCompId,appliedToolJSON)
	if err != nil {
		return nil, fmt.Errorf("could not put asset in ledger: %v", err)
	}

	//TODO remove dependencies

	return oldAppliedTool, nil
}

func (s *SmartContract) GetAppliedTool(ctx contractapi.TransactionContextInterface, appliedToolId string) (*AssetRepresentation.AppliedTool, error) {
	oldAppliedTool, _, err := GetAppliedToolIntern(s, ctx, appliedToolId)
	if err != nil {
		return nil, fmt.Errorf("could not get applied tool: %v", err)
	}
	return oldAppliedTool, nil
}


// Aux Functions

func (s *SmartContract) CheckUserAdmin(ctx contractapi.TransactionContextInterface) (bool, error) {
	role, found, err := ctx.GetClientIdentity().GetAttributeValue("hf.Type")
	if err != nil {
		return false, fmt.Errorf("failed to read user attributes: %v", err)
	}
	if !found {
		return false, fmt.Errorf("could not find type attribute in user")
	}

	return role == "admin", nil
}

func (s *SmartContract) GetSubmittingClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {

	b64ID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed to read clientID: %v", err)
	}
	decodeID, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
	}
	return string(decodeID), nil
}

func (s *SmartContract) ReadableToCompKey(ctx contractapi.TransactionContextInterface, key string) (string, error) {
	components := strings.Split(key, "_")

	compKey, err := ctx.GetStub().CreateCompositeKey(components[0], components[1:])
	if err != nil {
		return "", fmt.Errorf("could not create composite key from readable key: %v", err)
	}
	return compKey, nil
}

func (s *SmartContract) CompToReadableKey(ctx contractapi.TransactionContextInterface, key string) (string, error) {
	s1, ss, err := ctx.GetStub().SplitCompositeKey(key)
	if err != nil {
		return "", fmt.Errorf("could not split composite key: %v", err)
	}

	outKey := s1
	for _, keyPart := range ss {
		outKey += "_" + keyPart
	}
	return outKey, nil
}

func GetAssetIntern(s *SmartContract, ctx contractapi.TransactionContextInterface, assetId string) (*AssetRepresentation.Asset, string, error){
	compositeId, err := s.ReadableToCompKey(ctx, assetId)
	if err != nil {
		return nil, "", fmt.Errorf("could not parse asset ID: %v", err)
	}

	assetJson, err := ctx.GetStub().GetState(compositeId)
	if err != nil {
		return nil, "", fmt.Errorf("could not get world state: %v", err)
	}
	if assetJson == nil {
		return nil, "", fmt.Errorf("asset not found")
	}


	asset := AssetRepresentation.Asset{}
	err = json.Unmarshal(assetJson, &asset)
	if err != nil {
		return nil, "", fmt.Errorf("could not unmarshal asset: %v", err)
	}

	return &asset, compositeId, nil
}

func GetAppliedToolIntern(s *SmartContract, ctx contractapi.TransactionContextInterface, appliedToolId string) (*AssetRepresentation.AppliedTool, string, error){
	compositeId, err := s.ReadableToCompKey(ctx, appliedToolId)
	if err != nil {
		return nil, "", fmt.Errorf("could not parse applied tool ID: %v", err)
	}

	appliedToolJson, err := ctx.GetStub().GetState(compositeId)
	if err != nil {
		return nil, "", fmt.Errorf("could not get world state: %v", err)
	}
	if appliedToolJson == nil {
		return nil, "", fmt.Errorf("applied tool not found")
	}


	appliedTool := AssetRepresentation.AppliedTool{}
	err = json.Unmarshal(appliedToolJson, &appliedTool)
	if err != nil {
		return nil, "", fmt.Errorf("could not unmarshal applied tool: %v", err)
	}

	return &appliedTool, compositeId, nil
}


const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
func (s *SmartContract) RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (s *SmartContract) Contains(list []string, elem string) bool {
	for _, s := range list {
		if s == elem {
			return true
		}
	}
	return false
}

