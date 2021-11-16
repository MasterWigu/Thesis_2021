package APIs

import "github.com/MasterWigu/Thesis/appServer/AssetRepresentation"

type LoginResp struct {
	LedgerSessId string `json:"sess_id"`
	Username     string `json:"username"`
	Perms        string `json:"perms"`
}

type AssetTypesResp struct {
	AssetTypes []string `json:"asset_types"`
}

type AssetResp struct {
	Asset *AssetRepresentation.Asset `json:"asset"`
}

type AssetListResp struct {
	AssetList []*AssetRepresentation.Asset `json:"asset_list"`
}

type DependencyListResp struct {
	DependencyList []*AssetRepresentation.DependencyRelation `json:"dependency_list"`
}

type DependantListResp struct {
	DependantList []*AssetRepresentation.DependantRelation `json:"dependant_list"`
}

type AppliedToolResp struct {
	/*Targets []string `json:"targets"`*/
	AppliedTool *AssetRepresentation.AppliedTool `json:"applied_tool"`
}




