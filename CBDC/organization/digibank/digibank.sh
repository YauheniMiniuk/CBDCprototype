#!/bin/bash
#
# SPDX-License-Identifier: Apache-2.0
#
TMPFILE=`mktemp`
shopt -s extglob

function _exit(){
    printf "Exiting:%s\n" "$1"
    exit -1
}

: ${CHANNEL_NAME:="mychannel"}
: ${DELAY:="3"}
: ${MAX_RETRY:="5"}
: ${VERBOSE:="false"}

# Where am I?
DIR=${PWD}

# Locate the test network
cd "${DIR}/../test-network"
env | sort > $TMPFILE

OVERRIDE_ORG="1"
. ./scripts/envVar.sh

parsePeerConnectionParameters 1 2

# set the fabric config path
export FABRIC_CFG_PATH="${DIR}/../config"
export PATH="${DIR}/../bin:${PWD}:$PATH"

env | sort | comm -1 -3 $TMPFILE - | sed -E 's/(.*)=(.*)/export \1="\2"/'

rm $TMPFILE

cd "${DIR}"

peer lifecycle chaincode package cp.tar.gz --lang golang --path ${DIR}/chaincode-tea-go --label cp_0
peer lifecycle chaincode install cp.tar.gz
PACKAGE_ID="$(peer lifecycle chaincode queryinstalled | grep -oP 'cp_0:.*(?=,)')"
echo "Exported $PACKAGE_ID"
peer lifecycle chaincode approveformyorg --orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name cbdc -v 0 --package-id $PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA

peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} --channelID mychannel --name cbdc -v 0 --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent

peer lifecycle chaincode querycommitted --channelID mychannel --name cbdc --cafile ${PWD}/../test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

# peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/../test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n cbdc --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/../test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/../test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"Initialize","Args":["CBR", "CBR", "2"]}'

# peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/../test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n cbdc --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/../test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/../test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"Mint","Args":["5000"]}'