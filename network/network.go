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

	configtx := ""
	for i, line := range lines {
		fmt.Printf("\n%d %s", i, line)
		if strings.Contains(line, "&OrdererOrg") {
			configtx += line + "\n"
			configtx += fmt.Sprintf("        Name: OrdererOrg\n")
			configtx += fmt.Sprintf("        ID: OrdererOrg\n")
			configtx += fmt.Sprintf("        MSPDir: %s/ordererOrganizations/%s/msp\n", MSPBaseDir, comName)
			ordMSP := "OrdererOrg"
			configtx += fmt.Sprintf("        Policies:\n")
			configtx += fmt.Sprintf("            Readers:\n")
			configtx += fmt.Sprintf("                Type: Signature\n")
			configtx += fmt.Sprintf("                Rule: \"OR('%s.member')\"\n", ordMSP)
			configtx += fmt.Sprintf("            Writers:\n")
			configtx += fmt.Sprintf("                Type: Signature\n")
			configtx += fmt.Sprintf("                Rule: \"OR('%s.member')\"\n", ordMSP)
			configtx += fmt.Sprintf("            Admins:\n")
			configtx += fmt.Sprintf("                Type: Signature\n")
			configtx += fmt.Sprintf("                Rule: \"OR('%s.admin')\"\n", ordMSP)

			continue
		}
		configtx += line
		configtx += "\n"
	}

	err = ioutil.WriteFile("configtx.yaml", []byte(configtx), 0644)
	if err != nil {
		return false, err
	}

	return true, nil
}
