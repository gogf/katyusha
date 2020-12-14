package metrics

import (
	"github.com/gogf/gf/frame/g"
	"time"
)

// Log outputs each metric in the given registry periodically using the given logger.
func Log(r Registry, freq time.Duration) {
	LogScaled(r, freq, time.Nanosecond)
}

// LogOnCue outputs each metric in the given registry on demand through the channel
// using the given logger
func LogOnCue(r Registry, ch chan interface{}) {
	LogScaledOnCue(r, ch, time.Nanosecond)
}

// LogScaled outputs each metric in the given registry periodically using the given
// logger. Print timings in `scale` units (eg time.Millisecond) rather than nanos.
func LogScaled(r Registry, freq time.Duration, scale time.Duration) {
	ch := make(chan interface{})
	go func(channel chan interface{}) {
		tw.Add(freq, func() {
			channel <- struct{}{}
		})
	}(ch)
	LogScaledOnCue(r, ch, scale)
}

// LogScaledOnCue outputs each metric in the given registry on demand through the channel
// using the given logger. Print timings in `scale` units (eg time.Millisecond) rather
// than nanos.
func LogScaledOnCue(r Registry, ch chan interface{}, scale time.Duration) {
	du := float64(scale)
	duSuffix := scale.String()[1:]

	for _ = range ch {
		r.Each(func(name string, i interface{}) {
			switch metric := i.(type) {
			case ICounter:
				g.Log().Infof("counter %s", name)
				g.Log().Infof("  count:       %9d", metric.Count())
			case IGauge:
				g.Log().Infof("gauge %s", name)
				g.Log().Infof("  value:       %9d", metric.Value())
			case IGaugeFloat64:
				g.Log().Infof("gauge %s", name)
				g.Log().Infof("  value:       %f", metric.Value())
			case IHealthcheck:
				metric.Check()
				g.Log().Infof("healthcheck %s", name)
				g.Log().Infof("  error:       %v", metric.Error())
			case IHistogram:
				h := metric.Snapshot()
				ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
				g.Log().Infof("histogram %s", name)
				g.Log().Infof("  count:       %9d", h.Count())
				g.Log().Infof("  min:         %9d", h.Min())
				g.Log().Infof("  max:         %9d", h.Max())
				g.Log().Infof("  mean:        %12.2f", h.Mean())
				g.Log().Infof("  stddev:      %12.2f", h.StdDev())
				g.Log().Infof("  median:      %12.2f", ps[0])
				g.Log().Infof("  75%%:         %12.2f", ps[1])
				g.Log().Infof("  95%%:         %12.2f", ps[2])
				g.Log().Infof("  99%%:         %12.2f", ps[3])
				g.Log().Infof("  99.9%%:       %12.2f", ps[4])
			case IMeter:
				m := metric.Snapshot()
				g.Log().Infof("meter %s", name)
				g.Log().Infof("  count:       %9d", m.Count())
				g.Log().Infof("  1-min rate:  %12.2f", m.Rate1())
				g.Log().Infof("  5-min rate:  %12.2f", m.Rate5())
				g.Log().Infof("  15-min rate: %12.2f", m.Rate15())
				g.Log().Infof("  mean rate:   %12.2f", m.RateMean())
			case ITimer:
				t := metric.Snapshot()
				ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
				g.Log().Infof("timer %s", name)
				g.Log().Infof("  count:       %9d", t.Count())
				g.Log().Infof("  min:         %12.2f%s", float64(t.Min())/du, duSuffix)
				g.Log().Infof("  max:         %12.2f%s", float64(t.Max())/du, duSuffix)
				g.Log().Infof("  mean:        %12.2f%s", t.Mean()/du, duSuffix)
				g.Log().Infof("  stddev:      %12.2f%s", t.StdDev()/du, duSuffix)
				g.Log().Infof("  median:      %12.2f%s", ps[0]/du, duSuffix)
				g.Log().Infof("  75%%:         %12.2f%s", ps[1]/du, duSuffix)
				g.Log().Infof("  95%%:         %12.2f%s", ps[2]/du, duSuffix)
				g.Log().Infof("  99%%:         %12.2f%s", ps[3]/du, duSuffix)
				g.Log().Infof("  99.9%%:       %12.2f%s", ps[4]/du, duSuffix)
				g.Log().Infof("  1-min rate:  %12.2f", t.Rate1())
				g.Log().Infof("  5-min rate:  %12.2f", t.Rate5())
				g.Log().Infof("  15-min rate: %12.2f", t.Rate15())
				g.Log().Infof("  mean rate:   %12.2f", t.RateMean())
			}
		})
	}
}
