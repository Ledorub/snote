package encdec

import (
	gojson "encoding/json"
	"io"
)

type JSONEncoder struct{}

func (enc *JSONEncoder) Encode(data map[string]any) ([]byte, error) {
	return gojson.Marshal(data)
}

func NewJSONEncoder() *JSONEncoder {
	return &JSONEncoder{}
}

type JSONDecoder struct{}

func (enc *JSONDecoder) Decode(data io.Reader, dst any) error {
	return nil
}

func NewJSONDecoder() *JSONDecoder {
	return &JSONDecoder{}
}
