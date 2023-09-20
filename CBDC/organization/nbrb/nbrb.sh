#!/bin/bash
#
# SPDX-License-Identifier: Apache-2.0
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

# Locate the test-network
cd "${DIR}/../test-network/"
env | sort > $TMPFILE

OVERRIDE_ORG="2"
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