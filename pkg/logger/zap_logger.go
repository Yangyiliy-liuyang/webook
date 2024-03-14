package logger

import "go.uber.org/zap"

type ZapLogger struct {
	l *zap.Logger
}

func NewZapLogger(l *zap.Logger) Logger {
	return &ZapLogger{
		l: l,
	}
}
func (z *ZapLogger) Debug(msg string, args ...Field) {
	z.l.Debug(msg, z.toZapField(args)...)
}

func (z *ZapLogger) Info(msg string, args ...Field) {
	z.l.Info(msg, z.toZapField(args)...)
}

func (z *ZapLogger) Warn(msg string, args ...Field) {
	z.l.Warn(msg, z.toZapField(args)...)
}

func (z *ZapLogger) Error(msg string, args ...Field) {
	z.l.Error(msg, z.toZapField(args)...)
}

func (z *ZapLogger) toZapField(args []Field) []zap.Field {
	res := make([]zap.Field, len(args))
	for _, arg := range args {
		res = append(res, zap.Any(arg.Key, arg.Value))
	}
	return res
}
