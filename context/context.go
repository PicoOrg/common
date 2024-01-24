package context

import (
	"context"

	"github.com/picoorg/common/logger"
)

type Context interface {
	SetValue(key any, value interface{}) (ctx Context)
	GetValue(key any) (value interface{})
	GetRawContext() (ctx context.Context)
	GetRawLogger() (logger logger.Logger)
	LogWithField(key string, value any) (ctx Context)
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
	Panic(msg string)
}
