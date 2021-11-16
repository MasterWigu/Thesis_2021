package main

import (
	"fmt"
	"github.com/MasterWigu/Thesis/appServer/APIs"
	"mime/multipart"
)

var processor *Processor

type Processor struct {
}

func createProcessor() {
	processor = new(Processor)
}

func (p *Processor) GetAssetTypes(user *User) (*APIs.AssetTypesResp, error) {
	types, err := ledger.getAssetTypes(user)
	if err != nil {
		return nil, fmt.Errorf("unable to get asset types: %v", err)
	}

	ATR := APIs.AssetTypesResp{
		AssetTypes: types.AssetTypes,
	}
	return &ATR, nil
}

func (p *Processor) AddAssetType(user *User, newType string) (*APIs.AssetTypesResp, error) {
	types, err := ledger.addAssetType(user, newType)
	if err != nil {
		return nil, fmt.Errorf("unable to create asset type: %v", err)
	}

	ATR := APIs.AssetTypesResp{
		AssetTypes: types.AssetTypes,
	}
	return &ATR, nil
}


func (p *Processor) GetAssetsByType(user *User, assetType string) (*APIs.AssetListResp, error) {
	assets, err := ledger.getAssetsByType(user, assetType)
	if err != nil {
		return nil, fmt.Errorf("unable to create asset type: %v", err)
	}

	ALR := APIs.AssetListResp{
		AssetList: assets.AssetList,
	}
	return &ALR, nil
}

func (p *Processor) GetAssetById(user *User, assetId string) (*APIs.AssetResp, error) {
	asset, err := ledger.getAssetById(user, assetId)
	if err != nil {
		return nil, fmt.Errorf("unable to create asset type: %v", err)
	}

	return asset, nil
}


func (p *Processor) RegisterAsset(user *User, newAsset *APIs.AssetResp) (*APIs.AssetResp, error) {

	asset, err := ledger.registerAsset(user, newAsset.Asset)
	if err != nil {
		return nil, fmt.Errorf("unable to create asset: %v", err)
	}

	return asset, nil
}

func (p *Processor) ModifyAsset(user *User, newAsset *APIs.AssetResp) (*APIs.AssetResp, error) {

	asset, err := ledger.modifyAsset(user, newAsset.Asset)
	if err != nil {
		return nil, fmt.Errorf("unable to modify asset: %v", err)
	}

	return asset, nil
}


func (p *Processor) RemoveAsset(user *User, assetId string) error {

	_, err := ledger.deleteAsset(user, assetId)
	if err != nil {
		return fmt.Errorf("unable to remove asset: %v", err)
	}

	return nil
}




func (p *Processor) ExecuteAction(user *User, tool string, file *multipart.FileHeader) (*APIs.ExecuteResp, error) {

	resp, err := modules.executeAction(tool, file)
	if err != nil {
		return nil, fmt.Errorf("unable to execute action: %v", err)
	}

	if resp.AssetList != nil && resp.AssetList.AssetList != nil {
		for _, asset := range resp.AssetList.AssetList {
			if asset == nil {
				return nil, fmt.Errorf("asset is nil")
			}
			if asset.ID != "" {
				assetResp, err := ledger.modifyAsset(user, asset)
				if err != nil {
					return nil, fmt.Errorf("error modifying asset in ledger: %v", err)
				}
				asset = assetResp.Asset
			} else {
				assetResp, err := ledger.registerAsset(user, asset)
				if err != nil {
					return nil, fmt.Errorf("error registering asset in ledger: %v", err)
				}
				asset = assetResp.Asset
			}
		}
	}

	appTool, err := ledger.registerAppliedTool(user, resp.AppliedTool.AppliedTool)
	if err != nil {
		return nil, fmt.Errorf("error registering applied tool in ledger: %v", err)
	}

	resp.AppliedTool = appTool

	err = modules.informActionIdOnLedger(tool, resp.Id, resp.AppliedTool.AppliedTool.ID)
	if err != nil {
		return nil, fmt.Errorf("error registering aplied tool id on module storage: %v", err)
	}

	resp.Id = resp.AppliedTool.AppliedTool.ID

	return resp, nil
}

func (p *Processor) PlanAction(user *User, tool string, file *multipart.FileHeader) (*APIs.PlanResp, error) {

	resp, err := modules.planAction(tool, file)
	if err != nil {
		return nil, fmt.Errorf("unable to plan action: %v", err)
	}

	if resp.AssetList != nil && resp.AssetList.AssetList != nil {
		for _, asset := range resp.AssetList.AssetList {
			if asset == nil {
				return nil, fmt.Errorf("asset is nil")
			}
			if asset.ID != "" {
				assetResp, err := ledger.modifyAssetPlan(user, asset)
				if err != nil {
					return nil, fmt.Errorf("error modifying asset in ledger: %v", err)
				}
				asset = assetResp.Asset
			} else {
				assetResp, err := ledger.registerAssetPlan(user, asset)
				if err != nil {
					return nil, fmt.Errorf("error registering asset in ledger: %v", err)
				}
				asset = assetResp.Asset
			}
		}
	}

	appTool, err := ledger.registerAppliedToolPlan(user, resp.AppliedTool.AppliedTool)
	if err != nil {
		return nil, fmt.Errorf("error registering applied tool in ledger: %v", err)
	}

	appTool.AppliedTool.ID = "" // delete the generated id because is doesn't really exist in the ledger
	resp.AppliedTool = appTool

	return resp, nil
}

func (p *Processor) ConfirmAction(user *User, tool string, actionId string) (*APIs.ExecuteResp, error) {

	resp, err := modules.confirmAction(tool, actionId)

	if err != nil {
		return nil, fmt.Errorf("unable to confirm action: %v", err)
	}

	if resp.AssetList != nil && resp.AssetList.AssetList != nil {
		for _, asset := range resp.AssetList.AssetList {
			if asset == nil {
				return nil, fmt.Errorf("asset is nil")
			}
			if asset.ID != "" {
				assetResp, err := ledger.modifyAsset(user, asset)
				if err != nil {
					return nil, fmt.Errorf("error modifying asset in ledger: %v", err)
				}
				asset = assetResp.Asset
			} else {
				assetResp, err := ledger.registerAsset(user, asset)
				if err != nil {
					return nil, fmt.Errorf("error registering asset in ledger: %v", err)
				}
				asset = assetResp.Asset
			}
		}
	}

	appTool, err := ledger.registerAppliedTool(user, resp.AppliedTool.AppliedTool)
	if err != nil {
		return nil, fmt.Errorf("error registering applied tool in ledger: %v", err)
	}

	resp.AppliedTool = appTool

	return resp, nil

}


