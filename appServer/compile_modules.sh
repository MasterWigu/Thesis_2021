#!/bin/bash


sudo iptables -I FORWARD -i + -o + -j ACCEPT
sudo ip route add 10.10.200.0/24 via 10.10.50.2

sudo touch /etc/profile.d/go_git_vars.sh
echo 'export GOPRIVATE=github.com/MasterWigu/Thesis' | sudo tee -a /etc/profile.d/go_git_vars.sh
source /etc/profile

cp /home/vagrant/certshare/git_configs/.gitconfig /home/vagrant/.gitconfig
cp /home/vagrant/certshare/git_configs/known_hosts /home/vagrant/.ssh/
cp /home/vagrant/certshare/git_configs/id_rsa /home/vagrant/certshare/git_configs/id_rsa.pub /home/vagrant/.ssh
chmod 600 /home/vagrant/.ssh/id_rsa

eval "$(ssh-agent -s)"
ssh-add /home/vagrant/.ssh/id_rsa


rm /home/vagrant/appServer/brokerModule/goapp/broker
rm /home/vagrant/appServer/fabricModule/goapp/fabric
rm /home/vagrant/appServer/ansibleModule/goapp/ansible

cd /home/vagrant/appServer/brokerModule/goapp
go build

cd /home/vagrant/appServer/fabricModule/goapp
go build

cd /home/vagrant/appServer/ansibleModule/goapp
go build




cd /home/vagrant/appServer/docker_images
sudo docker-compose run -d fabric-module
sudo docker-compose run -d broker-module
sudo docker-compose run -d ansible-module

