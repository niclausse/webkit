package logx

import (
	"io"

	"github.com/sirupsen/logrus"
	"github.com/uniplaces/carbon"
)

func logrusWithWriter(writers ...io.Writer) *logrus.Logger {

	logger := logrus.StandardLogger()

	logger.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true, TimestampFormat: carbon.DefaultFormat})
	logger.SetOutput(io.MultiWriter(writers...))

	return logger

}

func logrusEntryWithWriter(writers ...io.Writer) *logrus.Entry {

	logger := logrus.StandardLogger()

	logger.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true, TimestampFormat: carbon.DefaultFormat})
	logger.SetOutput(io.MultiWriter(writers...))

	entry := logrus.NewEntry(logger)

	return entry

}
