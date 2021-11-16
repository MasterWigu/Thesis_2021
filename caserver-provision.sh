#!/bin/bash

source /home/vagrant/shared/scripts/cleanup.sh
source /home/vagrant/shared/scripts/prints.sh

infolnbig "CA Servers provision file"

infoln "cleanup of hosts file"
sudo sed -i '/tls-ca/d' /etc/hosts
sudo sed -i '/org0-ca/d' /etc/hosts
sudo sed -i '/org1-ca/d' /etc/hosts
successln "cleanup of hosts file"

### Adding the hostnames of the containers
echo '10.10.200.5     tls-ca' | sudo tee -a /etc/hosts     #guide port 7052
echo '10.10.200.10    org0-ca' | sudo tee -a /etc/hosts    #guide port 7053-5
echo '10.10.200.11    org1-ca' | sudo tee -a /etc/hosts    #guide port 7053-5
#echo '10.10.200.5     tls-ca-server' | sudo tee -a /etc/hosts
#echo '10.10.200.5     tls-ca-server' | sudo tee -a /etc/hosts
#echo '10.10.200.5     tls-ca-server' | sudo tee -a /etc/hosts
#echo '10.10.200.5     tls-ca-server' | sudo tee -a /etc/hosts


infoln "cleanup of old keys and files"
mkdir -p /home/vagrant/docker #if the folder does not exist, create (just to suppress error on cd)
cd /home/vagrant/docker
sudo docker-compose down
cd /

sudo rm -rf /home/vagrant/docker
sudo rm -rf /home/vagrant/fabric-ca-client
sudo rm -rf /home/vagrant/shared/caserver/tls-ca-server/admin/*
sudo rm -rf /home/vagrant/shared/caserver/tls-ca-server/crypto/*
sudo mkdir -p /home/vagrant/shared/caserver/tls-ca-server/admin
sudo mkdir -p /home/vagrant/shared/caserver/tls-ca-server/crypto

sudo rm -rf /home/vagrant/shared/caserver/org0-ca-server/crypto/*
sudo rm -rf /home/vagrant/shared/caserver/org1-ca-server/crypto/*
mkdir -p /home/vagrant/shared/caserver/org0-ca-server/crypto/temp
mkdir -p /home/vagrant/shared/caserver/org1-ca-server/crypto/temp
sudo rm -rf /home/vagrant/shared/caserver/org0-ca-server/admin/*
sudo rm -rf /home/vagrant/shared/caserver/org1-ca-server/admin/*
sudo rm -rf /home/vagrant/shared/caserver/org0/admin/*
sudo rm -rf /home/vagrant/shared/caserver/org1/admin/*
cleanCA
successln "cleanup of old keys and files"


infolnbig "TLS CA SERVER INITIALIZATION AND START"

infoln "Create docker folder and copy docker compose file"
mkdir /home/vagrant/docker
cp /home/vagrant/shared/caserver/docker-compose.yml /home/vagrant/docker
successln "Create docker folder and copy docker compose file"

infoln "Pre-make the TLS CA container directory and place the config file"
mkdir -p /home/vagrant/docker/tls/ca

cp /home/vagrant/shared/caserver/tls-ca-server/fabric-ca-server-config.yaml /home/vagrant/docker/tls/ca
successln "Pre-make the TLS CA container directory and place the config file"

infoln "Run the TLS CA container"
cd /home/vagrant/docker
sudo docker-compose run -d tls-ca

#wait for server start (just to be cautious)
sleep 5
successln "Run the TLS CA container"

cp /home/vagrant/docker/tls/ca/ca-cert.pem /home/vagrant/shared/caserver/tls-ca-server/crypto/tls-ca-cert.pem
successln "Copy the ca certificate to the shared folder and rename to more appropriate tls-ca-cert.pem(for further redistribution)"

successlnbig "TLS CA SERVER INITIALIZATION AND START"

infolnbig "CA CLIENT INIT"

infoln "Create folder structure"
cd /home/vagrant
mkdir fabric-ca-client
successln "Create folder structure"

infoln "Get client binary"
cd /home/vagrant/fabric-ca-client
wget -nv https://github.com/hyperledger/fabric-ca/releases/download/v1.4.9/hyperledger-fabric-ca-linux-amd64-1.4.9.tar.gz -O ca-client.tar.gz
tar -xf ca-client.tar.gz
mv bin/fabric-ca-client ./fabric-ca-client
rm -rf ca-client.tar.gz bin

#ensure the binary is executable
sudo chmod +x fabric-ca-client
successln "Get client binary"

#set FABRIC_CA_CLIENT_HOME environment variable
cleanCA
sudo touch /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_HOME=/home/vagrant/shared/caserver/tls-ca-server/admin' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_TLS_CERTFILES=/home/vagrant/shared/caserver/tls-ca-server/crypto/tls-ca-cert.pem' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_MSPDIR=msp' | sudo tee -a /etc/profile.d/ca_client_vars.sh
source /etc/profile
successln "Set fabric-ca-client environment variables"


### USER MANAGEMENT
infolnbig "Enroll tls ca admin and register entities on the TLS CA"
infoln "Enroll the admin user on the tls CA"
./fabric-ca-client enroll -u https://tlscaadmin:tlscaadminpw@tls-ca:7054 --enrollment.profile tls --csr.hosts '*'
successln "Enroll tls ca admin on the TLS CA"

./fabric-ca-client register --id.name peer1-org1 --id.secret peer1PW --id.type peer -u https://tls-ca:7054
./fabric-ca-client register --id.name peer2-org1 --id.secret peer2PW --id.type peer -u https://tls-ca:7054
./fabric-ca-client register --id.name orderer1-org0 --id.secret orderer1PW --id.type orderer -u https://tls-ca:7054
./fabric-ca-client register --id.name orderer2-org0 --id.secret orderer2PW --id.type orderer -u https://tls-ca:7054
./fabric-ca-client register --id.name orderer3-org0 --id.secret orderer3PW --id.type orderer -u https://tls-ca:7054
./fabric-ca-client register --id.name rcaorg0admin --id.secret rcaorg0adminpw -u https://tls-ca:7054
./fabric-ca-client register --id.name rcaorg1admin --id.secret rcaorg1adminpw -u https://tls-ca:7054
./fabric-ca-client register --id.name admin-org0 --id.secret org0AdminPW --id.type admin --id.affiliation org0.sysadmins -u https://tls-ca:7054
./fabric-ca-client register --id.name admin-org1 --id.secret org1AdminPW --id.type admin --id.affiliation org1.sysadmins -u https://tls-ca:7054
successln "Register entities on the TLS CA"

successlnbig "Enroll tls ca admin and register entities on the TLS CA"

infolnbig "Enroll org0 admin on the TLS CA"
cleanCA
sudo touch /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_HOME=/home/vagrant/shared/caserver/org0-ca-server/crypto/temp' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_TLS_CERTFILES=/home/vagrant/shared/caserver/tls-ca-server/crypto/tls-ca-cert.pem' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_MSPDIR=msp' | sudo tee -a /etc/profile.d/ca_client_vars.sh
source /etc/profile

./fabric-ca-client enroll -u https://rcaorg0admin:rcaorg0adminpw@tls-ca:7054 --enrollment.profile tls --csr.hosts '*'
successln "Enroll org0 admin on the TLS CA"

#Rename the rca private key file to key.pem (this assumes the only file in the folder is the key)
mv /home/vagrant/shared/caserver/org0-ca-server/crypto/temp/msp/keystore/* /home/vagrant/shared/caserver/org0-ca-server/crypto/temp/msp/keystore/key.pem
successln "Rename the rca private key file to key.pem"

#copy rcaorg0admin cert and key to org0 ca crypto folder
cp /home/vagrant/shared/caserver/org0-ca-server/crypto/temp/msp/keystore/* /home/vagrant/shared/caserver/org0-ca-server/crypto/key.pem
cp /home/vagrant/shared/caserver/org0-ca-server/crypto/temp/msp/signcerts/cert.pem /home/vagrant/shared/caserver/org0-ca-server/crypto/cert.pem
rm -rf /home/vagrant/shared/caserver/org0-ca-server/crypto/temp
successln "Copy rcaorg0admin cert and key to org0 ca crypto folder"

successlnbig "Enroll org0 admin on the TLS CA"


infolnbig "Enroll org1 admin on the TLS CA"
cleanCA
sudo touch /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_HOME=/home/vagrant/shared/caserver/org1-ca-server/crypto/temp' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_TLS_CERTFILES=/home/vagrant/shared/caserver/tls-ca-server/crypto/tls-ca-cert.pem' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_MSPDIR=msp' | sudo tee -a /etc/profile.d/ca_client_vars.sh
source /etc/profile

./fabric-ca-client enroll -u https://rcaorg1admin:rcaorg1adminpw@tls-ca:7054 --enrollment.profile tls --csr.hosts '*'
successln "Enroll org1 admin on the TLS CA"

#rename the rca private key file to key.pem (this assumes the only file in the folder is the key)
mv /home/vagrant/shared/caserver/org1-ca-server/crypto/temp/msp/keystore/* /home/vagrant/shared/caserver/org1-ca-server/crypto/temp/msp/keystore/key.pem
successln "Rename the rca private key file to key.pem"

#copy rcaorg1admin cert and key to org1 ca crypto folder
cp /home/vagrant/shared/caserver/org1-ca-server/crypto/temp/msp/keystore/* /home/vagrant/shared/caserver/org1-ca-server/crypto/key.pem
cp /home/vagrant/shared/caserver/org1-ca-server/crypto/temp/msp/signcerts/cert.pem /home/vagrant/shared/caserver/org1-ca-server/crypto/cert.pem
rm -rf /home/vagrant/shared/caserver/org1-ca-server/crypto/temp
successln "Copy rcaorg1admin cert and key to org1 ca crypto folder"

successlnbig "Enroll org1 admin on the TLS CA"


infolnbig "Create Org0 CA"
#relocate ourselves
cd /home/vagrant/docker

mkdir -p /home/vagrant/docker/org0/ca/tls
successln "Create org ca folders"

#copy rcaorg0admin credentials to a folder accessible by the org ca
cp /home/vagrant/shared/caserver/org0-ca-server/crypto/cert.pem /home/vagrant/docker/org0/ca/tls
cp /home/vagrant/shared/caserver/org0-ca-server/crypto/key.pem /home/vagrant/docker/org0/ca/tls
successln "Copy rcaorg0admin credentials to a folder accessible by the org ca"

#copy the config file to the correct folder
cp /home/vagrant/shared/caserver/org0-ca-server/fabric-ca-server-config.yaml /home/vagrant/docker/org0/ca
successln "Copy the config file to the correct folder"

infoln "Start the org0 CA server container"
sudo docker-compose run -d org0-ca

sleep 5
successln "Start the org0 CA server container"

#copy TLS CA cert to org crypto folder
cp /home/vagrant/docker/tls/ca/ca-cert.pem /home/vagrant/shared/caserver/org0-ca-server/crypto/tls-ca-cert.pem

#copy ORG CA cert to org crypto folder
cp /home/vagrant/docker/org0/ca/ca-cert.pem /home/vagrant/shared/caserver/org0-ca-server/crypto/ca-cert.pem
successln "Copy TLS and Org0 CAs certs to org crypto folder"
successlnbig "Create Org0 CA"

infolnbig "Create Org1 CA"
#relocate ourselves
cd /home/vagrant/docker

#create org ca folders
mkdir -p /home/vagrant/docker/org1/ca/tls
successln "Create org ca folders"

#copy rcaorg1admin credentials to a folder accessible by the org ca
cp /home/vagrant/shared/caserver/org1-ca-server/crypto/cert.pem /home/vagrant/docker/org1/ca/tls
cp /home/vagrant/shared/caserver/org1-ca-server/crypto/key.pem /home/vagrant/docker/org1/ca/tls
successln "Copy rcaorg1admin credentials to a folder accessible by the org ca"

#copy the config file to the correct folder
cp /home/vagrant/shared/caserver/org1-ca-server/fabric-ca-server-config.yaml /home/vagrant/docker/org1/ca
successln "Copy the config file to the correct folder"

infoln "Start the org1 CA server container"
sudo docker-compose run -d org1-ca

sleep 5
successln "Start the org0 CA server container"

#copy TLS CA cert to org crypto folder
cp /home/vagrant/docker/tls/ca/ca-cert.pem /home/vagrant/shared/caserver/org1-ca-server/crypto/tls-ca-cert.pem

#copy ORG CA cert to org crypto folder
cp /home/vagrant/docker/org1/ca/ca-cert.pem /home/vagrant/shared/caserver/org1-ca-server/crypto/ca-cert.pem
successln "Copy TLS and Org1 CAs certs to org crypto folder"
successlnbig "Create Org1 CA"



infolnbig "Enroll org0 ca admin and register entities on the Org0 CA"
infoln "Enroll the admin user on the org0 CA"
cd /home/vagrant/fabric-ca-client

cleanCA
sudo touch /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_HOME=/home/vagrant/shared/caserver/org0-ca-server/admin' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_TLS_CERTFILES=/home/vagrant/shared/caserver/org0-ca-server/crypto/ca-cert.pem' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_MSPDIR=msp' | sudo tee -a /etc/profile.d/ca_client_vars.sh
source /etc/profile


./fabric-ca-client enroll -u https://rcaorg0admin:rcaorg0adminpw@org0-ca:7054
successln "Enroll the admin user on the org0 CA"

./fabric-ca-client register --id.name orderer1-org0 --id.secret orderer1PW --id.type orderer -u https://org0-ca:7054
./fabric-ca-client register --id.name orderer2-org0 --id.secret orderer2PW --id.type orderer -u https://org0-ca:7054
./fabric-ca-client register --id.name orderer3-org0 --id.secret orderer3PW --id.type orderer -u https://org0-ca:7054
./fabric-ca-client register --id.name admin-org0 --id.secret org0AdminPW --id.type admin --id.affiliation org0.sysadmins --id.attrs "hf.Registrar.Roles=admin,hf.Registrar.Attributes=*,hf.AffiliationMgr=true,hf.Revoker=true,hf.GenCRL=true,admin=true:ecert,abac.init=true:ecert"  -u https://org0-ca:7054
successln "Register entities on the Org0 CA"
successlnbig "Enroll org0 ca admin and register entities on the Org0 CA"



infolnbig "Enroll org1 ca admin and register entities on the Org1 CA"
infoln "Enroll the admin user on the org1 CA"
cd /home/vagrant/fabric-ca-client

cleanCA
sudo touch /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_HOME=/home/vagrant/shared/caserver/org1-ca-server/admin' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_TLS_CERTFILES=/home/vagrant/shared/caserver/org1-ca-server/crypto/ca-cert.pem' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_MSPDIR=msp' | sudo tee -a /etc/profile.d/ca_client_vars.sh
source /etc/profile

./fabric-ca-client enroll -u https://rcaorg1admin:rcaorg1adminpw@org1-ca:7054
successln "Enroll the admin user on the org1 CA"

./fabric-ca-client register --id.name peer1-org1 --id.secret peer1PW --id.type peer -u https://org1-ca:7054
./fabric-ca-client register --id.name peer2-org1 --id.secret peer2PW --id.type peer -u https://org1-ca:7054
./fabric-ca-client register --id.name admin-org1 --id.secret org1AdminPW --id.type admin --id.affiliation org1.sysadmins --id.attrs "hf.Registrar.Roles=admin,hf.Registrar.Attributes=*,hf.AffiliationMgr=true,hf.Revoker=true,hf.GenCRL=true,admin=true:ecert,abac.init=true:ecert"  -u https://org1-ca:7054
./fabric-ca-client register --id.name user1-org1 --id.secret org1User1PW --id.type client --id.affiliation org1.sysusers -u https://org1-ca:7054
successln "Register entities on the Org1 CA"
successlnbig "Enroll org1 ca admin and register entities on the Org1 CA"

successlnbig "CA Servers provision file"
