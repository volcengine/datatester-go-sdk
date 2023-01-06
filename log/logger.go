/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package log

import "github.com/sirupsen/logrus"

type Fields map[string]interface{}

type Logger interface {
	WithField(key string, value interface{}) Logger
	WithFields(fields Fields) Logger
	Trace(args ...interface{})
	TraceF(format string, args ...interface{})
	Debug(args ...interface{})
	DebugF(format string, args ...interface{})
	Info(args ...interface{})
	InfoF(format string, args ...interface{})
	Warn(args ...interface{})
	WarnF(format string, args ...interface{})
	Error(args ...interface{})
	ErrorF(format string, args ...interface{})
	Panic(args ...interface{})
	PanicF(format string, args ...interface{})
	Fatal(args ...interface{})
	FatalF(format string, args ...interface{})
}

var logger Logger

func WithField(key string, value interface{}) Logger {
	return logger.WithField(key, value)
}
func WithFields(fields Fields) Logger {
	return logger.WithFields(fields)
}
func Trace(args ...interface{}) {
	logger.Trace(args...)
}
func TraceF(format string, args ...interface{}) {
	logger.TraceF(format, args...)
}
func Debug(args ...interface{}) {
	logger.Debug(args...)
}
func DebugF(format string, args ...interface{}) {
	logger.DebugF(format, args...)
}
func Info(args ...interface{}) {
	logger.Info(args...)
}
func InfoF(format string, args ...interface{}) {
	logger.InfoF(format, args...)
}
func Warn(args ...interface{}) {
	logger.Warn(args...)
}
func WarnF(format string, args ...interface{}) {
	logger.WarnF(format, args...)
}
func Error(args ...interface{}) {
	logger.Error(args...)
}
func ErrorF(format string, args ...interface{}) {
	logger.ErrorF(format, args...)
}
func Panic(args ...interface{}) {
	logger.Panic(args...)
}
func PanicF(format string, args ...interface{}) {
	logger.PanicF(format, args...)
}
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}
func Fatalf(format string, args ...interface{}) {
	logger.FatalF(format, args...)
}

func InitGlobalLogger(l Logger) {
	logger = l
}

func InitDefaultLogger() {
	if logger != nil {
		return
	}
	defaultLogger := logrus.New()
	defaultLogger.SetLevel(logrus.WarnLevel)
	logger = NewLogrusAdapt(defaultLogger)
}
