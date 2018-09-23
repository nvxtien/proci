package network

import (
	"log"
	"testing"
	"fmt"
	"os"
)

var g = &generator{}

func TestGenerator(t *testing.T) {

	baseDir := fmt.Sprintf("%s/src/github.com/proci/fixtures/crypto-config", os.Getenv("GOPATH"))

	if src, err := os.Stat(baseDir); os.IsExist(err) && src.IsDir() {
		err = os.RemoveAll(baseDir)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	g.NumberOfOrg(1).
		OrdererType("kafka").
		Company("hf.nvxtien.io").
		Profile("").
		MSPBaseDir(baseDir).
		PeersPerOrg(2).
		NumberOfOrderer(2).
		NumberOfChannel(1).
		NumberOfCa(1).
		NumberOfZookeeper(3). // Z will either be 3, 5, or 7.
		KafkaReplications(3).
		NumberOfKafka(4) // At a minimum, K should be set to 4.

	t.Run("GenConfigTx", func(t *testing.T) {
		g.GenerateConfigTx()
	})

	t.Run("GenerateCryptoCfg", func(t *testing.T) {
		g.GenerateCryptoCfg()
	})
	
	t.Run("ExecuteCryptogen", func(t *testing.T) {
		g.ExecuteCryptogen()
	})

	t.Run("CreateOrderGenesisBlock", func(t *testing.T) {
		g.CreateOrderGenesisBlock()
	})

	t.Run("CreateChannel", func(t *testing.T) {
		g.CreateChannels()
	})

	t.Run("CreateAnchorPeers", func(t *testing.T) {
		g.CreateAnchorPeers()
	})
}

func TestGenConfigTx(t *testing.T) {
	// 2, "kafka", 2, "test","./crypto-config", "example.com"
	g.NumberOfOrg(2).
	OrdererType("kafka").
	Company("nvxtien.com").
	Profile("test").
	MSPBaseDir("./crypto-config").
	PeersPerOrg(2).
	NumberOfOrderer(2)

	g.GenerateConfigTx()
}

func TestGenerateCryptoCfg(t *testing.T) {
	g.GenerateCryptoCfg()
}

func TestExecuteCryptogen(t *testing.T) {
	g.ExecuteCryptogen()
}

func TestCreateOrderGenesisBlock(t *testing.T) {
	g.CreateOrderGenesisBlock()
}

func TestGenerator_CreateChannels(t *testing.T) {
	g.CreateChannels()
}

func TestGenerator_CreateAnchorPeers(t *testing.T) {
	g.CreateAnchorPeers()
}

func BenchmarkGenerateCryptoCfg(b *testing.B) {
	//GenerateCryptoCfg()
}