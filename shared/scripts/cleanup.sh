#!/bin/bash

#clean env vars
function cleanCA() {
	sudo rm /etc/profile.d/ca_client_vars.sh
	sudo touch /etc/profile.d/ca_client_vars.sh
	echo 'unset FABRIC_CA_CLIENT_HOME' | sudo tee -a /etc/profile.d/ca_client_vars.sh
	echo 'unset FABRIC_CA_CLIENT_TLS_CERTFILES' | sudo tee -a /etc/profile.d/ca_client_vars.sh
	echo 'unset FABRIC_CA_CLIENT_MSPDIR' | sudo tee -a /etc/profile.d/ca_client_vars.sh
	source /etc/profile
	sudo rm /etc/profile.d/ca_client_vars.sh
}


function cleanPeer() {
	sudo rm /etc/profile.d/peer_cli_vars.sh
	sudo touch /etc/profile.d/peer_cli_vars.sh
	echo 'unset FABRIC_LOGGING_SPEC' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
	echo 'unset CORE_PEER_ADDRESS' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
	echo 'unset CORE_PEER_LOCALMSPID' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
	echo 'unset CORE_PEER_TLS_ENABLED' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
	echo 'unset CORE_PEER_TLS_ROOTCERT_FILE' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
	echo 'unset CORE_PEER_MSPCONFIGPATH' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
	echo 'unset FABRIC_CFG_PATH' | sudo tee -a /etc/profile.d/peer_cli_vars.sh
	source /etc/profile
	sudo rm /etc/profile.d/peer_cli_vars.sh
}

export -f cleanCA
export -f cleanPeer
