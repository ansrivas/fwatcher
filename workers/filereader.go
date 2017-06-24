package workers

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ansrivas/fwatcher/messages"
	"github.com/karrick/goavro"
)

type fileReadActor struct {
	kproducer Producer
	avroCodec *goavro.Codec
}

// publishFileToKafka publishes a file to Kafka in a avro-serialized manner.
// Currently the schema is hardcoded but can be easily extend to be passed from outside.
func (state *fileReadActor) publishFileToKafka(msg *messages.ReadFile, context actor.Context) {

	data, err := readFile(msg.Filename)
	if err != nil {
		log.Println(err.Error())
		return
	}

	mymap := make(map[string]interface{})
	for num, line := range strings.Split(string(data), "\n") {
		row := strings.Split(line, ";")
		if len(row) != 3 {
			log.Printf("Error: Unreadable row from file %s in  lineno %d, line %s", msg.Filename, num+1, line)
			continue
		}

		utcValue, err := convertToUTC(row[0])
		if err != nil {
			log.Printf("Error parsing the timestamp: file %s in  lineno %d, line %s", msg.Filename, num+1, line)
			continue
		}
		datapointValue, err := toFloat(row[2])
		if err != nil {
			log.Printf("Illegal value format file %s in  lineno %d, line %s", msg.Filename, num+1, line)
			continue
		}

		mymap["timestamp"] = utcValue
		mymap["datapoint"] = row[1]
		mymap["value"] = datapointValue

		// textual, err := state.avroCodec.TextualFromNative(nil, mymap)
		textual, err := state.avroCodec.BinaryFromNative(nil, mymap)
		if err != nil {
			fmt.Println("Conversion error", err)
			continue
		}

		// Figure out topic and partition from this.
		state.kproducer.Produce(textual, msg.Topic)
	}

	context.Parent().Tell(&messages.PublishAck{Filename: msg.Filename})

}
func initAvroDecoder(schema string) (*goavro.Codec, error) {
	return goavro.NewCodec(schema)
}

// CreateFileReaderProps create and spawn and child here
// brokers is a list of kafka brokers
func CreateFileReaderProps(context actor.Context, brokers []string) *actor.PID {

	avroEncoder, err := initAvroDecoder(avroSchema)
	if err != nil {
		log.Fatalln("Unable to create an avro decoder. Check your schema file.")
	}

	fileActor := &fileReadActor{
		kproducer: NewProducer(brokers),
		avroCodec: avroEncoder,
	}

	go func() {
		for err := range fileActor.kproducer.kafkaProducer.Errors() {
			log.Println("Failed to write access log entry:", err)
		}
	}()

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
