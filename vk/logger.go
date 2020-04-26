package vk

import (
	"log"
	"os"
)

type DefaultLogType struct{}

var defaultLogger = log.New(os.Stderr, "", log.LstdFlags)

type Logger interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
}

func (l *DefaultLogType) Error(args ...interface{}) {
	defaultLogger.SetPrefix("[Error] ")
	defaultLogger.Println(args...)
}

func (l *DefaultLogType) Errorf(format string, args ...interface{}) {
	defaultLogger.SetPrefix("[Error] ")
	defaultLogger.Printf(format, args...)
}

func (l *DefaultLogType) Warn(args ...interface{}) {
	defaultLogger.SetPrefix("[Warn] ")
	defaultLogger.Println(args...)
}

func (l *DefaultLogType) Warnf(format string, args ...interface{}) {
	defaultLogger.SetPrefix("[Warn] ")
	defaultLogger.Printf(format, args...)
}

func (l *DefaultLogType) Info(args ...interface{}) {
	defaultLogger.SetPrefix("[Info] ")
	defaultLogger.Println(args...)
}

func (l *DefaultLogType) Infof(format string, args ...interface{}) {
	defaultLogger.SetPrefix("[Info] ")
	defaultLogger.Printf(format, args...)
}

func (l *DefaultLogType) Debug(args ...interface{}) {
	defaultLogger.SetPrefix("[Debug] ")
	defaultLogger.Println(args...)
}

func (l *DefaultLogType) Debugf(format string, args ...interface{}) {
	defaultLogger.SetPrefix("[Debug] ")
	defaultLogger.Printf(format, args...)
}
