package logx

import (
	"context"
	"github.com/niclausse/webkit/consts"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
)

const EmptyLogPath = ""

var std Logger

func StandardLogger() Logger {
	if std == nil {
		std = New(LevelDebug, EmptyLogPath)
	}
	return std
}

func Infof(format string, args ...interface{}) {
	StandardLogger().Infof(format, args)
}

func Warnf(format string, args ...interface{}) {
	StandardLogger().Warnf(format, args)
}

func Errorf(format string, args ...interface{}) {
	StandardLogger().Errorf(format, args)
}

func Info(args ...interface{}) {
	StandardLogger().Info(args)
}

func Warn(args ...interface{}) {
	StandardLogger().Warn(args)
}

func Error(args ...interface{}) {
	StandardLogger().Error(args)
}

func SetLevel(level Level) {
	StandardLogger().SetLevel(level)
}

func GetLevel() Level {
	return StandardLogger().GetLevel()
}

func SetOutput(output io.Writer) {
	StandardLogger().SetOutput(output)
}

func GetOuter() io.Writer {
	return StandardLogger().GetOuter()
}

func WithField(key string, value interface{}) Entry {
	return StandardLogger().WithField(key, value)
}

func WithContext(ctx context.Context) Entry {
	return StandardLogger().WithContext(ctx)
}

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

// logger implement Logger
type logger struct {
	*logrus.Logger
}

func (g *logger) GetOuter() io.Writer {
	return g.Logger.Out
}

func (g *logger) WithField(key string, value interface{}) Entry {
	field := g.Logger.WithField(key, value)
	return &entryX{field}
}

func (g *logger) WithContext(ctx context.Context) Entry {
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

func (g *logger) SetLevel(level Level) {
	g.Logger.SetLevel(levelMap[level])
}

func (g *logger) GetLevel() Level {
	return levelLogrusMap[g.Logger.GetLevel()]
}

func (g *logger) SetFormatter(formatter Formatter) {
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
	return &logger{Logger: logrus.New()}
}

// New return a logger
func New(level Level, logPath string) Logger {
	_logger := &logger{
		Logger: logrus.New(),
	}

	_logger.SetLevel(level)
	_logger.SetFormatter(FormatterJSON)
	_logger.SetReportCaller(true)

	if len(logPath) > 0 {
		_logger.AddHook(newLfsSizeHook(logPath))
	}

	return _logger
}
