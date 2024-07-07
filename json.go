package gonertia

import (
	"encoding/json"
	"io"
)

// JSONMarshaller is marshaller which use for marshal/unmarshal JSON.
type JSONMarshaller interface {
	Marshal(v any) ([]byte, error)
	Decode(r io.Reader, v interface{}) error
}

type jsonDefaultMarshaller struct{}

func (j jsonDefaultMarshaller) Decode(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func (j jsonDefaultMarshaller) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
