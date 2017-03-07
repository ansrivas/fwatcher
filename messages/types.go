package messages

//FileModified message informs the actor about reading the file and sends the filename
type FileModified struct{ Filename string }

//ReadFile message is to read a given filename with path
type ReadFile struct{ Filename string }

//FileContent sends the content of the file in the message
type FileContent struct{ Content string }

//Publisher publishes a given message to a kafka topic
type Publisher struct{ Content string }
