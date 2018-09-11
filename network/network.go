package network

import (
	"strings"
	"os"
	"bufio"
	"io/ioutil"
	"fmt"
	"log"
	"os/exec"
	"github.com/proci"
)

type Generator interface {
	NumberOfOrg(int) Generator
	OrdererType(proci.OrdererType) Generator
	NumberOfKafka(int) Generator
	PeersPerOrg(int) Generator
	NumberOfOrderer(int) Generator
	NumberOfChannel(int) Generator
	NumberOfCa(int) Generator
	Profile(string) Generator
	MSPBaseDir(string) Generator
	Company(string) Generator

	GenerateConfigTx() (bool, error)
	GenerateCryptoCfg()
	ExecuteCryptogen()
	CreateOrderGenesisBlock()
	CreateChannels()
}

type generator struct {
	numberOfOrg		int
	ordererType		proci.OrdererType
	numberOfKafka	int
	peersPerOrg		int
	numberOfOrderer	int
	numberOfChannel	int
	numberOfCa		int
	profile 		string
	mspBaseDir		string
	company 		string
}

func (g *generator) NumberOfOrg(numberOfOrg int) Generator {
	g.numberOfOrg = numberOfOrg
	return g
}

func (g *generator) OrdererType(ordererType proci.OrdererType) Generator {
	g.ordererType = ordererType
	return g
}

func (g *generator) NumberOfKafka(numberOfKafka int) Generator {
	g.numberOfKafka = numberOfKafka
	return g
}

func (g *generator) PeersPerOrg(peersPerOrg int) Generator {
	g.peersPerOrg = peersPerOrg
	return g
}

func (g *generator) NumberOfOrderer(numberOfOrderer int) Generator {
	g.numberOfOrderer = numberOfOrderer
	return g
}

func (g *generator) NumberOfCa(numberOfCa int) Generator {
	g.numberOfCa = numberOfCa
	return g
}

func (g *generator) NumberOfChannel(numberOfChannel int) Generator {
	g.numberOfChannel = numberOfChannel
	return g
}

func (g *generator) Profile(profile string) Generator {
	g.profile = profile
	return g
}

func (g *generator) MSPBaseDir(mspBaseDir string) Generator {
	g.mspBaseDir = mspBaseDir
	return g
}

func (g *generator) Company(company string) Generator {
	g.company = company
	return g
}

func New() Generator {
	return &generator{}
}

// create configtx.yaml.
func (g *generator) GenerateConfigTx() (bool, error) {
	log.Println( " - generate configtx.yaml ...")

	numberOfOrg := g.numberOfOrg
	ordererType := g.ordererType
	numberOfKafka := g.numberOfKafka
	profile := g.profile
	MSPBaseDir := g.mspBaseDir
	comName := g.company


	if _, err := os.Stat(os.Getenv("GOPATH") + "/src/github.com/proci/configtx.yaml"); os.IsExist(err) {
		err := os.Remove(os.Getenv("GOPATH") + "/src/github.com/proci/configtx.yaml")
		if err != nil {
			log.Fatalf("Can not delete configtx.yaml")
		}
	}

	template, err := os.Open(os.Getenv("GOPATH") + "/src/github.com/proci/configtx.yaml-in")
	if err != nil {
		return false, err
	}

	defer template.Close()

	scanner := bufio.NewScanner(template)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	configtx := []string{""}
	for _, line := range lines {
		if !strings.Contains(line, "#") {
			if strings.Contains(line, "&OrdererOrg") {
				configtx = append(configtx, createOrdererOrg(line, MSPBaseDir, comName)...)
				//continue

			} else if strings.Contains(line, "&PeerOrg") {
				for i := 1; i <= numberOfOrg; i++ {
					configtx = append(configtx, createPeerOrg(line, i, MSPBaseDir, comName)...)
				}
				//continue

			} else if strings.Contains(line, "OrdererType:") {
				configtx = append(configtx, fmt.Sprintf("    OrdererType: %s\n", ordererType))
				//continue

			} else if strings.Contains(line, "Addresses:") {
				configtx = append(configtx, line)
				configtx = append(configtx, "\n")
				for i := 1; i <= numberOfOrg; i++ {
					configtx = append(configtx, fmt.Sprintf("         - orderer%d.%s:%d\n", i - 1, comName, 5005 + (i - 1)))
				}
				//continue

			} else if strings.Contains(line, "Brokers") {
				configtx = append(configtx, line)
				configtx = append(configtx, "\n")
				for i := 0; i < numberOfKafka; i++ {
					configtx = append(configtx, fmt.Sprintf("             - kafka%d:%d\n", i, 9092))
				}

			} else if strings.Contains(line, "&ProfileOrderString") {
				configtx = append(configtx, fmt.Sprintf("    %sOrgsOrdererGenesis:\n", profile))
				configtx = append(configtx, fmt.Sprintf("        <<: *ChannelDefaults\n"))

			} else if strings.Contains(line, "&ProfileOrgString") {
				configtx = append(configtx, fmt.Sprintf("    %sorgschannel:\n", profile))

			} else if strings.Contains(line, "*PeerOrg") {
				for i := 1; i <= numberOfOrg; i++ {
					configtx = append(configtx, fmt.Sprintf("                - *PeerOrg%d\n", i))
				}

			} else {
				configtx = append(configtx, line)
				configtx = append(configtx, "\n")
			}

		} else {
			configtx = append(configtx, line)
			configtx = append(configtx, "\n")
		}
	}

	err = ioutil.WriteFile(os.Getenv("GOPATH") + "/src/github.com/proci/configtx.yaml", []byte(strings.Join(configtx, "")), 0644)
	if err != nil {
		return false, err
	}

	return true, nil
}

func createPeerOrg(line string, i int, MSPBaseDir string, comName string) []string {
	var configtx []string
	configtx = append(configtx, fmt.Sprintf("%s%d\n", line, i))
	configtx = append(configtx, fmt.Sprintf("        Name: PeerOrg%d\n", i))
	configtx = append(configtx, fmt.Sprintf("        ID: PeerOrg%d\n", i))
	configtx = append(configtx, fmt.Sprintf("        MSPDir: %s/peerOrganizations/org%d.%s/msp\n", MSPBaseDir, i, comName))
	configtx = append(configtx, fmt.Sprintf("        ID: PeerOrg%d\n", i))
	configtx = append(configtx, fmt.Sprintf("        ID: PeerOrg%d\n", i))
	configtx = append(configtx, "        Policies:\n")
	configtx = append(configtx, "            Readers:\n")
	configtx = append(configtx, "                Type: Signature\n")
	configtx = append(configtx, fmt.Sprintf("                Rule: \"OR('PeerOrg%d.admin', 'PeerOrg%d.peer')\"\n", i, i))
	configtx = append(configtx, "            Writers:\n")
	configtx = append(configtx, "                Type: Signature\n")
	configtx = append(configtx, fmt.Sprintf("                Rule: \"OR('PeerOrg%d.admin', 'PeerOrg%d.client')\"\n", i, i))
	configtx = append(configtx, "            Admins:\n")
	configtx = append(configtx, "                Type: Signature\n")
	configtx = append(configtx, fmt.Sprintf("                Rule: \"OR('PeerOrg%d.admin')\"\n", i))
	configtx = append(configtx, "\n")
	configtx = append(configtx, "        AnchorPeers:\n")
	configtx = append(configtx, fmt.Sprintf("            - Host: peer0.org%d.%s\n", i, comName))
	configtx = append(configtx, fmt.Sprintf("              Port: %d\n", 7060 + 2 *(i - 1) + 1))

	configtx = append(configtx, "\n")
	return configtx
}

func createOrdererOrg(line string, MSPBaseDir string, comName string) []string {
	var str []string
	str = append(str, line)
	str = append(str, "\n")
	str = append(str, "        Name: OrdererOrg\n")
	str = append(str, "        ID: OrdererOrg\n")
	str = append(str, fmt.Sprintf("        MSPDir: %s/ordererOrganizations/%s/msp\n", MSPBaseDir, comName))

	ordMSP := "OrdererOrg"
	str = append(str, "        Policies:\n")
	str = append(str, "            Readers:\n")
	str = append(str, "                Type: Signature\n")
	str = append(str, fmt.Sprintf("                Rule: \"OR('%s.member')\"\n", ordMSP))
	str = append(str, "            Writers:\n")
	str = append(str, "                Type: Signature\n")
	str = append(str, fmt.Sprintf("                Rule: \"OR('%s.member')\"\n", ordMSP))
	str = append(str, "            Admins:\n")
	str = append(str, "                Type: Signature\n")
	str = append(str, fmt.Sprintf("                Rule: \"OR('%s.admin')\"\n", ordMSP))

	return str
	//return strings.Join(str, "")
}

func (g *generator) GenerateCryptoCfg() {
	log.Println("************* generate crypto-config.yaml *************")

	if _, err := os.Stat(os.Getenv("GOPATH") + "/src/github.com/proci" + "/crypto-config.yaml"); os.IsExist(err) {
		err := os.Remove(os.Getenv("GOPATH") + "/src/github.com/proci" + "/crypto-config.yaml")
		if err != nil {
			log.Fatalf("Can not delete crypto-config.yaml")
		}
	}

	nOrg := g.numberOfOrg
	nOrderer := g.numberOfOrderer
	peersPerOrg := g.peersPerOrg
	log.Printf("peersPerOrg %d",peersPerOrg)
	comName := g.company
	cryptocfg := []string{""}

	cryptocfg = append(cryptocfg, "OrdererOrgs:\n")
	cryptocfg = append(cryptocfg, "    - Name: Orderer\n")
	cryptocfg = append(cryptocfg, fmt.Sprintf("      Domain: %s\n", comName))
	cryptocfg = append(cryptocfg, "      Template:\n")
	cryptocfg = append(cryptocfg, fmt.Sprintf("        Count: %d\n", nOrderer))

	cryptocfg = append(cryptocfg, "PeerOrgs:\n")
	for i := 1; i <= nOrg; i++ {
		cryptocfg = append(cryptocfg, fmt.Sprintf("    - Name: PeerOrg%d\n", i))
		cryptocfg = append(cryptocfg, fmt.Sprintf("      Domain: org%d.%s\n", i, comName))
		cryptocfg = append(cryptocfg, "      EnableNodeOUs: true\n")
		cryptocfg = append(cryptocfg, "      Template:\n")
		cryptocfg = append(cryptocfg, fmt.Sprintf("        Count: %d\n", peersPerOrg))
		cryptocfg = append(cryptocfg, "      Users:\n")
		cryptocfg = append(cryptocfg, "        Count: 1\n")
	}

	err := ioutil.WriteFile(os.Getenv("GOPATH") + "/src/github.com/proci" + "/crypto-config.yaml", []byte(strings.Join(cryptocfg, "")), 0644)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (g *generator) ExecuteCryptogen() {
	log.Println("************* execute cryptogen *************")

	if _, err := os.Stat(os.Getenv("GOPATH") + "/src/github.com/proci" + "/crypto-config"); os.IsExist(err) {
		err := os.Remove(os.Getenv("GOPATH") + "/src/github.com/proci" + "/crypto-config")
		if err != nil {
			log.Fatalf("Can not delete crypto-config")
		}
	}

	// cryptogen should be in $GOPATH/bin ???
	path, err := exec.LookPath("cryptogen")
	if err != nil {
		log.Fatalf("Please install cryptogen: %s", err)
	}
	log.Printf("cryptogen is available at %s\n", path)

	cmd := exec.Command("cryptogen", "generate", "--output="+ os.Getenv("GOPATH") + "/src/github.com/proci" +  "/crypto-config", "--config=" + os.Getenv("GOPATH") + "/src/github.com/proci" +  "/crypto-config.yaml")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(stdoutStderr)
	}
}

func (g *generator) CreateOrderGenesisBlock() {
	log.Println(" - create orderer genesis block ...")

	profile := g.profile
	ordererDir := os.Getenv("GOPATH") + "/src/github.com/proci" + "/crypto-config/ordererOrganizations"
	ordBlock := fmt.Sprintf("%s/%sorderer.block", ordererDir, profile)
	//testChannel := "channel"

	//CHAN_PROFILE := fmt.Sprintf("%sChannel", profile)
	ORDERER_PROFILE := fmt.Sprintf("%sOrgsOrdererGenesis", profile)
	//ORG_PROFILE := fmt.Sprintf("%sorgschannel", profile)

	path, err := exec.LookPath("configtxgen")
	if err != nil {
		log.Fatal("Please install configtxgen")
	}

	log.Printf("configtxgen is available at %s\n", path)

	if _, err := os.Stat(ordererDir); os.IsNotExist(err) {
		//err := os.Remove("crypto-config")
		if err != nil {
			log.Fatalf("Can not find crypto-config")
		}
		log.Printf("Can not find crypto-config")
	}

	// configtxgen -profile "testOrgsOrdererGenesis" -channelID "channel" -outputBlock "./crypto-config/ordererOrganizations/orderer.block"

	// configtxgen -profile "testOrgsOrdererGenesis" -channelID "channel" -outputBlock "./crypto-config/ordererOrganizations/orderer.block"

	projectDir := os.Getenv("GOPATH") + "/src/github.com/proci"

	cmd := exec.Command("configtxgen",
		fmt.Sprintf("-configPath=%s", projectDir),
		fmt.Sprintf("-profile=%s", ORDERER_PROFILE),
		fmt.Sprintf("-channelID=%s", "channel"),
		fmt.Sprintf("-outputBlock=%s", ordBlock))

	//cmd := exec.Command("configtxgen", "-profile=testOrgsOrdererGenesis", "-channelID=channel", "-outputBlock=/home/tiennv14/devenv/gopath/src/github.com/proci/network/crypto-config/ordererOrganizations/orderer.block")

	log.Println(cmd.Args)

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(stdoutStderr)
	}

	//fmt.Printf("%s\n", stdoutStderr)
}

func (g *generator) CreateChannels() {
	log.Println("create channels ...")

	path, err := exec.LookPath("configtxgen")
	if err != nil {
		log.Fatal("Please install configtxgen")
	}

	log.Printf("configtxgen is available at %s\n", path)

	nChannel := g.numberOfChannel
	ordererDir := os.Getenv("GOPATH") + "/src/github.com/proci" + "/crypto-config/ordererOrganizations"
	ORG_PROFILE := fmt.Sprintf("%sorgschannel", g.profile)
	//for (( i=1; i<=$nChannel; i++ ))
	//do
	//channelTx=$ordererDir"/"$ORG_PROFILE$i".tx"
	//echo "$CFGEXE -profile $ORG_PROFILE -channelID $ORG_PROFILE"$i" -outputCreateChannelTx $channelTx"
	//$CFGEXE -profile $ORG_PROFILE -channelID $ORG_PROFILE"$i" -outputCreateChannelTx $channelTx
	//done

	projectDir := os.Getenv("GOPATH") + "/src/github.com/proci"

	for i := 1; i <= nChannel; i++ {
		channelTx := fmt.Sprintf("%s/%s%d.tx", ordererDir, ORG_PROFILE, i)
		fmt.Printf("%s", channelTx)
		cmd := exec.Command("configtxgen",
			fmt.Sprintf("-configPath=%s", projectDir),
			fmt.Sprintf("-profile=%s", ORG_PROFILE),
			fmt.Sprintf("-channelID=%s%d", ORG_PROFILE, i),
			fmt.Sprintf("-outputCreateChannelTx=%s", channelTx))

		//cmd := exec.Command("configtxgen", "-profile=testOrgsOrdererGenesis", "-channelID=channel", "-outputBlock=/home/tiennv14/devenv/gopath/src/github.com/proci/network/crypto-config/ordererOrganizations/orderer.block")

		log.Println(cmd.Args)

		stdoutStderr, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatal(stdoutStderr)
		}
	}
}

func (g *generator) CreateAnchorPeers() {
	log.Printf("Create anchor peer ...")

	path, err := exec.LookPath("configtxgen")
	if err != nil {
		log.Fatal("Please install configtxgen")
	}

	log.Printf("configtxgen is available at %s\n", path)


	//for (( i=1; i<=$nOrg; i++ ))
	//do
	//orgMSP="PeerOrg"$i
	//OrgMSP=$ordererDir"/"$orgMSP"anchors.tx"
	//echo "$CFGEXE -profile $ORG_PROFILE -outputAnchorPeersUpdate $OrgMSP -channelID $ORG_PROFILE"$i" -asOrg $orgMSP"
	//$CFGEXE -profile $ORG_PROFILE -outputAnchorPeersUpdate $OrgMSP -channelID $ORG_PROFILE"$i" -asOrg $orgMSP
	//done

	ordererDir := os.Getenv("GOPATH") + "/src/github.com/proci" + "/crypto-config/ordererOrganizations"
	ORG_PROFILE := fmt.Sprintf("%sorgschannel", g.profile)
	projectDir := os.Getenv("GOPATH") + "/src/github.com/proci"

	for i:=1; i<=g.numberOfOrg; i++ {
		orgMSP := fmt.Sprintf("PeerOrg%d", i)
		OrgMSP := fmt.Sprintf("%s/%sanchors.tx", ordererDir, orgMSP)
		cmd := exec.Command("configtxgen",
			fmt.Sprintf("-configPath=%s", projectDir),
			fmt.Sprintf("-profile=%s", ORG_PROFILE),
			fmt.Sprintf("-outputAnchorPeersUpdate=%s", OrgMSP),
			fmt.Sprintf("-channelID=%s%d", ORG_PROFILE, i),
			fmt.Sprintf("-asOrg=%s", orgMSP))

		//cmd := exec.Command("configtxgen", "-profile=testOrgsOrdererGenesis", "-channelID=channel", "-outputBlock=/home/tiennv14/devenv/gopath/src/github.com/proci/network/crypto-config/ordererOrganizations/orderer.block")

		log.Println(cmd.Args)

		stdoutStderr, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatal(stdoutStderr)
		}
	}
}