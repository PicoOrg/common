package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	debugFile = "data/log/debug.log"
	infoFile  = "data/log/info.log"
	errorFile = "data/log/error.log"
)

func New() Logger {
	return &implement{
		zapLogger: newZapLogger(1),
	}
}

func NewCtxLogger() Logger {
	return &implement{
		zapLogger: newZapLogger(2),
	}
}

func NewWithSkip(skip int) Logger {
	return &implement{
		zapLogger: newZapLogger(skip),
	}
}

type implement struct {
	zapLogger *zap.Logger
}

func (l *implement) WithField(name string, value any) Logger {
	return &implement{
		zapLogger: l.zapLogger.With(zap.Any(name, value)),
	}
}

func (l *implement) Debug(msg string) {
	l.zapLogger.Debug(msg)
}

func (l *implement) Info(msg string) {
	l.zapLogger.Info(msg)
}

func (l *implement) Warn(msg string) {
	l.zapLogger.Warn(msg)
}

func (l *implement) Error(msg string) {
	l.zapLogger.Error(msg)
}

func (l *implement) Fatal(msg string) {
	l.zapLogger.Fatal(msg)
}

func (l *implement) Panic(msg string) {
	l.zapLogger.Panic(msg)
}

func newZapLogger(addCallerSkip int) *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	makeLoggerFile(debugFile)
	debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})
	debugLumberJackLogger := zapcore.AddSync(&lumberjack.Logger{
		Filename:   debugFile,
		MaxSize:    32, // megabytes
		MaxBackups: 30,
		MaxAge:     30,   // days
		Compress:   true, // disabled by default
	})

	makeLoggerFile(infoFile)
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
	infoLumberJackLogger := zapcore.AddSync(&lumberjack.Logger{
		Filename:   infoFile,
		MaxSize:    32, // megabytes
		MaxBackups: 30,
		MaxAge:     30,   // days
		Compress:   true, // disabled by default
	})

	makeLoggerFile(errorFile)
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	errorLumberJackLogger := zapcore.AddSync(&lumberjack.Logger{
		Filename:   errorFile,
		MaxSize:    32, // megabytes
		MaxBackups: 30,
		MaxAge:     30,   // days
		Compress:   true, // disabled by default
	})

	debugConsoleLogger := zapcore.Lock(&wrapStdout{})

	cores := []zapcore.Core{
		zapcore.NewCore(encoder, debugLumberJackLogger, debugLevel),
		zapcore.NewCore(encoder, infoLumberJackLogger, infoLevel),
		zapcore.NewCore(encoder, errorLumberJackLogger, errorLevel),
		zapcore.NewCore(encoder, debugConsoleLogger, debugLevel),
	}
	return zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel), zap.AddCallerSkip(addCallerSkip))
}

type wrapStdout struct {
}

func (w *wrapStdout) Write(p []byte) (n int, err error) {
	os.Stdout.Write([]byte("\033[1K\r"))
	return os.Stdout.Write(p)
}

func (w *wrapStdout) Sync() error {
	return os.Stdout.Sync()
}

func makeLoggerFile(name string) {
	if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
		panic(fmt.Errorf("failed to create log directory: %w", err))
	}
}
