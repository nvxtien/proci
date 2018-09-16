package network

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Container struct {
	Image			string `yaml:"image""`
	Environment		[]string `yaml:"environment,omitempty"`
	Ports			[]string `yaml:"ports,omitempty"`
	Command			string `yaml:"command,omitempty"`
	Volumes			[]string `yaml:"volumes,omitempty"`
	Depends_on 		[]string `yaml:"depends_on,omitempty"`
	Container_name	string `yaml:"container_name,omitempty"`
	Working_dir		string `yaml:"working_dir,omitempty"`
}

type Topo struct {
	//CaAddress string `json:"caAddress"`
	//CaPort string `json:"caPort"`
	//OrdererAddress string `json:"ordererAddress"`
	//OrdererPort string `json:"ordererPort"`
	//KafkaAddress string `json:"kafkaAddress"`
	//KafkaPort string `json:"kafkaPort"`
	//CouchdbAddress string `json:"couchdbAddress"`
	//CouchdbPort string `json:"couchdbPort"`
	//Vp0Address string `json:"vp0Address"`
	//Vp0Port string `json:"vp0port"`
	//EvtAddress string `json:"evtAddress"`
	//EvtPort string `json:"evtort"`
	Services *services	`yaml:"services"`
}

type services struct {
	Ca 			[]*Container	`yaml:"ca,omitempty"`
	Zookeeper 	*Container	`yaml:"zookeeper,omitempty"`
	Kafka 		*Container	`yaml:"kafka,omitempty"`
	Orderer 	*Container	`yaml:"orderer,omitempty"`
	Peer 		*Container	`yaml:"peer,omitempty"`
	Couchdb 	*Container	`yaml:"couchdb,omitempty"`
}
var s_ca =
`
  ca%d:
    image: hyperledger/fabric-ca
    environment: 
      - FABRIC_CA_HOME=/etc/hyperledger/fabric-ca-server
      - FABRIC_CA_SERVER_CA_NAME=ca%d
      - FABRIC_CA_SERVER_CA_CERTFILE=%s/ca.org%d.%s-cert.pem
      - FABRIC_CA_SERVER_CA_KEYFILE=%s/%s
      - FABRIC_CA_SERVER_TLS_ENABLED=true
      - FABRIC_CA_SERVER_TLS_CERTFILE=%s/ca.org%d.%s-cert.pem
      - FABRIC_CA_SERVER_TLS_KEYFILE=%s/%s
    ports: 
      - %d:7054
    command: sh -c 'fabric-ca-server start --cfg.identities.allowremove --cfg.affiliations.allowremove -b admin:adminpw -d'
    volumes: 
      - %s/peerOrganizations/org%d.%s/ca/:/etc/hyperledger/fabric-ca-server-config
    container_name: ca%d
`
var caPort = 7054


var s_zookeeper =
`
  zookeeper%d:
    image: hyperledger/fabric-zookeeper
    environment: 
      - ZOO_MY_ID=%d
      - ZOO_PORT=%d
      - ZOO_SERVERS=%s
    ports:%s
    container_name: zookeeper%d
`
var zooPort = 2181

var s_kafka =
`
  kafka%d:
    image: hyperledger/fabric-kafka
    environment: 
      - KAFKA_BROKER_ID=%d
      - KAFKA_DEFAULT_REPLICATION_FACTOR=%d
      - KAFKA_MESSAGE_MAX_BYTES=103809024
      - KAFKA_REPLICA_FETCH_MAX_BYTES=103809024
      - KAFKA_ZOOKEEPER_CONNECT=%s
      - KAFKA_MIN_INSYNC_REPLICAS=2
      - KAFKA_UNCLEAN_LEADER_ELECTION_ENABLE=false
    depends_on:%s
    container_name: kafka%d
    ports: 
      - %d:9092
`

var kafkaPort = 9092

var s_orderer =
`
  orderer%d.%s:
    image: hyperledger/fabric-orderer
    environment: 
      - ORDERER_GENERAL_LOGLEVEL=ERROR
      - ORDERER_GENERAL_LISTENADDRESS=0.0.0.0
      - ORDERER_GENERAL_LISTENPORT=%d
      - ORDERER_GENERAL_GENESISMETHOD=file
      - ORDERER_GENERAL_GENESISFILE=/opt/hyperledger/fabric/msp/crypto-config/ordererOrganizations/orderer.block
      - ORDERER_GENERAL_LOCALMSPID=OrdererOrg
      - ORDERER_GENERAL_LOCALMSPDIR=%s/msp
      - ORDERER_GENERAL_TLS_ENABLED=true
      - ORDERER_GENERAL_TLS_PRIVATEKEY=%s/tls/server.key
      - ORDERER_GENERAL_TLS_CERTIFICATE=%s/tls/server.crt
      - ORDERER_GENERAL_TLS_ROOTCAS=[%s/tls/ca.crt]
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric
    command: orderer
    volumes: 
      - %s:/opt/hyperledger/fabric/msp/crypto-config
    ports: 
      - %d:%d
    container_name: orderer%d.%s
    depends_on:%s
`

var ordererPort = 5005

var s_peer =
`
  %s:
    image: hyperledger/fabric-peer
    environment: 
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      - CORE_LOGGING_LEVEL=INFO
      - CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=nl_default
      - CORE_PEER_PROFILE_ENABLED=true
      - CORE_PEER_GOSSIP_USELEADERELECTION=true
      - CORE_PEER_GOSSIP_ORGLEADER=false
      - CORE_PEER_GOSSIP_ENDPOINT=%s:%d
      - CORE_PEER_LISTENADDRESS=%s:%d
      - CORE_PEER_ID=%s
      - CORE_PEER_EVENTS_ADDRESS=%s:%d%s
      - CORE_PEER_MSPCONFIGPATH=%s/msp
      - CORE_LEDGER_STATE_STATEDATABASE=CouchDB
      - CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS=couchdb%d:5984
      - CORE_PEER_LOCALMSPID=PeerOrg%d
      - CORE_PEER_ADDRESS=%s:%d
      - CORE_PEER_GOSSIP_EXTERNALENDPOINT=%s:%d
      - CORE_PEER_TLS_ENABLED=true
      - CORE_PEER_TLS_CERT_FILE=%s/tls/server.crt
      - CORE_PEER_TLS_KEY_FILE=%s/tls/server.key
      - CORE_PEER_TLS_ROOTCERT_FILE=%s/tls/ca.crt
    volumes: 
      - /var/run/:/host/var/run/
      - %s:/opt/hyperledger/fabric/msp/crypto-config
    ports: 
      - %d:%d
      - %d:%d
    depends_on: 
      - orderer%d.%s
      - couchdb%d%s
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer
    command: peer node start
    container_name: peer%d.org%d.%s
`
//- CORE_PEER_TLS_CLIENTAUTHREQUIRED=true
//- CORE_PEER_TLS_CLIENTROOTCAS_FILES=%s

var vp0Port = 7051
var event0Port = 6051

var s_couchdb =
`
  couchdb%d:
    image: hyperledger/fabric-couchdb
    ports: 
      - %d:5984
    container_name: couchdb%d
`

var couchdbPort = 5984


func (g *generator) CreateDockerCompose(filename string) {

	fmt.Printf("\nnumber of Orderer: %d", g.numberOfOrderer)
	fmt.Printf("\nnumber of Ca: %d", g.numberOfCa)
	fmt.Printf("\nnumber of Organization: %d", g.numberOfOrg)
	fmt.Printf("\nnumber of Peer per Org: %d", g.peersPerOrg)
	fmt.Printf("\nnumber of Zookeeper: %d", g.numberOfZookeeper)
	fmt.Printf("\nnumber of Kafka Broker: %d", g.numberOfKafka)
	fmt.Printf("\nnumber of Kafka Replication: %d", g.kafkaReplications)
	fmt.Printf("Database type: %s", "couchdb")

	//nPeerPerOrg:  2
	//	number of Orderer:  1
	//	number of Kafka Broker:  1
	//	number of Kafka Replication:  1
	//	number of zookeepers:  1
	//	number of orgs:  2
	//	number of peer:  4
	//	DB type:  couchdb
	//addCA:  2
	//	ca address: 0.0.0.0
	//	ca Port: 7054
	//	orderer address: 0.0.0.0
	//	orderer Port: 5005
	//	couchdb address: 0.0.0.0
	//	couchdb Port: 5984
	//	peer0 address: 0.0.0.0
	//	peer0 Port: 7061
	//	kafka address: 0.0.0.0
	//	kafka Port: 9092
	//	event Hub address: 0.0.0.0
	//	event Hub Port: 6051
	//zport:zport1=2181:2182


	var compose = "version: '2'\n\nservices:"

	compose += writeCa(g.numberOfOrg, g.numberOfCa, caPort, g.mspBaseDir, g.company)
	compose += writeZookeeper(g.numberOfZookeeper, zooPort)
	compose += writeKafka(g.numberOfKafka, g.kafkaReplications, kafkaPort, g.numberOfZookeeper, zooPort)
	compose += writeOrderer(g.numberOfOrderer, ordererPort, g.numberOfKafka, g.company, g.mspBaseDir)
	compose += writePeer(g.numberOfOrderer, g.numberOfOrg, g.peersPerOrg, vp0Port, event0Port, g.company, g.mspBaseDir)
	compose += writeCouchdb(g.numberOfOrg, g.peersPerOrg, couchdbPort)

	fmt.Println(compose)

	err := ioutil.WriteFile(os.Getenv("GOPATH") + "/src/github.com/proci" + "/docker-compose.yaml", []byte(compose), 0644)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func expose(numberOfZookeeper, baseZooPort int) (result string) {
	var port = baseZooPort
	for i:=0; i<=numberOfZookeeper-1; i++ {
		result += fmt.Sprintf("\n      - %d:%d", port, port)
		port += 1
	}
	return
}

func writeCa(numberOfOrg, numberOfCa, baseCaPort int, mspBaseDir, company string) (result string) {
	port := baseCaPort
	config := "/etc/hyperledger/fabric-ca-server-config"

	org := 1
	sk := ""

	for i:=0; i<=numberOfCa-1; i++ {
		dir := fmt.Sprintf("%s/peerOrganizations/org%d.%s/ca", mspBaseDir, org, company)
		filepath.Walk(dir, func(path string, info os.FileInfo, e error) error {
			if e != nil {
				return e
			}
			// check if it is a regular file (not dir)
			if info.Mode().IsRegular() &&  strings.LastIndex(info.Name(), "sk") > 0 {
				sk = info.Name()
				fmt.Println("file name:", info.Name())
				//fmt.Println("file path:", path)
			}
			return nil
		})
		result += fmt.Sprintf(s_ca, i, i, config, org, company, config, sk, config, org, company, config, sk, port + i*100, mspBaseDir, org, company, i)
		// What happen if number of CA is greater than number of Org
		org++
		if org == numberOfOrg+1 {
			org = 1
		}
	}
	return
}

func writeZookeeper(numberOfZookeeper, baseZooPort int) (result string) {
	var servers = ""
	var port = baseZooPort
	for i:=0; i<=numberOfZookeeper-1; i++ {
		servers += fmt.Sprintf("server.%d=zookeeper%d:%d:%d:participant ", i+1, i, port+1, port+2)
		port += 1000
	}

	fmt.Println(servers)

	port = baseZooPort
	for i:=0; i<=numberOfZookeeper-1; i++ {
		var expose = expose(numberOfZookeeper, port)
		result += fmt.Sprintf(s_zookeeper, i, i+1, port, strings.TrimSpace(servers), expose, i)
		port += 1000
	}
	return
}

func writeKafka(numberOfKafka, kafkaReplications, baseKafkaPort, numberOfZookeeper, baseZooPort int) (result string) {

	var zookeeperConnect = fmt.Sprintf("zookeeper%d:%d", 0, baseZooPort)
	var port = baseZooPort
	for i:=1; i<=numberOfZookeeper-1; i++ {
		port += 1000
		zookeeperConnect += fmt.Sprintf(",zookeeper%d:%d", i, port)
	}

	var dependsonZookeeper = dependsonZookeepr(numberOfZookeeper)
	for i:=0; i<=numberOfKafka-1; i++ {
		result += fmt.Sprintf(s_kafka, i, i, kafkaReplications, zookeeperConnect, dependsonZookeeper, i, baseKafkaPort+i)
	}

	return
}

func dependsonZookeepr(numberOfZookeeper int) (result string) {
	for i:=0; i<=numberOfZookeeper-1; i++ {
		result += fmt.Sprintf("\n      - zookeeper%d", i)
	}
	return
}

func writeOrderer(numberOfOrderer, ordererPort, numberOfKafka int, company, mspBaseDir string) (result string) {

	var dependsonKafka = dependsonKafka(numberOfKafka)
	var port = ordererPort
	for i:=0; i<=numberOfOrderer-1; i++ {
		ordererDir := fmt.Sprintf("/opt/hyperledger/fabric/msp/crypto-config/ordererOrganizations/%s/orderers/orderer%d.%s", company, i, company)
		result += fmt.Sprintf(s_orderer, i, company, port + i, ordererDir, ordererDir, ordererDir, ordererDir, mspBaseDir, port+i, port+i, i, company, dependsonKafka)
	}
	return
}

func dependsonKafka(numberOfKafka int) (result string) {
	for i:=0; i<=numberOfKafka-1; i++ {
		result += fmt.Sprintf("\n      - kafka%d", i)
	}
	return
}


func writePeer(numberOfOrderer, numberOfOrg, peersPerOrg, vp0Port, event0Port int, company, mspBaseDir string) (result string) {
	vpPort := vp0Port
	vpPortOut := vp0Port + 1000
	eventPort := event0Port

	//clientauth

	//CLIENTROOTCAS := fmt.Sprintf("/opt/hyperledger/fabric/msp/crypto-config/ordererOrganizations/trade.com/orderers/orderer0.%s/tls/ca.crt", company)
	//for i:=1; i<=numberOfOrderer-1; i++ {
	//	CLIENTROOTCAS += fmt.Sprintf(" /opt/hyperledger/fabric/msp/crypto-config/ordererOrganizations/trade.com/orderers/orderer%d.%s/tls/ca.crt", i, company)
	//}
	//CLIENTROOTCAS += fmt.Sprintf(" /opt/hyperledger/fabric/msp/crypto-config/peerOrganizations/org1.%s/ca/ca.org1.%s-cert.pem", company, company)
	//for i:=2; i<=numberOfOrg; i++ {
	//	CLIENTROOTCAS += fmt.Sprintf(" /opt/hyperledger/fabric/msp/crypto-config/peerOrganizations/org%d.%s/ca/ca.org%d.%s-cert.pem", i, company, i, company)
	//}

	couchdb := 0
	ord := 0

	peer0depend := ""
	boostrap := ""

	for i:=1; i<=numberOfOrg; i++ {
		if ord == numberOfOrderer {
			ord = 0
		}
		for j:=0; j<=peersPerOrg-1; j++ {
			peer := fmt.Sprintf("peer%d.org%d.%s", j, i, company)
			orderer := fmt.Sprintf("org%d.%s", i, company)
			peerDir := fmt.Sprintf("/opt/hyperledger/fabric/msp/crypto-config/peerOrganizations/%s/peers/%s", orderer, peer)


			if j != 0 {
				peer0Port := vp0Port + (i - 1)* peersPerOrg
				peer0 := fmt.Sprintf("peer0.org%d.%s", i, company)
				peer0depend = fmt.Sprintf("\n      - peer0.org%d.%s", i, company)
				boostrap = fmt.Sprintf("\n      - CORE_PEER_GOSSIP_BOOTSTRAP=%s:%d", peer0, peer0Port)
			} else {
				peer0depend = ""
				boostrap = ""
			}

			//result += fmt.Sprintf(s_peer, peer, peer, vpPort, peer, vpPort, peer, peer, eventPort, boostrap, peerDir,
			//	peer, vpPort, peer, vpPort, peerDir, peerDir, peerDir, CLIENTROOTCAS, mspBaseDir, vpPort, vpPort,
			//	eventPort, eventPort, ord, company, couchdb, peer0depend, j, i, company)

			result += fmt.Sprintf(s_peer, peer, peer, vpPort, peer, vpPort, peer, peer, eventPort, boostrap, peerDir, couchdb, i,
				peer, vpPort, peer, vpPort, peerDir, peerDir, peerDir, mspBaseDir, vpPortOut, vpPort,
				eventPort, eventPort, ord, company, couchdb, peer0depend, j, i, company)

			//vpPort++
			vpPortOut++
			eventPort++
			couchdb++
		}
		ord++
	}
	return
}

func writeCouchdb(numberOfOrg, peersPerOrg, couchdbPort int) (result string) {
	port := couchdbPort
	n := numberOfOrg * peersPerOrg
	for i:=0; i<=n-1; i++ {
		result += fmt.Sprintf(s_couchdb, i, port, i)
		port++
	}
	return
}