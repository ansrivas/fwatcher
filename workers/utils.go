package workers

import (
	"strconv"
	"strings"

	"github.com/jinzhu/now"
)

func init() {
	// TODO: Initialize this w.r.t a config file
	now.TimeFormats = append(now.TimeFormats, "2006-01-02T15:04:05")
}

func toFloat(input string) (float32, error) {
	// var value float32
	val, err := strconv.ParseFloat(strings.TrimSpace(input), 32)
	if err != nil {
		return 0, err
	}
	return float32(val), err
}

func convertToUTC(ts string) (int64, error) {
	parsed, err := now.Parse(ts)
	if err != nil {
		return 0, err
	}
	return parsed.UTC().Unix(), nil
}
