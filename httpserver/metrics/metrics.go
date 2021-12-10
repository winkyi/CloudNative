package metrics

import (
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

const MetricsNamespace = "httpserver"

var (
	functionLatency = CreateExcutionTimeMetric(MetricsNamespace, "Time spent")
)

func Register() {
	if err := prometheus.Register(functionLatency); err != nil {
		glog.Error(err)
	}
}

// ExcutionTimer
type ExcutionTimer struct {
	histo *prometheus.HistogramVec
	start time.Time
	last  time.Time
}

// NewExcutionTimer
func NewExcutionTimer(histo *prometheus.HistogramVec) *ExcutionTimer {
	now := time.Now()
	return &ExcutionTimer{
		histo: histo,
		start: now,
		last:  now,
	}
}

func NewTimer() *ExcutionTimer {
	return NewExcutionTimer(functionLatency)
}

func (t *ExcutionTimer) ObserveTotal() {
	(*t.histo).WithLabelValues("total").Observe(time.Now().Sub(t.start).Seconds())
}

func CreateExcutionTimeMetric(namespace, help string) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "latency_seconds",
			Help:      help,
			Buckets:   prometheus.ExponentialBuckets(0.001, 2, 15),
		}, []string{"step"})
}
