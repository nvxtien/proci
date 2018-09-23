package network

import (
	"fmt"
	"github.com/btcsuitereleases/btcutil/base58"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/test/integration"
	"github.com/hyperledger/fabric-sdk-go/test/metadata"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"log"
	"os"

	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
)

type Client interface {
	CreateChannel()
	JoinChannel()
	CreateChainCode()
}

type channel struct {

} 

const (
	channelID      = "mychannel"
	orgName        = "Org1"
	orgAdmin       = "Admin"
	ordererOrgName = "OrdererOrg"
)

func (c *channel) CreateChannel() {

	projectDir := fmt.Sprintf("%s/src/github.com/proci", os.Getenv("GOPATH"))
	configOpt := config.FromFile(projectDir + "/config_e2e.yaml")
	//if configOpt == nil {
	//	log.Fatal("configOpt must be not nil")
	//}

	sdk, err := fabsdk.New(configOpt)
	if err != nil || sdk == nil {
		log.Fatal("sdk is not nil", err)
	}

	//clientContext allows creation of transactions using the supplied identity as the credential.
	clientContext := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(ordererOrgName))

	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		log.Fatalf("Failed to create channel management client: %s", err)
	}

	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(orgName))
	if err != nil {
		log.Fatal(err.Error())
	}

	adminIdentity, err := mspClient.GetSigningIdentity(orgAdmin)

	sk, _ := adminIdentity.PrivateKey().Bytes()

	fmt.Println(">>>>>>>>>>>>>>>>>>> " + base58.Encode(sk))

	if err != nil {
		log.Fatal(err.Error())
	}
	req := resmgmt.SaveChannelRequest{ChannelID: channelID,
		ChannelConfigPath: integration.GetChannelConfigPath(channelID + ".tx"),
		SigningIdentities: []msp.SigningIdentity{adminIdentity}}
	txID, err := resMgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("orderer0.nvxtien.com"))
	if txID.TransactionID == "" {
		log.Fatalf("txID can not be empty")
	}

	fmt.Println(txID)
}

func (c *channel) JoinChannel()  {

	configOpt := config.FromFile("config_e2e.yaml")
	sdk, err := fabsdk.New(configOpt)

	//prepare context
	adminContext := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(orgName))

	// Org resource management client
	orgResMgmt, err := resmgmt.New(adminContext)
	if err != nil {
		log.Fatalf("Failed to create new resource management client: %s", err)
	}

	// Org peers join channel
	if err = orgResMgmt.JoinChannel(channelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("orderer0.nvxtien.com")); err != nil {
		log.Fatalf("Org peers failed to JoinChannel: %s", err)
	}
}

var (
	ccID = "example_cc_e2e" + metadata.TestRunID
)

func (c *channel) CreateChainCode() {

	configOpt := config.FromFile("config_e2e.yaml")
	sdk, err := fabsdk.New(configOpt)

	//prepare context
	adminContext := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(orgName))

	// Org resource management client
	orgResMgmt, err := resmgmt.New(adminContext)
	if err != nil {
		log.Fatalf("Failed to create new resource management client: %s", err)
	}


	ccPkg, err := packager.NewCCPackage("github.com/example_cc", integration.GetDeployPath())
	if err != nil {
		log.Fatalf(err.Error())
	}
	// Install example cc to org peers
	installCCReq := resmgmt.InstallCCRequest{Name: ccID, Path: "github.com/example_cc", Version: "0", Package: ccPkg}
	_, err = orgResMgmt.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Fatalf(err.Error())
	}
	// Set up chaincode policy
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"Org1MSP"})
	// Org resource manager will instantiate 'example_cc' on channel
	resp, err := orgResMgmt.InstantiateCC(
		channelID,
		resmgmt.InstantiateCCRequest{Name: ccID, Path: "github.com/example_cc", Version: "0", Args: integration.ExampleCCInitArgs(), Policy: ccPolicy},
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
	)

	fmt.Println(resp)
}