package logx

import (
	"context"
	"io"

	"github.com/penglin1995/webflow/consts"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
)

type Logger interface {
	WithField(key string, value interface{}) Entry
	WithContext(ctx context.Context) Entry
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	SetLevel(level Level)
	GetLevel() Level
	SetOutput(output io.Writer)
	GetOuter() io.Writer
}

type Entry interface {
	WithField(key string, value interface{}) Entry
	WithContext(ctx context.Context) Entry
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

// entryX implement Entry
type entryX struct {
	*logrus.Entry
}

func (e *entryX) WithField(key string, value interface{}) Entry {
	field := e.Entry.WithField(key, value)
	return &entryX{field}
}

func (e *entryX) WithContext(ctx context.Context) Entry {
	field := e.Entry.WithContext(ctx)
	return &entryX{field}
}

// gLog implement Logger
type gLog struct {
	*logrus.Logger
}

func (g *gLog) GetOuter() io.Writer {
	return g.Logger.Out
}

func (g *gLog) WithField(key string, value interface{}) Entry {
	field := g.Logger.WithField(key, value)
	return &entryX{field}
}

func (g *gLog) WithContext(ctx context.Context) Entry {
	var (
		traceID         string
		spanID          string
		fromContextSpan opentracing.Span
	)
	if ctx == nil {
		goto Final
	}

	fromContextSpan = opentracing.SpanFromContext(ctx)
	if fromContextSpan != nil {

		var spanContext = fromContextSpan.Context()
		switch assertData := spanContext.(type) {
		case jaeger.SpanContext:
			traceID = assertData.TraceID().String()
			spanID = assertData.SpanID().String()
			goto Final
		}
	}

Final:
	return g.WithField(consts.OpenTraceTraceID.String(), traceID).
		WithField(consts.OpenTraceSpanID.String(), spanID)

}

func (g *gLog) SetLevel(level Level) {
	g.Logger.SetLevel(levelMap[level])
}

func (g *gLog) GetLevel() Level {
	return levelLogrusMap[g.Logger.GetLevel()]
}

func (g *gLog) SetFormatter(formatter Formatter) {
	if formatter == "" {
		return
	}

	if formatter == FormatterJSON {
		g.Logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		g.Logger.SetFormatter(&logrus.TextFormatter{})
	}
}

// Default return a default logger
func Default() Logger {
	GLog := &gLog{
		logrus.New(),
	}
	return GLog
}

// New return a logger
func New(level Level, logPath string) Logger {
	GLog := &gLog{
		logrus.New(),
	}

	GLog.SetLevel(level)
	GLog.SetFormatter(FormatterJSON)
	GLog.SetReportCaller(true)
	GLog.AddHook(newLfsSizeHook(logPath))

	return GLog

}
