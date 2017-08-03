package databus_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"

	"errors"
	"testing"

	"github.com/datamountaineer/schema-registry"
)

func TestDatabus(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Databus Suite", []Reporter{junitReporter})
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
