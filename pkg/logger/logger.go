package logger

import (
	"context"
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	globalLogger *zap.SugaredLogger
	defaultLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
)

func init() {
	// Initialize the global logger with the default error level
	SetLogger(New(defaultLevel))
}

// New returns a new *zap.SugaredLogger instance that outputs JSON logs.
// If no logging level is provided, it defaults to zap.ErrorLevel.
func New(level zapcore.LevelEnabler, opts ...zap.Option) *zap.SugaredLogger {
	return NewWithOutput(level, os.Stdout, opts...)
}

// NewWithOutput returns a new *zap.SugaredLogger instance with JSON encoding,
// writing logs to the given sink (io.Writer). If level is nil, the default level is used.
func NewWithOutput(level zapcore.LevelEnabler, sink io.Writer, opts ...zap.Option) *zap.SugaredLogger {
	if level == nil {
		level = defaultLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(sink),
		level,
	)

	return zap.New(core, opts...).Sugar()
}

// SetLogger assigns the provided logger as the global logger instance.
// Note: This function is not safe for concurrent use.
func SetLogger(loggerInstance *zap.SugaredLogger) {
	globalLogger = loggerInstance
}

// Debug logs a message at Debug level.
func Debug(ctx context.Context, args ...interface{}) {
	globalLogger.Debug(args...)
}

// Debugf logs a formatted message at Debug level.
func Debugf(ctx context.Context, format string, args ...interface{}) {
	globalLogger.Debugf(format, args...)
}

// DebugKV logs a message with key-value pairs at Debug level.
func DebugKV(ctx context.Context, msg string, kvs ...interface{}) {
	globalLogger.Debugw(msg, kvs...)
}

// Info logs a message at Info level.
func Info(ctx context.Context, args ...interface{}) {
	globalLogger.Info(args...)
}

// Infof logs a formatted message at Info level.
func Infof(ctx context.Context, format string, args ...interface{}) {
	globalLogger.Infof(format, args...)
}

// InfoKV logs a message with key-value pairs at Info level.
func InfoKV(ctx context.Context, msg string, kvs ...interface{}) {
	globalLogger.Infow(msg, kvs...)
}

// Warn logs a message at Warn level.
func Warn(ctx context.Context, args ...interface{}) {
	globalLogger.Warn(args...)
}

// Warnf logs a formatted message at Warn level.
func Warnf(ctx context.Context, format string, args ...interface{}) {
	globalLogger.Warnf(format, args...)
}

// WarnKV logs a message with key-value pairs at Warn level.
func WarnKV(ctx context.Context, msg string, kvs ...interface{}) {
	globalLogger.Warnw(msg, kvs...)
}

// Error logs a message at Error level.
func Error(ctx context.Context, args ...interface{}) {
	globalLogger.Error(args...)
}

// Errorf logs a formatted message at Error level.
func Errorf(ctx context.Context, format string, args ...interface{}) {
	globalLogger.Errorf(format, args...)
}

// ErrorKV logs a message with key-value pairs at Error level.
func ErrorKV(ctx context.Context, msg string, kvs ...interface{}) {
	globalLogger.Errorw(msg, kvs...)
}

// Fatal logs a message at Fatal level and then terminates the program.
func Fatal(ctx context.Context, args ...interface{}) {
	globalLogger.Fatal(args...)
}

// Fatalf logs a formatted message at Fatal level and then terminates the program.
func Fatalf(ctx context.Context, format string, args ...interface{}) {
	globalLogger.Fatalf(format, args...)
}

// FatalKV logs a message with key-value pairs at Fatal level and then terminates the program.
func FatalKV(ctx context.Context, msg string, kvs ...interface{}) {
	globalLogger.Fatalw(msg, kvs...)
}

// Panic logs a message at Panic level and then panics.
func Panic(ctx context.Context, args ...interface{}) {
	globalLogger.Panic(args...)
}

// Panicf logs a formatted message at Panic level and then panics.
func Panicf(ctx context.Context, format string, args ...interface{}) {
	globalLogger.Panicf(format, args...)
}

// PanicKV logs a message with key-value pairs at Panic level and then panics.
func PanicKV(ctx context.Context, msg string, kvs ...interface{}) {
	globalLogger.Panicw(msg, kvs...)
}
