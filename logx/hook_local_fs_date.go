package logx

import (
	"io"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
)

func newOneRotateFileWriter(level log.Level, logDir string) (io.Writer, error) {

	setLogDir(logDir)

	switch level {
	case log.TraceLevel:
		fallthrough
	case log.DebugLevel:
		fallthrough
	case log.InfoLevel:
		return rotatelogs.New(
			trace2InfoFile+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(trace2InfoFile), // 设置软连接
			//rotatelogs.WithRotationSize(2),        // 设置日志文件滚动大小
			rotatelogs.WithRotationTime(time.Hour*24), // 设置滚动时间
			rotatelogs.WithRotationCount(20),          // 设置日志文件最大保留数量
		)
	case log.WarnLevel:
		return rotatelogs.New(
			warnFile+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(warnFile),
			rotatelogs.WithRotationTime(time.Hour*24),
			rotatelogs.WithRotationCount(20),
		)
	case log.ErrorLevel:
		return rotatelogs.New(
			errorFile+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(errorFile),
			rotatelogs.WithRotationTime(time.Hour*24),
			rotatelogs.WithRotationCount(20),
		)
	}
	return rotatelogs.New(
		fatal2PanicFile+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(fatal2PanicFile),
		rotatelogs.WithRotationTime(time.Hour*24),
		rotatelogs.WithRotationCount(20),
	)
}

// newLfsDateHook 创建一个本地日志写入及按日期切割的hook
func newLfsDateHook(logDir string) log.Hook {

	infoWriter, err := newOneRotateFileWriter(log.InfoLevel, logDir)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("new trace2info IO writer")
	}

	warnWriter, err := newOneRotateFileWriter(log.WarnLevel, logDir)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("new warn IO writer")
	}

	errorWriter, err := newOneRotateFileWriter(log.ErrorLevel, logDir)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("new error IO writer")
	}

	fatal2panicWriter, err := newOneRotateFileWriter(log.FatalLevel, logDir)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("new fatal2panic IO writer")
	}

	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		log.TraceLevel: infoWriter,
		log.DebugLevel: infoWriter,
		log.InfoLevel:  infoWriter,
		log.WarnLevel:  warnWriter,
		log.ErrorLevel: errorWriter,
		log.FatalLevel: fatal2panicWriter,
		log.PanicLevel: fatal2panicWriter,
	}, &log.JSONFormatter{})

	return lfsHook
}
