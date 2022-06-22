package logx

import "github.com/sirupsen/logrus"

type (
	Level     string
	Formatter string
)

func (l Level) Str() string {
	return string(l)
}

func (f Formatter) Str() string {
	return string(f)
}

const (
	LevelPanic Level = "panic"
	LevelFatal Level = "fatal"
	LevelError Level = "error"
	LevelWarn  Level = "warn"
	LevelInfo  Level = "info"
	LevelDebug Level = "debug"
	LevelTrace Level = "trace"

	FormatterJSON Formatter = "JSON"
	FormatterText Formatter = "TEXT"

	DefaultLogPath string = "./runtime/log"
)

var (
	levelMap       map[Level]logrus.Level
	levelLogrusMap map[logrus.Level]Level
)

func init() {
	levelMap = make(map[Level]logrus.Level)
	levelMap[LevelPanic] = logrus.PanicLevel
	levelMap[LevelFatal] = logrus.FatalLevel
	levelMap[LevelError] = logrus.ErrorLevel
	levelMap[LevelWarn] = logrus.WarnLevel
	levelMap[LevelInfo] = logrus.InfoLevel
	levelMap[LevelDebug] = logrus.DebugLevel
	levelMap[LevelTrace] = logrus.TraceLevel

	levelLogrusMap = make(map[logrus.Level]Level)
	levelLogrusMap[logrus.PanicLevel] = LevelPanic
	levelLogrusMap[logrus.FatalLevel] = LevelFatal
	levelLogrusMap[logrus.ErrorLevel] = LevelError
	levelLogrusMap[logrus.WarnLevel] = LevelWarn
	levelLogrusMap[logrus.InfoLevel] = LevelInfo
	levelLogrusMap[logrus.DebugLevel] = LevelDebug
	levelLogrusMap[logrus.TraceLevel] = LevelTrace
}
