package network

import (
	"testing"
)

func TestGenConfigTx(t *testing.T) {

	GenerateConfigTx("./crypto-config", "example.com")
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
	GenerateCryptoCfg()
}