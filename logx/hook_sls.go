package logx

import (
	"io"
	"log"
	"sync"

	"github.com/sirupsen/logrus"
)

// logging to sls service, strip colors to make the output more readable.
var defaultFormatter = &logrus.TextFormatter{DisableColors: true}

// sLSHook is a hook to handle writing to aliyun sls log service.
type sLSHook struct {
	Writer io.Writer

	lock *sync.Mutex

	formatter logrus.Formatter

	defaultWriter io.Writer

	hasDefaultWriter bool
}

func newSLSHook(writer io.Writer, formatter logrus.Formatter) *sLSHook {
	log.SetOutput(writer)
	hook := &sLSHook{
		lock: new(sync.Mutex),
	}

	hook.SetFormatter(formatter)
	hook.SetDefaultWriter(writer)

	return hook
}

// SetFormatter sets the format that will be used by hook.
// If using text formatter, this method will disable color output to make the log file more readable.
func (hook *sLSHook) SetFormatter(formatter logrus.Formatter) {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	if formatter == nil {
		formatter = defaultFormatter
	} else {
		switch formatter.(type) {
		case *logrus.TextFormatter:
			textFormatter := formatter.(*logrus.TextFormatter)
			textFormatter.DisableColors = true
		}
	}

	hook.formatter = formatter
}

// SetDefaultWriter sets default writer for levels that don't have any defined writer.
func (hook *sLSHook) SetDefaultWriter(defaultWriter io.Writer) {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	hook.defaultWriter = defaultWriter
	hook.hasDefaultWriter = true
}

// Fire writes the log file to defined path or using the defined writer.
// User who run this function needs write permissions to the file or directory if the file does not yet exist.
func (hook *sLSHook) Fire(entry *logrus.Entry) error {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	if hook.Writer != nil || hook.hasDefaultWriter {
		return hook.ioWrite(entry)
	}

	return nil
}

// Write a log line to an io.Writer.
func (hook *sLSHook) ioWrite(entry *logrus.Entry) error {
	var (
		writer io.Writer
		msg    []byte
		err    error
	)

	writer = hook.Writer

	msg, err = hook.formatter.Format(entry)

	if err != nil {
		log.Printf("failed to generate string for entry: %v", err)
		return err
	}
	_, err = writer.Write(msg)
	return err
}

// Levels returns configured log levels.
func (hook *sLSHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
