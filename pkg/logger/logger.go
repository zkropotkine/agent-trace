package logger

import (
	"context"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/zkropotkine/agent-trace/config"
)

// Logger defines a lightweight interface for structured logging.
type Logger interface {
	WithField(key string, value interface{}) *logrus.Entry
	WithFields(fields logrus.Fields) *logrus.Entry
	WithError(err error) *logrus.Entry

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})
}

// NewLogRusLogger initializes a new Logger based on config.Log.
func NewLogRusLogger(cfg config.Log) *logrus.Entry {
	logger := logrus.New()

	lvl, err := logrus.ParseLevel(strings.ToLower(cfg.Level))
	if err != nil {
		lvl = logrus.InfoLevel
	}
	logger.SetLevel(lvl)

	switch strings.ToLower(cfg.Format) {
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{
			//DisableColors: true,
			FullTimestamp:   true,
			ForceColors:     true, // forces ANSI colors
			DisableQuote:    true, // cleaner output
			TimestampFormat: "2006-01-02 15:04:05",
		})
	default:
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	return logrus.NewEntry(logger)
}

// DefaultLogger fallback logger instance
var DefaultLogger = NewLogRusLogger(config.Log{Level: "info", Format: "json"})

// ctxKeyType is an unexported type used as the key for storing loggers in context.
type ctxKeyType struct{}

var ctxKey = ctxKeyType{}

// WithLogger stores a Logger in the provided context.
func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, ctxKey, logger)
}

// FromContext retrieves the Logger from context, falling back to defaultLogger.
func FromContext(ctx context.Context) Logger {
	if log, ok := ctx.Value(ctxKey).(Logger); ok {
		return log
	}

	return DefaultLogger
}
