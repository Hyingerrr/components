package logger

import (
	"io"

	"github.com/sirupsen/logrus"
)

var (
	LogClient = New()
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
	LogClient = l
}

// todo  需要其他方法时，可以添加

func SetOutput(output io.Writer) {
	LogClient.SetOutput(output)
}

func Info(args ...interface{}) {
	LogClient.Info(args...)
}

func Infof(format string, args ...interface{}) {
	LogClient.Infof(format, args...)
}

func Warn(args ...interface{}) {
	LogClient.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	LogClient.Warnf(format, args...)
}

func Error(args ...interface{}) {
	LogClient.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	LogClient.Errorf(format, args...)
}

func Panic(args ...interface{}) {
	LogClient.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	LogClient.Panicf(format, args...)
}

func WithField(key string, value interface{}) Logger {
	return LogClient.WithField(key, value)
}

func WithFields(fields logrus.Fields) Logger {
	return LogClient.WithFields(fields)
}

func WithError(err error) Logger {
	return LogClient.WithError(err)
}

func Debug(args ...interface{}) {
	LogClient.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	LogClient.Debugf(format, args...)
}