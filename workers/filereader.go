package workers

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/fwatcher/messages"
)

type fileReadActor struct {
	kproducer Producer
}

//CreateFileReaderProps create and spawn and child here
func CreateFileReaderProps(context actor.Context, bootstrapServers string) *actor.PID {

	fileActor := &fileReadActor{kproducer: NewProducer(bootstrapServers)}

	fileReadActorProps := actor.FromInstance(fileActor)
	return context.Spawn(fileReadActorProps)
}

func (state *fileReadActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {

	case *actor.Started:
	case *messages.ReadFile:

		// Need to be sure if this is an okay practice to run a coroutine in an actor
		go func() {
			data, err := readFile(msg.Filename)
			if err != nil {
				log.Println(err.Error())
				return
			}
			state.kproducer.Produce(data)
			context.Parent().Tell(&messages.PublishAck{Filename: msg.Filename})
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
