package zenkit

// key is the type used to store internal values in the context.
type key int

const (
	serviceKey key = iota + 1
	metricsKey
	identityKey
)
