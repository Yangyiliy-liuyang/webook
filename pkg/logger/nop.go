package logger

import "go.uber.org/zap"

type NopLogger struct {
	Logger zap.Logger
}

func NewNopLogger() *NopLogger {
	return &NopLogger{}
}
func (z *NopLogger) Debug(msg string, args ...Field) {

}

func (z *NopLogger) Info(msg string, args ...Field) {

}

func (z *NopLogger) Warn(msg string, args ...Field) {

}

func (z *NopLogger) Error(msg string, args ...Field) {

}
