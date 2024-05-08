package encdec

import (
	"errors"
	"github.com/goccy/go-yaml"
	"io"
)

type YAMLEncoder struct{}

func (e *YAMLEncoder) Encode(data any) ([]byte, error) {
	return yaml.Marshal(data)
}

func NewYAMLEncoder() *YAMLEncoder {
	return &YAMLEncoder{}
}

type YAMLDecoder struct{}

func (d *YAMLDecoder) Decode(data io.Reader, dst any) error {
	dec := yaml.NewDecoder(data, yaml.Strict())
	err := dec.Decode(dst)

	if errors.Is(err, io.EOF) {
		return errors.New("data must not be empty")
	}
	return err
}

func NewYAMLDecoder() *YAMLDecoder {
	return &YAMLDecoder{}
}
