package network

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"testing"
)

func TestStart(t *testing.T) {
	//projectDir := fmt.Sprintf("%s/src/github.com/proci", os.Getenv("GOPATH"))
	//Start(projectDir)

	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}
	imgs, err := client.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		panic(err)
	}
	for _, img := range imgs {
		fmt.Println("ID: ", img.ID)
		fmt.Println("RepoTags: ", img.RepoTags)
		fmt.Println("Created: ", img.Created)
		fmt.Println("Size: ", img.Size)
		fmt.Println("VirtualSize: ", img.VirtualSize)
		fmt.Println("ParentId: ", img.ParentID)
	}

	//client.CreateContainer()
}