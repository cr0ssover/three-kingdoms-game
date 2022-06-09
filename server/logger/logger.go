package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var errorLogger *zap.SugaredLogger

func init() {
	encoder := zapcore.NewConsoleEncoder(
		zapcore.EncoderConfig{
			MessageKey:  "msg",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalColorLevelEncoder,
			TimeKey:     "ts",
			EncodeTime: func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
				pae.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
			},
			CallerKey:      "caller_line",
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
		},
	)

	level := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level), //打印到控制台
	)
	log := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	errorLogger = log.Sugar()
}

func Debug(arg ...interface{}) {
	errorLogger.Debug(arg...)
}

func Info(arg ...interface{}) {
	errorLogger.Info(arg...)
}

func Warn(arg ...interface{}) {
	errorLogger.Warn(arg...)
}

func Error(arg ...interface{}) {
	errorLogger.Error(arg...)
}

func Panic(arg ...interface{}) {
	errorLogger.Panic(arg...)
}

func Debugf(template string, arg ...interface{}) {
	errorLogger.Debugf(template, arg...)
}

func Infof(template string, arg ...interface{}) {
	errorLogger.Infof(template, arg...)
}

func Warnf(template string, arg ...interface{}) {
	errorLogger.Warnf(template, arg...)
}

func Errorf(template string, arg ...interface{}) {
	errorLogger.Errorf(template, arg...)
}

func Panicf(template string, arg ...interface{}) {
	errorLogger.Panicf(template, arg...)
}

func ZapEngine() *zap.SugaredLogger {
	return errorLogger
}
