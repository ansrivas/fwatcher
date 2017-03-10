package workers

import (
	"fmt"
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/fwatcher/messages"
)

// CoordinatorActor to supervise other actors
type CoordinatorActor struct {
	fileReader       *actor.PID
	BootStrapServers string
}

// Receive on a coordinator actor
func (parent *CoordinatorActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		fmt.Println("Parent actor started now ...")
		parent.fileReader = CreateFileReaderProps(context, parent.BootStrapServers)
	case *messages.FileModified:
		parent.fileReader.Tell(&messages.ReadFile{Filename: msg.Filepath})
		// context.Request(parent.fileReader, &messages.ReadFile{Filename: msg.Filepath})
	case *messages.PublishAck:
		log.Printf("File %v has been successfully published to kafka", msg.Filename)
	}
}
