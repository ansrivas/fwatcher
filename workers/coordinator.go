package workers

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/fwatcher/messages"
)

type coordinatorActor struct {
	fileReader *actor.PID
}

func (parent *coordinatorActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		fmt.Println("Parent actor started now ...")
		parent.fileReader = CreateFileReaderProps(context)
	case *messages.FileModified:
		parent.fileReader.Tell(&messages.ReadFile{Filename: msg.Filepath})
		// context.Request(parent.fileReader, &messages.ReadFile{Filename: msg.Filepath})
	case *messages.PublishAck:
		fmt.Println("File has been successfully published to kafka")
	}
}

//NewCoordinatorActor instantiates a new coordinator actor which is responsible for initializing all worker actors
func NewCoordinatorActor() actor.Actor {
	return &coordinatorActor{}
}
