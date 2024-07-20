package json

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
	"io"
)

func init() {
	extra.RegisterFuzzyDecoders()
}

var JSON = jsoniter.ConfigCompatibleWithStandardLibrary

func Marshal(v interface{}) ([]byte, error) {
	return JSON.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return JSON.Unmarshal(data, v)
}

func SafeMarshal(v interface{}) []byte {
	b, err := JSON.Marshal(v)
	if err != nil {
		return []byte("{}")
	}
	return b
}

func SafeMarshalString(v interface{}) string {
	return string(SafeMarshal(v))
}

func NewDecoder(reader io.Reader) *jsoniter.Decoder {
	return JSON.NewDecoder(reader)
}
