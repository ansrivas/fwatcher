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

func (state *fileReadActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {

	case *actor.Started:
		state.p = NewProducer()
	case *messages.ReadFile:

		// Need to be sure if this is an okay practice to run a coroutine in an actor
		go func() {
			data, err := readFile(msg.Filename)
			if err != nil {
				log.Println(err.Error())
				return
			}
			state.p.Produce(data)
			context.Parent().Tell(&messages.PublishAck{})
		}()

		// context.Sender().Tell(&messages.FileContent{Content: data})
		//Testing inform self
		// context.Self().Tell(&messages.FileContent{Content: data})

		//Testing poison pill
		// context.Self().Tell(&actor.PoisonPill{})
	case *messages.PublishAck:
		fmt.Println("File has been successfully published")
	}
}

func readFile(filename string) (string, error) {
	log.Println("Now reading file: ", filename)
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("Unable to read file: %s", filename)
	}
	return string(dat), nil
}
