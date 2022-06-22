package logx

import (
	"io"

	"gopkg.in/natefinch/lumberjack.v2"

	log "github.com/sirupsen/logrus"

	"github.com/rifflock/lfshook"
)

var (
	trace2InfoFile, warnFile, errorFile, fatal2PanicFile string
)

func setLogDir(logDir string) {
	trace2InfoFile = logDir + "/info/trace2info.log"
	warnFile = logDir + "/warn/warn.log"
	errorFile = logDir + "/error/error.log"
	fatal2PanicFile = logDir + "/fatal/fatal2panic.log"
}

// newLfsSizeHook 创建一个本地日志写入及按日志文件大小切割的hook
func newLfsSizeHook(logDir string) log.Hook {

	trace2Info := newOneLumberjackWriter(string(LevelDebug), logDir, 100, 30, true, false)
	warn := newOneLumberjackWriter(string(LevelWarn), logDir, 100, 30, true, false)
	err := newOneLumberjackWriter(string(LevelError), logDir, 100, 30, true, false)
	fatal2Panic := newOneLumberjackWriter(string(LevelFatal), logDir, 100, 30, true, false)

	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		log.TraceLevel: trace2Info,
		log.DebugLevel: trace2Info,
		log.InfoLevel:  trace2Info,
		log.WarnLevel:  warn,
		log.ErrorLevel: err,
		log.FatalLevel: fatal2Panic,
		log.PanicLevel: fatal2Panic,
	}, &log.JSONFormatter{})

	return lfsHook
}

func newOneLumberjackWriter(level, logDir string, maxSize, maxAge int, localTime, compress bool) io.Writer {

	setLogDir(logDir)

	switch level {
	case string(LevelTrace):
		fallthrough
	case string(LevelDebug):
		fallthrough
	case string(LevelInfo):
		return &lumberjack.Logger{
			Filename:  trace2InfoFile,
			MaxSize:   maxSize,
			MaxAge:    maxAge,
			LocalTime: localTime,
			Compress:  compress,
		}
	case string(LevelWarn):
		return &lumberjack.Logger{
			Filename:  warnFile,
			MaxSize:   maxSize,
			MaxAge:    maxAge,
			LocalTime: localTime,
			Compress:  compress,
		}
	case string(LevelError):
		return &lumberjack.Logger{
			Filename:  errorFile,
			MaxSize:   maxSize,
			MaxAge:    maxAge,
			LocalTime: localTime,
			Compress:  compress,
		}
	}
	return &lumberjack.Logger{
		Filename:  fatal2PanicFile,
		MaxSize:   maxSize,
		MaxAge:    maxAge,
		LocalTime: localTime,
		Compress:  compress,
	}
}
