package databus

import (
	schemaregistry "github.com/datamountaineer/schema-registry"
	"github.com/linkedin/goavro"
)

func GetCodec(client schemaregistry.Client, subject string) (int, goavro.Codec, error) {
	schema, err := client.GetLatestSchema(subject)
	if err != nil {
		return 0, nil, err
	}
	codec, err := goavro.NewCodec(schema.Schema)
	if err != nil {
		return 0, nil, err
	}
	return schema.Id, codec, nil
}
