package network

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
      - FABRIC_CA_HOME:/etc/hyperledger/fabric-ca-server
      - FABRIC_CA_SERVER_CA_NAME:ca%d
    ports: 
      - %d:7054
    command: sh -c 'fabric-ca-server start --cfg.identities.allowremove --cfg.affiliations.allowremove -b admin:adminpw -d'
    volumes: 
      - %s/peerOrganizations/org1.%s/ca/:/etc/hyperledger/fabric-ca-server-config
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
    expose:%s
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

func (g *generator) CreateDockerCompose(filename string) {

	var compose = "version: '2'\n\nservices:"

	for i:=0; i<=g.numberOfCa-1; i++ {
		compose += fmt.Sprintf(s_ca, i, i, caPort + i, g.mspBaseDir, g.company, i)
	}

	compose += writeZookeeper(g.numberOfZookeeper, zooPort)
	compose += writeKafka(g.numberOfKafka, g.kafkaReplications, kafkaPort, g.numberOfZookeeper, zooPort)
	compose += writeOrderer(g.numberOfOrderer, ordererPort, g.numberOfKafka, g.company, g.mspBaseDir)

	fmt.Println(compose)

	err := ioutil.WriteFile(os.Getenv("GOPATH") + "/src/github.com/proci" + "/docker-compose.yaml", []byte(compose), 0644)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func expose(numberOfZookeeper, baseZooPort int) (result string) {
	var port = baseZooPort
	for i:=0; i<=numberOfZookeeper-1; i++ {
		result += fmt.Sprintf("\n      - \"%d\"", port)
		port += 1
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
		result += fmt.Sprintf(s_zookeeper, i, i+1, port, servers, expose, i)
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