package network

import (
	"testing"
	"fmt"
	"os"
)

var g = &generator{}

func TestGenerator(t *testing.T) {

	baseDir := fmt.Sprintf("%s/src/github.com/proci/crypto-config", os.Getenv("GOPATH"))

	g.NumberOfOrg(2).
		OrdererType("kafka").
		Company("nvxtien.com").
		Profile("test").
		MSPBaseDir(baseDir).
		PeersPerOrg(2).
		NumberOfOrderer(3).
		NumberOfChannel(2)

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