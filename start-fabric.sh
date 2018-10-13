#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error
set -e

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1
starttime=$(date +%s)
CC_SRC_PATH=github.com/proci/chaincode

# launch network; create channel and join peer to channel
set -ev
# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1

# Now launch the CLI container in order to install, instantiate chaincode
# and prime the ledger with our 10 cars
#docker-compose -f ./docker-compose.yml up -d cli
#sleep 10

#cli=$(docker ps --format "table {{.ID}}" -f "label=com.docker.stack.namespace=ekyc-cli" | tail -1)

#create channel
docker exec -e "CORE_PEER_LOCALMSPID=PeerOrg1" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.hf.nvxtien.io/users/Admin@org1.hf.nvxtien.io/msp" -e "CORE_PEER_ADDRESS=peer0.org1.hf.nvxtien.io:7051" cli peer channel create -o orderer0.hf.nvxtien.io:7050 -c orgschannel1 -f /etc/hyperledger/configtx/orgschannel1.tx

#join peer0
docker exec -e "CORE_PEER_LOCALMSPID=PeerOrg1" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.hf.nvxtien.io/users/Admin@org1.hf.nvxtien.io/msp" -e "CORE_PEER_ADDRESS=peer0.org1.hf.nvxtien.io:7051" cli peer channel join -b order.genesis.block
#join peer1
docker exec -e "CORE_PEER_LOCALMSPID=PeerOrg1" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.hf.nvxtien.io/users/Admin@org1.hf.nvxtien.io/msp" -e "CORE_PEER_ADDRESS=peer1.org1.hf.nvxtien.io:7051" cli peer channel join -b order.genesis.block

#update anchor
docker exec -e "CORE_PEER_LOCALMSPID=PeerOrg1" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.hf.nvxtien.io/users/Admin@org1.hf.nvxtien.io/msp" -e "CORE_PEER_ADDRESS=peer0.org1.hf.nvxtien.io:7051" cli peer channel update -o orderer0.hf.nvxtien.io:7050 -c ekycchannel -f /etc/hyperledger/configtx/PeerOrg1anchors.tx

#install chain code on peer0 & peer1
docker exec -e "CORE_PEER_LOCALMSPID=PeerOrg1" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.hf.nvxtien.io/users/Admin@org1.hf.nvxtien.io/msp" -e "CORE_PEER_ADDRESS=peer0.org1.hf.nvxtien.io:7051" cli peer chaincode install -n ekyc -v 1.0 -p "$CC_SRC_PATH"
docker exec -e "CORE_PEER_LOCALMSPID=PeerOrg1" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.hf.nvxtien.io/users/Admin@org1.hf.nvxtien.io/msp" -e "CORE_PEER_ADDRESS=peer1.org1.hf.nvxtien.io:7051" cli peer chaincode install -n ekyc -v 1.0 -p "$CC_SRC_PATH"
sleep 30

#instantiate chaincode on peer0
docker exec -e "CORE_PEER_LOCALMSPID=PeerOrg1" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.hf.nvxtien.io/users/Admin@org1.hf.nvxtien.io/msp" -e "CORE_PEER_ADDRESS=peer0.org1.hf.nvxtien.io:7051" cli peer chaincode instantiate -o orderer0.hf.nvxtien.io:7050 -C ekycchannel -n ekyc -v 1.0 -c '{"function":"init","Args":[""]}' -P "OR ('PeerOrg1.member')"

sleep 30

printf "\nTotal setup execution time : $(($(date +%s) - starttime)) secs ...\n\n\n"
