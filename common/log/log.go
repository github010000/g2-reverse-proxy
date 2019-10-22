package log

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func SetUp() error {
	logrus.SetOutput(os.Stdout)

	setLogLevel("debug")

	return nil
}

func setLogLevel(logLevel string) {
	switch strings.ToLower(logLevel) {
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func Info(args ...interface{}) {
	logrus.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func InfoWithFields(msg string, fields map[string]interface{}) {
	logrus.WithFields(fields).Info(msg)
}

func Debug(args ...interface{}) {
	logrus.Debug(args...)
}

func DebugWithFields(msg string, fields map[string]interface{}) {
	logrus.WithFields(fields).Debug(msg)
}

func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func Warn(args ...interface{}) {
	logrus.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

func Error(args ...interface{}) {
	logrus.Error(args...)
}

func ErrorWithFields(msg string, fields map[string]interface{}) {
	logrus.WithFields(fields).Error(msg)
}

func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}
