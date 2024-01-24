package mocks

import (
	"log"

	"github.com/picoorg/common/logger"
)

func NewLogger() logger.Logger {
	return &loggerModel{}
}

type loggerModel struct {
}

func (l *loggerModel) WithField(name string, value any) logger.Logger {
	return l
}

func (l *loggerModel) Debug(msg string) {
	log.Println(msg)
}

func (l *loggerModel) Info(msg string) {
	log.Println(msg)
}

func (l *loggerModel) Warn(msg string) {
	log.Println(msg)
}

func (l *loggerModel) Error(msg string) {
	log.Println(msg)
}

func (l *loggerModel) Fatal(msg string) {
	log.Fatalln(msg)
}

func (l *loggerModel) Panic(msg string) {
	log.Panicln(msg)
}
