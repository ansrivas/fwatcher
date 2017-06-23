package workers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TimeParsing(t *testing.T) {

	assert := assert.New(t)
	actual := "2017-06-23T11:58:38"
	expected := int64(1498211918)

	utc, _ := convertToUTC(actual)
	assert.Equal(utc, expected)

	actual = "2017-06-23T11:58:99"
	_, err := convertToUTC(actual)

	assert.NotNil(err, "Should throw an error in case of illegal dates")
}

func Test_ToFloat(t *testing.T) {

	assert := assert.New(t)
	actual := "12"
	expected := float64(12)

	fvalue, _ := toFloat(actual)
	assert.Equal(fvalue, expected)

	actual = "a2"

	_, err := toFloat(actual)
	assert.NotNil(err, "Should throw error while parsing a2 as float")
}
