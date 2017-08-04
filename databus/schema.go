package databus

import (
	schemaregistry "github.com/datamountaineer/schema-registry"
	"github.com/karrick/goavro"
	"github.com/pkg/errors"
)

// SchemaDecoder decodes text encodings according to a schema into native Go
// types.
type SchemaDecoder interface {
	// NativeFromTextual converts textual (JSON) data into native Go types
	// according to a schema.
	NativeFromTextual([]byte) (interface{}, []byte, error)
}

// SchemaEncoder encodes native Go types according to a schema to a textual
// representation.
type SchemaEncoder interface {
	// TextualFromNative encodes Go native types to a JSON representation of
	// a schema.
	TextualFromNative([]byte, interface{}) ([]byte, error)
}

// SchemaHaver is an object with a schema.
type SchemaHaver interface {
	// Schema retrieves the schema for this object.
	Schema() string
}

// SchemaCodec is a schema encoder/decoder.
type SchemaCodec interface {
	SchemaHaver
	SchemaEncoder
	SchemaDecoder
}

// GetCodec retrieves the Avro schema with the subject specified from a schema
// registry. It returns a codec that can be used to decode from binary or text
// to Go native types.
func GetCodec(client schemaregistry.Client, subject string) (int, SchemaCodec, error) {
	schema, err := client.GetLatestSchema(subject)
	if err != nil {
		return 0, nil, errors.Wrapf(err, "failed to get latest schema for subject %s", subject)
	}
	codec, err := goavro.NewCodec(schema.Schema)
	if err != nil {
		return 0, nil, errors.Wrap(err, "failed to create codec")
	}
	return schema.Id, codec, nil
}
