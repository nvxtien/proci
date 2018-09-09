package network

import (
	"testing"
	"fmt"
	"os"
)

func TestCreateDockerCompose(t *testing.T) {
	projectDir := fmt.Sprintf("%s/src/github.com/proci", os.Getenv("GOPATH"))
	filename := fmt.Sprintf("%s/network.json", projectDir)

	CreateDockerCompose(filename)
}