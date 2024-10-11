package gonertia

import (
	"encoding/json"
	"io"
)

// JSONMarshaller is marshaller which use for marshal/unmarshal JSON.
type JSONMarshaller interface {
	Marshal(v any) ([]byte, error)
	Decode(r io.Reader, v any) error
}

type jsonDefaultMarshaller struct{}

func (j jsonDefaultMarshaller) Decode(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}

func (j jsonDefaultMarshaller) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}
