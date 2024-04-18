package encdec

import (
	gojson "encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type JSONEncoder struct{}

func (enc *JSONEncoder) Encode(data any) ([]byte, error) {
	return gojson.Marshal(data)
}

func NewJSONEncoder() *JSONEncoder {
	return &JSONEncoder{}
}

type JSONDecoder struct{}

func (enc *JSONDecoder) Decode(data io.Reader, dst any) error {
	dec := gojson.NewDecoder(data)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var (
			syntaxError            *gojson.SyntaxError
			unmarshalTypeError     *gojson.UnmarshalTypeError
			invalidUnmarshallError *gojson.InvalidUnmarshalError
		)

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("data contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("data contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			// https://github.com/golang/go/issues/25956
			return errors.New("data contains badly-formed JSON")
		case errors.Is(err, io.EOF):
			return errors.New("data must not be empty")
		case errors.As(err, &invalidUnmarshallError):
			panic(err)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("json contains unknown key %s", fieldName)
		default:
			return err
		}
	}

	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}

func NewJSONDecoder() *JSONDecoder {
	return &JSONDecoder{}
}
