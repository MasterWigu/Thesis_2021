package main

import (
	"archive/zip"
	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"time"
)

var LoginTimeout = 2 * time.Minute
var UsersList []*User

type User struct {
	sessionId    string
	ledgerSessId string
	userPerm     string //user or admin
	lastUse      time.Time
	username     string
	credZip      multipart.File
	cert         *x509.Certificate
	prvKey       *ecdsa.PrivateKey
}

func createUsersBase() {
	UsersList = make([]*User, 0)
}


// checkLoggedIn returns a pointer to the user if it is logged in, or nil if not and a bool that represents "not logged in"
func checkLoggedIn(sessId string) (*User, bool) {
	for _, user := range UsersList {
		if user.sessionId == sessId {
			if time.Now().Sub(user.lastUse) < LoginTimeout {
				user.lastUse = time.Now()
				return user, false
			}
		}

	}
	return nil, true
}

func parsePrivateKey(der []byte) (*ecdsa.PrivateKey, error) {
	if parsedKey, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		key, ok := parsedKey.(*ecdsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("private Key is not the right type")
		}
		return key, nil
	}
	return nil, fmt.Errorf("private Key is not the right type")
}

func unzipCredFile(credFile io.ReaderAt, credFileSize int64) (*x509.Certificate, *ecdsa.PrivateKey, error) {
	var cert *x509.Certificate
	var key *ecdsa.PrivateKey

	var r, err = zip.NewReader(credFile, credFileSize)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read zip file: %v", err)
	}

	// Iterate through the files in the archive
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return nil, nil, err
		}
		buf := bytes.NewBuffer(nil)
		_, err = buf.ReadFrom(rc)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read files in zip: %v", err)

		}
		block, _ := pem.Decode(buf.Bytes())
		if block == nil {
			return nil, nil, fmt.Errorf("failed to parse certificate PEM: %v", err)
		}

		if f.Name == "cert.pem" {
			cert, err = x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to parse certificate: %v", err)
			}
		}
		if f.Name == "key.pem" {
			key, err = parsePrivateKey(block.Bytes)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to parse privatekey: %v", err)
			}
		}
		err = rc.Close()
		if err != nil {
			log.Println("failed to parse privateKey: " + err.Error())
		}

	}
	return cert, key, nil
}

// loginUser creates a new session in the internal state and logs the user on the ledger too for verification
// @returns bool logged in
// @returns error
func loginUser(credFile multipart.File, credFileSize int64, sessId string) (bool, error) {

	cert, key, err := unzipCredFile(credFile, credFileSize)
	if err != nil {
		return false, fmt.Errorf("could not unzip file: %v", err)
	}

	// all users are created with permission = user  and without name, it will be inserted when ledger module responds
	newUser := User{sessId, "", "user", time.Now(), "", credFile, cert, key}

	loggedIn, err := ledger.loginUser(&newUser)
	if err != nil || !loggedIn {
		return loggedIn, err
	}

	UsersList = append(UsersList, &newUser)

	return true, nil
}
