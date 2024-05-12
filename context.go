package log

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func Debugf(c context.Context, format string, args ...any) {
	Logger(c).Debugf(format, args...)
}

func Errorf(c context.Context, format string, args ...any) {
	Logger(c).Errorf(format, args...)
}

func Fatalf(c context.Context, format string, args ...any) {
	Logger(c).Fatalf(format, args...)
}

func Infof(c context.Context, format string, args ...any) {
	Logger(c).Infof(format, args...)
}

func Panicf(c context.Context, format string, args ...any) {
	Logger(c).Panicf(format, args...)
}

func Warnf(c context.Context, format string, args ...any) {
	Logger(c).Warnf(format, args...)
}

func Debug(c context.Context, args ...any) {
	Logger(c).Debug(args...)
}

func Error(c context.Context, args ...any) {
	Logger(c).Error(args...)
}

func Fatal(c context.Context, args ...any) {
	Logger(c).Fatal(args...)
}

func Info(c context.Context, args ...any) {
	Logger(c).Info(args...)
}

func Panic(c context.Context, args ...any) {
	Logger(c).Panic(args...)
}

func Warn(c context.Context, args ...any) {
	Logger(c).Warn(args...)
}

func IsLevelEnabled(c context.Context, level logrus.Level) bool {
	switch l := Logger(c).(type) {
	case *logrus.Logger:
		return l.IsLevelEnabled(level)
	case *logrus.Entry:
		return l.Logger.IsLevelEnabled(level)
	default:
		return true
	}
}

type logKey struct{}

func Logger(c context.Context) logrus.FieldLogger {
	if l, ok := c.Value(logKey{}).(logrus.FieldLogger); ok {
		return l
	}
	return logrus.StandardLogger()
}

func WithField(c context.Context, key string, value any) context.Context {
	return WithLogger(c, Logger(c).WithField(key, value))
}

func WithFields(c context.Context, fields logrus.Fields) context.Context {
	return WithLogger(c, Logger(c).WithFields(fields))
}

func WithPath(c context.Context, path string) context.Context {
	return WithField(c, pathKey, path)
}

func WithLogger(c context.Context, l logrus.FieldLogger) context.Context {
	return context.WithValue(c, logKey{}, l)
}

// NewForwarder creates a new logger that logs unformatted output to the given io.Writer.
func NewForwarder(out io.Writer, level logrus.Level) logrus.FieldLogger {
	return &logrus.Logger{
		Out:       out,
		Formatter: PlainFormatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     level,
		ExitFunc:  os.Exit,
	}
}

// NewLogger creates a new logger that logs its output to a file named by logFileName or to stdout, in case
// the logFileName is the empty string.
func NewLogger(logFileName string) (logrus.FieldLogger, context.CancelFunc, error) {
	var out io.Writer
	var cancel context.CancelFunc
	if logFileName != "" {
		err := os.MkdirAll(filepath.Dir(logFileName), 0755)
		if err != nil {
			return nil, nil, err
		}
		f, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, nil, err
		}
		cancel = func() {
			_ = f.Close()
		}
		out = f
	} else {
		out = os.Stdout
		cancel = func() {}
	}
	return &logrus.Logger{
		Out:          out,
		Formatter:    Formatter("15:04:05.0000"),
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.DebugLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}, cancel, nil
}
