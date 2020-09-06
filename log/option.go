package log

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Option represents the option of log instance
type Option struct {
	DisableStacktrace *bool  `json:"disableStacktrace"`
	LogLevel          string `json:"logLevel"`
	LogFile           string `json:"logFile"`
}

func unifyConfig(opts ...*Option) *zap.Config {
	var (
		level, file       string
		disableStacktrace bool
	)
	cfg := zap.NewProductionConfig()

	// Get DisableStacktrace
	if opts != nil && opts[0].DisableStacktrace != nil {
		disableStacktrace = *opts[0].DisableStacktrace
	} else {
		disableStacktrace = viper.GetBool("LOG_DISABLE_STACKTRACE")
	}

	// Get log level
	if opts != nil && opts[0].LogLevel != "" {
		level = opts[0].LogLevel
	} else {
		level = viper.GetString("LOG_LEVEL")
		if level == "" {
			level = viper.GetString("logLevel")
		}
	}

	// Get log file
	if opts != nil && opts[0].LogFile != "" {
		file = opts[0].LogFile
	} else {
		file = viper.GetString("LOG_FILE")
		if file == "" {
			file = viper.GetString("logFile")
		}
	}

	// Set DisableStacktrace
	cfg.DisableStacktrace = disableStacktrace
	// Set log level
	cfg.Level.SetLevel(convertLogLevel(level))
	// Set output path
	var paths = []string{
		"stdout",
	}
	if file != "" {
		paths = append(paths, file)
	}
	cfg.OutputPaths = paths

	// Set encoder
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.EncodeTime = timeEncoder
	cfg.EncoderConfig.EncodeDuration = milliSecondsDurationEncoder

	return &cfg
}

func convertLogLevel(level string) zapcore.Level {
	m := map[string]zapcore.Level{
		"debug":   zap.DebugLevel,
		"info":    zap.InfoLevel,
		"warning": zap.WarnLevel,
		"error":   zap.ErrorLevel,
		"dpanic":  zap.DPanicLevel,
		"panic":   zap.PanicLevel,
		"fatal":   zap.FatalLevel,
	}
	if v, ok := m[level]; ok {
		return v
	}
	// If level string is invalid, return debug level for default
	return zap.DebugLevel
}
