#!/usr/bin/env bash
echo
echo "##########################################################"
echo "##### Generate certificates using cryptogen tool #########"
echo "##########################################################"

# Define channel name
export CHANNEL_NAME="exchange"
# Please change value as below (all value below is example)
ORG_NAME="sale"

NAME="exchange"

DOMAIN="$NAME.io"

CA_DOMAIN="ca.$ORG_NAME.$DOMAIN"

echo $CA_DOMAIN

ORG1_DOMAIN="$ORG_NAME.$DOMAIN"

ORDERER_DOMAIN="orderer.$DOMAIN"

ORG_MSP="SaleMSP"

GENESIS_PROFILE="Exchange"
NETWORK_ID="exchange"

if [ -d "crypto-config" ]; then
rm -Rf crypto-config

cp config.yamlt config.yaml
sed -i "s/DOMAIN/$DOMAIN/g" config.yaml
sed -i "s/ORG_NAME/$ORG_NAME/g" config.yaml
sed -i "s/GENESIS_PROFILE/$GENESIS_PROFILE/g" config.yaml

cp configtx.yamlt configtx.yaml
sed -i "s/DOMAIN/$DOMAIN/g" configtx.yaml
sed -i "s/ORG_NAME/$ORG_NAME/g" configtx.yaml

cp crypto-config.yamlt crypto-config.yaml
sed -i "s/DOMAIN/$DOMAIN/g" crypto-config.yaml
sed -i "s/ORG_NAME/$ORG_NAME/g" crypto-config.yaml

cp docker-compose.yamlt docker-compose.yaml
sed -i "s/DOMAIN/$DOMAIN/g" docker-compose.yaml
sed -i "s/ORG_NAME/$ORG_NAME/g" docker-compose.yaml
sed -i "s/NETWORK_ID/$NETWORK_ID/g" ../docker-compose.yaml

#CURRENT_DIR=$PWD

#cd ..


cd CURRENT_DIR

fi
set -x
cryptogen generate --config=./crypto-config.yaml
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
configtxgen -profile ChainHero -outputBlock ./artifacts/orderer.genesis.block
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
configtxgen -profile ChainHero -outputCreateChannelTx ./artifacts/"$NAME".channel.tx -channelID $CHANNEL_NAME
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
#configtxgen -profile ChainHero -outputAnchorPeersUpdate ./config/${NAME_ORG}anchors.tx -channelID $CHANNEL_NAME -asOrg $NAME_ORG
configtxgen -profile ChainHero -outputAnchorPeersUpdate ./artifacts/"$ORG_MSP"."$NAME".anchors.tx -channelID $CHANNEL_NAME -asOrg SaleOrg

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

CURRENT_DIR=$PWD

echo $CURRENT_DIR
echo $ORG1_DOMAIN

cd ./crypto-config/peerOrganizations/$ORG1_DOMAIN/ca/
echo $PWD
PRIV_KEY=$(ls *_sk)
cd "$CURRENT_DIR"
echo $PWD
sed -i "s/CA1_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose.yaml