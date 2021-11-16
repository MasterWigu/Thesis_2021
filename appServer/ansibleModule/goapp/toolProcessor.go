package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/MasterWigu/Thesis/appServer/APIs"
	"github.com/MasterWigu/Thesis/appServer/AssetRepresentation"
	"mime/multipart"
)

var processor *Processor

type Processor struct {
}


func createProcessor() {
	processor = new(Processor)
}

func (p *Processor) planAction(files *multipart.FileHeader) (*APIs.PlanResp, error) {
	// Generate 8 byte session id
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("failed to generate action id")
	}
	actionId := "temp_" + base64.URLEncoding.EncodeToString(b)

	zipHash, err := fileHandler.SaveActionFiles(files, actionId)
	if err != nil {
		return nil, fmt.Errorf("failed to save temp action files")
	}

	err = fileHandler.PrepareRunningDir(actionId)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare running dir for temp action files")
	}

	assetIds, err := fileHandler.ProcessHosts(actionId)

	_, err = tool.DryRunAction(actionId)
	if err != nil {
		return nil, fmt.Errorf("failed to run tool: %v", err)
	}

	appliedTool := AssetRepresentation.AppliedTool{
		ID:                     "",
		AppliedTo:              assetIds,
		AssociatedDependencies: nil,
		ToolName:               "ansible",
		FileName:               actionId,
		FileHash:               zipHash,
		Finished:               false,
		FinalState:             "",
		Reverted:               "",
	}

	appliedToolResp := APIs.AppliedToolResp{
		AppliedTool: &appliedTool,
	}

	planResp := APIs.PlanResp{
		AssetList:   nil,
		AppliedTool: &appliedToolResp,
	}

	fileHandler.DeleteRunningFiles(actionId)
	fileHandler.DeleteActionFiles(actionId)

	return &planResp, nil
}

func (p *Processor) execAction(files *multipart.FileHeader) (*APIs.ExecuteResp, error) {
	// Generate 8 byte session id
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("failed to generate action id: %v", err)
	}
	actionId := base64.URLEncoding.EncodeToString(b)

	zipHash, err := fileHandler.SaveActionFiles(files, actionId)
	if err != nil {
		return nil, fmt.Errorf("failed to save temp action files: %v", err)
	}

	err = fileHandler.PrepareRunningDir(actionId)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare running dir for temp action files: %v", err)
	}

	assetIds, err := fileHandler.ProcessHosts(actionId)

	_, err = tool.DryRunAction(actionId)
	if err != nil {
		return nil, fmt.Errorf("failed to run tool: %v", err)
	}

	appliedTool := AssetRepresentation.AppliedTool{
		ID:                     "",
		AppliedTo:              assetIds,
		AssociatedDependencies: nil,
		ToolName:               "ansible",
		FileName:               actionId,
		FileHash:               zipHash,
		Finished:               false,
		FinalState:             "",
		Reverted:               "",
	}

	fileHandler.DeleteRunningFiles(actionId)


	appliedToolResp := APIs.AppliedToolResp{
		AppliedTool: &appliedTool,
	}

	executeResp := APIs.ExecuteResp{
		Id:          actionId,
		AssetList:   nil,
		AppliedTool: &appliedToolResp,
	}


	return &executeResp, nil
}

func (p *Processor) ConfirmAction(actionId string)  (*APIs.ExecuteResp, bool, error) {
	// bool represents if the action id was found
	if !fileHandler.RunningDirExists(actionId) {
		if !fileHandler.ActionZipExists(actionId) {
			return nil, false, nil
		}

		err := fileHandler.PrepareRunningDir(actionId)
		if err != nil {
			return nil, true, fmt.Errorf("failed to prepare running dir for temp action files: %v", err)
		}
	}
	assetIds, err := fileHandler.ProcessHosts(actionId)
	if err != nil {
		return nil, true, fmt.Errorf("failed to process target hosts: %v", err)
	}

	toolOutput, err := tool.ExecAction(actionId)
	if err != nil {
		return nil, true, fmt.Errorf("failed to run tool: %v", err)
	}

	appliedTool := AssetRepresentation.AppliedTool{
		ID:                     actionId,
		AppliedTo:              assetIds,
		AssociatedDependencies: nil,
		ToolName:               "ansible",
		FileName:               actionId,
		FileHash:               "",
		Finished:               true,
		FinalState:             toolOutput,
		Reverted:               "",
	}



	fileHandler.DeleteRunningFiles(actionId)

	appliedToolResp := APIs.AppliedToolResp{
		AppliedTool: &appliedTool,
	}

	executeResp := APIs.ExecuteResp{
		Id:          actionId,
		AssetList:   nil,
		AppliedTool: &appliedToolResp,
	}

	return &executeResp, true, nil
}

func (p *Processor) BindLedgerId(actionId string, ledgerId string) error {
	return fileHandler.changeActionId(actionId, ledgerId)
}


