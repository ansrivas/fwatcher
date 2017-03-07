package main

import (
	"fmt"
	"log"

	"github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	conf "github.com/fwatcher/internal"
	flag "github.com/spf13/pflag"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "path to a configuration file")
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
	if configPath == "" {
		log.Fatalf("No config file provided... ")
	}
	config, err := conf.GetConfig(configPath)
	if err != nil {
		log.Fatalf(err.Error())
	}
	hosts := config.GetString("kafka.hosts")
	dirToWatch := config.GetString("app.dir")
	log.Println(hosts, dirToWatch)
	props := actor.FromInstance(&helloActor{})
	pid := actor.Spawn(props)
	pid.Tell(&hello{Who: "Roger"})
	console.ReadLine()
}
