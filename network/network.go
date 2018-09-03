package network

import (
	"strings"
	"os"
	"bufio"
	"io/ioutil"
	"fmt"
)

// create configtx.yaml.
func GenerateConfigTx(MSPBaseDir, comName string) (bool, error) {

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
	for i, line := range lines {
		fmt.Printf("\n%d %s", i, line)

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
				configtx = append(configtx, line)
				configtx = append(configtx, "\n")
				configtx = append(configtx, fmt.Sprintf("    %sOrgsOrdererGenesis:", profile))
				configtx = append(configtx, fmt.Sprintf("        <<: *ChannelDefaults"))

			} else if strings.Contains(line, "&ProfileOrgString") {
				configtx = append(configtx, fmt.Sprintf("    %sorgschannel", profile))

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
