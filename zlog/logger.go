package zlog

import (
	"github.com/niclausse/webkit/consts"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	zapLogger     *zap.Logger
	sugaredLogger *zap.SugaredLogger
)

// log文件后缀类型
const (
	txtLogNormal    = "normal"    // 正常的日志：info、debug
	txtLogWarnFatal = "warnfatal" // 异常的日志： warn、error、fatal
	txtLogStdout    = "stdout"
)

func newLogger() *zap.Logger {
	var stdLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= logConfig.ZapLevel && lvl >= DebugLevel
	})

	var errLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= logConfig.ZapLevel && lvl >= WarnLevel
	})

	var infoLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= logConfig.ZapLevel && lvl <= InfoLevel
	})

	var zapCores []zapcore.Core
	if logConfig.Log2Stdout {
		var encoder zapcore.Encoder
		if logConfig.Mode == consts.DevelopMode {
			encoder = getConsoleEncoder()
		} else {
			encoder = getJsonEncoder()
		}
		zapCores = append(zapCores, zapcore.NewCore(encoder, getLogWriter(txtLogStdout), stdLevel))
	}

	if logConfig.Log2File {
		zapCores = append(zapCores, zapcore.NewCore(getJsonEncoder(), getLogWriter(txtLogNormal), infoLevel))
		zapCores = append(zapCores, zapcore.NewCore(getJsonEncoder(), getLogWriter(txtLogWarnFatal), errLevel))
	}

	core := zapcore.NewTee(zapCores...)

	return zap.New(core, zap.AddCaller(), zap.Fields(), zap.Development())
}

func getConsoleEncoder() zapcore.Encoder {
	// time字段编码器
	timeEncoder := zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.999999")

	encoderCfg := zapcore.EncoderConfig{
		LevelKey:       "level",
		TimeKey:        "time",
		CallerKey:      "file",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径编码器
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	return zapcore.NewConsoleEncoder(encoderCfg)
}

func getJsonEncoder() zapcore.Encoder {
	// time字段编码器
	timeEncoder := zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.999999")

	encoderCfg := zapcore.EncoderConfig{
		LevelKey:       "level",
		TimeKey:        "time",
		CallerKey:      "file",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径编码器
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	return zapcore.NewJSONEncoder(encoderCfg) // todo: custom webkit json_encoder
}

func getLogWriter(loggerType string) (ws zapcore.WriteSyncer) {
	var w io.Writer
	if loggerType == txtLogStdout {
		// stdOut
		w = os.Stdout
	} else {
		// 打印到 name.log[.wf] 中
		var err error
		filename := filepath.Join(strings.TrimSuffix(logConfig.Path, "/"), appendLogFileTail(logConfig.AppName, loggerType))
		w, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			panic("open log file error: " + err.Error())
		}
	}

	if !logConfig.BufferSwitch {
		return zapcore.AddSync(w)
	}

	// 开启缓冲区
	ws = &zapcore.BufferedWriteSyncer{
		WS:            zapcore.AddSync(w),
		Size:          logConfig.BufferSize,
		FlushInterval: logConfig.BufferFlushInterval,
		Clock:         nil,
	}
	return ws
}

// genFilename 拼装完整文件名
func appendLogFileTail(appName, loggerType string) string {
	var tailFixed string
	switch loggerType {
	case txtLogNormal:
		tailFixed = ".log"
	case txtLogWarnFatal:
		tailFixed = ".log.wf"
	default:
		tailFixed = ".log"
	}
	return appName + tailFixed
}
