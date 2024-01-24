package context

import (
	"context"

	"github.com/picoorg/common/logger"
)

func New(ctx context.Context, log logger.Logger) Context {
	return &implement{
		context: ctx,
		logger:  log,
	}
}

type implement struct {
	context context.Context
	logger  logger.Logger
}

func (m *implement) clone() *implement {
	return &implement{
		context: m.context,
		logger:  m.logger,
	}
}

func (m *implement) SetValue(key any, value interface{}) Context {
	ctx := m.clone()
	ctx.context = context.WithValue(ctx.context, key, value)
	return ctx
}

func (m *implement) GetValue(key any) interface{} {
	return m.context.Value(key)
}

func (m *implement) GetRawContext() context.Context {
	return m.context
}

func (m *implement) GetRawLogger() logger.Logger {
	return m.logger
}

func (m *implement) LogWithField(key string, value any) Context {
	ctx := m.clone()
	ctx.logger = m.logger.WithField(key, value)
	return ctx
}

func (m *implement) Debug(msg string) {
	m.logger.Debug(msg)
}

func (m *implement) Info(msg string) {
	m.logger.Info(msg)
}

func (m *implement) Warn(msg string) {
	m.logger.Warn(msg)
}

func (m *implement) Error(msg string) {
	m.logger.Error(msg)
}

func (m *implement) Fatal(msg string) {
	m.logger.Fatal(msg)
}

func (m *implement) Panic(msg string) {
	m.logger.Panic(msg)
}
