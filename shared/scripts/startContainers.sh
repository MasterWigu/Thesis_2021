iptables -I FORWARD -i + -o + -j ACCEPT

cd /home/vagrant/docker

docker-compose run -d tls-ca
docker-compose run -d org0-ca
docker-compose run -d org1-ca
sleep 5

docker-compose up -d peer1-org1
docker-compose up -d peer2-org1

sleep 5
docker-compose up -d orderer1-org0
docker-compose up -d orderer2-org0
docker-compose up -d orderer3-org0