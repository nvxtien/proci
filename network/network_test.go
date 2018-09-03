package network

import (
	"testing"
)

func TestGenConfigTx(t *testing.T) {

	GenerateConfigTx("/root/gopath/src/github.com/hyperledger/fabric-test/fabric/common/tools/cryptogen/crypto-config", "example.com")
}