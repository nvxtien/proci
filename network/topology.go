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
      - ZOO_PORT=2181
      - ZOO_SERVERS=server.1=zookeeper0:2182:2183:participant server.2=zookeeper1:3182:3183:participant server.3=zookeeper2:4182:4183:participant
    expose:
      - "2181"
      - "2182"
      - "2183"
    container_name: zookeeper0
`
var zooPort = 2181

//       - ZOO_SERVERS=server.1=zookeeper0:2182:2183:participant server.2=zookeeper1:3182:3183:participant server.3=zookeeper2:4182:4183:participant
//    expose:
//      - "2181"
//      - "2182"
//      - "2183"

//"zookeeper": {
//"image": "hyperledger/fabric-zookeeper",
//"environment": {
//"ZOO_MY_ID": "0",
//"ZOO_PORT": "zookeeper:2181",
//"ZOO_SERVERS": "ZOO_SERVERS=server.1=zookeeper0:2182:2183:participant"
//},
//"ports": [
//"zooPort:2181"
//],
//"container_name": "zookeeper"
//}

func (g *generator) CreateDockerCompose(filename string) {

	var compose = "version: '2'\n\nservices:"

	for i:=0; i<=g.numberOfCa-1; i++ {
		compose += fmt.Sprintf(s_ca, i, i, caPort + i, g.mspBaseDir, g.company, i)
	}

	//ca := Container{}
	//err := yaml.Unmarshal([]byte(s_ca), &ca)
	//if err != nil {
	//	log.Fatalf(err.Error())
	//}
	//
	//topo := Topo{&services{}}
	//topo.Services.Ca = append(topo.Services.Ca, &ca)
	//topo.Services.Ca = append(topo.Services.Ca, &ca)

	//def, err := ioutil.ReadFile(filename)
	//if err != nil {
	//	log.Fatalf(err.Error())
	//}

	//fmt.Printf("%s", string(def))


	//fmt.Printf("%s\n", t.Services.Ca.Command)

	//err = yaml.Unmarshal(def, &t)
	//fmt.Printf("%s\n", compose)

	//compose, _ := yaml.Marshal(&topo)
	//fmt.Printf("%s\n", compose)

	err := ioutil.WriteFile(os.Getenv("GOPATH") + "/src/github.com/proci" + "/docker-compose.yaml", []byte(compose), 0644)
	if err != nil {
		log.Fatalf(err.Error())
	}

	//err = yaml.Unmarshal(compose, &t)
	//compose, _ = yaml.Marshal(&t)
	//fmt.Printf("%s\n", t.Services.Peer.Environment)


}