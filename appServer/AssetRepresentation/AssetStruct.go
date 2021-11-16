package AssetRepresentation

type AppliedTool struct {
	ID       string `json:"id"`
	AppliedTo []string `json:"applied_to"`
	AssociatedDependencies []string `json:"associated_dependencies"`

	ToolName string `json:"tool_name"`
	FileName string `json:"file_name"`
	FileHash string `json:"file_hash"`
	Finished bool   `json:"finished"`
	FinalState string `json:"final_state"`
	Reverted string `json:"reverted"`
}

type DependencyRelation struct {
	Dependency string `json:"dependency"`
	OriginID string `json:"origin"`
}

type DependantRelation struct {
	Dependant string `json:"dependant"`
	OriginID string `json:"origin"`
}

// Asset is the main asset type struct (represents all infrastructure elements)
type Asset struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Location string `json:"location"`
	//Privacy int `json:"privacy"`
	Owner string `json:"owner"`

	SpecRamGB    int `json:"spec_ram_gb"`
	SpecCpuCores int `json:"spec_cpu_cores"`

	AppliedTools []string `json:"applied_tools"`

	IpAddrs []string `json:"ip_addrs"`

	Dependencies []*DependencyRelation `json:"dependencies"`
	Dependants   []*DependantRelation `json:"dependants"`

	Implemented bool `json:"implemented"`
}



type TypeTracker struct {
	AssetTypes []string `json:"asset_types"`
}

type AssetList struct {
	Assets []*Asset `json:"assets"`
}

type DependencyList struct {
	Dependencies []*DependencyRelation `json:"dependencies"`
}

type DependantList struct {
	Dependants   []*DependantRelation `json:"dependants"`
}
