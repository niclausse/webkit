package resource

import (
	"context"
	"errors"
	"fmt"
	"github.com/niclausse/webkit/utils"
	"github.com/niclausse/webkit/zlog"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	ormUtil "gorm.io/gorm/utils"
	"time"
)

type ormLogger struct {
	Service  string
	Addr     string
	Database string
	logger   *zlog.Logger
}

func newLogger(serviceName, addr, database string) *ormLogger {
	if serviceName == "" {
		serviceName = database
	}

	return &ormLogger{
		Service:  serviceName,
		Addr:     addr,
		Database: database,
		logger:   zlog.GetZapLogger().WithOptions(zlog.AddCallerSkip(2)),
	}
}

// Print just for mysql go-sql-driver error log
func (l *ormLogger) Print(args ...interface{}) {
	l.logger.Error(fmt.Sprint(args...), l.commonFields(nil)...)
}

// LogMode log mode
func (l *ormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

// Info print info
func (l ormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	m := fmt.Sprintf(msg, append([]interface{}{ormUtil.FileWithLineNum()}, data...)...)
	// 非trace日志改为debug级别输出
	l.logger.Debug(m, l.commonFields(ctx)...)
}

// Warn print warn messages
func (l ormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	m := fmt.Sprintf(msg, append([]interface{}{ormUtil.FileWithLineNum()}, data...)...)
	l.logger.Warn(m, l.commonFields(ctx)...)
}

// Error print error messages
func (l ormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	m := fmt.Sprintf(msg, append([]interface{}{ormUtil.FileWithLineNum()}, data...)...)
	l.logger.Error(m, l.commonFields(ctx)...)
}

// Trace print sql message
func (l ormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	end := time.Now()
	elapsed := end.Sub(begin)
	cost := float64(elapsed.Nanoseconds()/1e4) / 100.0

	// 请求是否成功
	msg := "mysql do success"
	ralCode := -0
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// 没有找到记录不统计在请求错误中
		msg = err.Error()
		ralCode = -1
	}

	sql, rows := fc()
	fileLineNum := ormUtil.FileWithLineNum()
	fields := l.commonFields(ctx)
	fields = append(fields,
		zlog.Int64("affectedRow", rows),
		zlog.String("requestEndTime", utils.GetFormatRequestTime(end)),
		zlog.String("requestStartTime", utils.GetFormatRequestTime(begin)),
		zlog.String("fileLine", fileLineNum),
		zlog.Float64("cost", cost),
		zlog.Int("ralCode", ralCode),
		zlog.String("sql", sql),
	)

	l.logger.Info(msg, fields...)
}

func (l ormLogger) commonFields(ctx context.Context) []zlog.Field {
	requestId, _ := ctx.Value(zlog.ContextKeyRequestID).(string)
	uri, _ := ctx.Value(zlog.ContextKeyURI).(string)
	fields := []zlog.Field{
		zlog.String("requestId", requestId),
		zlog.String("uri", uri),
		zlog.String("service", l.Service),
		zlog.String("addr", l.Addr),
		zlog.String("db", l.Database),
	}
	return fields
}
