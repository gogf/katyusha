package metrix

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"reflect"
	"strings"
	"sync"
)

// DuplicateMetric is the error returned by Registry.Register when a metric
// already exists.  If you mean to Register that metric you must first
// Unregister the existing metric.
type DuplicateMetric string

func (err DuplicateMetric) Error() string {
	return fmt.Sprintf("duplicate metric: %s", string(err))
}

// Stoppable defines the metrics which has to be stopped.
type IStoppable interface {
	Stop()
}

// A Registry holds references to a set of metrics by name and can iterate
// over them, calling callback functions provided by the user.
//
// This is an interface so as to encourage other structs to implement
// the Registry API as appropriate.
type Registry interface {

	// Call the given function for each registered metric.
	Each(func(string, interface{}))

	// Get the metric by the given name or nil if none is registered.
	Get(string) interface{}

	// GetAll metrics in the Registry.
	GetAll() map[string]g.MapStrAny

	// Gets an existing metric or registers the given one.
	// The interface can be the metric to register if not found in registry,
	// or a function returning the metric for lazy instantiation.
	GetOrRegister(string, interface{}) interface{}

	// Register the given metric under the given name.
	Register(string, interface{}) error

	// Run all registered healthchecks.
	RunHealthchecks()

	// Unregister the metric with the given name.
	Unregister(string)

	// Unregister all metrics.  (Mostly for testing.)
	UnregisterAll()
}

//The standard implementation of a Registry is a mutex-protected map
//of names to metrics.
type StandardRegistry struct {
	metrics g.MapStrAny
	mutex   sync.RWMutex
}

func (s *StandardRegistry) register(name string, i interface{}) error {
	if _, ok := s.metrics[name]; ok {
		return DuplicateMetric(name)
	}
	switch i.(type) {
	case ICounter, IGauge, IGaugeFloat64, IHealthcheck, IHistogram, IMeter, ITimer:
		s.metrics[name] = i
	}
	return nil
}

type metricKV struct {
	name  string
	value interface{}
}

func (s *StandardRegistry) registered() []metricKV {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	metrics := make([]metricKV, 0, len(s.metrics))
	for name, i := range s.metrics {
		metrics = append(metrics, metricKV{
			name:  name,
			value: i,
		})
	}
	return metrics
}

func (r *StandardRegistry) stop(name string) {
	if i, ok := r.metrics[name]; ok {
		if s, ok := i.(IStoppable); ok {
			s.Stop()
		}
	}
}

func (s *StandardRegistry) Each(f func(string, interface{})) {
	metrics := s.registered()
	for i := range metrics {
		kv := &metrics[i]
		f(kv.name, kv.value)
	}
}

func (s *StandardRegistry) Get(name string) interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.metrics[name]
}

// Gets an existing metric or creates and registers a new one. Threadsafe
// alternative to calling Get and Register on failure.
// The interface can be the metric to register if not found in registry,
// or a function returning the metric for lazy instantiation.
func (s *StandardRegistry) GetOrRegister(name string, i interface{}) interface{} {
	// access the read lock first which should be re-entrant
	s.mutex.RLock()
	metric, ok := s.metrics[name]
	s.mutex.RUnlock()
	if ok {
		return metric
	}

	// only take the write lock if we'll be modifying the metrics map
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if metric, ok := s.metrics[name]; ok {
		return metric
	}
	if v := reflect.ValueOf(i); v.Kind() == reflect.Func {
		i = v.Call(nil)[0].Interface()
	}
	s.register(name, i)
	return i
}

func (s *StandardRegistry) Register(name string, i interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.register(name, i)
}

func (s *StandardRegistry) RunHealthchecks() {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for _, i := range s.metrics {
		if h, ok := i.(IHealthcheck); ok {
			h.Check()
		}
	}
}

func (s *StandardRegistry) GetAll() map[string]g.MapStrAny {
	data := make(map[string]g.MapStrAny)
	s.Each(func(name string, i interface{}) {
		values := make(map[string]interface{})
		switch metric := i.(type) {
		case ICounter:
			values["count"] = metric.Count()
		case IGauge:
			values["value"] = metric.Value()
		case IGaugeFloat64:
			values["value"] = metric.Value()
		case IHealthcheck:
			values["error"] = nil
			metric.Check()
			if err := metric.Error(); nil != err {
				values["error"] = metric.Error().Error()
			}
		case IHistogram:
			h := metric.Snapshot()
			ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			values["count"] = h.Count()
			values["min"] = h.Min()
			values["max"] = h.Max()
			values["mean"] = h.Mean()
			values["stddev"] = h.StdDev()
			values["median"] = ps[0]
			values["75%"] = ps[1]
			values["95%"] = ps[2]
			values["99%"] = ps[3]
			values["99.9%"] = ps[4]
		case IMeter:
			m := metric.Snapshot()
			values["count"] = m.Count()
			values["1m.rate"] = m.Rate1()
			values["5m.rate"] = m.Rate5()
			values["15m.rate"] = m.Rate15()
			values["mean.rate"] = m.RateMean()
		case ITimer:
			t := metric.Snapshot()
			ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			values["count"] = t.Count()
			values["min"] = t.Min()
			values["max"] = t.Max()
			values["mean"] = t.Mean()
			values["stddev"] = t.StdDev()
			values["median"] = ps[0]
			values["75%"] = ps[1]
			values["95%"] = ps[2]
			values["99%"] = ps[3]
			values["99.9%"] = ps[4]
			values["1m.rate"] = t.Rate1()
			values["5m.rate"] = t.Rate5()
			values["15m.rate"] = t.Rate15()
			values["mean.rate"] = t.RateMean()
		}
		data[name] = values
	})
	return data
}

func (s *StandardRegistry) Unregister(name string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.stop(name)
	delete(s.metrics, name)
}

func (s *StandardRegistry) UnregisterAll() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for name, _ := range s.metrics {
		s.stop(name)
		delete(s.metrics, name)
	}
}

// Create a new registry.
func NewRegistry() Registry {
	return &StandardRegistry{metrics: g.MapStrAny{}}
}

var DefaultRegistry Registry = NewRegistry()

type PrefixedRegistry struct {
	underlying Registry
	prefix     string
}

func NewPrefixedRegistry(prefix string) Registry {
	return &PrefixedRegistry{
		underlying: NewRegistry(),
		prefix:     prefix,
	}
}

func NewPrefixedChildRegistry(parent Registry, prefix string) Registry {
	return &PrefixedRegistry{
		underlying: parent,
		prefix:     prefix,
	}
}

// Call the given function for each registered metric.
func (r *PrefixedRegistry) Each(fn func(string, interface{})) {
	wrappedFn := func(prefix string) func(string, interface{}) {
		return func(name string, iface interface{}) {
			if strings.HasPrefix(name, prefix) {
				fn(name, iface)
			} else {
				return
			}
		}
	}

	baseRegistry, prefix := findPrefix(r, "")
	baseRegistry.Each(wrappedFn(prefix))
}

func findPrefix(registry Registry, prefix string) (Registry, string) {
	switch r := registry.(type) {
	case *PrefixedRegistry:
		return findPrefix(r.underlying, r.prefix+prefix)
	case *StandardRegistry:
		return r, prefix
	}
	return nil, ""
}

// Get the metric by the given name or nil if none is registered.
func (r *PrefixedRegistry) Get(name string) interface{} {
	realName := r.prefix + name
	return r.underlying.Get(realName)
}

// Gets an existing metric or registers the given one.
// The interface can be the metric to register if not found in registry,
// or a function returning the metric for lazy instantiation.
func (r *PrefixedRegistry) GetOrRegister(name string, metric interface{}) interface{} {
	realName := r.prefix + name
	return r.underlying.GetOrRegister(realName, metric)
}

// Register the given metric under the given name. The name will be prefixed.
func (r *PrefixedRegistry) Register(name string, metric interface{}) error {
	realName := r.prefix + name
	return r.underlying.Register(realName, metric)
}

// Run all registered healthchecks.
func (r *PrefixedRegistry) RunHealthchecks() {
	r.underlying.RunHealthchecks()
}

// GetAll metrics in the Registry
func (r *PrefixedRegistry) GetAll() map[string]map[string]interface{} {
	return r.underlying.GetAll()
}

// Unregister the metric with the given name. The name will be prefixed.
func (r *PrefixedRegistry) Unregister(name string) {
	realName := r.prefix + name
	r.underlying.Unregister(realName)
}

// Unregister all metrics.  (Mostly for testing.)
func (r *PrefixedRegistry) UnregisterAll() {
	r.underlying.UnregisterAll()
}

// Call the given function for each registered metric.
func Each(f func(string, interface{})) {
	DefaultRegistry.Each(f)
}

// Get the metric by the given name or nil if none is registered.
func Get(name string) interface{} {
	return DefaultRegistry.Get(name)
}

// Gets an existing metric or creates and registers a new one. Threadsafe
// alternative to calling Get and Register on failure.
func GetOrRegister(name string, i interface{}) interface{} {
	return DefaultRegistry.GetOrRegister(name, i)
}

// Register the given metric under the given name.  Returns a DuplicateMetric
// if a metric by the given name is already registered.
func Register(name string, i interface{}) error {
	return DefaultRegistry.Register(name, i)
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
	DefaultRegistry.RunHealthchecks()
}

// Unregister the metric with the given name.
func Unregister(name string) {
	DefaultRegistry.Unregister(name)
}
