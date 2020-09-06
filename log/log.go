package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger represents the struct of logger
type Logger struct {
	*zap.SugaredLogger
	Level zapcore.Level
}

// NewLogger create an instance of zap SugarLogger with custom config
func NewLogger(opts ...*Option) *Logger {
	cfg := unifyConfig(opts...)

	log, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return &Logger{
		log.Sugar(),
		cfg.Level.Level(),
	}
}

// Log implements go-kit logger
func (sl Logger) Log(kv ...interface{}) error {
	level := sl.Level
	forGokitLog(sl.SugaredLogger, level)
	return nil
}

// FlushLogger encapsule Sync function
func (sl *Logger) FlushLogger() {
	if sl != nil {
		_ = sl.Sync()
	}
}

// OriginLogger represents the struct of origin logger
type OriginLogger struct {
	*zap.Logger
	Level zapcore.Level
}

// NewOriginLogger create an instance of zap logger with custom config
func NewOriginLogger(opts ...*Option) *OriginLogger {
	cfg := unifyConfig(opts...)

	log, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return &OriginLogger{
		log,
		cfg.Level.Level(),
	}
}

// Log implements go-kit logger
func (l OriginLogger) Log(kv ...interface{}) error {
	sugarLogger := l.Sugar()
	level := l.Level
	forGokitLog(sugarLogger, level)
	return nil
}

// FlushLogger encapsule Sync function
func (l *OriginLogger) FlushLogger() {
	if l != nil {
		_ = l.Sync()
	}
}

func forGokitLog(sugarLogger *zap.SugaredLogger, level zapcore.Level, kv ...interface{}) {
	switch level {
	case zapcore.DebugLevel:
		sugarLogger.Debugw("", kv...)
	case zapcore.InfoLevel:
		sugarLogger.Infow("", kv...)
	case zapcore.WarnLevel:
		sugarLogger.Warnw("", kv...)
	case zapcore.ErrorLevel:
		sugarLogger.Errorw("", kv...)
	case zapcore.DPanicLevel:
		sugarLogger.DPanicw("", kv...)
	case zapcore.PanicLevel:
		sugarLogger.Panicw("", kv...)
	case zapcore.FatalLevel:
		sugarLogger.Fatalw("", kv...)
	default:
		sugarLogger.Infow("", kv...)
	}
}
