package log

import "github.com/sirupsen/logrus"

var logger = logrus.New()

type Fields logrus.Fields

func SetLogLevel(level logrus.Level)  {
	logger.SetLevel(level)
}

func SetLogFormatter(formatter logrus.Formatter)  {
	logger.Formatter = formatter
}

// Debug
func Debug(args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Debug(args)
	}
}
// 带有field的Debug
func DebugWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Debug(l)
	}
}
// Info
func Info(args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Info(args...)
	}
}
// 带有field的Info
func InfoWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Info(l)
	}
}
// Warn
func Warn(args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Warn(args...)
	}
}
// 带有Field的Warn
func WarnWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Warn(l)
	}
}
// Error
func Error(args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Error(args...)
	}
}
// 带有Fields的Error
func ErrorWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.ErrorLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Error(l)
	}
}
// Fatal
func Fatal(args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Fatal(args...)
	}
}
// 带有Field的Fatal
func FatalWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.FatalLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Fatal(l)
	}
}
// Panic
func Panic(args ...interface{}) {
	if logger.Level >= logrus.PanicLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Panic(args...)
	}
}
// 带有Field的Panic
func PanicWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.PanicLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Panic(l)
	}
}