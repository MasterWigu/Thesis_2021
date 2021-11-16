package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/MasterWigu/Thesis/appServer/APIs"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

var modules *Modules

type Modules struct {
	modules map[string]string
}

func initModules() {
	modules = new(Modules)


	modules.modules = make(map[string]string)

	file, err := os.Open("/brokerModule/resources/modules.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 5 {
			continue
		}
		components := strings.Split(scanner.Text(), " ")
		modules.modules[components[0]] = components[1] + ":" + components[2]
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}


	return
}


func (m *Modules) planAction(tool string, files *multipart.FileHeader) (*APIs.PlanResp, error) {
	url := m.modules[tool]


	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("toolResources", "toolResources.zip")

	if err != nil {
		return nil, err
	}

	filesO, err := files.Open()
	if err != nil {
		return nil, err
	}

	io.Copy(part, filesO)
	writer.Close()
	request, err := http.NewRequest("POST", "https://" + url + "/plan", body)

	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())

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
		return nil, fmt.Errorf("internal tool error")
	}

	jsonResp := APIs.PlanResp{}
	err = json.Unmarshal(content, &jsonResp)
	if err != nil {
		return nil, err
	}


	return &jsonResp, nil
}

func (m *Modules) executeAction(tool string, files *multipart.FileHeader) (*APIs.ExecuteResp, error) {
	url := m.modules[tool]


	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("toolResources", "toolResources.zip")

	if err != nil {
		return nil, err
	}

	filesO, err := files.Open()
	if err != nil {
		return nil, err
	}

	io.Copy(part, filesO)
	writer.Close()
	request, err := http.NewRequest("POST", "https://" + url + "/execute", body)

	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())

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
		return nil, fmt.Errorf("internal tool error")
	}

	jsonResp := APIs.ExecuteResp{}
	err = json.Unmarshal(content, &jsonResp)
	if err != nil {
		return nil, err
	}


	return &jsonResp, nil
}

func (m *Modules) confirmAction(tool string, id string) (*APIs.ExecuteResp, error) {
	url := m.modules[tool]

	request, err := http.NewRequest("POST", "https://" + url + "/execute/" + id + "/confirm", nil)
	if err != nil {
		return nil, err
	}

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
		return nil, fmt.Errorf("internal tool error")
	}

	jsonResp := APIs.ExecuteResp{}
	err = json.Unmarshal(content, &jsonResp)
	if err != nil {
		return nil, err
	}


	return &jsonResp, nil
}

func (m *Modules) informActionIdOnLedger(tool string, actionId string, ledgerId string) error {
	url := m.modules[tool]

	request, err := http.NewRequest("POST", "https://" + url + "/actions/" + actionId + "/ledgerId/" + ledgerId, nil)
	if err != nil {
		return err
	}

	response, err := client.Do(request)

	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("internal tool error")
	}
	return nil
}