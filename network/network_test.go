package network

import (
	"testing"
)

var g = &generator{}

func TestGenConfigTx(t *testing.T) {
	// 2, "kafka", 2, "test","./crypto-config", "example.com"
	//g.NumberOfOrg(2)
	//g.OrdererType("kafka")
	//g.Company("nvxtien.com")
	//g.Profile("test")
	//g.MSPBaseDir("./crypto-config")

	g.GenerateConfigTx()
}

func TestGenerateCryptoCfg(t *testing.T) {
	GenerateCryptoCfg()
}

func TestExecuteCryptogen(t *testing.T) {
	ExecuteCryptogen()
}

func TestCreateOrderGenesisBlock(t *testing.T) {
	CreateOrderGenesisBlock()
}

func BenchmarkGenerateCryptoCfg(b *testing.B) {
	//GenerateCryptoCfg()
}