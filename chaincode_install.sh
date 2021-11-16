#!/bin/bash

source /home/vagrant/shared/scripts/cleanup.sh
source /home/vagrant/shared/scripts/prints.sh


#setup git for go get from private repo
sudo touch /etc/profile.d/go_git_vars.sh
echo 'export GOPRIVATE=github.com/MasterWigu/Thesis' | sudo tee -a /etc/profile.d/go_git_vars.sh
source /etc/profile

cp /home/vagrant/certshare/git_configs/.gitconfig /home/vagrant/.gitconfig
cp /home/vagrant/certshare/git_configs/known_hosts /home/vagrant/.ssh/
cp /home/vagrant/certshare/git_configs/id_rsa /home/vagrant/certshare/git_configs/id_rsa.pub /home/vagrant/.ssh
chmod 600 /home/vagrant/.ssh/id_rsa

eval "$(ssh-agent -s)"
ssh-add /home/vagrant/.ssh/id_rsa



cd /home/vagrant/chaincode/asset-transfer-basic
go mod vendor

cd /home/vagrant/chaincode/user-check
go mod vendor

cd /home/vagrant/chaincode/inventory-management
go mod vendor

mkdir /home/vagrant/chaincode/packaged
cd /home/vagrant/chaincode/packaged


cleanPeer
sudo touch /etc/profile.d/peer_cli_vars.sh
#echo 'export FABRIC_LOGGING_SPEC=DEBUG' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
echo 'export CORE_PEER_ADDRESS=peer1-org1:7051' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
echo 'export CORE_PEER_LOCALMSPID="org1MSP"' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
echo 'export CORE_PEER_TLS_ENABLED=true' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
echo 'export CORE_PEER_TLS_ROOTCERT_FILE=/home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
echo 'export CORE_PEER_MSPCONFIGPATH=/home/vagrant/shared/nodesserver/org1/admin/msp' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
echo 'export FABRIC_CFG_PATH=/home/vagrant/fabric-tools/config/' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
source /etc/profile


/home/vagrant/fabric-tools/bin/peer lifecycle chaincode package basic.tar.gz --path /home/vagrant/chaincode/asset-transfer-basic/ --lang golang --label basic_1.0

/home/vagrant/fabric-tools/bin/peer lifecycle chaincode package userCheck.tar.gz --path /home/vagrant/chaincode/user-check/ --lang golang --label userCheck_1.0
/home/vagrant/fabric-tools/bin/peer lifecycle chaincode package inventoryMgt.tar.gz --path /home/vagrant/chaincode/inventory-management/ --lang golang --label inventoryMgt_1.0


/home/vagrant/fabric-tools/bin/peer lifecycle chaincode install basic.tar.gz
/home/vagrant/fabric-tools/bin/peer lifecycle chaincode install userCheck.tar.gz
/home/vagrant/fabric-tools/bin/peer lifecycle chaincode install inventoryMgt.tar.gz


cleanPeer
sudo touch /etc/profile.d/peer_cli_vars.sh
#echo 'export FABRIC_LOGGING_SPEC=DEBUG' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
echo 'export CORE_PEER_ADDRESS=peer2-org1:7051' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
echo 'export CORE_PEER_LOCALMSPID="org1MSP"' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
echo 'export CORE_PEER_TLS_ENABLED=true' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
echo 'export CORE_PEER_TLS_ROOTCERT_FILE=/home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
echo 'export CORE_PEER_MSPCONFIGPATH=/home/vagrant/shared/nodesserver/org1/admin/msp' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
echo 'export FABRIC_CFG_PATH=/home/vagrant/fabric-tools/config/' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
source /etc/profile

/home/vagrant/fabric-tools/bin/peer lifecycle chaincode install basic.tar.gz
/home/vagrant/fabric-tools/bin/peer lifecycle chaincode install userCheck.tar.gz
/home/vagrant/fabric-tools/bin/peer lifecycle chaincode install inventoryMgt.tar.gz


CC_PACKAGE_ID=$(/home/vagrant/fabric-tools/bin/peer lifecycle chaincode queryinstalled)
CC_PACKAGE_ID=${CC_PACKAGE_ID#*ID: basic}
CC_PACKAGE_ID=${CC_PACKAGE_ID%, Label: basic*} 
CC_PACKAGE_ID=basic${CC_PACKAGE_ID}
/home/vagrant/fabric-tools/bin/peer lifecycle chaincode approveformyorg -o orderer1-org0:7050 --channelID channel1 --name basic --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile /home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem


CC_PACKAGE_ID=$(/home/vagrant/fabric-tools/bin/peer lifecycle chaincode queryinstalled)
CC_PACKAGE_ID=${CC_PACKAGE_ID#*ID: userCheck}
CC_PACKAGE_ID=${CC_PACKAGE_ID%, Label: userCheck*} 
CC_PACKAGE_ID=userCheck${CC_PACKAGE_ID}
/home/vagrant/fabric-tools/bin/peer lifecycle chaincode approveformyorg -o orderer1-org0:7050 --channelID channel1 --name userCheck --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile /home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem

CC_PACKAGE_ID=$(/home/vagrant/fabric-tools/bin/peer lifecycle chaincode queryinstalled)
CC_PACKAGE_ID=${CC_PACKAGE_ID#*ID: inventoryMgt}
CC_PACKAGE_ID=${CC_PACKAGE_ID%, Label: inventoryMgt*} 
CC_PACKAGE_ID=inventoryMgt${CC_PACKAGE_ID}
/home/vagrant/fabric-tools/bin/peer lifecycle chaincode approveformyorg -o orderer1-org0:7050 --channelID channel1 --name inventoryMgt --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile /home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem




/home/vagrant/fabric-tools/bin/peer lifecycle chaincode checkcommitreadiness --channelID channel1 --name basic --version 1.0 --sequence 1 --tls --cafile /home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem --output json
/home/vagrant/fabric-tools/bin/peer lifecycle chaincode checkcommitreadiness --channelID channel1 --name userCheck --version 1.0 --sequence 1 --tls --cafile /home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem --output json
/home/vagrant/fabric-tools/bin/peer lifecycle chaincode checkcommitreadiness --channelID channel1 --name inventoryMgt --version 1.0 --sequence 1 --tls --cafile /home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem --output json


/home/vagrant/fabric-tools/bin/peer lifecycle chaincode commit -o orderer1-org0:7050 --channelID channel1 --name basic --version 1.0 --sequence 1 --tls --cafile /home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem --peerAddresses peer1-org1:7051 --tlsRootCertFiles /home/vagrant/shared/nodesserver/org1/peers/peer1/tls/tlscacerts/tls-tls-ca-7054.pem
/home/vagrant/fabric-tools/bin/peer lifecycle chaincode commit -o orderer1-org0:7050 --channelID channel1 --name userCheck --version 1.0 --sequence 1 --tls --cafile /home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem --peerAddresses peer1-org1:7051 --tlsRootCertFiles /home/vagrant/shared/nodesserver/org1/peers/peer1/tls/tlscacerts/tls-tls-ca-7054.pem
/home/vagrant/fabric-tools/bin/peer lifecycle chaincode commit -o orderer1-org0:7050 --channelID channel1 --name inventoryMgt --version 1.0 --sequence 1 --tls --cafile /home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem --peerAddresses peer1-org1:7051 --tlsRootCertFiles /home/vagrant/shared/nodesserver/org1/peers/peer1/tls/tlscacerts/tls-tls-ca-7054.pem


/home/vagrant/fabric-tools/bin/peer chaincode invoke -o orderer1-org0:7050 --tls --cafile /home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem -C channel1 -n basic --peerAddresses peer1-org1:7051 --tlsRootCertFiles /home/vagrant/shared/nodesserver/org1/peers/peer1/tls/tlscacerts/tls-tls-ca-7054.pem -c '{"function":"InitLedger","Args":[]}'

/home/vagrant/fabric-tools/bin/peer chaincode query -o orderer1-org0:7050 --tls --cafile /home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem -C channel1 -n userCheck --peerAddresses peer1-org1:7051 --tlsRootCertFiles /home/vagrant/shared/nodesserver/org1/peers/peer1/tls/tlscacerts/tls-tls-ca-7054.pem -c '{"function":"GetPerms","Args":[]}'
/home/vagrant/fabric-tools/bin/peer chaincode invoke -o orderer1-org0:7050 --tls --cafile /home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem -C channel1 -n inventoryMgt --peerAddresses peer1-org1:7051 --tlsRootCertFiles /home/vagrant/shared/nodesserver/org1/peers/peer1/tls/tlscacerts/tls-tls-ca-7054.pem -c '{"function":"InitLedger","Args":[]}'


sleep 5
/home/vagrant/fabric-tools/bin/peer chaincode query -C channel1 -n basic -c '{"Args":["GetAllAssets"]}'
/home/vagrant/fabric-tools/bin/peer chaincode query -C channel1 -n inventoryMgt -c '{"Args":["GetAssetTypes"]}'
