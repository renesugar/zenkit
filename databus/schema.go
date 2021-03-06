package databus

import (
	schemaregistry "github.com/datamountaineer/schema-registry"
	"github.com/linkedin/goavro"
	"github.com/pkg/errors"
)

// SchemaDecoder decodes text encodings according to a schema into native Go
// types.
type SchemaDecoder interface {
	// NativeFromBinary converts binary Avro data into native Go types
	// according to a schema.
	NativeFromBinary([]byte) (interface{}, []byte, error)
}

// SchemaEncoder encodes native Go types according to a schema to a textual
// representation.
type SchemaEncoder interface {
	// BinaryFromNative encodes Go native types to a binary-encoded
	// representation of a schema.
	BinaryFromNative([]byte, interface{}) ([]byte, error)
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
