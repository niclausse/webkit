package logx

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
)

type traceHook struct {
}

func (t *traceHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
func (t *traceHook) Fire(entry *logrus.Entry) error {

	if entry.Context == nil {
		return nil
	}

	fromContextSpan := opentracing.SpanFromContext(entry.Context)
	if fromContextSpan != nil {
		var traceID string
		var spanID string
		var spanContext = fromContextSpan.Context()
		switch assertData := spanContext.(type) {
		case jaeger.SpanContext:
			traceID = assertData.TraceID().String()
			spanID = assertData.SpanID().String()
			entry.WithField("X-Trace-ID", traceID)
			entry.WithField("X-Span-ID", spanID)

		}
	}

	return nil
}
