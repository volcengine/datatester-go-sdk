/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package log

import "github.com/sirupsen/logrus"

type loggerAdaptor struct {
	l *logrus.Logger
}

func (a loggerAdaptor) WithField(key string, value interface{}) Logger {
	return newFieldAdapt(a.l.WithField(key, value))
}

func (a loggerAdaptor) WithFields(fields Fields) Logger {
	return newFieldAdapt(a.l.WithFields(logrus.Fields(fields)))
}

func (a loggerAdaptor) Trace(args ...interface{}) {
	a.l.Trace(args...)
}

func (a loggerAdaptor) TraceF(format string, args ...interface{}) {
	a.l.Tracef(format, args...)
}

func (a loggerAdaptor) Debug(args ...interface{}) {
	a.l.Debug(args...)
}

func (a loggerAdaptor) DebugF(format string, args ...interface{}) {
	a.l.Debugf(format, args...)
}

func (a loggerAdaptor) Info(args ...interface{}) {
	a.l.Info(args...)
}

func (a loggerAdaptor) InfoF(format string, args ...interface{}) {
	a.l.Infof(format, args...)
}

func (a loggerAdaptor) Warn(args ...interface{}) {
	a.l.Warn(args...)
}

func (a loggerAdaptor) WarnF(format string, args ...interface{}) {
	a.l.Warnf(format, args...)
}

func (a loggerAdaptor) Error(args ...interface{}) {
	a.l.Error(args...)
}

func (a loggerAdaptor) ErrorF(format string, args ...interface{}) {
	a.l.Errorf(format, args...)
}

func (a loggerAdaptor) Panic(args ...interface{}) {
	a.l.Panic(args...)
}

func (a loggerAdaptor) PanicF(format string, args ...interface{}) {
	a.l.Panicf(format, args...)
}

func (a loggerAdaptor) Fatal(args ...interface{}) {
	a.l.Fatal(args...)
}

func (a loggerAdaptor) FatalF(format string, args ...interface{}) {
	a.l.Fatalf(format, args...)
}

type fieldAdapt struct {
	e *logrus.Entry
}

func (f fieldAdapt) WithField(key string, value interface{}) Logger {
	return newFieldAdapt(f.e.WithField(key, value))
}

func (f fieldAdapt) WithFields(fields Fields) Logger {
	return newFieldAdapt(f.e.WithFields(logrus.Fields(fields)))
}

func (f fieldAdapt) Trace(args ...interface{}) {
	f.e.Trace(args...)
}

func (f fieldAdapt) TraceF(format string, args ...interface{}) {
	f.e.Tracef(format, args...)
}

func (f fieldAdapt) Debug(args ...interface{}) {
	f.e.Debug(args...)
}

func (f fieldAdapt) DebugF(format string, args ...interface{}) {
	f.e.Debugf(format, args...)
}

func (f fieldAdapt) Info(args ...interface{}) {
	f.e.Info(args...)
}

func (f fieldAdapt) InfoF(format string, args ...interface{}) {
	f.e.Infof(format, args...)
}

func (f fieldAdapt) Warn(args ...interface{}) {
	f.e.Warn(args...)
}

func (f fieldAdapt) WarnF(format string, args ...interface{}) {
	f.e.Warnf(format, args...)
}

func (f fieldAdapt) Error(args ...interface{}) {
	f.e.Error(args...)
}

func (f fieldAdapt) ErrorF(format string, args ...interface{}) {
	f.e.Errorf(format, args...)
}

func (f fieldAdapt) Panic(args ...interface{}) {
	f.e.Panic(args...)
}

func (f fieldAdapt) PanicF(format string, args ...interface{}) {
	f.e.Panicf(format, args...)
}

func (f fieldAdapt) Fatal(args ...interface{}) {
	f.e.Fatal(args...)
}

func (f fieldAdapt) FatalF(format string, args ...interface{}) {
	f.e.Fatalf(format, args...)
}

func newFieldAdapt(e *logrus.Entry) Logger {
	return fieldAdapt{e}
}

func NewLogrusAdapt(l *logrus.Logger) Logger {
	return &loggerAdaptor{
		l: l,
	}
}
