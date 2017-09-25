package databus_test

import (
	"context"

	schemaregistry "github.com/datamountaineer/schema-registry"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/logging/logrus"
	"github.com/pkg/errors"
	"github.com/zenoss/zenkit/test"
	. "github.com/zenoss/zenkit/databus"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type mockFactory struct {
	registry SchemaRegistry
}

func (mock *mockFactory) NewSchemaRegistry(registryURI string) SchemaRegistry {
	return mock.registry
}

var _ = Describe("Registry", func() {

	var (
		ctx context.Context
		factory SchemaRegistryFactory
		client *schemaregistry.MockClient
		registry SchemaRegistry
		newClientFunc NewSchemaRegistryClientFunc
		logger   = test.TestLogger()

		testSubject, testSchema string
		isRegisteredError, registrationError error

		isRegisteredMock = func(subject, schema string) (bool, schemaregistry.Schema, error) {
			unused := schemaregistry.Schema{}
			if subject == testSubject && schema == testSchema {
				return false, unused, isRegisteredError
			}
			return true, unused, nil
		}

		registerNewSchemaMock = func(subject, schema string) (int, error) {
			Ω(subject).Should(Equal(testSubject))
			Ω(schema).Should(Equal(testSchema))
			return 0, registrationError
		}
	)

	BeforeEach(func() {
		ctx = goa.WithLogger(context.Background(), goalogrus.New(logger))
	})

	Context("when a schema registry factory is added to a context", func() {
		JustBeforeEach(func() {
			factory = DefaultSchemaRegistryFactory(ctx)
			ctx = WithSchemaRegistryFactory(ctx, factory)
		})
		It("it can be retrieved from such context", func() {
			result, err := ContextSchemaRegistryFactory(ctx)
			Ω(err).Should(BeNil())
			Ω(result).Should(Equal(factory))
		})
	})

	Context("when a schema registry factory is not added to a context", func() {
		It("retrieving it from the context fails", func() {
			result, err := ContextSchemaRegistryFactory(ctx)
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(Equal(ErrFactoryNotFound))
			Ω(result).Should(BeNil())
		})
	})

	Context("when a schema registry client cannot be created", func() {
		expectedError := errors.New("new client func failed")
		JustBeforeEach(func() {
			newClientFunc =  func(baseurl string) (schemaregistry.Client, error) {
				return nil, expectedError
			}
			factory = BuildSchemaRegistryFactory(ctx, newClientFunc)
			Ω(factory).ShouldNot(BeNil())
			registry = factory.NewSchemaRegistry("registryURI")
		})
		It("schema registration should fail", func() {
			err := registry.Register("unused", "unused", "unused", "unused")
			Ω(err).ShouldNot(BeNil())
			Ω(errors.Cause(err)).Should(Equal(expectedError))
		})
	})

	Context("when a key schema is registered", func() {
		JustBeforeEach(func() {
			testSubject = "keySubject"
			testSchema = "keySchema"
			client = &schemaregistry.MockClient{
				IsRegisteredFn:      isRegisteredMock,
				RegisterNewSchemaFn: registerNewSchemaMock,
			}
			newClientFunc =  func(baseurl string) (schemaregistry.Client, error) {
				return client, nil
			}
			factory = BuildSchemaRegistryFactory(ctx, newClientFunc)
			Ω(factory).ShouldNot(BeNil())

			registry = factory.NewSchemaRegistry("registryURI")
		})

		Context("and it has not been registered before", func() {
			JustBeforeEach(func() {
				isRegisteredError = errors.New("40401")
				registrationError = nil
			})
			It("it should be added to the registry", func() {
				err := registry.Register(testSubject, testSchema, "unused", "unused")
				Ω(err).Should(BeNil())
			})
		})

		Context("and the registry check fails", func() {
			JustBeforeEach(func() {
				isRegisteredError = errors.New("check failed")
			})
			It("it should not be added to the registry", func() {
				err := registry.Register(testSubject, testSchema, "unused", "unused")
				Ω(err).ShouldNot(BeNil())
				Ω(errors.Cause(err)).Should(Equal(isRegisteredError))
			})
		})

		Context("and registration fails", func() {
			JustBeforeEach(func() {
				isRegisteredError = errors.New("40401")
				registrationError = errors.New("registration failed")
			})
			It("it should not be added to the registry", func() {
				err := registry.Register(testSubject, testSchema, "unused", "unused")
				Ω(err).ShouldNot(BeNil())
				Ω(errors.Cause(err)).Should(Equal(registrationError))
			})
		})
	})

	Context("when a value schema is registered", func() {
		JustBeforeEach(func() {
			testSubject = "valueSubject"
			testSchema = "valueSchema"
			client = &schemaregistry.MockClient{
				IsRegisteredFn:      isRegisteredMock,
				RegisterNewSchemaFn: registerNewSchemaMock,
			}
			newClientFunc =  func(baseurl string) (schemaregistry.Client, error) {
				return client, nil
			}
			factory = BuildSchemaRegistryFactory(ctx, newClientFunc)
			Ω(factory).ShouldNot(BeNil())

			registry = factory.NewSchemaRegistry("registryURI")
		})

		Context("and it has not been registered before", func() {
			JustBeforeEach(func() {
				isRegisteredError = errors.New("40401")
				registrationError = nil
			})
			It("it should be added to the registry", func() {
				err := registry.Register("unused", "unused", testSubject, testSchema)
				Ω(err).Should(BeNil())
			})
		})

		Context("and the registry check fails", func() {
			JustBeforeEach(func() {
				isRegisteredError = errors.New("check failed")
			})
			It("it should not be added to the registry", func() {
				err := registry.Register("unused", "unused", testSubject, testSchema)
				Ω(err).ShouldNot(BeNil())
				Ω(errors.Cause(err)).Should(Equal(isRegisteredError))
			})
		})

		Context("and registration fails", func() {
			JustBeforeEach(func() {
				isRegisteredError = errors.New("40401")
				registrationError = errors.New("registration failed")
			})
			It("it should not be added to the registry", func() {
				err := registry.Register("unused", "unused", testSubject, testSchema)
				Ω(err).ShouldNot(BeNil())
				Ω(errors.Cause(err)).Should(Equal(registrationError))
			})
		})
	})
})


