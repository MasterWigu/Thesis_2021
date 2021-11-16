package main

import (
	"archive/zip"
	"bufio"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

var fileHandler *FileHandler

type FileHandler struct {
	actionStoragePath string
	runningDirPath string
}


func createFileHandler() {
	fileHandler = new(FileHandler)
	fileHandler.actionStoragePath = "/terraformModule/resources/actionStorage/"
	fileHandler.runningDirPath = "/terraformModule/resources/runningDir/"
}

func (f *FileHandler) DeleteActionFiles(actionId string) {
	err := os.Remove(f.actionStoragePath + actionId + ".zip")
	if err != nil {
		log.Println("failed to delete action files: " + err.Error())
	}
}

func (f *FileHandler) DeleteRunningFiles(actionId string) {
	err := os.RemoveAll(f.runningDirPath + actionId)
	if err != nil {
		log.Println("failed to delete action files: " + err.Error())
	}
}

func (f *FileHandler) SaveActionFiles(file *multipart.FileHeader, actionId string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}

	h := sha256.New()
	if _, err = io.Copy(h, src); err != nil {
		return "", fmt.Errorf("could not get hash of zip file: %v", err)
	}


	out, err := os.Create(f.actionStoragePath + actionId + ".zip")
	if err != nil {
		f.DeleteActionFiles(actionId)
		return "", fmt.Errorf("failed to create new file in fs for zip file: %v", err)
	}

	_, err = io.Copy(out, src)
	if err != nil {
		f.DeleteActionFiles(actionId)
		return "", fmt.Errorf("failed to copy contents of zip to file in fs: %v", err)
	}

	err = out.Close()
	if err != nil {
		log.Println("failed to close fs file: " + err.Error())
	}
	err = src.Close()
	if err != nil {
		log.Println("failed to close memory file: " + err.Error())
	}

	return base64.URLEncoding.EncodeToString(h.Sum(nil)), nil
}

// PrepareRunningDir returns list of asset ids and error
func (f *FileHandler) PrepareRunningDir(actionId string) error {
	dstPath := f.runningDirPath + actionId + "/"
	err := os.Mkdir(f.runningDirPath + actionId, 0755)
	if err != nil {
		return fmt.Errorf("failed to create running dir: %v", err)
	}


	r, err := zip.OpenReader(f.actionStoragePath + actionId + ".zip")
	if err != nil {
		return fmt.Errorf("failed to read zip file: %v", err)
	}

	// Iterate through the files in the archive
	for _, fp := range r.File {
		rc, err := fp.Open()
		if err != nil {
			f.DeleteRunningFiles(actionId)
			return fmt.Errorf("failed to open file inside zip file: %v", err)
		}

		out, err := os.Create(dstPath + fp.Name)
		if err != nil {
			f.DeleteRunningFiles(actionId)
			return fmt.Errorf("failed to create new file in fs: %v", err)
		}

		_, err = io.Copy(out, rc)
		if err != nil {
			f.DeleteRunningFiles(actionId)
			return fmt.Errorf("failed to copy contents to file in fs: %v", err)
		}
		err = rc.Close()
		if err != nil {
			log.Println("failed to close memory file: " + err.Error())
		}
		err = out.Close()
		if err != nil {
			log.Println("failed to close fs file: " + err.Error())
		}

	}

	return nil
}

func (f *FileHandler) ProcessHosts(actionId string) ([]string, error) {
	/* reads the hosts file in the running dir, copies the pairs (ip, host) to the system's hosts file
	and returns the asset id's of the hosts
	Assumes the file has three params per line (ip, host, assetId)
	 */
	assetIds := make([]string, 0)


	hostsFile, err := os.OpenFile("/etc/hosts", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open systems hosts file: %v", err)
	}
	defer hostsFile.Close()

	file, err := os.Open(f.runningDirPath + actionId + "/hosts")
	if err != nil {
		return nil, fmt.Errorf("failed to open running config hosts file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 5 {
			continue
		}
		components := strings.Split(scanner.Text(), " ")
		assetIds = append(assetIds, components[2])

		_, err = hostsFile.WriteString(components[0] + "\t" + components[1])
		if err != nil {
			return nil, fmt.Errorf("failed to write new entries to hosts file: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}


	// get the inventory and extract asset id list
	return assetIds, nil
}

func (f *FileHandler) ActionZipExists(actionId string) bool {

	if pathAbs, err := filepath.Abs(f.actionStoragePath +  actionId + ".zip"); err != nil {
		return false
	} else if _, err := os.Stat(pathAbs); os.IsNotExist(err) {
		return false
	}
	return true

}


func (f *FileHandler) RunningDirExists(actionId string) bool {
	path := f.runningDirPath + actionId
	if pathAbs, err := filepath.Abs(path); err != nil {
		return false
	} else if fileInfo, err := os.Stat(pathAbs); os.IsNotExist(err) || !fileInfo.IsDir() {
		return false
	}
	return true
}

func (f *FileHandler) changeActionId(actionId string, ledgerId string) error {
	if f.ActionZipExists(actionId) {
		src, err := filepath.Abs(f.actionStoragePath +  actionId + ".zip")
		if err != nil {
			return err
		}
		dst, err := filepath.Abs(f.actionStoragePath +  ledgerId + ".zip")
		if err != nil {
			return err
		}

		err = os.Rename(src, dst)
		if err != nil {
			return err
		}
	}

	if f.RunningDirExists(actionId) {
		src, err := filepath.Abs(f.runningDirPath +  actionId + ".zip")
		if err != nil {
			return err
		}
		dst, err := filepath.Abs(f.runningDirPath +  ledgerId + ".zip")
		if err != nil {
			return err
		}

		err = os.Rename(src, dst)
		if err != nil {
			return err
		}
	}

	return nil
}

