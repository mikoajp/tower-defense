package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var L *zap.SugaredLogger

func Init(level string) error {
	cfg := zap.NewProductionConfig()
	cfg.Encoding = "json"
	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	switch level {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	logger, err := cfg.Build()
	if err != nil { return err }
	L = logger.Sugar()
	return nil
}

func Sync() { if L != nil { _ = L.Sync() } }

func Debugw(msg string, keysAndValues ...interface{}) { if L != nil { L.Debugw(msg, keysAndValues...) } }
func Infow(msg string, keysAndValues ...interface{})   { if L != nil { L.Infow(msg, keysAndValues...) } }
func Warnw(msg string, keysAndValues ...interface{})  { if L != nil { L.Warnw(msg, keysAndValues...) } }
func Errorw(msg string, keysAndValues ...interface{}) { if L != nil { L.Errorw(msg, keysAndValues...) } }
