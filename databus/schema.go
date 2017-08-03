package databus

import (
	schemaregistry "github.com/datamountaineer/schema-registry"
	"github.com/karrick/goavro"
	"github.com/pkg/errors"
)

func GetCodec(client schemaregistry.Client, subject string) (int, *goavro.Codec, error) {
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
