package main

import (
	"fmt"
	"log"

	"github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	flag "github.com/spf13/pflag"
)

var configName string

func init() {
	flag.StringVar(&configName, "config", "", "path to a configuration file")
}

type hello struct{ Who string }
type helloActor struct{}

func (state *helloActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *hello:
		fmt.Printf("Hello %v\n", msg.Who)
	}
}

func main() {
	flag.Parse()
	if configName == "" {
		log.Fatalf("No config file provided... ")
	}
	fmt.Println("config file is:", configName)
	props := actor.FromInstance(&helloActor{})
	pid := actor.Spawn(props)
	pid.Tell(&hello{Who: "Roger"})
	console.ReadLine()
}
