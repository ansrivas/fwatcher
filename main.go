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
	"github.com/fsnotify/fsnotify"
	conf "github.com/fwatcher/internal"
	"github.com/fwatcher/messages"
	"github.com/fwatcher/workers"
	flag "github.com/spf13/pflag"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "path to a configuration file")
}

func watchDirectory(ctx context.Context, dirToWatch string, pid *actor.PID) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(dirToWatch)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event := <-watcher.Events:

			log.Println(event.Op)
			if event.Op == fsnotify.Create {
				log.Println("File created...")
				pid.Tell(&messages.FileModified{Filepath: event.Name})
			}

		case err := <-watcher.Errors:
			log.Println("error:", err)

		}
	}
}
func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	backoffWindow := time.Duration(time.Second * 10)
	initialBackoff := time.Duration(time.Second * 3)

	supervisor := actor.NewExponentialBackoffStrategy(backoffWindow, initialBackoff)

	props := actor.
		FromProducer(workers.NewCoordinatorActor).
		WithSupervisor(supervisor)

	pid := actor.Spawn(props)

	go watchDirectory(ctx, dirToWatch, pid)

	<-sigchan
	cancel()
	fmt.Printf("Terminating the program successfully\n")

}
