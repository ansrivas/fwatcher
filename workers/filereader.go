package workers

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ansrivas/fwatcher/messages"
	"github.com/karrick/goavro"
)

type fileReadActor struct {
	kproducer Producer
}

// publishFileToKafka publishes a file to Kafka in a avro-serialized manner.
// Currently the schema is hardcoded but can be easily extend to be passed from outside.
func (state *fileReadActor) publishFileToKafka(msg *messages.ReadFile, context actor.Context) {

	data, err := readFile(msg.Filename)
	if err != nil {
		log.Println(err.Error())
		return
	}
	codec, err := goavro.NewCodec(avro_schema)
	if err != nil {
		fmt.Println("Decoding error", err)
	}
	mymap := make(map[string]interface{})
	for num, line := range strings.Split(string(data), "\n") {
		row := strings.Split(line, ";")
		if len(row) != 3 {
			log.Printf("Error: Unreadable row from file %s in  lineno %d, line %s", msg.Filename, num, line)
			continue
		}
		mymap["timestamp"] = row[0]
		mymap["datapoint"] = row[1]

		value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 32)
		if err != nil {
			log.Printf("Illegal value format %v", row[2])
			continue
		}
		mymap["value"] = value

		textual, err := codec.BinaryFromNative(nil, mymap)
		if err != nil {
			fmt.Println("Conversion error", err)
			continue
		}
		state.kproducer.Produce(textual)
	}

	context.Parent().Tell(&messages.PublishAck{Filename: msg.Filename})

}

//CreateFileReaderProps create and spawn and child here
func CreateFileReaderProps(context actor.Context, bootstrapServers string) *actor.PID {

	fileActor := &fileReadActor{kproducer: NewProducer(bootstrapServers)}
	fileReadActorProps := actor.FromInstance(fileActor)
	return context.Spawn(fileReadActorProps)
}

func (state *fileReadActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {

	case *messages.ReadFile:

		go state.publishFileToKafka(msg, context)

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
