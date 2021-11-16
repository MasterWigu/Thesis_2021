package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/MasterWigu/Thesis/appServer/APIs"
	"github.com/MasterWigu/Thesis/appServer/AssetRepresentation"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
)

var ledger *Ledger

type Ledger struct {
	url string
}

var client *http.Client

func initLedger() {
	ledger = new(Ledger)
	ledger.url = "https://fabric-module:8090"


	caCert, err := ioutil.ReadFile("/certs/tls-cacert.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}


	return
}

/*	r.POST("/loginUser", loginHandler)

	r.GET("/assets/types", getAssetTypesHandler)
	r.POST("/assets/types/", addAssetTypeHandler)

	r.GET("/assets/:type/", getAssetsByTypeHandler)
	r.GET("/asset/:id/", getAssetHandler)
	r.POST("/asset/", registerAssetHandler)
	r.POST("/asset/modify/", modifyAssetHandler)
	r.PUT("/asset/:id/confirm/", confirmAssetHandler)
	r.DELETE("/asset/:id/", removeAssetHandler)
	r.GET("/asset/:id/dependencies/", listDependenciesHandler)
	r.GET("/asset/:id/dependants/", listDependantsHandler)
	r.POST("/asset/:assetId/dependencies/:dependencyId/", addDependencyHandler)
	//r.DELETE("/asset/:assetId/dependencies/:dependencyId/", removeDependencyHandler)

	r.GET("/appliedTool/:id", getAppliedToolHandler)
	r.POST("/appliedTool", createAppliedToolHandler)
	r.POST("/appliedTool/finish", finishAppliedToolHandler)
	r.POST("/appliedTool/revert", revertAppliedToolHandler)

 */


func (l *Ledger) loginUser(user *User) (bool, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("identityFile", "certs.zip")

	if err != nil {
		return false, err
	}

	io.Copy(part, user.credZip)
	writer.Close()
	request, err := http.NewRequest("POST", l.url + "/loginUser", body)

	if err != nil {
		return false, err
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())

	response, err := client.Do(request)

	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return false, err
	}

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusUnauthorized {
			return false, nil
		}
		return false, fmt.Errorf("internal ledger error")
	}

	fmt.Println(string(content))

	jsonResp := APIs.LoginResp{}
	err = json.Unmarshal(content, &jsonResp)
	if err != nil {
		return false, err
	}

	user.username = jsonResp.Username
	user.ledgerSessId = jsonResp.LedgerSessId
	user.userPerm = jsonResp.Perms

	return true, nil
}

func (l *Ledger) getAssetTypes(user *User) (*APIs.AssetTypesResp, error) {
	request, err := http.NewRequest("GET", l.url + "/assets/types", nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Cookie", "SESSION_ID="+user.ledgerSessId)

	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusUnauthorized {
			return nil, nil
		}
		return nil, fmt.Errorf("internal ledger error")
	}

	jsonResp := APIs.AssetTypesResp{}
	err = json.Unmarshal(content, &jsonResp)
	if err != nil {
		return nil, err
	}

	return &jsonResp, nil
}

func (l *Ledger) addAssetType(user *User, newType string) (*APIs.AssetTypesResp, error) {

	request, err := http.NewRequest("POST", l.url + "/assets/types?newType=" + newType, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Cookie", "SESSION_ID="+user.ledgerSessId)

	//request.Header.Add("Content-Type", "application/json")

	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusUnauthorized {
			return nil, nil
		}
		return nil, fmt.Errorf("internal ledger error")
	}

	jsonResp := APIs.AssetTypesResp{}
	err = json.Unmarshal(content, &jsonResp)
	if err != nil {
		return nil, err
	}

	return &jsonResp, nil
}

func (l *Ledger) getAssetsByType(user *User, typeName string) (*APIs.AssetListResp, error) {
	request, err := http.NewRequest("GET", l.url + "/assets/" + typeName, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Cookie", "SESSION_ID="+user.ledgerSessId)

	//request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusUnauthorized {
			return nil, nil
		}
		return nil, fmt.Errorf("internal ledger error")
	}

	jsonResp := APIs.AssetListResp{}
	err = json.Unmarshal(content, &jsonResp)
	if err != nil {
		return nil, err
	}

	return &jsonResp, nil
}

func (l *Ledger) getAssetById(user *User, assetId string) (*APIs.AssetResp, error) {
	request, err := http.NewRequest("GET", l.url + "/asset/" + assetId, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Cookie", "SESSION_ID="+user.ledgerSessId)

	//request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusUnauthorized {
			return nil, nil
		}
		return nil, fmt.Errorf("internal ledger error")
	}

	jsonResp := APIs.AssetResp{}
	err = json.Unmarshal(content, &jsonResp)
	if err != nil {
		return nil, err
	}

	return &jsonResp, nil
}

func (l *Ledger) registerAsset(user *User, asset *AssetRepresentation.Asset) (*APIs.AssetResp, error) {
	assetRequest := APIs.AssetResp{
		Asset: asset,
	}
	reqJSON, err := json.Marshal(assetRequest)
	if err != nil {
		return  nil, err
	}

	request, err := http.NewRequest("POST", l.url + "/asset/", bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, err
	}


	request.Header.Set("Cookie", "SESSION_ID="+user.ledgerSessId)
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusUnauthorized {
			return nil, nil
		}
		return nil, fmt.Errorf("internal ledger error")
	}

	jsonResp := APIs.AssetResp{}
	err = json.Unmarshal(content, &jsonResp)
	if err != nil {
		return nil, err
	}

	return &jsonResp, nil
}

func (l *Ledger) registerAssetPlan(user *User, asset *AssetRepresentation.Asset) (*APIs.AssetResp, error) {
	assetRequest := APIs.AssetResp{
		Asset: asset,
	}
	reqJSON, err := json.Marshal(assetRequest)
	if err != nil {
		return  nil, err
	}

	request, err := http.NewRequest("POST", l.url + "/asset/plan", bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, err
	}


	request.Header.Set("Cookie", "SESSION_ID="+user.ledgerSessId)
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusUnauthorized {
			return nil, nil
		}
		return nil, fmt.Errorf("internal ledger error")
	}

	jsonResp := APIs.AssetResp{}
	err = json.Unmarshal(content, &jsonResp)
	if err != nil {
		return nil, err
	}

	return &jsonResp, nil
}


func (l *Ledger) modifyAsset(user *User, asset *AssetRepresentation.Asset) (*APIs.AssetResp, error) {
	assetRequest := APIs.AssetResp{
		Asset: asset,
	}
	reqJSON, err := json.Marshal(assetRequest)
	if err != nil {
		return  nil, err
	}

	request, err := http.NewRequest("POST", l.url + "/asset/modify", bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, err
	}


	request.Header.Set("Cookie", "SESSION_ID="+user.ledgerSessId)
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusUnauthorized {
			return nil, nil
		}
		return nil, fmt.Errorf("internal ledger error")
	}

	jsonResp := APIs.AssetResp{}
	err = json.Unmarshal(content, &jsonResp)
	if err != nil {
		return nil, err
	}

	return &jsonResp, nil
}

func (l *Ledger) modifyAssetPlan(user *User, asset *AssetRepresentation.Asset) (*APIs.AssetResp, error) {
	assetRequest := APIs.AssetResp{
		Asset: asset,
	}
	reqJSON, err := json.Marshal(assetRequest)
	if err != nil {
		return  nil, err
	}

	request, err := http.NewRequest("POST", l.url + "/asset/modify/plan", bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, err
	}


	request.Header.Set("Cookie", "SESSION_ID="+user.ledgerSessId)
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusUnauthorized {
			return nil, nil
		}
		return nil, fmt.Errorf("internal ledger error")
	}

	jsonResp := APIs.AssetResp{}
	err = json.Unmarshal(content, &jsonResp)
	if err != nil {
		return nil, err
	}

	return &jsonResp, nil
}

func (l *Ledger) confirmAsset(user *User, assetId string) (*APIs.AssetResp, error) {
	request, err := http.NewRequest("PUT", l.url + "/asset/" + assetId + "/confirm", nil)
	if err != nil {
		return nil, err
	}


	request.Header.Set("Cookie", "SESSION_ID="+user.ledgerSessId)
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusUnauthorized {
			return nil, nil
		}
		return nil, fmt.Errorf("internal ledger error")
	}

	jsonResp := APIs.AssetResp{}
	err = json.Unmarshal(content, &jsonResp)
	if err != nil {
		return nil, err
	}

	return &jsonResp, nil
}


// bool success
func (l *Ledger) deleteAsset(user *User, assetId string) (bool, error) {
	request, err := http.NewRequest("DELETE", l.url + "/asset/" + assetId, nil)
	if err != nil {
		return false, err
	}


	request.Header.Set("Cookie", "SESSION_ID="+user.ledgerSessId)
	response, err := client.Do(request)

	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusUnauthorized {
			return false, nil
		}
		return false, fmt.Errorf("internal ledger error")
	}
	return true, nil
}


func (l *Ledger) getDependencies(user *User, assetId string) (*APIs.DependencyListResp, error) {
	request, err := http.NewRequest("GET", l.url + "/asset/" + assetId + "/dependencies", nil)
	if err != nil {
		return nil, err
	}


	request.Header.Set("Cookie", "SESSION_ID="+user.ledgerSessId)
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusUnauthorized {
			return nil, nil
		}
		return nil, fmt.Errorf("internal ledger error")
	}

	jsonResp := APIs.DependencyListResp{}
	err = json.Unmarshal(content, &jsonResp)
	if err != nil {
		return nil, err
	}

	return &jsonResp, nil
}

func (l *Ledger) getDependants(user *User, assetId string) (*APIs.DependantListResp, error) {
	request, err := http.NewRequest("GET", l.url + "/asset/" + assetId + "/dependants", nil)
	if err != nil {
		return nil, err
	}


	request.Header.Set("Cookie", "SESSION_ID="+user.ledgerSessId)
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusUnauthorized {
			return nil, nil
		}
		return nil, fmt.Errorf("internal ledger error")
	}

	jsonResp := APIs.DependantListResp{}
	err = json.Unmarshal(content, &jsonResp)
	if err != nil {
		return nil, err
	}

	return &jsonResp, nil
}


// bool success
func (l *Ledger) addDependency(user *User, assetId string, dependencyId string) (bool, error) {
	request, err := http.NewRequest("POST", l.url + "/asset/" + assetId + "/dependency/" + dependencyId, nil)
	if err != nil {
		return false, err
	}


	request.Header.Set("Cookie", "SESSION_ID="+user.ledgerSessId)
	response, err := client.Do(request)

	if err != nil {
		return false, err
	}
	defer response.Body.Close()


	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusUnauthorized {
			return false, nil
		}
		return false, fmt.Errorf("internal ledger error")
	}
	return true, nil
}

// bool success
func (l *Ledger) addDependencyWithOrigin(user *User, assetId string, dependencyId string, originId string) (bool, error) {
	request, err := http.NewRequest("POST", l.url + "/asset/" + assetId + "/dependency/" + dependencyId + "?originId=" + originId, nil)
	if err != nil {
		return false, err
	}


	request.Header.Set("Cookie", "SESSION_ID="+user.ledgerSessId)
	response, err := client.Do(request)

	if err != nil {
		return false, err
	}
	defer response.Body.Close()


	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusUnauthorized {
			return false, nil
		}
		return false, fmt.Errorf("internal ledger error")
	}
	return true, nil
}

// bool success
func (l *Ledger) removeDependency(user *User, assetId string, dependencyId string) (bool, error) {
	request, err := http.NewRequest("DELETE", l.url + "/asset/" + assetId + "/dependency/" + dependencyId, nil)
	if err != nil {
		return false, err
	}


	request.Header.Set("Cookie", "SESSION_ID="+user.ledgerSessId)
	response, err := client.Do(request)

	if err != nil {
		return false, err
	}
	defer response.Body.Close()


	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusUnauthorized {
			return false, nil
		}
		return false, fmt.Errorf("internal ledger error")
	}
	return true, nil
}


/* APPLIED TOOLS */

func (l *Ledger) getAppliedTool(user *User, appliedToolId string) (*APIs.AppliedToolResp, error) {
	request, err := http.NewRequest("GET", l.url + "/appliedTool/" + appliedToolId, nil)
	if err != nil {
		return nil, err
	}


	request.Header.Set("Cookie", "SESSION_ID="+user.ledgerSessId)
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusUnauthorized {
			return nil, nil
		}
		return nil, fmt.Errorf("internal ledger error")
	}

	jsonResp := APIs.AppliedToolResp{}
	err = json.Unmarshal(content, &jsonResp)
	if err != nil {
		return nil, err
	}

	return &jsonResp, nil
}

func (l *Ledger) registerAppliedTool(user *User, appliedTool *AssetRepresentation.AppliedTool) (*APIs.AppliedToolResp, error) {
	appToolReq := APIs.AppliedToolResp{
		AppliedTool: appliedTool,
	}

	reqJSON, err := json.Marshal(appToolReq)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", l.url + "/appliedTool", bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Cookie", "SESSION_ID="+user.ledgerSessId)
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusUnauthorized {
			return nil, nil
		}
		return nil, fmt.Errorf("internal ledger error")
	}

	jsonResp := APIs.AppliedToolResp{}
	err = json.Unmarshal(content, &jsonResp)
	if err != nil {
		return nil, err
	}

	return &jsonResp, nil
}

func (l *Ledger) registerAppliedToolPlan(user *User, appliedTool *AssetRepresentation.AppliedTool) (*APIs.AppliedToolResp, error) {
	appToolReq := APIs.AppliedToolResp{
		AppliedTool: appliedTool,
	}

	reqJSON, err := json.Marshal(appToolReq)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", l.url + "/appliedTool/plan", bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Cookie", "SESSION_ID="+user.ledgerSessId)
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusUnauthorized {
			return nil, nil
		}
		return nil, fmt.Errorf("internal ledger error")
	}

	jsonResp := APIs.AppliedToolResp{}
	err = json.Unmarshal(content, &jsonResp)
	if err != nil {
		return nil, err
	}

	return &jsonResp, nil
}


func (l *Ledger) finishAppliedTool(user *User, appliedTool *AssetRepresentation.AppliedTool) (*APIs.AppliedToolResp, error) {
	appToolReq := APIs.AppliedToolResp{
		AppliedTool: appliedTool,
	}

	reqJSON, err := json.Marshal(appToolReq)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", l.url + "/appliedTool", bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Cookie", "SESSION_ID="+user.ledgerSessId)
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusUnauthorized {
			return nil, nil
		}
		return nil, fmt.Errorf("internal ledger error")
	}

	jsonResp := APIs.AppliedToolResp{}
	err = json.Unmarshal(content, &jsonResp)
	if err != nil {
		return nil, err
	}

	return &jsonResp, nil
}

func (l *Ledger) revertAppliedTool(user *User, appliedTool *AssetRepresentation.AppliedTool) (*APIs.AppliedToolResp, error) {
	appToolReq := APIs.AppliedToolResp{
		AppliedTool: appliedTool,
	}

	reqJSON, err := json.Marshal(appToolReq)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", l.url + "/appliedTool", bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Cookie", "SESSION_ID="+user.ledgerSessId)
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusUnauthorized {
			return nil, nil
		}
		return nil, fmt.Errorf("internal ledger error")
	}

	jsonResp := APIs.AppliedToolResp{}
	err = json.Unmarshal(content, &jsonResp)
	if err != nil {
		return nil, err
	}

	return &jsonResp, nil
}
