package artemis

import (
	"crypto/tls"
	"net/http/httptrace"
	"time"
)

type HttpTracer struct {
	RequestStartTime time.Time

	GetConnTime           time.Time
	DNSStartTime          time.Time
	DNSHost               string
	ConnectStartTime      time.Time
	TLSHandshakeStartTime time.Time

	Metrics *TracingMetrics

	ReqMethod string
	ReqURL    string
}

func NewHttpTracer(requestStart time.Time, metrics *TracingMetrics, method, reqUrl string) *HttpTracer {
	return &HttpTracer{
		RequestStartTime: requestStart,
		Metrics:          metrics,
		ReqMethod:        method,
		ReqURL:           reqUrl,
	}
}

func (ht *HttpTracer) GetConn(hostPort string) {
	ht.GetConnTime = time.Now()
}

func (ht *HttpTracer) GotConn(info httptrace.GotConnInfo) {
	ht.Metrics.GetConnectionDurationSecondsMetric(ht.GetConnTime, ht.ReqMethod, ht.ReqURL)

	if info.Reused {
		ht.Metrics.ReuseConnectionsMetric()
	}

	if info.WasIdle {
		ht.Metrics.ReuseIdleConnectionsMetric()
	}
}

func (ht *HttpTracer) GotFirstResponseByte() {
	ht.Metrics.FirstByteReceiveDurationSecondsMetric(ht.RequestStartTime, ht.ReqMethod, ht.ReqURL)
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
	ht.Metrics.ConnectionHandshakeDurationSecondsMetric(ht.ConnectStartTime, ht.ReqMethod, ht.ReqURL)
}

func (ht *HttpTracer) TLSHandshakeStart() {
	ht.TLSHandshakeStartTime = time.Now()
}

func (ht *HttpTracer) TLSHandshakeDone(tls.ConnectionState, error) {
	ht.Metrics.TLSHandshakeDurationSecondsMetric(ht.TLSHandshakeStartTime, ht.ReqMethod, ht.ReqURL)
}

func (ht *HttpTracer) WroteHeaders() {
	ht.Metrics.HeaderWriteDrurationSecondsMetric(ht.RequestStartTime, ht.ReqMethod, ht.ReqURL)
}

func (ht *HttpTracer) WroteRequest(info httptrace.WroteRequestInfo) {
	ht.Metrics.RequestWriteDurationSecondsMetric(ht.RequestStartTime, ht.ReqMethod, ht.ReqURL)
}
