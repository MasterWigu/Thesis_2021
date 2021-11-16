#!/bin/bash

source /home/vagrant/shared/scripts/prints.sh


infolnbig "Starting General provision script"

infoln "Installing Golang"
wget -cq https://golang.org/dl/go1.16.linux-amd64.tar.gz -O - | sudo tar -xz -C /usr/local
sudo touch /etc/profile.d/golang_vars.sh
echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile.d/golang_vars.sh
echo 'export GOPATH=$HOME/go' | sudo tee -a /etc/profile.d/golang_vars.sh
source /etc/profile
successln "Installing Golang"


infoln "Installing Base Packages"
sudo apt-get update
sudo apt-get install -y libtool libltdl-dev software-properties-common build-essential
successln "Installing Base Packages"

infoln "Installing Docker"
sudo apt-get install -y apt-transport-https ca-certificates curl gnupg
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io
successln "Installing Docker"

infoln "Installing Docker-compose"
sudo curl -fsSL "https://github.com/docker/compose/releases/download/1.28.5/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
docker-compose --version
successln "Installing Docker-compose"


# infoln "Installing npm"
# sudo apt-get install -y npm
# successln "Installing npm"

infoln "Allow external hosts to access containers"
sudo iptables -I FORWARD -i + -o + -j ACCEPT
successln "Allow external hosts to access containers"

successlnbig "General Provision script"
