package consts

type TraceKey string

func (t TraceKey) String() string {
	return string(t)
}

const (
	OpenTracer          TraceKey = "Tracer"
	OpenTraceParentSpan TraceKey = "parentSpan"
	OpenTraceTraceID    TraceKey = "X-Trace-ID"
	OpenTraceSpanID     TraceKey = "X-Span-ID"
)
