package messages

//FileModified message informs the actor about reading the file and sends the filename
type FileModified struct{ Filepath string }

//ReadFile message is to read a given filename with path
type ReadFile struct{ Filename string }

//PublishAck sends the ack about publishing to kafka
type PublishAck struct{ Filename string }

//Publisher publishes a given message to a kafka topic
type Publisher struct{ Content string }
