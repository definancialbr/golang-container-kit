package prometheus

import (
	"net/http"
	"time"

	"github.com/GrooveDEF/golang-container-kit/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
)

func InterfaceSliceToMetricOptionSlice(interfaceOptions []interface{}) []MetricOption {

	options := make([]MetricOption, len(interfaceOptions))

	for i, interfaceOption := range interfaceOptions {
		options[i] = interfaceOption.(MetricOption)
	}

	return options

}

func StringSliceToPrometheusLabels(keysAndValues []string) prometheus.Labels {

	labels := make(prometheus.Labels)

	for i := 0; i < len(keysAndValues); i += 2 {
		labels[keysAndValues[i]] = keysAndValues[i+1]
	}

	return labels

}

var (
	DefaultHttpHandlerOptions = promhttp.HandlerOpts{
		Timeout:             10 * time.Second,
		DisableCompression:  true,
		MaxRequestsInFlight: 10,
		EnableOpenMetrics:   true,
	}
)

type MetricServiceOption func(*MetricService)

type MetricService struct {
	registry           *prometheus.Registry
	httpHandlerOptions promhttp.HandlerOpts
	pusher             *push.Pusher
	pusherGroupings    map[string]string
}

func WithHttpHandlerOptions(httpHandlerOptions promhttp.HandlerOpts) MetricServiceOption {
	return func(m *MetricService) {
		m.httpHandlerOptions = httpHandlerOptions
	}
}

func WithPusher(url, job string) MetricServiceOption {
	return func(m *MetricService) {
		m.pusher = push.New(url, job)
	}
}

func WithPusherGrouping(key, value string) MetricServiceOption {
	return func(m *MetricService) {
		m.pusherGroupings[key] = value
	}
}

func NewMetricService(options ...MetricServiceOption) *MetricService {

	m := &MetricService{
		registry:           prometheus.NewRegistry(),
		httpHandlerOptions: DefaultHttpHandlerOptions,
		pusherGroupings:    make(map[string]string),
	}

	for _, option := range options {
		option(m)
	}

	if m.pusher != nil {

		m.pusher = m.pusher.Gatherer(m.registry)

		for key, value := range m.pusherGroupings {
			m.pusher = m.pusher.Grouping(key, value)
		}

	}

	return m

}

func (m *MetricService) Counter(options ...interface{}) metrics.Counter {

	options = append(options, WithRegistry(m.registry))
	counter := NewCounter(InterfaceSliceToMetricOptionSlice(options)...)

	return counter

}

func (m *MetricService) Gauge(options ...interface{}) metrics.Gauge {

	options = append(options, WithRegistry(m.registry))
	gauge := NewGauge(InterfaceSliceToMetricOptionSlice(options)...)

	return gauge

}

func (m *MetricService) Histogram(options ...interface{}) metrics.Histogram {

	options = append(options, WithRegistry(m.registry))
	histogram := NewHistogram(InterfaceSliceToMetricOptionSlice(options)...)

	return histogram

}

func (m *MetricService) Handler() http.Handler {
	return promhttp.InstrumentMetricHandler(m.registry, promhttp.HandlerFor(m.registry, m.httpHandlerOptions))
}

func (m *MetricService) Push() error {
	return m.pusher.Push()
}

type MetricOption func(*MetricOptionSet)

type MetricOptionSet struct {
	prometheus.Opts
	Buckets  []float64
	Labels   []string
	registry *prometheus.Registry
}

func WithNamespace(namespace string) MetricOption {
	return func(optionSet *MetricOptionSet) {
		optionSet.Namespace = namespace
	}
}

func WithName(name string) MetricOption {
	return func(optionSet *MetricOptionSet) {
		optionSet.Name = name
	}
}

func WithSubsystem(subsystem string) MetricOption {
	return func(optionSet *MetricOptionSet) {
		optionSet.Subsystem = subsystem
	}
}

func WithHelp(help string) MetricOption {
	return func(optionSet *MetricOptionSet) {
		optionSet.Help = help
	}
}

func WithBuckets(buckets []float64) MetricOption {
	return func(optionSet *MetricOptionSet) {
		optionSet.Buckets = buckets
	}
}

func WithLabels(labels []string) MetricOption {
	return func(optionSet *MetricOptionSet) {
		optionSet.Labels = labels
	}
}

func WithRegistry(registry *prometheus.Registry) MetricOption {
	return func(optionSet *MetricOptionSet) {
		optionSet.registry = registry
	}
}

func NewMetricOptionSet(options ...MetricOption) *MetricOptionSet {

	o := &MetricOptionSet{}

	for _, option := range options {
		option(o)
	}

	return o

}

func (o *MetricOptionSet) AsCounterOpts() prometheus.CounterOpts {
	return prometheus.CounterOpts(o.Opts)
}

func (o *MetricOptionSet) AsGaugeOpts() prometheus.GaugeOpts {
	return prometheus.GaugeOpts(o.Opts)
}

func (o *MetricOptionSet) AsHistogramOpts() prometheus.HistogramOpts {
	return prometheus.HistogramOpts{
		Namespace: o.Namespace,
		Subsystem: o.Subsystem,
		Name:      o.Name,
		Help:      o.Help,
		Buckets:   o.Buckets,
	}
}

type Counter struct {
	vec     *prometheus.CounterVec
	counter prometheus.Counter
}

func NewCounter(options ...MetricOption) *Counter {

	o := NewMetricOptionSet(options...)
	c := &Counter{}

	if len(o.Labels) > 0 {
		c.vec = prometheus.NewCounterVec(o.AsCounterOpts(), o.Labels)
		o.registry.MustRegister(c.vec)
		return c
	}

	c.counter = prometheus.NewCounter(o.AsCounterOpts())
	o.registry.MustRegister(c.counter)

	return c

}

func (c *Counter) WithLabels(keysAndValues ...string) metrics.Counter {

	labels := StringSliceToPrometheusLabels(keysAndValues)

	return &Counter{
		vec:     c.vec,
		counter: c.vec.With(labels),
	}

}

func (c *Counter) WithLabelValues(values ...string) metrics.Counter {

	return &Counter{
		vec:     c.vec,
		counter: c.vec.WithLabelValues(values...),
	}

}

func (c *Counter) Add(value float64) {
	c.counter.Add(value)
}

type Gauge struct {
	vec   *prometheus.GaugeVec
	gauge prometheus.Gauge
}

func NewGauge(options ...MetricOption) *Gauge {

	o := NewMetricOptionSet(options...)
	g := &Gauge{}

	if len(o.Labels) > 0 {
		g.vec = prometheus.NewGaugeVec(o.AsGaugeOpts(), o.Labels)
		o.registry.MustRegister(g.vec)
		return g
	}

	g.gauge = prometheus.NewGauge(o.AsGaugeOpts())
	o.registry.MustRegister(g.gauge)

	return g

}

func (g *Gauge) WithLabels(keysAndValues ...string) metrics.Gauge {

	labels := StringSliceToPrometheusLabels(keysAndValues)

	return &Gauge{
		vec:   g.vec,
		gauge: g.vec.With(labels),
	}

}

func (g *Gauge) WithLabelValues(values ...string) metrics.Gauge {

	return &Gauge{
		vec:   g.vec,
		gauge: g.vec.WithLabelValues(values),
	}

}

func (g *Gauge) Add(value float64) {
	g.gauge.Add(value)
}

func (g *Gauge) Sub(value float64) {
	g.gauge.Sub(value)
}

func (g *Gauge) Set(value float64) {
	g.gauge.Set(value)
}

type Histogram struct {
	vec      *prometheus.HistogramVec
	observer prometheus.Observer
}

func NewHistogram(options ...MetricOption) *Histogram {

	o := NewMetricOptionSet(options...)
	h := &Histogram{}

	if len(o.Labels) > 0 {
		h.vec = prometheus.NewHistogramVec(o.AsHistogramOpts(), o.Labels)
		o.registry.MustRegister(h.vec)
		return h
	}

	histogram := prometheus.NewHistogram(o.AsHistogramOpts())
	h.observer = histogram
	o.registry.MustRegister(histogram)

	return h

}

func (g *Histogram) WithLabels(keysAndValues ...string) metrics.Histogram {

	labels := StringSliceToPrometheusLabels(keysAndValues)

	return &Histogram{
		vec:      g.vec,
		observer: g.vec.With(labels),
	}

}

func (g *Histogram) WithLabelValues(values ...string) metrics.Histogram {

	return &Histogram{
		vec:      g.vec,
		observer: g.vec.WithLabelValues(values...),
	}

}

func (g *Histogram) Observe(value float64) {
	g.observer.Observe(value)
}
