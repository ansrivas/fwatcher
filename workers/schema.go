package workers

const avroSchema = `{
  "type" : "record",
  "name" : "test_schema",
  "namespace" : "com.srivastava.avro",
  "fields" : [ {
    "name" : "timestamp",
    "type" : "int",
    "doc" : "Datetime format with timezone information"
  }, {
    "name" : "datapoint",
    "type" : "string",
    "doc" : "The name of the datapoint"
  }, {
    "name" : "value",
    "type" : "float",
    "doc" : "Value of the datapoint"
  } ],
  "doc:" : "A basic schema for storing iot messages"
}`
