#!/bin/bash

source /home/vagrant/shared/scripts/cleanup.sh
source /home/vagrant/shared/scripts/prints.sh


infolnbig "Create Certificates and Keys for inter-module TLS communication"

#go to the fabric-ca-client dir
cd /home/vagrant/fabric-ca-client

#make ourselves the TLS CA admin
cleanCA
sudo touch /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_HOME=/home/vagrant/shared/caserver/tls-ca-server/admin' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_TLS_CERTFILES=/home/vagrant/shared/caserver/tls-ca-server/crypto/tls-ca-cert.pem' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_MSPDIR=msp' | sudo tee -a /etc/profile.d/ca_client_vars.sh
source /etc/profile


infoln "Registering the modules on the TLS CA"
./fabric-ca-client register --id.name broker-module --id.secret brmodule --id.type client -u https://tls-ca:7054
./fabric-ca-client register --id.name fabric-module --id.secret fbmodule --id.type client -u https://tls-ca:7054
./fabric-ca-client register --id.name webserver-module --id.secret wsmodule --id.type client -u https://tls-ca:7054
./fabric-ca-client register --id.name ansible-module --id.secret asmodule --id.type client -u https://tls-ca:7054
./fabric-ca-client register --id.name terraform-module --id.secret tfmodule --id.type client -u https://tls-ca:7054

successln "Registering the modules on the TLS CA"

#create temp folder for the modules msps
mkdir -p /home/vagrant/certshare/temp

#copy tls ca cert to certshare folder
cp /home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem /home/vagrant/certshare/tls-cacert.pem






infoln "Enroll the broker module"
cleanCA
sudo touch /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_HOME=/home/vagrant/certshare/temp' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_TLS_CERTFILES=/home/vagrant/shared/caserver/tls-ca-server/crypto/tls-ca-cert.pem' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_MSPDIR=msp' | sudo tee -a /etc/profile.d/ca_client_vars.sh
source /etc/profile

./fabric-ca-client enroll -u https://broker-module:brmodule@tls-ca:7054 --enrollment.profile tls --csr.hosts 'broker-module'

successln "Enroll the broker module"

infoln "Copying files"

rm -rf /home/vagrant/certshare/broker-module
mkdir -p /home/vagrant/certshare/broker-module

cp /home/vagrant/certshare/temp/msp/keystore/* /home/vagrant/certshare/broker-module/key.pem
cp /home/vagrant/certshare/temp/msp/signcerts/* /home/vagrant/certshare/broker-module/cert.pem
cp /home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem /home/vagrant/certshare/broker-module/tls-cacert.pem

rm -rf /home/vagrant/certshare/temp/*

successln "Copying files"
successln "Enroll the broker module"





infoln "Enroll the fabric1 module"
cleanCA
sudo touch /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_HOME=/home/vagrant/certshare/temp' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_TLS_CERTFILES=/home/vagrant/shared/caserver/tls-ca-server/crypto/tls-ca-cert.pem' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_MSPDIR=msp' | sudo tee -a /etc/profile.d/ca_client_vars.sh
source /etc/profile

./fabric-ca-client enroll -u https://fabric-module:fbmodule@tls-ca:7054 --enrollment.profile tls --csr.hosts 'fabric-module'

successln "Enroll the fabric1 module"

infoln "Copying files"

rm -rf /home/vagrant/certshare/fabric-module
mkdir -p /home/vagrant/certshare/fabric-module

cp /home/vagrant/certshare/temp/msp/keystore/* /home/vagrant/certshare/fabric-module/key.pem
cp /home/vagrant/certshare/temp/msp/signcerts/* /home/vagrant/certshare/fabric-module/cert.pem
cp /home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem /home/vagrant/certshare/fabric-module/tls-cacert.pem

rm -rf /home/vagrant/certshare/temp/*

successln "Copying files"
successln "Enroll the fabric1 module"






infoln "Enroll the webserver1 module"
cleanCA
sudo touch /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_HOME=/home/vagrant/certshare/temp' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_TLS_CERTFILES=/home/vagrant/shared/caserver/tls-ca-server/crypto/tls-ca-cert.pem' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_MSPDIR=msp' | sudo tee -a /etc/profile.d/ca_client_vars.sh
source /etc/profile

./fabric-ca-client enroll -u https://webserver-module:wsmodule@tls-ca:7054 --enrollment.profile tls --csr.hosts 'webserver-module'

successln "Enroll the webserver1 module"

infoln "Copying files"

rm -rf /home/vagrant/certshare/webserver-module
mkdir -p /home/vagrant/certshare/webserver-module

cp /home/vagrant/certshare/temp/msp/keystore/* /home/vagrant/certshare/webserver-module/key.pem
cp /home/vagrant/certshare/temp/msp/signcerts/* /home/vagrant/certshare/webserver-module/cert.pem
cp /home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem /home/vagrant/certshare/webserver-module/tls-cacert.pem

rm -rf /home/vagrant/certshare/temp/*

successln "Copying files"
successln "Enroll the webserver1 module"




infoln "Enroll the ansible1 module"
cleanCA
sudo touch /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_HOME=/home/vagrant/certshare/temp' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_TLS_CERTFILES=/home/vagrant/shared/caserver/tls-ca-server/crypto/tls-ca-cert.pem' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_MSPDIR=msp' | sudo tee -a /etc/profile.d/ca_client_vars.sh
source /etc/profile

./fabric-ca-client enroll -u https://ansible-module:asmodule@tls-ca:7054 --enrollment.profile tls --csr.hosts 'ansible-module'

successln "Enroll the ansible1 module"

infoln "Copying files"

rm -rf /home/vagrant/certshare/ansible-module
mkdir -p /home/vagrant/certshare/ansible-module

cp /home/vagrant/certshare/temp/msp/keystore/* /home/vagrant/certshare/ansible-module/key.pem
cp /home/vagrant/certshare/temp/msp/signcerts/* /home/vagrant/certshare/ansible-module/cert.pem
cp /home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem /home/vagrant/certshare/ansible-module/tls-cacert.pem

rm -rf /home/vagrant/certshare/temp/*

successln "Copying files"
successln "Enroll the ansible1 module"



infoln "Enroll the terraform1 module"
cleanCA
sudo touch /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_HOME=/home/vagrant/certshare/temp' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_TLS_CERTFILES=/home/vagrant/shared/caserver/tls-ca-server/crypto/tls-ca-cert.pem' | sudo tee -a /etc/profile.d/ca_client_vars.sh
echo 'export FABRIC_CA_CLIENT_MSPDIR=msp' | sudo tee -a /etc/profile.d/ca_client_vars.sh
source /etc/profile

./fabric-ca-client enroll -u https://terraform-module:tfmodule@tls-ca:7054 --enrollment.profile tls --csr.hosts 'terraform-module'

successln "Enroll the terraform1 module"

infoln "Copying files"

rm -rf /home/vagrant/certshare/terraform-module
mkdir -p /home/vagrant/certshare/terraform-module

cp /home/vagrant/certshare/temp/msp/keystore/* /home/vagrant/certshare/terraform-module/key.pem
cp /home/vagrant/certshare/temp/msp/signcerts/* /home/vagrant/certshare/terraform-module/cert.pem
cp /home/vagrant/shared/nodesserver/org1/msp/tlscacerts/tls-cert.pem /home/vagrant/certshare/terraform-module/tls-cacert.pem

rm -rf /home/vagrant/certshare/temp/*

successln "Copying files"
successln "Enroll the terraform1 module"


#delete temp folder
rm -rf /home/vagrant/certshare/temp

successlnbig "Create Certificates and Keys for inter-module TLS communication"

