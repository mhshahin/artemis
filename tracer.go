package artemis

import (
	"net/http/httptrace"
	"time"
)

type HttpTracer struct {
	RequestStartTime time.Time

	GetConnTime      time.Time
	DNSStartTime     time.Time
	DNSHost          string
	ConnectStartTime time.Time

	Metrics *TracingMetrics
}

func NewHttpTracer(requestStart time.Time, metrics *TracingMetrics) *HttpTracer {
	return &HttpTracer{
		RequestStartTime: requestStart,
		Metrics:          metrics,
	}
}

func (ht *HttpTracer) GetConn(hostPort string) {
	ht.GetConnTime = time.Now()
}

func (ht *HttpTracer) GotConn(info httptrace.GotConnInfo) {
	ht.Metrics.GetConnectionDurationSecondsMetric(ht.GetConnTime)

	if info.Reused {
		ht.Metrics.ReuseConnectionsMetric()
	}

	if info.WasIdle {
		ht.Metrics.ReuseIdleConnectionsMetric()
	}
}

func (ht *HttpTracer) GotFirstResponseByte() {
	ht.Metrics.FirstByteReceiveDurationSecondsMetric(ht.RequestStartTime)
}

func (ht *HttpTracer) DNSStart(info httptrace.DNSStartInfo) {
	ht.DNSStartTime = time.Now()
	ht.DNSHost = info.Host
}

func (ht *HttpTracer) DNSDone(info httptrace.DNSDoneInfo) {
	ht.Metrics.DNSLookupDurationSecondsMetric(ht.DNSStartTime, ht.DNSHost)

	if info.Coalesced {
		ht.Metrics.DNSCoalescedMetric(ht.DNSHost)
	}
}

func (ht *HttpTracer) ConnStart(network, addr string) {
	ht.ConnectStartTime = time.Now()
}

func (ht *HttpTracer) ConnDone(network, addr string, err error) {
	ht.Metrics.ConnectionHandshakeDurationSecondsMetric(ht.ConnectStartTime)
}

func (ht *HttpTracer) WroteHeaders() {
	ht.Metrics.HeaderWriteDrurationSecondsMetric(ht.RequestStartTime)
}

func (ht *HttpTracer) WroteRequest(info httptrace.WroteRequestInfo) {
	ht.Metrics.RequestWriteDurationSecondsMetric(ht.RequestStartTime)
}
