package network

import (
	"log"
	"os"
	"os/exec"
)

func Start(dir string) {
	log.Println("************* start network *************")
	filepath := dir + "/docker-compose.yaml"
	log.Println(filepath)
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		log.Fatalf("Not found docker-compose.yaml")
	}

	path, err := exec.LookPath("docker-compose")
	if err != nil {
		log.Fatalf("Please install docker-compose: %s", err)
	}
	log.Printf("docker-compose is available at %s\n", path)

	cmd, _ := exec.Command("docker-compose", "-f" + filepath, "up").Output()
	log.Println(cmd)
	//stdoutStderr, err := cmd.CombinedOutput()
	//if err != nil {
	//	log.Fatal(string(cmd))
	//}
}

