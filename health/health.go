package health

import (
	"golang.org/x/sync/syncmap"
)

// Status describes the health state
type Status string

const (
	OK       Status = "OK"
	DEGRADED Status = "DEGRADED"
	CRITICAL Status = "CRITICAL"
)

var r *Registry

func init() {
	r = New()
}

// Result describes the state of the health check
type Result struct {
	Name   string
	Status Status
	Err    error
}

// StatusChecker returns the status of the health check
type StatusChecker interface {
	CheckStatus() (Status, error)
}

// Registry is a health check registry that maintains the set of health checks.
type Registry struct {
	statusCheckers *syncmap.Map
}

// New returns an initialized Registry instance.
func New() *Registry {
	return &Registry{
		statusCheckers: &syncmap.Map{},
	}
}

// Intended for testing, will reset all to default settings.
// In the public interface for the health package so applications can use it in
// their testing as well.
func Reset() {
	r = New()
}

// Register adds a new health check to the registry
func Register(name string, statusChecker StatusChecker) { r.Register(name, statusChecker) }
func (r *Registry) Register(name string, statusChecker StatusChecker) {
	r.statusCheckers.Store(name, statusChecker)
}

// Execute runs all health checks in parallel
func Execute() []*Result { return r.Execute() }
func (r *Registry) Execute() []*Result {
	ch := make(chan *Result)
	count := 0

	r.statusCheckers.Range(func(key, value interface{}) bool {
		name, statusChecker := key.(string), value.(StatusChecker)
		count++
		go func() {
			stat, err := statusChecker.CheckStatus()
			ch <- &Result{Name: name, Status: stat, Err: err}
		}()
		return true
	})

	results := make([]*Result, count)
	for i := 0; i < count; i++ {
		results[i] = <-ch
	}
	return results
}
