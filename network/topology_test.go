package network

import (
	"testing"
	"fmt"
	"os"
)

func TestCreateDockerCompose(t *testing.T) {
	projectDir := fmt.Sprintf("%s/src/github.com/proci", os.Getenv("GOPATH"))
	filename := fmt.Sprintf("%s/network.json", projectDir)

	baseDir := fmt.Sprintf("%s/src/github.com/proci/crypto-config", os.Getenv("GOPATH"))

	//g.NumberOfOrg(6).
	//	OrdererType("kafka").
	//	Company("trade.com").
	//	Profile("test").
	//	MSPBaseDir(baseDir).
	//	PeersPerOrg(2).
	//	NumberOfOrderer(3).
	//	NumberOfChannel(3).
	//	NumberOfCa(6).
	//	NumberOfZookeeper(3).
	//	KafkaReplications(3).
	//	NumberOfKafka(3)


	g.NumberOfOrg(1).
		OrdererType("kafka").
		Company("trade.com").
		Profile("test").
		MSPBaseDir(baseDir).
		PeersPerOrg(2).
		NumberOfOrderer(1).
		NumberOfChannel(3).
		NumberOfCa(1).
		NumberOfZookeeper(3). // Z will either be 3, 5, or 7.
		KafkaReplications(3).
		NumberOfKafka(4) // At a minimum, K should be set to 4.

	g.CreateDockerCompose(filename)
}