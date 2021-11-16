package main

import (
	"archive/zip"
	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"io"
	"log"
	"path/filepath"
	"time"
)

var LoginTimeout = 5 * time.Minute
var Wallet *gateway.Wallet
var UsersList []*User

type User struct {
	sessionId string
	userPerm  string //user or admin
	lastUse   time.Time
	username  string
	network   *gateway.Network
}

func createUsersBase() {
	Wallet = gateway.NewInMemoryWallet()
	UsersList = make([]*User, 0)
}

// checkLoggedIn returns a pointer to the user if it is logged in, or nil if not and a bool that represents "not logged in"
func checkLoggedIn(sessId string) (*User, bool) {
	for i, user := range UsersList {
		if user.sessionId == sessId {
			if time.Now().Sub(user.lastUse) < LoginTimeout {
				user.lastUse = time.Now()
				return user, false
			} else {
				// remove old entry from users list
				UsersList = append(UsersList[:i], UsersList[i+1:]...)
				err := Wallet.Remove(sessId)
				if err != nil {
					fmt.Printf("error removing old entry from wallet: %v", err)
				}
			}
		}

	}
	return nil, true
}

// unzipCredFile receives the zip file containing the cert.pem and key.pem of the user (and its size) and extracts
// both to a gateway.X509Identity and the user id
func unzipCredFile(credFile io.ReaderAt, credFileSize int64) (string, *gateway.X509Identity, error) {
	var cert string
	var key string
	var userId string

	var r, err = zip.NewReader(credFile, credFileSize)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read zip file: %v", err)
	}

	// Iterate through the files in the archive
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return "", nil, err
		}
		buf := bytes.NewBuffer(nil)
		_, err = buf.ReadFrom(rc)
		if err != nil {
			return "", nil, fmt.Errorf("failed to read files in zip: %v", err)
		}

		block, _ := pem.Decode(buf.Bytes())
		if block == nil {
			return "", nil, fmt.Errorf("failed to parse certificate PEM: %v", err)
		}

		if f.Name == "cert.pem" {
			cert = string(buf.Bytes())
			cert2, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return "", nil, fmt.Errorf("failed to parse certificate: %v", err)
			}
			userId = cert2.Subject.CommonName

		}
		if f.Name == "key.pem" {
			key = string(buf.Bytes())
		}
		err = rc.Close()
		if err != nil {
			log.Println("failed to parse privateKey: " + err.Error())
		}
	}

	identity := gateway.NewX509Identity("org1MSP", cert, key)

	return userId, identity, nil
}


// verifyUserCerts verifies if the certificate and private key match
// @returns bool is they match
func verifyUserCerts(identity *gateway.X509Identity) (bool, error) {
	certString := identity.Certificate()
	keyString := identity.Key()

	certBlock, _ := pem.Decode([]byte(certString))
	if certBlock == nil {
		return false, fmt.Errorf("failed to parse certificate PEM")
	}

	keyBlock, _ := pem.Decode([]byte(keyString))
	if keyBlock == nil {
		return false, fmt.Errorf("failed to parse key PEM")
	}

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return false, fmt.Errorf("failed to parse certificate: %v", err)
	}

	parsedKey, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		return false, fmt.Errorf("failed to parse private key: %v", err)
	}
	key, ok := parsedKey.(*ecdsa.PrivateKey)
	if !ok {
		return false, fmt.Errorf("failed to parse private key")
	}

	if !key.PublicKey.Equal(cert.PublicKey) {
		return false, nil
	}
	return true, nil
}

// loginUser handles the user login, creating the user identity in the list of users to subsequent authentications
func loginUser(credFile io.ReaderAt, credFileSize int64, sessId string) (*User, error) {
	userId, identity, err := unzipCredFile(credFile, credFileSize)
	if err != nil {
		return nil, fmt.Errorf("failed to unzip file: %v", err)
	}

	match, err := verifyUserCerts(identity)
	if err != nil {
		return nil, fmt.Errorf("could not locally verify cert and key: %v", err)
	}
	if !match {
		return nil, nil
	}

	if !Wallet.Exists(sessId) {
		err := Wallet.Put(sessId, identity)
		if err != nil {
			return nil, fmt.Errorf("failed to populate wallet contents: %v", err)
		}
	}

	ccpPath := "/fabricModule/resources/connection-org1.yaml"

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(Wallet, sessId),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gateway: %v", err)
	}
	network, err := gw.GetNetwork("channel1")
	if err != nil {
		log.Printf("Failed to get network: %v\n", err)
		return nil, err
	}

	// all users are created with permission = user  and without name, it will be inserted when ledger module responds
	newUser := User{sessId, "user", time.Now(), userId, network}

	perm, err := getPermissions(&newUser)
	if err != nil {
		log.Printf("Failed to get perms: %v\n", err)
		return nil, err
	}
	newUser.userPerm = perm

	UsersList = append(UsersList, &newUser)

	return &newUser, nil
}
