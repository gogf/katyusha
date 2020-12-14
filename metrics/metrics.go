package metrics

import (
	"github.com/gogf/gf/os/gtimer"
	"time"
)

// UseNilMetrics is checked by the constructor functions for all of the
// standard metrics.  If it is true, the metric returned is a stub.
//
// This global kill-switch helps quantify the observer effect and makes
// for less cluttered pprof profiles.
var (
	UseNilMetrics bool = false
	tw                 = gtimer.New(100, 10*time.Millisecond, 6)
)
