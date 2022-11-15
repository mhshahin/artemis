package artemis

import (
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Tracer struct {
	Metrics *TracingMetrics
}

func NewTracer(namespace string) *Tracer {
	metrics := NewTracingMetrics(namespace)

	collectors := metrics.GetCollectors()

	prometheus.MustRegister(collectors...)

	return &Tracer{
		Metrics: metrics,
	}
}

func (t *Tracer) RequestWithTracer(request *http.Request) *http.Request {
	requestStart := time.Now()

	httpTracer := NewHttpTracer(requestStart, t.Metrics)

	clientTrace := &httptrace.ClientTrace{
		GetConn:              httpTracer.GetConn,
		GotConn:              httpTracer.GotConn,
		GotFirstResponseByte: httpTracer.GotFirstResponseByte,
		DNSStart:             httpTracer.DNSStart,
		DNSDone:              httpTracer.DNSDone,
		ConnectStart:         httpTracer.ConnStart,
		ConnectDone:          httpTracer.ConnDone,
		WroteHeaders:         httpTracer.WroteHeaders,
		WroteRequest:         httpTracer.WroteRequest,
	}

	clientTraceCtx := httptrace.WithClientTrace(request.Context(), clientTrace)

	tracerRequest := request.WithContext(clientTraceCtx)

	return tracerRequest
}
