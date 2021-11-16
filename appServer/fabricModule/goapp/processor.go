package main

import (
	"fmt"
	"github.com/MasterWigu/Thesis/appServer/APIs"
)

var processor *Processor

type Processor struct {

}

func CreateProcessor() {
	processor = new(Processor)
}

func (p *Processor) GetAssetTypes(user *User) (*APIs.AssetTypesResp, error) {
	types, err := GetAssetTypes(user)
	if err != nil {
		return nil, fmt.Errorf("unable to get asset types: %v", err)
	}
	
	ATR := APIs.AssetTypesResp{
		AssetTypes: types.AssetTypes,
	}
	return &ATR, nil
}

func (p *Processor) AddAssetType(user *User, newType string) (*APIs.AssetTypesResp, error) {
	types, err := AddAssetType(user, newType)
	if err != nil {
		return nil, fmt.Errorf("unable to create asset type: %v", err)
	}

	ATR := APIs.AssetTypesResp{
		AssetTypes: types.AssetTypes,
	}
	return &ATR, nil
}


func (p *Processor) GetAssetsByType(user *User, assetType string) (*APIs.AssetListResp, error) {
	assets, err := GetAssetsByType(user, assetType)
	if err != nil {
		return nil, fmt.Errorf("unable to create asset type: %v", err)
	}

	ALR := APIs.AssetListResp{
		AssetList: assets.Assets,
	}
	return &ALR, nil
}

func (p *Processor) GetAssetById(user *User, assetId string) (*APIs.AssetResp, error) {
	asset, err := GetAsset(user, assetId)
	if err != nil {
		return nil, fmt.Errorf("unable to create asset type: %v", err)
	}

	AR := APIs.AssetResp{
		Asset: asset,
	}
	return &AR, nil
}


func (p *Processor) RegisterAsset(user *User, newAsset *APIs.AssetResp) (*APIs.AssetResp, error) {

	asset, err := RegisterAsset(user, newAsset.Asset)
	if err != nil {
		return nil, fmt.Errorf("unable to create asset: %v", err)
	}

	AR := APIs.AssetResp{
		Asset: asset,
	}
	return &AR, nil
}

func (p *Processor) RegisterAssetPlan(user *User, newAsset *APIs.AssetResp) (*APIs.AssetResp, error) {

	asset, err := RegisterAssetPlan(user, newAsset.Asset)
	if err != nil {
		return nil, fmt.Errorf("unable to create asset: %v", err)
	}

	AR := APIs.AssetResp{
		Asset: asset,
	}
	return &AR, nil
}

func (p *Processor) ModifyAsset(user *User, newAsset *APIs.AssetResp) (*APIs.AssetResp, error) {

	asset, err := ModifyAsset(user, newAsset.Asset)
	if err != nil {
		return nil, fmt.Errorf("unable to modify asset: %v", err)
	}

	AR := APIs.AssetResp{
		Asset: asset,
	}
	return &AR, nil
}

func (p *Processor) ModifyAssetPlan(user *User, newAsset *APIs.AssetResp) (*APIs.AssetResp, error) {

	asset, err := ModifyAssetPlan(user, newAsset.Asset)
	if err != nil {
		return nil, fmt.Errorf("unable to modify asset: %v", err)
	}

	AR := APIs.AssetResp{
		Asset: asset,
	}
	return &AR, nil
}

func (p *Processor) ConfirmAsset(user *User, assetId string) (*APIs.AssetResp, error) {

	asset, err := ConfirmAsset(user, assetId)
	if err != nil {
		return nil, fmt.Errorf("unable to confirm asset: %v", err)
	}

	AR := APIs.AssetResp{
		Asset: asset,
	}
	return &AR, nil
}

func (p *Processor) RemoveAsset(user *User, assetId string) error {

	err := RemoveAsset(user, assetId)
	if err != nil {
		return fmt.Errorf("unable to remove asset: %v", err)
	}

	return nil
}

func (p *Processor) ListDependencies(user *User, assetId string) (*APIs.DependencyListResp, error) {

	deps, err := ListDependencies(user, assetId)
	if err != nil {
		return nil, fmt.Errorf("unable to list dependencies: %v", err)
	}

	DLR := APIs.DependencyListResp{
		DependencyList: deps.Dependencies,
	}
	return &DLR, nil
}

func (p *Processor) ListDependants(user *User, assetId string) (*APIs.DependantListResp, error) {

	deps, err := ListDependants(user, assetId)
	if err != nil {
		return nil, fmt.Errorf("unable to list dependendants: %v", err)
	}

	DLR := APIs.DependantListResp{
		DependantList: deps.Dependants,
	}
	return &DLR, nil
}

func (p *Processor) AddDependency(user *User, assetId string, dependencyId string, originId string) error {

	err := AddDependency(user, assetId, dependencyId, originId)
	if err != nil {
		return fmt.Errorf("unable to add dependencyt: %v", err)
	}

	return nil
}

func (p *Processor) RemoveDependency(user *User, assetId string, dependencyId string) error {

	err := RemoveDependency(user, assetId, dependencyId)
	if err != nil {
		return fmt.Errorf("unable to remove dependency: %v", err)
	}

	return nil
}

/* Applied Tools */

func (p *Processor) CreateAppliedTool(user *User, appliedTool *APIs.AppliedToolResp) (*APIs.AppliedToolResp, error) {
	appTool, err := CreateAppliedTool(user, appliedTool.AppliedTool)
	if err != nil {
		return nil, fmt.Errorf("unable to create applied tool: %v", err)
	}

	APR := APIs.AppliedToolResp{
		AppliedTool: appTool,
	}

	return &APR, nil
}

func (p *Processor) CreateAppliedToolPlan(user *User, appliedTool *APIs.AppliedToolResp) (*APIs.AppliedToolResp, error) {
	appTool, err := CreateAppliedToolPlan(user, appliedTool.AppliedTool)
	if err != nil {
		return nil, fmt.Errorf("unable to create applied tool: %v", err)
	}

	APR := APIs.AppliedToolResp{
		AppliedTool: appTool,
	}

	return &APR, nil
}

func (p *Processor) FinishAppliedTool(user *User, appliedTool *APIs.AppliedToolResp, appliedTo []string) (*APIs.AppliedToolResp, error) {
	appTool, err := FinishAppliedTool(user, appliedTool.AppliedTool, appliedTo)
	if err != nil {
		return nil, fmt.Errorf("unable to create applied tool: %v", err)
	}

	APR := APIs.AppliedToolResp{
		AppliedTool: appTool,
	}

	return &APR, nil
}

func (p *Processor) RevertAppliedTool(user *User, appliedTool *APIs.AppliedToolResp) (*APIs.AppliedToolResp, error) {
	appTool, err := RevertAppliedTool(user, appliedTool.AppliedTool)
	if err != nil {
		return nil, fmt.Errorf("unable to create applied tool: %v", err)
	}

	APR := APIs.AppliedToolResp{
		AppliedTool: appTool,
	}

	return &APR, nil
}

func (p *Processor) GetAppliedTool(user *User, appliedToolId string) (*APIs.AppliedToolResp, error) {
	appTool, err := GetAppliedTool(user, appliedToolId)
	if err != nil {
		return nil, fmt.Errorf("unable to create applied tool: %v", err)
	}

	APR := APIs.AppliedToolResp{
		AppliedTool: appTool,
	}

	return &APR, nil
}