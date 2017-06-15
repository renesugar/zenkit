package zenkit

// key is the type used to store internal values in the context.
type key int

const (
	metricsKey key = iota + 1
	identityKey
)
