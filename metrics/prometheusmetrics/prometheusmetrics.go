package prometheusmetrics

import (
	"fmt"
	"github.com/gogf/katyusha/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"strings"
	"sync"
	"time"
)

// PrometheusConfig provides a container with config parameters for the
// Prometheus Exporter
type PrometheusConfig struct {
	namespace        string
	subsystem        string
	promRegistry     prometheus.Registerer //Prometheus registry
	FlushInterval    time.Duration         //interval to update prom metrics
	gauges           map[string]prometheus.Gauge
	customMetrics    map[string]*CustomCollector
	histogramBuckets []float64
	timerBuckets     []float64
	mutex            *sync.Mutex
}

// NewPrometheusProvider returns a Provider that produces Prometheus metrics.
// Namespace and subsystem are applied to all produced metrics.
func NewPrometheusProvider(namespace string, subsystem string, promRegistry prometheus.Registerer, FlushInterval time.Duration) *PrometheusConfig {
	return &PrometheusConfig{
		namespace:        namespace,
		subsystem:        subsystem,
		promRegistry:     promRegistry,
		FlushInterval:    FlushInterval,
		gauges:           make(map[string]prometheus.Gauge),
		customMetrics:    make(map[string]*CustomCollector),
		histogramBuckets: []float64{0.05, 0.1, 0.25, 0.50, 0.75, 0.9, 0.95, 0.99},
		timerBuckets:     []float64{0.50, 0.95, 0.99, 0.999},
		mutex:            new(sync.Mutex),
	}
}

func (c *PrometheusConfig) WithHistogramBuckets(b []float64) *PrometheusConfig {
	c.histogramBuckets = b
	return c
}

func (c *PrometheusConfig) WithTimerBuckets(b []float64) *PrometheusConfig {
	c.timerBuckets = b
	return c
}

func (c *PrometheusConfig) flattenKey(key string) string {
	key = strings.Replace(key, " ", "_", -1)
	key = strings.Replace(key, ".", "_", -1)
	key = strings.Replace(key, "-", "_", -1)
	key = strings.Replace(key, "=", "_", -1)
	key = strings.Replace(key, "/", "_", -1)
	return key
}

func (c *PrometheusConfig) createKey(name string) string {
	return fmt.Sprintf("%s_%s_%s", c.namespace, c.subsystem, name)
}

func (c *PrometheusConfig) gaugeFromNameAndValue(name string, val float64) {
	key := c.createKey(name)
	g, ok := c.gauges[key]
	if !ok {
		g = prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: c.flattenKey(c.namespace),
			Subsystem: c.flattenKey(c.subsystem),
			Name:      c.flattenKey(name),
			Help:      name,
		})
		c.promRegistry.Register(g)
		c.gauges[key] = g
	}
	g.Set(val)
}

func (c *PrometheusConfig) histogramFromNameAndMetric(name string, goMetric interface{}, buckets []float64) {
	key := c.createKey(name)

	collector, ok := c.customMetrics[key]
	if !ok {
		collector = NewCustomCollector(c.mutex)
		c.promRegistry.MustRegister(collector)
		c.customMetrics[key] = collector
	}

	var ps []float64
	var count uint64
	var sum float64
	var typeName string

	switch metric := goMetric.(type) {
	case metrics.IHistogram:
		snapshot := metric.Snapshot()
		ps = snapshot.Percentiles(buckets)
		count = uint64(snapshot.Count())
		sum = float64(snapshot.Sum())
		typeName = "histogram"
	case metrics.ITimer:
		snapshot := metric.Snapshot()
		ps = snapshot.Percentiles(buckets)
		count = uint64(snapshot.Count())
		sum = float64(snapshot.Sum())
		typeName = "timer"
	default:
		panic(fmt.Sprintf("unexpected metric type %T", goMetric))
	}

	bucketVals := make(map[float64]uint64)
	for ii, bucket := range buckets {
		bucketVals[bucket] = uint64(ps[ii])
	}

	desc := prometheus.NewDesc(
		prometheus.BuildFQName(
			c.flattenKey(c.namespace),
			c.flattenKey(c.subsystem),
			fmt.Sprintf("%s_%s", c.flattenKey(name), typeName),
		),
		c.flattenKey(name),
		[]string{},
		map[string]string{},
	)

	if constHistogram, err := prometheus.NewConstHistogram(
		desc,
		count,
		sum,
		bucketVals,
	); err == nil {
		c.mutex.Lock()
		collector.metric = constHistogram
		c.mutex.Unlock()
	}
}

func (c *PrometheusConfig) UpdatePrometheusMetrics() {
	for _ = range time.Tick(c.FlushInterval) {
		c.UpdatePrometheusMetricsOnce()
	}
}

func (c *PrometheusConfig) UpdatePrometheusMetricsOnce() error {
	metrics.Each(func(name string, i interface{}) {
		switch metric := i.(type) {
		case metrics.ICounter:
			c.gaugeFromNameAndValue(name, float64(metric.Count()))
		case metrics.IGauge:
			c.gaugeFromNameAndValue(name, float64(metric.Value()))
		case metrics.IGaugeFloat64:
			c.gaugeFromNameAndValue(name, float64(metric.Value()))
		case metrics.IHistogram:
			samples := metric.Snapshot().Sample().Values()
			if len(samples) > 0 {
				lastSample := samples[len(samples)-1]
				c.gaugeFromNameAndValue(name, float64(lastSample))
			}

			c.histogramFromNameAndMetric(name, metric, c.histogramBuckets)
		case metrics.IMeter:
			lastSample := metric.Snapshot().Rate1()
			c.gaugeFromNameAndValue(name, float64(lastSample))
		case metrics.ITimer:
			lastSample := metric.Snapshot().Rate1()
			c.gaugeFromNameAndValue(name, float64(lastSample))

			c.histogramFromNameAndMetric(name, metric, c.timerBuckets)
		}
	})
	return nil
}

// for collecting prometheus.constHistogram objects
type CustomCollector struct {
	prometheus.Collector

	metric prometheus.Metric
	mutex  *sync.Mutex
}

func NewCustomCollector(mutex *sync.Mutex) *CustomCollector {
	return &CustomCollector{
		mutex: mutex,
	}
}

func (c *CustomCollector) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock()
	if c.metric != nil {
		val := c.metric
		ch <- val
	}
	c.mutex.Unlock()
}

func (p *CustomCollector) Describe(ch chan<- *prometheus.Desc) {
	// empty method to fulfill prometheus.Collector interface
}
