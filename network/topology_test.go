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

	g.NumberOfOrg(2).
		OrdererType("kafka").
		Company("nvxtien.com").
		Profile("test").
		MSPBaseDir(baseDir).
		PeersPerOrg(2).
		NumberOfOrderer(2).
		NumberOfChannel(2).
		NumberOfCa(2).
		NumberOfZookeeper(3)

	g.CreateDockerCompose(filename)
}