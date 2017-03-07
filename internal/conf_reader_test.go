package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	assert := assert.New(t)
	conf, err := GetConfig("config.yaml")
	expected := "world"
	assert.Nil(err, "File should be properly read")
	assert.Equal(expected, conf.GetString("field.hello"), "expected to read world from the config")

	conf, err = GetConfig(".")
	assert.NotNil(err, "Error should be not nil")
}
