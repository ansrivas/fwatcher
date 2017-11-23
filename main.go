package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"syscall"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ansrivas/fwatcher/db"
	conf "github.com/ansrivas/fwatcher/internal"

	"github.com/ansrivas/fwatcher/workers"

	flag "github.com/spf13/pflag"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "path to a configuration file")
}

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if configPath == "" {
		log.Fatalf("No config file provided... ")
	}
	config, err := conf.GetConfig(configPath)
	checkAndExit(err)

	hosts := config.GetString("kafka.hosts")
	dirToWatch := config.GetString("app.dir")
	allowedExtensions := config.GetStringSlice("app.filetypes")
	log.Println(hosts, dirToWatch)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	backoffWindow := time.Duration(time.Second * 10)
	initialBackoff := time.Duration(time.Second * 3)

	supervisor := actor.NewExponentialBackoffStrategy(backoffWindow, initialBackoff)

	props := actor.
		FromInstance(&workers.CoordinatorActor{BootStrapServers: hosts}).
		WithSupervisor(supervisor)

	pid := actor.Spawn(props)

	db.NewDb()

	go watchDirectory(ctx, dirToWatch, allowedExtensions, pid)

	<-sigchan
	cancel()
	pid.Stop()
	fmt.Printf("Terminating the program successfully\n")

}
