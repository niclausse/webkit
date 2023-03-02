package zlog

import (
	"fmt"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type Level = zapcore.Level

const (
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	WarnLevel  = zapcore.WarnLevel
	ErrorLevel = zapcore.ErrorLevel
	FatalLevel = zapcore.FatalLevel
)

const (
	defaultLogPath = "./log"
	defaultAppName = "app"
)

type config struct {
	ZapLevel zapcore.Level
	AppName  string

	// 以下变量仅对开发环境生效
	Log2Stdout bool
	Log2File   bool
	Path       string

	// 缓冲区
	BufferSwitch        bool
	BufferSize          int
	BufferFlushInterval time.Duration
}

// 全局配置 仅限InitLogConfig函数进行变更
var logConfig *config

func init() {
	logConfig = &config{
		ZapLevel: zapcore.InfoLevel,
		AppName:  defaultAppName,

		Log2Stdout: true,
		Log2File:   true,
		Path:       defaultLogPath,

		// 缓冲区，如果不配置默认使用以下配置
		BufferSwitch:        true,
		BufferSize:          256 * 1024, // 256kb
		BufferFlushInterval: 5 * time.Second,
	}

	InitLogConfig()
}

type Option func()

func WithAppName(appName string) Option {
	return func() {
		logConfig.AppName = appName
	}
}

func WithLevel(level Level) Option {
	return func() {
		logConfig.ZapLevel = level
	}
}

func WithLog2Stdout(stdout bool) Option {
	return func() {
		logConfig.Log2Stdout = stdout
	}
}

func WithLog2File(file bool) Option {
	return func() {
		logConfig.Log2File = file
	}
}

func WithLogDirPath(path string) Option {
	return func() {
		if len(path) > 0 {
			logConfig.Path = path
		}
	}
}

func InitLogConfig(opts ...Option) {
	for _, opt := range opts {
		opt()
	}

	// 目录不存在则先创建目录
	if _, err := os.Stat(logConfig.Path); os.IsNotExist(err) {
		err = os.MkdirAll(logConfig.Path, 0777)
		if err != nil {
			panic(fmt.Errorf("log conf err: create log dir '%s' error: %s", logConfig.Path, err))
		}
	}
}
