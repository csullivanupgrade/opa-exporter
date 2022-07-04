// Package log provides methods for creating the logger
package log

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ContextKey is the context key type.
type ContextKey string

const (
	// ContextLogger is the context logger key.
	ContextLogger = ContextKey("logger")
)

func setLevel(level string) (zapcore.Level, bool) {
	switch level {
	case "debug":
		return zap.DebugLevel, false
	case "info":
		return zap.InfoLevel, false
	case "warn":
		return zap.WarnLevel, false
	case "error":
		return zap.ErrorLevel, false
	case "dpanic":
		return zap.DPanicLevel, false
	case "panic":
		return zap.PanicLevel, false
	case "fatal":
		return zap.FatalLevel, false
	}
	return zap.DebugLevel, true
}

// NewLogger creates a zap logger based on the mode you specify, defaulting to production if no mode was matched.
func NewLogger(level string, mode string) *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "file",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	zapLevel, warnBadLevelConfig := setLevel(level)
	atomicLevel := zap.NewAtomicLevelAt(zapLevel)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		atomicLevel,
	)

	var logger *zap.Logger
	switch mode {
	case "production":
		logger = zap.New(core)
	case "nop":
		logger = zap.NewNop()
	case "development":
		caller := zap.AddCaller()
		development := zap.Development()
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
			atomicLevel,
		)
		logger = zap.New(core, caller, development)
	default:
		logger = zap.New(core)
		logger.Warn("log mode not specified - defaulting to production")
	}

	if warnBadLevelConfig {
		logger.Warn("log level not specified - defaulting to info")
	}

	logger.Info("logger successfully initialized")
	return logger
}

// SetContext sets the given logger in a new child context.
func SetContext(ctx context.Context, contextLogger *zap.Logger) context.Context {
	return context.WithValue(ctx, ContextLogger, contextLogger)
}

// FromContext gets the logger from the given context.
func FromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(ContextLogger).(*zap.Logger); ok {
		return logger
	}

	logger := NewLogger("info", "production")
	logger.Warn("logger could not be retrieved from context")
	return logger
}
