package tools

import (
	"github.com/mitchellh/mapstructure"
	"strconv"
	"strings"
)

func IsStringInSlice(s string, slc []string) bool {
	for _, v := range slc {
		if v == s {
			return true
		}
	}

	return false
}

func MapToStruct(input interface{}, output interface{}) error {
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   output,
		TagName:  "json",
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}

func SliceOfIntsToString(a []int, sep string) string {
	if len(a) == 0 {
		return ""
	}

	b := make([]string, len(a))
	for i, v := range a {
		b[i] = strconv.Itoa(v)
	}
	return strings.Join(b, sep)
}
