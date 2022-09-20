package runtime

import (
	"context"
	"encoding/json"
	"io"
)

// JSONBuiltin is a Marshaler which marshals/unmarshals into/from JSON
// with the standard "encoding/json" package of Golang.
// Although it is generally faster for simple proto messages than JSONPb,
// it does not support advanced features of protobuf, e.g. map, oneof, ....
//
// The NewEncoder and NewDecoder types return *json.Encoder and
// *json.Decoder respectively.
type JSONBuiltin struct{}

// ContentType always Returns "application/json".
func (*JSONBuiltin) ContentType(_ interface{}) string {
	return "application/json"
}

// Marshal marshals "v" into JSON
func (j *JSONBuiltin) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal unmarshals JSON data into "v".
func (j *JSONBuiltin) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// NewDecoder returns a Decoder which reads JSON stream from "r".
func (j *JSONBuiltin) NewDecoder(r io.Reader) Decoder {
	return &jsonDecoder{r}
}

// NewEncoder returns an Encoder which writes JSON stream into "w".
func (j *JSONBuiltin) NewEncoder(w io.Writer) Encoder {
	return &jsonEncoder{w}
}

// Delimiter for newline encoded JSON streams.
func (j *JSONBuiltin) Delimiter() []byte {
	return []byte("\n")
}

type jsonEncoder struct {
	w io.Writer
}

func (j *jsonEncoder) Encode(_ context.Context, v interface{}) error {
	return json.NewEncoder(j.w).Encode(v)
}

type jsonDecoder struct {
	r io.Reader
}

func (j *jsonDecoder) Decode(_ context.Context, v interface{}) error {
	return json.NewDecoder(j.r).Decode(v)
}
