package artemis

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type TracingMetrics struct {
	GetConnectionDurationSeconds       *prometheus.HistogramVec
	ReuseConnections                   *prometheus.CounterVec
	ReuseIdleConnections               *prometheus.CounterVec
	FirstByteReceiveDurationSeconds    *prometheus.HistogramVec
	DNSLookupDurationSeconds           *prometheus.HistogramVec
	DNSCoalesced                       *prometheus.CounterVec
	ConnectionHandshakeDurationSeconds *prometheus.HistogramVec
	HeaderWriteDrurationSeconds        *prometheus.HistogramVec
	RequestWriteDurationSeconds        *prometheus.HistogramVec
}

var (
	SmallDurationBuckets = []float64{.0001, .0005, .001, .002, .005, .01, .05, .1, 1, 2.5, 5, 10}

	HttpLabels prometheus.Labels
	HostLabel  = "host"
)

func NewTracingMetrics(namespace string) *TracingMetrics {
	return &TracingMetrics{
		GetConnectionDurationSeconds: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "http_get_connection_duration_seconds",
			Help:      "HTTP Get Connection Duration",
			Buckets:   SmallDurationBuckets,
		}, []string{}),
		ReuseConnections: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "http_reuse_connections",
			Help:      "HTTP Connection Re-use Counter",
		}, []string{}),
		ReuseIdleConnections: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "http_reuse_idle_connections",
			Help:      "HTTP Idle Connection Re-use Counter",
		}, []string{}),
		FirstByteReceiveDurationSeconds: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "http_first_byte_response_duration_seconds",
			Help:      "HTTP Duration of Getting First Response Bytes",
			Buckets:   SmallDurationBuckets,
		}, []string{}),
		DNSLookupDurationSeconds: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "http_dns_lookup_duration_seconds",
			Help:      "HTTP DNS Lookup Duration",
			Buckets:   SmallDurationBuckets,
		}, []string{HostLabel}),
		DNSCoalesced: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "http_dns_coalesced_queries_counter",
			Help:      "HTTP DNS Query Coalesced Counter",
		}, []string{HostLabel}),
		ConnectionHandshakeDurationSeconds: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "http_connection_handshake_duration_seconds",
			Help:      "HTTP Connection Handshake Duration",
			Buckets:   SmallDurationBuckets,
		}, []string{}),
		HeaderWriteDrurationSeconds: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "http_header_write_duration_seconds",
			Help:      "HTTP Header Write Duration",
			Buckets:   SmallDurationBuckets,
		}, []string{}),
		RequestWriteDurationSeconds: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "http_request_write_duration_seconds",
			Help:      "HTTP Request Write Duration",
			Buckets:   SmallDurationBuckets,
		}, []string{}),
	}
}

func (tm *TracingMetrics) GetConnectionDurationSecondsMetric(getConnTime time.Time) {
	getConnDuration := time.Since(getConnTime)
	tm.GetConnectionDurationSeconds.With(HttpLabels).Observe(getConnDuration.Seconds())
}

func (tm *TracingMetrics) ReuseConnectionsMetric() {
	tm.ReuseConnections.With(HttpLabels).Inc()
}

func (tm *TracingMetrics) ReuseIdleConnectionsMetric() {
	tm.ReuseIdleConnections.With(HttpLabels).Inc()
}

func (tm *TracingMetrics) FirstByteReceiveDurationSecondsMetric(startTime time.Time) {
	getFirstByteDuration := time.Since(startTime)
	tm.FirstByteReceiveDurationSeconds.With(HttpLabels).Observe(getFirstByteDuration.Seconds())
}

func (tm *TracingMetrics) DNSLookupDurationSecondsMetric(dnsStartTime time.Time, dnsHost string) {
	dnsDuration := time.Since(dnsStartTime)
	tm.DNSLookupDurationSeconds.With(prometheus.Labels{HostLabel: dnsHost}).Observe(dnsDuration.Seconds())
}

func (tm *TracingMetrics) DNSCoalescedMetric(dnsHost string) {
	tm.DNSCoalesced.With(prometheus.Labels{HostLabel: dnsHost}).Inc()
}

func (tm *TracingMetrics) ConnectionHandshakeDurationSecondsMetric(connStartTime time.Time) {
	handshakeDuration := time.Since(connStartTime)
	tm.ConnectionHandshakeDurationSeconds.With(HttpLabels).Observe(handshakeDuration.Seconds())
}

func (tm *TracingMetrics) HeaderWriteDrurationSecondsMetric(startTime time.Time) {
	headerWriteDuration := time.Since(startTime)
	tm.HeaderWriteDrurationSeconds.With(HttpLabels).Observe(headerWriteDuration.Seconds())
}

func (tm *TracingMetrics) RequestWriteDurationSecondsMetric(startTime time.Time) {
	requestWriteDuration := time.Since(startTime)
	tm.RequestWriteDurationSeconds.With(HttpLabels).Observe(requestWriteDuration.Seconds())
}

func (tm *TracingMetrics) GetCollectors() []prometheus.Collector {
	return []prometheus.Collector{
		tm.GetConnectionDurationSeconds,
		tm.ReuseIdleConnections,
		tm.FirstByteReceiveDurationSeconds,
		tm.DNSLookupDurationSeconds,
		tm.DNSCoalesced,
		tm.ConnectionHandshakeDurationSeconds,
		tm.HeaderWriteDrurationSeconds,
		tm.RequestWriteDurationSeconds,
	}
}
