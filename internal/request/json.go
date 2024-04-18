package request

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type jsonDecoder interface {
	Decode(data io.Reader, dst any) error
}

type JSONRequestReader struct {
	logger  *log.Logger
	decoder jsonDecoder
}

func (reader *JSONRequestReader) Read(data io.Reader, dst any) error {
	err := reader.decoder.Decode(data, dst)
	var maxBytesError *http.MaxBytesError
	if errors.As(err, &maxBytesError) {
		err = fmt.Errorf("body is too large")
	}
	return err
}

func NewJSONReader(logger *log.Logger, decoder jsonDecoder) *JSONRequestReader {
	return &JSONRequestReader{logger: logger, decoder: decoder}
}
