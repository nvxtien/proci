#!/usr/bin/env bash
echo
echo "##########################################################"
echo "##### Generate certificates using cryptogen tool #########"
echo "##########################################################"

# Define channel name
export CHANNEL_NAME="exchange"
# Please change value as below (all value below is example)

# Refer to configtx.yaml
ORDER_PROFILE="OrgsOrdererGenesis"
CHANNEL_PROFILE="OrgsChannel"

ORG="PeerOrg1"

SERVICE_NAME="nvxtien"
DOMAIN="$SERVICE_NAME.io"

ORG_NAME="sale"
ORG1_DOMAIN="$ORG_NAME.$DOMAIN"

CA_DOMAIN="ca.$ORG_NAME.$DOMAIN"

echo $CA_DOMAIN


ORDERER_DOMAIN="orderer.$DOMAIN"

ORG_MSP="SaleMSP"

GENESIS_PROFILE="Exchange"
NETWORK_ID="exchange"

if [ -d "./fixtures/crypto-config" ]; then
rm -Rf ./fixtures/crypto-config/*
fi

if [ -d "./fixtures/artifacts" ]; then
rm -Rf ./fixtures/artifacts/*
fi

#cp configtx.yamlt configtx.yaml
#sed -i "s/DOMAIN/$DOMAIN/g" configtx.yaml
#sed -i "s/ORG_NAME/$ORG_NAME/g" configtx.yaml
#
#cp crypto-config.yamlt crypto-config.yaml
#sed -i "s/DOMAIN/$DOMAIN/g" crypto-config.yaml
#sed -i "s/ORG_NAME/$ORG_NAME/g" crypto-config.yaml
#
#cp docker-compose.yamlt docker-compose.yaml
#sed -i "s/DOMAIN/$DOMAIN/g" docker-compose.yaml
#sed -i "s/ORG_NAME/$ORG_NAME/g" docker-compose.yaml
#sed -i "s/NETWORK_ID/$NETWORK_ID/g" docker-compose.yaml
#
#cp ../config.yamlt ../config.yaml
#sed -i "s/DOMAIN/$DOMAIN/g" ../config.yaml
#sed -i "s/ORG_NAME/$ORG_NAME/g" ../config.yaml
#sed -i "s/GENESIS_PROFILE/$GENESIS_PROFILE/g" ../config.yaml

set -x
cryptogen generate --config=./crypto-config.yaml --output=./fixtures/crypto-config
res=$?
set +x
if [ $res -ne 0 ]; then
echo "Failed to generate certificates..."
exit 1
fi
echo

echo "##########################################################"
echo "#########  Generating Orderer Genesis block ##############"
echo "##########################################################"
set -x
configtxgen -profile $ORDER_PROFILE -outputBlock ./fixtures/artifacts/orderer.genesis.block -channelID $CHANNEL_NAME
res=$?
set +x
if [ $res -ne 0 ]; then
echo "Failed to generate orderer genesis block..."
exit 1
fi
echo
echo "#################################################################"
echo "### Generating channel configuration transaction 'channel.tx' ###"
echo "#################################################################"
set -x
configtxgen -profile $CHANNEL_PROFILE -outputCreateChannelTx ./fixtures/artifacts/"$SERVICE_NAME".channel.tx -channelID $CHANNEL_NAME
res=$?
set +x
if [ $res -ne 0 ]; then
echo "Failed to generate channel configuration transaction..."
exit 1
fi

echo
echo "#################################################################"
echo "#######    Generating anchor peer update for " + $ORG_NAME +"   ##########"
echo "#################################################################"
set -x
configtxgen -profile $CHANNEL_PROFILE -outputAnchorPeersUpdate ./fixtures/artifacts/"$ORG_MSP"."$NAME".anchors.tx -channelID $CHANNEL_NAME -asOrg $ORG

res=$?
set +x
if [ $res -ne 0 ]; then
echo "Failed to generate anchor peer update for " + $ORG_NAME + "..."
exit 1
fi

#cp start-template.sh start.sh

#sed -i "s/CA_DOMAIN/$CA_DOMAIN/g" start.sh
#sed -i "s/ORG1_DOMAIN/$ORG1_DOMAIN/g" start.sh
#sed -i "s/ORDERER_DOMAIN/$ORDERER_DOMAIN/g" start.sh
#sed -i "s/NAME_ORG/$NAME_ORG/g" start.sh
#sed -i "s/CHANNEL_NAME/$CHANNEL_NAME/g" start.sh

#cp docker-compose-template.yml docker-compose.yml

#sed -i "s/CA_DOMAIN/$CA_DOMAIN/g" docker-compose.yml
#sed -i "s/ORG1_DOMAIN/$ORG1_DOMAIN/g" docker-compose.yml
#sed -i "s/ORDERER_DOMAIN/$ORDERER_DOMAIN/g" docker-compose.yml
#sed -i "s/DOMAIN/$DOMAIN/g" docker-compose.yml

ARCH=$(uname -s | grep Darwin)
if [ "$ARCH" == "Darwin" ]; then
OPTS="-it"
else
OPTS="-i"
fi

CURRENT_DIR=$PWD

echo $CURRENT_DIR
echo $ORG1_DOMAIN

cd ./fixtures/crypto-config/peerOrganizations/$ORG1_DOMAIN/ca/
echo $PWD
PRIV_KEY=$(ls *_sk)
cd "$CURRENT_DIR"
echo $PWD
sed $OPTS "s/CA1_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose.yaml