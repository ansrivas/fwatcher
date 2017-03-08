package workers

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/fwatcher/messages"
)

type fileReadActor struct {
	p Producer
}

func newfileReadActor() actor.Actor {
	return &fileReadActor{}
}

//CreateFileReaderProps create and spawn and child here
func CreateFileReaderProps(context actor.Context) *actor.PID {
	fileReadActorProps := actor.FromProducer(newfileReadActor)
	return context.Spawn(fileReadActorProps)
}

//CreateFileReaderPropsGlobal creates an actor in current actor system, not with a given context
func CreateFileReaderPropsGlobal() *actor.PID {
	fileReadActorProps := actor.FromProducer(newfileReadActor)
	return actor.Spawn(fileReadActorProps)
}

func (state *fileReadActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {

	case *actor.Started:
		state.p = NewProducer()
	case *messages.ReadFile:

		data := readFile(msg.Filename)
		context.Parent().Tell(&messages.FileContent{Content: data})
		state.p.Produce(data)

		// context.Sender().Tell(&messages.FileContent{Content: data})
		//Testing inform self
		// context.Self().Tell(&messages.FileContent{Content: data})

		//Testing poison pill
		// context.Self().Tell(&actor.PoisonPill{})
	case *messages.FileContent:
		fmt.Println("File content in child:", msg.Content)
	}
}

func readFile(filename string) string {
	log.Println("Now reading file: ", filename)
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("Unable to read file: %s", filename))
	}
	return string(dat)
}
