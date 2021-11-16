#!/bin/bash

source /home/vagrant/shared/scripts/cleanup.sh
source /home/vagrant/shared/scripts/prints.sh

rm -rf /home/vagrant/chaincode

mkdir /home/vagrant/chaincode
cd /home/vagrant/chaincode

git clone https://github.com/hyperledger/fabric-samples

cd /home/vagrant/chaincode/fabric-samples/asset-transfer-basic/chaincode-go
go mod vendor


mkdir /home/vagrant/chaincode/fabric-samples/tests
cd /home/vagrant/chaincode/fabric-samples/tests


cleanPeer
sudo touch /etc/profile.d/peer_cli_vars.sh
#echo 'export FABRIC_LOGGING_SPEC=DEBUG' | sudo tee -a /etc/profile.d/peer_cli_vars.sh.sh
echo 'export CORE_PEER_ADDRESS=peer1-org1:7051' | sudo tee -a /etc/profile.d/peer_cli_vars.sh.sh
echo 'export CORE_PEER_LOCALMSPID="org1MSP"' | sudo tee -a /etc/profile.d/peer_cli_vars.sh.sh
echo 'export CORE_PEER_TLS_ENABLED=true' | sudo tee -a /etc/profile.d/peer_cli_vars.sh.sh
echo 'export CORE_PEER_TLS_ROOTCERT_FILE=/home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem' | sudo tee -a /etc/profile.d/peer_cli_vars.sh.sh
echo 'export CORE_PEER_MSPCONFIGPATH=/home/vagrant/shared/nodesserver/org1/admin/msp' | sudo tee -a /etc/profile.d/peer_cli_vars.sh.sh
echo 'export FABRIC_CFG_PATH=/home/vagrant/fabric-tools/config/' | sudo tee -a /etc/profile.d/peer_cli_vars.sh.sh
source /etc/profile


/home/vagrant/fabric-tools/bin/peer lifecycle chaincode package basic.tar.gz --path /home/vagrant/chaincode/fabric-samples/asset-transfer-basic/chaincode-go/ --lang golang --label basic_1.0


/home/vagrant/fabric-tools/bin/peer lifecycle chaincode install basic.tar.gz


cleanPeer
sudo touch /etc/profile.d/peer_cli_vars.sh
#echo 'export FABRIC_LOGGING_SPEC=DEBUG' | sudo tee -a /etc/profile.d/peer_cli_vars.sh.sh
echo 'export CORE_PEER_ADDRESS=peer2-org1:7051' | sudo tee -a /etc/profile.d/peer_cli_vars.sh.sh
echo 'export CORE_PEER_LOCALMSPID="org1MSP"' | sudo tee -a /etc/profile.d/peer_cli_vars.sh.sh
echo 'export CORE_PEER_TLS_ENABLED=true' | sudo tee -a /etc/profile.d/peer_cli_vars.sh.sh
echo 'export CORE_PEER_TLS_ROOTCERT_FILE=/home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem' | sudo tee -a /etc/profile.d/peer_cli_vars.sh.sh
echo 'export CORE_PEER_MSPCONFIGPATH=/home/vagrant/shared/nodesserver/org1/admin/msp' | sudo tee -a /etc/profile.d/peer_cli_vars.sh.sh
echo 'export FABRIC_CFG_PATH=/home/vagrant/fabric-tools/config/' | sudo tee -a /etc/profile.d/peer_cli_vars.sh.sh
source /etc/profile

/home/vagrant/fabric-tools/bin/peer lifecycle chaincode install basic.tar.gz


CC_PACKAGE_ID=$(/home/vagrant/fabric-tools/bin/peer lifecycle chaincode queryinstalled)
CC_PACKAGE_ID=${CC_PACKAGE_ID#*ID: }
CC_PACKAGE_ID=${CC_PACKAGE_ID%, L*} 


/home/vagrant/fabric-tools/bin/peer lifecycle chaincode approveformyorg -o orderer1-org0:7050 --channelID channel1 --name basic --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile /home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem


/home/vagrant/fabric-tools/bin/peer lifecycle chaincode checkcommitreadiness --channelID channel1 --name basic --version 1.0 --sequence 1 --tls --cafile /home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem --output json


/home/vagrant/fabric-tools/bin/peer lifecycle chaincode commit -o orderer1-org0:7050 --channelID channel1 --name basic --version 1.0 --sequence 1 --tls --cafile /home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem --peerAddresses peer1-org1:7051 --tlsRootCertFiles /home/vagrant/shared/nodesserver/org1/peers/peer1/tls/tlscacerts/tls-tls-ca-7054.pem

/home/vagrant/fabric-tools/bin/peer chaincode invoke -o orderer1-org0:7050 --tls --cafile /home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem -C channel1 -n basic --peerAddresses peer1-org1:7051 --tlsRootCertFiles /home/vagrant/shared/nodesserver/org1/peers/peer1/tls/tlscacerts/tls-tls-ca-7054.pem -c '{"function":"InitLedger","Args":[]}'


/home/vagrant/fabric-tools/bin/peer chaincode query -C channel1 -n basic -c '{"Args":["GetAllAssets"]}'