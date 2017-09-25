package databus

import (
	"context"
	"strings"

	schemaregistry "github.com/datamountaineer/schema-registry"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/zenoss/zenkit"
)

// key is the type used to store internal values in the context.
type key int

const (
	schemaRegistryFactoryKey key = iota + 1
)

var (
	// ErrFactoryNotFound occurs when the schema registry factory is not on the context
	ErrFactoryNotFound = errors.New("schema registry factory not found on context")
)

// WithSchemaRegistryFactory returns the context with a schema registry factory attached
// Example:
// service.Context = databus.WithSchemaRegistryFactory(service.Context, databus.DefaultSchemaRegistryFactory(service.Context))
func WithSchemaRegistryFactory(ctx context.Context, f SchemaRegistryFactory) context.Context {
	return context.WithValue(ctx, schemaRegistryFactoryKey, f)
}

// ContextSchemaRegistryFactory returns the schema registry factory attached to the context
func ContextSchemaRegistryFactory(ctx context.Context) (SchemaRegistryFactory, error) {
	if v := ctx.Value(schemaRegistryFactoryKey); v != nil {
		tf := v.(SchemaRegistryFactory)
		return tf, nil
	}
	return nil, ErrFactoryNotFound
}

// SchemaRegistry can be used to register kafka key and value schema's.
type SchemaRegistry interface {
	Register(KeySubject string, KeySchema string, ValueSubject string, ValueSchema string) error
}

// SchemaRegistryFactory can be used to generate new instances of SchemaRegistry
type SchemaRegistryFactory interface {
	NewSchemaRegistry(registryURI string) SchemaRegistry
}

type defaultRegistry struct {
	logger        *logrus.Logger
	registryURI   string
	newClientFunc NewSchemaRegistryClientFunc
}

type defaultRegistryFactory struct {
	newClientFunc NewSchemaRegistryClientFunc
	logger        *logrus.Logger
}

// DefaultSchemaRegistryFactory returns a default implementation of the SchemaRegistryFactory,
// which uses the Logger from the specified context.
func DefaultSchemaRegistryFactory(ctx context.Context) SchemaRegistryFactory {
	return BuildSchemaRegistryFactory(ctx, schemaregistry.NewClient)
}

// NewSchemaRegistryClient is a function that can be used to generate new clients for the schema registry
// Note this type primarily exists for unit-testing within this package.
type NewSchemaRegistryClientFunc func(baseurl string) (schemaregistry.Client, error)

// BuildSchemaRegistryFactory returns a default implementation of SchemaRegistryFactory
// with a caller-specified NewSchemaRegistryClientFunc.
// This method is primarily intended for unit-testing this package. Users of this package are encouraged
// to use the simpler DefaultSchemaRegistryFactory() method instead of this method.
func BuildSchemaRegistryFactory(ctx context.Context, newClientFunc NewSchemaRegistryClientFunc) SchemaRegistryFactory {
	return &defaultRegistryFactory{
		newClientFunc: newClientFunc,
		logger:        zenkit.ContextLogger(ctx).Logger,
	}
}

// NewSchemaRegistry returns a default implementation of SchemaRegistry
// which uses the Logger from defaultRegistryFactory
func (f *defaultRegistryFactory) NewSchemaRegistry(registryURI string) SchemaRegistry {
	return &defaultRegistry{
		logger:        f.logger,
		newClientFunc: f.newClientFunc,
		registryURI:   registryURI,
	}
}

func (r *defaultRegistry) Register(KeySubject string, KeySchema string, ValueSubject string, ValueSchema string) error {
	client, err := r.newClientFunc(r.registryURI)
	if err != nil {
		return errors.Wrap(err, "failed to create schema registry client")
	}

	err = r.register(client, KeySubject, KeySchema)
	if err != nil {
		return errors.Wrap(err, "failed to register key")
	}

	err = r.register(client, ValueSubject, ValueSchema)
	if err != nil {
		return errors.Wrap(err, "failed to register value")
	}

	return nil
}

func (r *defaultRegistry) register(client schemaregistry.Client, subject, schema string) error {
	logger := r.logger.WithFields(logrus.Fields{
		"registry": r.registryURI,
		"subject": subject,
		"schema": schema,
	})

	registered, _, err := client.IsRegistered(subject, schema)
	if err != nil {
		// If a schema for the subject has never been registered before, the IsRegistered method
		// will return an error.  In this case, we want to register the schema.  We don't have access
		// to the underlying error type to get the error code for this situation (which is 40401) so we
		// will check the message for the error code.
		if !strings.Contains(err.Error(), "40401") {
			logger.WithError(err).Error("Error occurred checking schema registration")
			return err
		}
	}

	if !registered {
		_, err = client.RegisterNewSchema(subject, schema)
		if err != nil {
			logger.WithError(err).Error("Error occurred registering new schema")
			return err
		}
		logger.Info("New schema registered")
	} else {
		logger.Debug("Schema already registered")
	}

	return nil
}
