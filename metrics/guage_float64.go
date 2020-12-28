package metrics

import (
	"math"
	"sync/atomic"
)

// GaugeFloat64s hold a float64 value that can be set arbitrarily.
type IGaugeFloat64 interface {
	Snapshot() IGaugeFloat64
	Update(float64)
	Value() float64
}

// GaugeFloat64 returns an existing GaugeFloat64 or constructs and registers to defaultRegistry
func GaugeFloat64(name string) IGaugeFloat64 {
	return GetOrRegisterGaugeFloat64(name, defaultRegistry)
}

// GetOrRegisterGaugeFloat64 returns an existing GaugeFloat64 or constructs and registers a
// new StandardGaugeFloat64.
func GetOrRegisterGaugeFloat64(name string, r Registry) IGaugeFloat64 {
	if nil == r {
		r = defaultRegistry
	}
	return r.GetOrRegister(name, NewGaugeFloat64()).(IGaugeFloat64)
}

// NewGaugeFloat64 constructs a new StandardGaugeFloat64.
func NewGaugeFloat64() IGaugeFloat64 {
	return &StandardGaugeFloat64{
		value: 0.0,
	}
}

// NewRegisteredGaugeFloat64 constructs and registers a new StandardGaugeFloat64.
func NewRegisteredGaugeFloat64(name string, r Registry) IGaugeFloat64 {
	c := NewGaugeFloat64()
	if nil == r {
		r = defaultRegistry
	}
	r.Register(name, c)
	return c
}

// NewFunctionalGauge constructs a new FunctionalGauge.
func NewFunctionalGaugeFloat64(f func() float64) IGaugeFloat64 {
	return &FunctionalGaugeFloat64{value: f}
}

// NewRegisteredFunctionalGauge constructs and registers a new StandardGauge.
func NewRegisteredFunctionalGaugeFloat64(name string, r Registry, f func() float64) IGaugeFloat64 {
	c := NewFunctionalGaugeFloat64(f)
	if nil == r {
		r = defaultRegistry
	}
	r.Register(name, c)
	return c
}

// GaugeFloat64Snapshot is a read-only copy of another GaugeFloat64.
type GaugeFloat64Snapshot float64

// Snapshot returns the snapshot.
func (g GaugeFloat64Snapshot) Snapshot() IGaugeFloat64 { return g }

// Update panics.
func (GaugeFloat64Snapshot) Update(float64) {
	panic("Update called on a GaugeFloat64Snapshot")
}

// Value returns the value at the time the snapshot was taken.
func (g GaugeFloat64Snapshot) Value() float64 { return float64(g) }

// StandardGaugeFloat64 is the standard implementation of a GaugeFloat64 and uses
// sync.Mutex to manage a single float64 value.
type StandardGaugeFloat64 struct {
	value uint64
}

// Snapshot returns a read-only copy of the gauge.
func (g *StandardGaugeFloat64) Snapshot() IGaugeFloat64 {
	return GaugeFloat64Snapshot(g.Value())
}

// Update updates the gauge's value.
func (g *StandardGaugeFloat64) Update(v float64) {
	atomic.StoreUint64(&g.value, math.Float64bits(v))
}

// Value returns the gauge's current value.
func (g *StandardGaugeFloat64) Value() float64 {
	return math.Float64frombits(atomic.LoadUint64(&g.value))
}

// FunctionalGaugeFloat64 returns value from given function
type FunctionalGaugeFloat64 struct {
	value func() float64
}

// Value returns the gauge's current value.
func (g FunctionalGaugeFloat64) Value() float64 {
	return g.value()
}

// Snapshot returns the snapshot.
func (g FunctionalGaugeFloat64) Snapshot() IGaugeFloat64 { return GaugeFloat64Snapshot(g.Value()) }

// Update panics.
func (FunctionalGaugeFloat64) Update(float64) {
	panic("Update called on a FunctionalGaugeFloat64")
}
