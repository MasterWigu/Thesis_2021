#!/bin/bash

source /home/vagrant/shared/scripts/cleanup.sh
source /home/vagrant/shared/scripts/prints.sh

infolnbig "Create channel and add orderers and peers to it"

infoln "Get fabric-tools and add to PATH"
cd /home/vagrant
mkdir -p fabric-tools
cd /home/vagrant/fabric-tools

#get the tools
wget -nv https://github.com/hyperledger/fabric/releases/download/v2.3.1/hyperledger-fabric-linux-amd64-2.3.1.tar.gz -O fabric-bins.tar.gz
tar -xf fabric-bins.tar.gz
rm fabric-bins.tar.gz

sudo touch /etc/profile.d/fabric-tools-path.sh
echo 'export PATH=$PATH:/home/vagrant/fabric-tools/bin/' | sudo tee -a /etc/profile.d/fabric-tools-path.sh
source /etc/profile
successln "Get fabric-tools and add to PATH"

infoln "Generate genesis block"
mkdir -p /home/vagrant/channelStart

cd /home/vagrant/channelStart

cp /home/vagrant/shared/nodesserver/configtx.yaml /home/vagrant/channelStart/configtx.yaml

# Create the genesis block
../fabric-tools/bin/configtxgen -profile SampleAppChannelEtcdRaft -outputBlock genesis_block.pb -channelID channel1
successln "Generate genesis block"


infoln "Join orderers to the channel"
cd /home/vagrant/fabric-tools/bin
./osnadmin channel join --channelID channel1  --config-block /home/vagrant/channelStart/genesis_block.pb -o orderer1-org0:9443 --ca-file /home/vagrant/shared/nodesserver/org0/msp/tlscacerts/tls-cert.pem --client-cert /home/vagrant/shared/nodesserver/org0/admin/tls/signcerts/cert.pem --client-key /home/vagrant/shared/nodesserver/org0/admin/tls/keystore/key.pem
./osnadmin channel join --channelID channel1  --config-block /home/vagrant/channelStart/genesis_block.pb -o orderer2-org0:9443 --ca-file /home/vagrant/shared/nodesserver/org0/msp/tlscacerts/tls-cert.pem --client-cert /home/vagrant/shared/nodesserver/org0/admin/tls/signcerts/cert.pem --client-key /home/vagrant/shared/nodesserver/org0/admin/tls/keystore/key.pem
./osnadmin channel join --channelID channel1  --config-block /home/vagrant/channelStart/genesis_block.pb -o orderer3-org0:9443 --ca-file /home/vagrant/shared/nodesserver/org0/msp/tlscacerts/tls-cert.pem --client-cert /home/vagrant/shared/nodesserver/org0/admin/tls/signcerts/cert.pem --client-key /home/vagrant/shared/nodesserver/org0/admin/tls/keystore/key.pem
successln "Join orderers to the channel"


infoln "Check if all orderers joined the channel"
./osnadmin channel list --channelID channel1 -o orderer1-org0:9443 --ca-file /home/vagrant/shared/nodesserver/org0/msp/tlscacerts/tls-cert.pem --client-cert /home/vagrant/shared/nodesserver/org0/admin/tls/signcerts/cert.pem --client-key /home/vagrant/shared/nodesserver/org0/admin/tls/keystore/key.pem
./osnadmin channel list --channelID channel1 -o orderer2-org0:9443 --ca-file /home/vagrant/shared/nodesserver/org0/msp/tlscacerts/tls-cert.pem --client-cert /home/vagrant/shared/nodesserver/org0/admin/tls/signcerts/cert.pem --client-key /home/vagrant/shared/nodesserver/org0/admin/tls/keystore/key.pem
./osnadmin channel list --channelID channel1 -o orderer3-org0:9443 --ca-file /home/vagrant/shared/nodesserver/org0/msp/tlscacerts/tls-cert.pem --client-cert /home/vagrant/shared/nodesserver/org0/admin/tls/signcerts/cert.pem --client-key /home/vagrant/shared/nodesserver/org0/admin/tls/keystore/key.pem
successln "Check if all orderers joined the channel"


#wait for the orderers to figure out themselves
sleep 10


infoln "Join peer 1 to the channel"
cleanPeer
sudo touch /etc/profile.d/peer_cli_vars.sh
#echo 'export FABRIC_LOGGING_SPEC=DEBUG' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
echo 'export CORE_PEER_ADDRESS=peer1-org1:7051' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
echo 'export CORE_PEER_LOCALMSPID=org1MSP' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
echo 'export CORE_PEER_TLS_ENABLED=true' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
echo 'export CORE_PEER_TLS_ROOTCERT_FILE=/home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
echo 'export CORE_PEER_MSPCONFIGPATH=/home/vagrant/shared/nodesserver/org1/admin/msp' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
echo 'export FABRIC_CFG_PATH=/home/vagrant/fabric-tools/config/' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
source /etc/profile




./peer channel fetch newest /home/vagrant/channelStart/channel1.block -c channel1 --orderer orderer1-org0:7050 --tls --cafile /home/vagrant/shared/nodesserver/org0/msp/tlscacerts/tls-cert.pem


#join peer1 to the channel
./peer channel join -b /home/vagrant/channelStart/channel1.block

./peer channel list
./peer channel getinfo -c channel1
successln "Join peer 1 to the channel"

infoln "Join peer 2 to the channel"
cleanPeer
sudo touch /etc/profile.d/peer_cli_vars.sh
#echo 'export FABRIC_LOGGING_SPEC=DEBUG' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
echo 'export CORE_PEER_ADDRESS=peer2-org1:7051' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
echo 'export CORE_PEER_LOCALMSPID=org1MSP' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
echo 'export CORE_PEER_TLS_ENABLED=true' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
echo 'export CORE_PEER_TLS_ROOTCERT_FILE=/home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
echo 'export CORE_PEER_MSPCONFIGPATH=/home/vagrant/shared/nodesserver/org1/admin/msp' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
echo 'export FABRIC_CFG_PATH=/home/vagrant/fabric-tools/config/' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
source /etc/profile

#join peer2 to the channel
./peer channel join -b /home/vagrant/channelStart/channel1.block

./peer channel list
./peer channel getinfo -c channel1
successln "Join peer 2 to the channel"


successlnbig "Create channel and add orderers and peers to it"
