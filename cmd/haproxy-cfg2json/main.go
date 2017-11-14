package main

import (
	"encoding/json"
	"fmt"
	"github.com/tkmgo/haproxyconfigparser"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("You need to specify the haproxy.cfg file.")
	}
	filename := os.Args[1]

	services, err := haproxyconfigparser.ParseFromFile(filename)
	if err != nil {
		log.Fatalf("Failed to parse the config: %s", err)
	}

	bytes, err := json.MarshalIndent(services, "", "  ")
	if err != nil {
		log.Fatalf("Failed to translate to JSON")
	}

	fmt.Println(string(bytes))
}
