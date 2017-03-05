package main

import (
	"fmt"

	flag "github.com/spf13/pflag"
)

var configName string

func init() {
	flag.StringVar(&configName, "config", "config.yaml", "path to a configuration file")
}

func main() {
	flag.Parse()
	fmt.Println("config file is:", configName)
}
