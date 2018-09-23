package main

import (
	"fmt"
	"github.com/proci/blockchain"
	"github.com/proci/web"
	"github.com/proci/web/controllers"
	"os"
)

func main() {
	// Definition of the Fabric SDK properties
	fSetup := blockchain.FabricSetup{
		// Network parameters
		OrdererID: "orderer0.hf.nvxtien.io",

		// Channel parameters
		ChannelID:     "orgschannel1",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/proci/fixtures/crypto-config/ordererOrganizations/orgschannel1.tx",

		// Chaincode parameters
		ChainCodeID:     "proci",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/proci/chaincode/",
		OrgAdmin:        "Admin",
		OrgName:         "org1",
		ConfigFile:      "config.yaml",

		// User parameters
		UserName: "User1",
	}

	// Initialization of the Fabric SDK from the previously set properties
	err := fSetup.Initialize()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
		return
	}
	// Close SDK
	defer fSetup.CloseSDK()

	// Install and instantiate the chaincode
	err = fSetup.InstallAndInstantiateCC()
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
		return
	}

	// Launch the web application listening
	app := &controllers.Application{
		Fabric: &fSetup,
	}
	web.Serve(app)
}
