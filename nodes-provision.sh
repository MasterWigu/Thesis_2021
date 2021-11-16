#!/bin/bash

source /home/vagrant/shared/scripts/cleanup.sh
source /home/vagrant/shared/scripts/prints.sh


infoln "Cleanup of hosts file"
sudo sed -i '/peer1-org1/d' /etc/hosts
sudo sed -i '/peer2-org1/d' /etc/hosts
sudo sed -i '/orderer1-org0/d' /etc/hosts
sudo sed -i '/orderer2-org0/d' /etc/hosts
sudo sed -i '/orderer3-org0/d' /etc/hosts
successln "Cleanup of hosts file"

infoln "Add the nodes containers to hosts file"
echo '10.10.200.15    peer1-org1' | sudo tee -a /etc/hosts
echo '10.10.200.16    peer2-org1' | sudo tee -a /etc/hosts
echo '10.10.200.20    orderer1-org0' | sudo tee -a /etc/hosts
echo '10.10.200.21    orderer2-org0' | sudo tee -a /etc/hosts
echo '10.10.200.22    orderer3-org0' | sudo tee -a /etc/hosts
successln "Add the nodes containers to hosts file"


sudo rm -rf /home/vagrant/shared/nodesserver/org0
sudo rm -rf /home/vagrant/shared/nodesserver/org1
successln "Delete orgs folders"

infoln "Create folders for peers and orderers"
mkdir -p /home/vagrant/shared/nodesserver/org0/msp
mkdir -p /home/vagrant/shared/nodesserver/org1/msp
mkdir -p /home/vagrant/shared/nodesserver/org0/admin
mkdir -p /home/vagrant/shared/nodesserver/org1/admin
mkdir -p /home/vagrant/shared/nodesserver/org1/user1
mkdir -p /home/vagrant/shared/nodesserver/org1/peers
mkdir -p /home/vagrant/shared/nodesserver/org1/peers/peer1/msp /home/vagrant/shared/nodesserver/org1/peers/peer1/tls
mkdir -p /home/vagrant/shared/nodesserver/org1/peers/peer2/msp /home/vagrant/shared/nodesserver/org1/peers/peer2/tls
mkdir -p /home/vagrant/shared/nodesserver/org0/orderers/orderer1/msp /home/vagrant/shared/nodesserver/org0/orderers/orderer1/tls
mkdir -p /home/vagrant/shared/nodesserver/org0/orderers/orderer2/msp /home/vagrant/shared/nodesserver/org0/orderers/orderer2/tls
mkdir -p /home/vagrant/shared/nodesserver/org0/orderers/orderer3/msp /home/vagrant/shared/nodesserver/org0/orderers/orderer3/tls
successln "Create folders for peers and orderers"

infoln "Create folder where the genesis block will be created"
rm -rf /home/vagrant/channelStart
mkdir -p /home/vagrant/channelStart
successln "Create folder where the genesis block will be created"

infoln "Make root certs avaliable for peers"
mkdir -p /home/vagrant/shared/nodesserver/org1/msp/tlscacerts /home/vagrant/shared/nodesserver/org1/msp/cacerts 
cp /home/vagrant/shared/caserver/tls-ca-server/crypto/tls-ca-cert.pem /home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem
cp /home/vagrant/shared/caserver/org1-ca-server/crypto/ca-cert.pem /home/vagrant/shared/nodesserver/org1/msp/cacerts/ca-cert.pem

mkdir -p /home/vagrant/shared/nodesserver/org0/msp/tlscacerts /home/vagrant/shared/nodesserver/org0/msp/cacerts 
cp /home/vagrant/shared/caserver/tls-ca-server/crypto/tls-ca-cert.pem /home/vagrant/shared/nodesserver/org0/msp/tlscacerts/tls-cert.pem
cp /home/vagrant/shared/caserver/org0-ca-server/crypto/ca-cert.pem /home/vagrant/shared/nodesserver/org0/msp/cacerts/ca-cert.pem
successln "Make root certs avaliable for peers"


cd /home/vagrant/fabric-ca-client

infolnbig "Enroll org0 admin"
cleanCA
sudo touch /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_HOME=/home/vagrant/shared/nodesserver/org0/admin' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_TLS_CERTFILES=/home/vagrant/shared/nodesserver/org0/msp/cacerts/ca-cert.pem' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_MSPDIR=msp' | sudo tee -a /etc/profile.d/ca_client_vars.sh
source /etc/profile

./fabric-ca-client enroll -u https://admin-org0:org0AdminPW@org0-ca:7054
mv /home/vagrant/shared/nodesserver/org0/admin/msp/keystore/* /home/vagrant/shared/nodesserver/org0/admin/msp/keystore/key.pem
cp /home/vagrant/shared/nodesserver/localMSPConf.yaml /home/vagrant/shared/nodesserver/org0/admin/msp/config.yaml
successln "Enroll org0 admin over the org0 CA"

cleanCA
sudo touch /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_HOME=/home/vagrant/shared/nodesserver/org0/admin' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_TLS_CERTFILES=/home/vagrant/shared/nodesserver/org0/msp/tlscacerts/tls-cert.pem' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_MSPDIR=tls' | sudo tee -a /etc/profile.d/ca_client_vars.sh
source /etc/profile

./fabric-ca-client enroll -u https://admin-org0:org0AdminPW@tls-ca:7054
mv /home/vagrant/shared/nodesserver/org0/admin/tls/keystore/* /home/vagrant/shared/nodesserver/org0/admin/tls/keystore/key.pem
successln "Enroll org0 admin over the tls CA"
successlnbig "Enroll org0 admin"

infolnbig "Enroll org1 admin"
cleanCA
sudo touch /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_HOME=/home/vagrant/shared/nodesserver/org1/admin' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_TLS_CERTFILES=/home/vagrant/shared/nodesserver/org1/msp/cacerts/ca-cert.pem' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_MSPDIR=msp' | sudo tee -a /etc/profile.d/ca_client_vars.sh
source /etc/profile

./fabric-ca-client enroll -u https://admin-org1:org1AdminPW@org1-ca:7054
mv /home/vagrant/shared/nodesserver/org1/admin/msp/keystore/* /home/vagrant/shared/nodesserver/org1/admin/msp/keystore/key.pem
cp /home/vagrant/shared/nodesserver/localMSPConf.yaml /home/vagrant/shared/nodesserver/org1/admin/msp/config.yaml
successln "Enroll org1 admin over the org1 CA"

cleanCA
sudo touch /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_HOME=/home/vagrant/shared/nodesserver/org1/admin' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_TLS_CERTFILES=/home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_MSPDIR=tls' | sudo tee -a /etc/profile.d/ca_client_vars.sh
source /etc/profile

./fabric-ca-client enroll -u https://admin-org1:org1AdminPW@tls-ca:7054
mv /home/vagrant/shared/nodesserver/org1/admin/tls/keystore/* /home/vagrant/shared/nodesserver/org1/admin/tls/keystore/key.pem
successln "Enroll org1 admin over the tls CA"
successlnbig "Enroll org1 admin"


infolnbig "Enroll peer1"
cleanCA
sudo touch /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_HOME=/home/vagrant/shared/nodesserver/org1/peers/peer1' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_TLS_CERTFILES=/home/vagrant/shared/nodesserver/org1/msp/cacerts/ca-cert.pem' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_MSPDIR=msp' | sudo tee -a /etc/profile.d/ca_client_vars.sh
source /etc/profile

./fabric-ca-client enroll -u https://peer1-org1:peer1PW@org1-ca:7054
successln "Enroll peer1 over the org1 CA"

cleanCA
sudo touch /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_HOME=/home/vagrant/shared/nodesserver/org1/peers/peer1' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_TLS_CERTFILES=/home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_MSPDIR=tls' | sudo tee -a /etc/profile.d/ca_client_vars.sh
source /etc/profile

./fabric-ca-client enroll -u https://peer1-org1:peer1PW@tls-ca:7054 --enrollment.profile tls --csr.hosts 'peer1-org1'
mv /home/vagrant/shared/nodesserver/org1/peers/peer1/tls/keystore/* /home/vagrant/shared/nodesserver/org1/peers/peer1/tls/keystore/key.pem

cp /home/vagrant/shared/nodesserver/localMSPConf.yaml /home/vagrant/shared/nodesserver/org1/peers/peer1/msp/config.yaml
cp /home/vagrant/shared/nodesserver/localMSPConf.yaml /home/vagrant/shared/nodesserver/org1/msp/config.yaml
successln "Enroll peer1 over the tls CA"
successlnbig "Enroll peer1"


infolnbig "Enroll peer2"
cleanCA
sudo touch /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_HOME=/home/vagrant/shared/nodesserver/org1/peers/peer2' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_TLS_CERTFILES=/home/vagrant/shared/nodesserver/org1/msp/cacerts/ca-cert.pem' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_MSPDIR=msp' | sudo tee -a /etc/profile.d/ca_client_vars.sh
source /etc/profile

./fabric-ca-client enroll -u https://peer2-org1:peer2PW@org1-ca:7054
successln "Enroll peer2 over the org1 CA"

cleanCA
sudo touch /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_HOME=/home/vagrant/shared/nodesserver/org1/peers/peer2' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_TLS_CERTFILES=/home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_MSPDIR=tls' | sudo tee -a /etc/profile.d/ca_client_vars.sh
source /etc/profile

./fabric-ca-client enroll -u https://peer2-org1:peer2PW@tls-ca:7054 --enrollment.profile tls --csr.hosts 'peer2-org1'
mv /home/vagrant/shared/nodesserver/org1/peers/peer2/tls/keystore/* /home/vagrant/shared/nodesserver/org1/peers/peer2/tls/keystore/key.pem

cp /home/vagrant/shared/nodesserver/localMSPConf.yaml /home/vagrant/shared/nodesserver/org1/peers/peer2/msp/config.yaml
cp /home/vagrant/shared/nodesserver/localMSPConf.yaml /home/vagrant/shared/nodesserver/org1/msp/config.yaml
successln "Enroll peer1 over the tls CA"
successlnbig "Enroll peer1"


infolnbig "Enroll orderer1"
cleanCA
sudo touch /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_HOME=/home/vagrant/shared/nodesserver/org0/orderers/orderer1' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_TLS_CERTFILES=/home/vagrant/shared/nodesserver/org0/msp/cacerts/ca-cert.pem' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_MSPDIR=msp' | sudo tee -a /etc/profile.d/ca_client_vars.sh
source /etc/profile

./fabric-ca-client enroll -u https://orderer1-org0:orderer1PW@org0-ca:7054
successln "Enroll orderer1 over the org0 CA"

cleanCA
sudo touch /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_HOME=/home/vagrant/shared/nodesserver/org0/orderers/orderer1' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_TLS_CERTFILES=/home/vagrant/shared/nodesserver/org0/msp/tlscacerts/tls-cert.pem' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_MSPDIR=tls' | sudo tee -a /etc/profile.d/ca_client_vars.sh
source /etc/profile

./fabric-ca-client enroll -u https://orderer1-org0:orderer1PW@tls-ca:7054 --enrollment.profile tls --csr.hosts 'orderer1-org0'

mv /home/vagrant/shared/nodesserver/org0/orderers/orderer1/tls/keystore/* /home/vagrant/shared/nodesserver/org0/orderers/orderer1/tls/keystore/key.pem
cp /home/vagrant/shared/nodesserver/localMSPConf.yaml /home/vagrant/shared/nodesserver/org0/orderers/orderer1/msp/config.yaml
cp /home/vagrant/shared/nodesserver/localMSPConf.yaml /home/vagrant/shared/nodesserver/org0/msp/config.yaml

successln "Enroll orderer1 over the tls CA"
successlnbig "Enroll orderer1"


infolnbig "Enroll orderer2"
cleanCA
sudo touch /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_HOME=/home/vagrant/shared/nodesserver/org0/orderers/orderer2' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_TLS_CERTFILES=/home/vagrant/shared/nodesserver/org0/msp/cacerts/ca-cert.pem' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_MSPDIR=msp' | sudo tee -a /etc/profile.d/ca_client_vars.sh
source /etc/profile

./fabric-ca-client enroll -u https://orderer2-org0:orderer2PW@org0-ca:7054
successln "Enroll orderer2 over the org0 CA"

cleanCA
sudo touch /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_HOME=/home/vagrant/shared/nodesserver/org0/orderers/orderer2' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_TLS_CERTFILES=/home/vagrant/shared/nodesserver/org0/msp/tlscacerts/tls-cert.pem' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_MSPDIR=tls' | sudo tee -a /etc/profile.d/ca_client_vars.sh
source /etc/profile

./fabric-ca-client enroll -u https://orderer2-org0:orderer2PW@tls-ca:7054 --enrollment.profile tls --csr.hosts 'orderer2-org0'

mv /home/vagrant/shared/nodesserver/org0/orderers/orderer2/tls/keystore/* /home/vagrant/shared/nodesserver/org0/orderers/orderer2/tls/keystore/key.pem
cp /home/vagrant/shared/nodesserver/localMSPConf.yaml /home/vagrant/shared/nodesserver/org0/orderers/orderer2/msp/config.yaml
cp /home/vagrant/shared/nodesserver/localMSPConf.yaml /home/vagrant/shared/nodesserver/org0/msp/config.yaml
successln "Enroll orderer2 over the tls CA"
successlnbig "Enroll orderer2"


infolnbig "Enroll orderer3"
cleanCA
sudo touch /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_HOME=/home/vagrant/shared/nodesserver/org0/orderers/orderer3' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_TLS_CERTFILES=/home/vagrant/shared/nodesserver/org0/msp/cacerts/ca-cert.pem' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_MSPDIR=msp' | sudo tee -a /etc/profile.d/ca_client_vars.sh
source /etc/profile

./fabric-ca-client enroll -u https://orderer3-org0:orderer3PW@org0-ca:7054
successln "Enroll orderer3 over the org0 CA"


cleanCA
sudo touch /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_HOME=/home/vagrant/shared/nodesserver/org0/orderers/orderer3' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_TLS_CERTFILES=/home/vagrant/shared/nodesserver/org0/msp/tlscacerts/tls-cert.pem' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_MSPDIR=tls' | sudo tee -a /etc/profile.d/ca_client_vars.sh
source /etc/profile

./fabric-ca-client enroll -u https://orderer3-org0:orderer3PW@tls-ca:7054 --enrollment.profile tls --csr.hosts 'orderer3-org0'

mv /home/vagrant/shared/nodesserver/org0/orderers/orderer3/tls/keystore/* /home/vagrant/shared/nodesserver/org0/orderers/orderer3/tls/keystore/key.pem
cp /home/vagrant/shared/nodesserver/localMSPConf.yaml /home/vagrant/shared/nodesserver/org0/orderers/orderer3/msp/config.yaml
cp /home/vagrant/shared/nodesserver/localMSPConf.yaml /home/vagrant/shared/nodesserver/org0/msp/config.yaml
successln "Enroll orderer3 over the tls CA"
successlnbig "Enroll orderer3"


infolnbig "Start peer1"
mkdir -p /home/vagrant/docker/org1/peer1/home
cp /home/vagrant/shared/nodesserver/peer1-core.yaml /home/vagrant/docker/org1/peer1/home/core.yaml

cd /home/vagrant/docker
sudo docker-compose up -d peer1-org1
successlnbig "Start peer1"


infolnbig "Start peer2"
mkdir -p /home/vagrant/docker/org1/peer2/home
cp /home/vagrant/shared/nodesserver/peer2-core.yaml /home/vagrant/docker/org1/peer2/home/core.yaml

cd /home/vagrant/docker
sudo docker-compose up -d peer2-org1
successlnbig "Start peer2"


infolnbig "Start orderer1"
mkdir -p /home/vagrant/docker/org0/orderer1/home
cp /home/vagrant/shared/nodesserver/orderer1.yaml /home/vagrant/docker/org0/orderer1/home/orderer.yaml

cd /home/vagrant/docker
sudo docker-compose up -d orderer1-org0
successlnbig "Start orderer1"

infolnbig "Start orderer2"
mkdir -p /home/vagrant/docker/org0/orderer2/home
cp /home/vagrant/shared/nodesserver/orderer2.yaml /home/vagrant/docker/org0/orderer2/home/orderer.yaml

cd /home/vagrant/docker
sudo docker-compose up -d orderer2-org0
successlnbig "Start orderer2"

infolnbig "Start orderer3"
mkdir -p /home/vagrant/docker/org0/orderer3/home
cp /home/vagrant/shared/nodesserver/orderer3.yaml /home/vagrant/docker/org0/orderer3/home/orderer.yaml

cd /home/vagrant/docker
sudo docker-compose up -d orderer3-org0
successlnbig "Start orderer3"



infolnbig "Enroll org1 user1"
cd /home/vagrant/fabric-ca-client

cleanCA
sudo touch /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_HOME=/home/vagrant/shared/nodesserver/org1/user1' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_TLS_CERTFILES=/home/vagrant/shared/nodesserver/org1/msp/cacerts/ca-cert.pem' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_MSPDIR=msp' | sudo tee -a /etc/profile.d/ca_client_vars.sh
source /etc/profile

./fabric-ca-client enroll -u https://user1-org1:org1User1PW@org1-ca:7054
mv /home/vagrant/shared/nodesserver/org1/user1/msp/keystore/* /home/vagrant/shared/nodesserver/org1/user1/msp/keystore/key.pem
cp /home/vagrant/shared/nodesserver/localMSPConf.yaml /home/vagrant/shared/nodesserver/org1/user1/msp/config.yaml
successln "Enroll org1 userv1 over the org1 CA"

successlnbig "Enroll org1 user1"
