package APIs

type PlanResp struct {
	AssetList *AssetListResp `json:"asset_list"` // only filled if the assets are new!!!
	AppliedTool *AppliedToolResp `json:"applied_tool"`
}

type ExecuteResp struct {
	Id string `json:"id"`
	AssetList *AssetListResp `json:"asset_list"` // only filled if the assets are new!!!
	AppliedTool *AppliedToolResp `json:"applied_tool"`
}
