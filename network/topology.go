package network

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type component struct {
	image			string `json:"image""`
	environment		[]string `json:"environment,omitempty"`
	ports			[]string `json:"ports,omitempty"`
	command			string `json:"command,omitempty"`
	volumes			[]string `json:"volumes,omitempty"`
	depends_on 		string `json:"depends_on,omitempty"`
	container_name	string `json:"container_name,omitempty"`
	working_dir		string `json:"working_dir,omitempty"`
}

type ca struct {
	component
}

func CreateDockerCompose(filename string) {

	def, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf(err.Error())
	}

	var dat map[string]interface{}

	if err := json.Unmarshal(def, &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat)

	//var caPort, ordererPort, vp0Port, kafkaPort, couchdbPort, evtPort int
	//var caAddress, ordererAddress, vp0Address, kafkaAddress, couchdbAddress, evtAddress string

	/*for k, v := range dat {
		fmt.Printf("%s %s\n", k, v)
	}*/

	services := dat["services"].(map[string]interface{})
	fmt.Println(services)

	ca := services["ca"].(map[string]interface{})
	fmt.Printf("%s \n", ca)

	/*var services map[string]interface{}
	if err := json.Unmarshal(strs, &services); err != nil {
		panic(err)
	}
	fmt.Println(services)*/

	//.([]interface{})
	//str1 := strs[0].(string)
	//fmt.Println(strs)
}