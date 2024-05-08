package encdec

import (
	"errors"
	"github.com/goccy/go-yaml"
	"io"
)

type YAMLEncoder struct{}

func (enc *YAMLEncoder) Encode(data any) ([]byte, error) {
	return yaml.Marshal(data)
}

func NewYAMLEncoder() *YAMLEncoder {
	return &YAMLEncoder{}
}

type YAMLDecoder struct{}

func (dec *YAMLDecoder) Decode(data io.Reader, dst any) error {
	d := yaml.NewDecoder(data, yaml.Strict())
	err := d.Decode(dst)

	if errors.Is(err, io.EOF) {
		return errors.New("data must not be empty")
	}
	return err
}

func NewYAMLDecoder() *YAMLDecoder {
	return &YAMLDecoder{}
}
