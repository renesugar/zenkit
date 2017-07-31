package databus_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"errors"
	"github.com/datamountaineer/schema-registry"
	"testing"
)

func TestDatabus(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Databus Suite")
}

func GetSchemaRegistryMockClient(schemas map[string]string, ids map[string]int) schemaregistry.Client {
	return &schemaregistry.MockClient{
		GetLatestSchemaFn: func(subject string) (schemaregistry.Schema, error) {
			var empty schemaregistry.Schema
			schema, ok := schemas[subject]
			if !ok {
				return empty, errors.New("Nope")
			}
			return schemaregistry.Schema{
				Id:      ids[subject],
				Schema:  schema,
				Subject: subject,
				Version: 1,
			}, nil
		},
	}
}
