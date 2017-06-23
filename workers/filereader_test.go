package workers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const testSchema = `{
  "type" : "record",
  "name" : "test_schema",
  "namespace" : "com.srivastava.avro",
  "fields" : [   {
    "name" : "value",
    "type" : "float",
    "doc" : "Value of the datapoint"
  } ],
  "doc:" : "A basic schema for storing iot messages"
}`

const testWrongSchema = `{
  "type" : "record",
  "name" : "test_schema",
  "namespace" : "com.srivastava.avro",
  "fields" : [
    "name" : "value",
    "type" : "float",
    "doc" : "Value of the datapoint"
  } ],
  "doc:" : "A basic schema for storing iot messages"
}`

func Test_AvroEncoder(t *testing.T) {
	assert := assert.New(t)

	codec, err := initAvroDecoder(testSchema)
	assert.NotNil(codec, "Decoder created successfully")
	assert.Nil(err, "Decoder shouldnt panic on successful parsing")

	_, err = initAvroDecoder(testWrongSchema)
	assert.NotNil(err, "Avro encoder/Decoder should throw an error on parsing wrong schema")
}
