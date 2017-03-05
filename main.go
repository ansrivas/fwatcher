package main

import (
	"fmt"
	"log"

	flag "github.com/spf13/pflag"
)

var configName string

func init() {
	flag.StringVar(&configName, "config", "", "path to a configuration file")
}

func main() {
	flag.Parse()
	if configName == "" {
		log.Fatalf("No config file provided... ")
	}
	fmt.Println("config file is:", configName)
}
