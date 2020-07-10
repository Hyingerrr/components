package logger

import (
	"io"

	"github.com/sirupsen/logrus"
)

var (
	logInstance = New()
)

type Logger interface {
	SetOutput(output io.Writer)                // 设置输出
	SetLevel(level logrus.Level)               // 设置log等级 	// 获取log等级
	AddHook(hook logrus.Hook)                  // 添加hook
	Debug(args ...interface{})                 // 记录 DebugLevel 级别的日志
	Debugf(format string, args ...interface{}) // 格式化并记录 DebugLevel 级别的日志
	Info(args ...interface{})                  // 记录 InfoLevel 级别的日志
	Infof(format string, args ...interface{})  // 格式化并记录 InfoLevel 级别的日志
	Warn(args ...interface{})                  // 记录 WarnLevel 级别的日志
	Warnf(format string, args ...interface{})  // 格式化并记录 WarnLevel 级别的日志
	Fatalf(format string, args ...interface{})
	Error(args ...interface{})                      // 记录 ErrorLevel 级别的日志
	Errorf(format string, args ...interface{})      // 格式化并记录 ErrorLevel 级别的日志
	Panic(args ...interface{})                      // 记录 PanicLevel 级别的日志
	Panicf(format string, args ...interface{})      // 格式化并记录 PanicLevel 级别的日志
	WithField(key string, value interface{}) Logger // 为日志添加一个上下文数据
	WithFields(fields logrus.Fields) Logger         // 为日志添加多个上下文数据
	WithError(err error) Logger                     // 为日志添加标准错误上下文数据

	// 兼容grpc logger；暂时没有实现
	Warningln(args ...interface{})
	Warningf(format string, args ...interface{})
	Warning(args ...interface{})
	Fatal(args ...interface{})
	Fatalln(args ...interface{})
	Infoln(args ...interface{})
	Errorln(args ...interface{})
	V(l int) bool
}

func SetLogger(l Logger) {
	logInstance = l
}

// todo  需要其他方法时，可以添加

func SetOutput(output io.Writer) {
	logInstance.SetOutput(output)
}

func Info(args ...interface{}) {
	logInstance.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logInstance.Infof(format, args...)
}

func Warn(args ...interface{}) {
	logInstance.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	logInstance.Warnf(format, args...)
}

func Error(args ...interface{}) {
	logInstance.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	logInstance.Errorf(format, args...)
}

func Panic(args ...interface{}) {
	logInstance.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	logInstance.Panicf(format, args...)
}

func WithField(key string, value interface{}) Logger {
	return logInstance.WithField(key, value)
}

func WithFields(fields logrus.Fields) Logger {
	return logInstance.WithFields(fields)
}

func WithError(err error) Logger {
	return logInstance.WithError(err)
}

func Debug(args ...interface{}) {
	logInstance.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	logInstance.Debugf(format, args...)
}
