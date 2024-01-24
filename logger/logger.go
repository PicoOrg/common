package logger

type Logger interface {
	Debug(msg string)
	Error(msg string)
	Fatal(msg string)
	Info(msg string)
	Panic(msg string)
	Warn(msg string)
	WithField(key string, value any) (logger Logger)
}
