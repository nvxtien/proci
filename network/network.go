package network

import (
	"strings"
	"os"
	"bufio"
	"io/ioutil"
	"fmt"
	"log"
	"os/exec"
)

// create configtx.yaml.
func GenerateConfigTx(MSPBaseDir, comName string) (bool, error) {
	log.Println( " - generate configtx.yaml ...")

	if _, err := os.Stat("./configtx.yaml"); os.IsExist(err) {
		err := os.Remove("./configtx.yaml")
		if err != nil {
			log.Fatalf("Can not delete configtx.yaml")
		}
	}

	template, err := os.Open("./configtx.yaml-in")
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

	numberOfOrg := 2
	ordererType := "kafka"
	numberOfKafka := 1
	profile := "test"

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

	err = ioutil.WriteFile("configtx.yaml", []byte(strings.Join(configtx, "")), 0644)
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
	configtx = append(configtx, fmt.Sprintf("            - Host: peer0.org%d.example.com\n", i))
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

func GenerateCryptoCfg() {
	log.Println("************* generate crypto-config.yaml *************")

	if _, err := os.Stat("./crypto-config.yaml"); os.IsExist(err) {
		err := os.Remove("./crypto-config.yaml")
		if err != nil {
			log.Fatalf("Can not delete crypto-config.yaml")
		}
	}

	nOrg := 2
	nOrderer := 2
	peersPerOrg := 2
	comName := "example.com"
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

	err := ioutil.WriteFile("crypto-config.yaml", []byte(strings.Join(cryptocfg, "")), 0644)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func ExecuteCryptogen() {
	log.Println("************* execute cryptogen *************")

	if _, err := os.Stat("./crypto-config"); os.IsExist(err) {
		err := os.Remove("crypto-config")
		if err != nil {
			log.Fatalf("Can not delete crypto-config")
		}
	}

	path, err := exec.LookPath("cryptogen")
	if err != nil {
		log.Fatal("Please install cryptogen")
	}
	log.Printf("cryptogen is available at %s\n", path)

	cmd := exec.Command("cryptogen", "generate", "--output=./crypto-config", "--config=./crypto-config.yaml")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", stdoutStderr)
}

func CreateOrderGenesisBlock() {
	log.Println(" - create orderer genesis block ...")

	ordererDir := "./crypto-config/ordererOrganizations"
	ordBlock := fmt.Sprintf("%s/orderer.block", ordererDir)
	profile := "test"
	//testChannel := "channel"

	//CHAN_PROFILE := fmt.Sprintf("%sChannel", profile)
	ORDERER_PROFILE := fmt.Sprintf("%sOrgsOrdererGenesis", profile)
	//ORG_PROFILE := fmt.Sprintf("%sorgschannel", profile)

	path, err := exec.LookPath("configtxgen")
	if err != nil {
		log.Fatal("Please install configtxgen")
	}

	log.Printf("configtxgen is available at %s\n", path)
	// configtxgen -profile "testOrgsOrdererGenesis" -channelID "channel" -outputBlock "./crypto-config/ordererOrganizations/orderer.block"

	// configtxgen -profile "testOrgsOrdererGenesis" -channelID "channel" -outputBlock "./crypto-config/ordererOrganizations/orderer.block"

	cmd := exec.Command("configtxgen",
		fmt.Sprintf("-profile=%s", ORDERER_PROFILE),
		fmt.Sprintf("-channelID=%s", "channel"),
		fmt.Sprintf("-outputBlock=%s", ordBlock))

	//cmd := exec.Command("configtxgen", "-profile=testOrgsOrdererGenesis", "-channelID=channel", "-outputBlock=/home/tiennv14/devenv/gopath/src/github.com/proci/network/crypto-config/ordererOrganizations/orderer.block")

	log.Println(cmd.Args)


	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", stdoutStderr)
}
