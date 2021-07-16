package metrics

import "net/http"

type MetricService interface {
	Counter(...interface{}) Counter
	Gauge(...interface{}) Gauge
	Histogram(...interface{}) Histogram
	Handler() http.Handler
	Push() error
}

type Counter interface {
	WithLabels(...string) Counter
	Add(float64)
}

type Gauge interface {
	WithLabels(...string) Gauge
	Add(float64)
	Sub(float64)
	Set(float64)
}

type Histogram interface {
	WithLabels(...string) Histogram
	Observe(float64)
}
