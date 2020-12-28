package metrics

import (
	"github.com/gogf/gf/os/gtimer"
	"time"
)

// This global kill-switch helps quantify the observer effect and makes
// for less cluttered pprof profiles.
var (
	tw = gtimer.New(100, 10*time.Millisecond, 6)
)

var defaultRegistry Registry = NewRegistry()

// Call the given function for each registered metric.
func Each(f func(string, interface{})) {
	defaultRegistry.Each(f)
}

// Get the metric by the given name or nil if none is registered.
func Get(name string) interface{} {
	return defaultRegistry.Get(name)
}

// Gets an existing metric or creates and registers a new one. Threadsafe
// alternative to calling Get and Register on failure.
func GetOrRegister(name string, i interface{}) interface{} {
	return defaultRegistry.GetOrRegister(name, i)
}

// Register the given metric under the given name.  Returns a DuplicateMetric
// if a metric by the given name is already registered.
func Register(name string, i interface{}) error {
	return defaultRegistry.Register(name, i)
}

// Register the given metric under the given name.  Panics if a metric by the
// given name is already registered.
func MustRegister(name string, i interface{}) {
	if err := Register(name, i); err != nil {
		panic(err)
	}
}

// Run all registered healthchecks.
func RunHealthchecks() {
	defaultRegistry.RunHealthchecks()
}

// Unregister the metric with the given name.
func Unregister(name string) {
	defaultRegistry.Unregister(name)
}
