package workers

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ansrivas/fwatcher/messages"
	"github.com/linkedin/goavro"
)

type fileReadActor struct {
	kproducer Producer
	avroCodec *goavro.Codec
}

// publishFileToKafka publishes a file to Kafka in a avro-serialized manner.
// Currently the schema is hardcoded but can be easily extend to be passed from outside.
func (state *fileReadActor) publishFileToKafka(msg *messages.ReadFile, context actor.Context) {
	log.Printf("Now reading file: %s", msg.Filename)
	fileHandle, err := os.Open(msg.Filename)
	defer fileHandle.Close()
	if err != nil {
		log.Printf("Unable to open the file %s", msg.Filename)
		return
	}
	fileReader := bufio.NewScanner(fileHandle)

	mymap := make(map[string]interface{})
	var num = 0
	for fileReader.Scan() {
		line := fileReader.Text()
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
		state.kproducer.Produce(textual)
		num++
	}

	context.Parent().Tell(&messages.PublishAck{Filename: msg.Filename})

}
func initAvroDecoder(schema string) (*goavro.Codec, error) {
	return goavro.NewCodec(schema)
}

// CreateFileReaderProps create and spawn and child here
func CreateFileReaderProps(context actor.Context, bootstrapServers string) *actor.PID {

	avroEncoder, err := initAvroDecoder(avroSchema)
	if err != nil {
		log.Fatalln("Unable to create an avro decoder. Check your schema file.")
	}

	fileActor := &fileReadActor{
		kproducer: NewProducer(bootstrapServers),
		avroCodec: avroEncoder,
	}

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
