package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	JsonFormatter = "json"
	TextFormatter = "text"
	errorKey      = "error"

	defaultOption = &Option{
		Output:           os.Stdout,
		Level:            logrus.InfoLevel,
		Formatter:        JsonFormatter,
		EnableHTMLEscape: false,
		ReportCaller:     true,
	}
)

type log struct {
	entry        *logrus.Entry
	depth        int // 上报的深度
	reportCaller bool
}

type Option struct {
	Output           io.Writer
	Level            logrus.Level
	Formatter        string
	EnableHTMLEscape bool
	ReportCaller     bool // 是否上报堆栈
}

// 创建默认的日志记录器
func New() Logger {
	return NewWithOption(defaultOption)
}

// 创建日志记录器
func NewWithOption(option *Option) *log {
	if option == nil {
		option = defaultOption
	}

	logger := logrus.New()

	// set level
	logger.SetLevel(option.Level)

	// set formatter
	if option.Formatter == JsonFormatter {
		logger.SetFormatter(&logrus.JSONFormatter{DisableHTMLEscape: !option.EnableHTMLEscape})
	}

	// set output
	logger.SetOutput(option.Output)

	// set no lock
	logger.SetNoLock()

	return &log{
		entry:        logrus.NewEntry(logger),
		depth:        1,
		reportCaller: option.ReportCaller,
	}
}

// 通往logrus的入口
func (l *log) log(level logrus.Level, args ...interface{}) {
	entry := l.entry
	if l.reportCaller {
		entry = l.entry.WithField("file", caller(l.depth+3))
	}
	entry.Log(level, args...)
}

// 设置日志等级
func (l *log) SetLevel(level logrus.Level) {
	l.entry.Logger.SetLevel(level)
}

func (l *log) Log(level logrus.Level, arg ...interface{}) {
	l.log(level, arg...)
}

func (l *log) Logf(level logrus.Level, format string, arg ...interface{}) {
	l.log(level, fmt.Sprintf(format, arg...))
}

// 记录一条 LevelInfo 级别的日志
func (l *log) Info(args ...interface{}) {
	l.Log(logrus.InfoLevel, args...)
}

// 格式化并记录一条 LevelInfo 级别的日志
func (l *log) Infof(format string, args ...interface{}) {
	l.Logf(logrus.InfoLevel, format, args...)
}

// 记录一条 LevelWarn 级别的日志
func (l *log) Warn(args ...interface{}) {
	l.Log(logrus.WarnLevel, args...)
}

// 格式化并记录一条 LevelWarn 级别的日志
func (l *log) Warnf(format string, args ...interface{}) {
	l.Logf(logrus.WarnLevel, format, args...)
}

// 记录一条 LevelError 级别的日志
func (l *log) Error(args ...interface{}) {
	l.Log(logrus.ErrorLevel, args...)
}

// 格式化并记录一条 LevelError 级别的日志
func (l *log) Errorf(format string, args ...interface{}) {
	l.Logf(logrus.ErrorLevel, format, args...)
}

// 记录一条 LevelPanic 级别的日志
func (l *log) Panic(args ...interface{}) {
	l.Log(logrus.PanicLevel, args...)
	panic(fmt.Sprint(args...))
}

// 格式化并记录一条 LevelPanic 级别的日志
func (l *log) Panicf(format string, args ...interface{}) {
	l.Logf(logrus.PanicLevel, format, args...)
	panic(fmt.Sprintf(format, args...))
}

// 记录一条 LevelDebug 级别的日志
func (l *log) Debug(args ...interface{}) {
	l.Log(logrus.DebugLevel, args...)
}

// 格式化并记录一条 LevelDebug 级别的日志
func (l *log) Debugf(format string, args ...interface{}) {
	l.Logf(logrus.DebugLevel, format, args...)
}

// 为当前日志附加一个上下文数据
func (l *log) WithField(key string, value interface{}) Logger {
	return l.WithFields(logrus.Fields{key: value})
}

// 为当前日志附加一组上下文数据
func (l *log) WithFields(fields logrus.Fields) Logger {
	if l.reportCaller {
		if err, ok := fields[errorKey].(interface {
			Stack() []string
		}); ok {
			fields["err.stack"] = strings.Join(err.Stack(), ";")
		}
	}
	return &log{entry: l.entry.WithFields(fields), reportCaller: l.reportCaller}
}

// 为当前日志附加一个错误
func (l *log) WithError(err error) Logger {
	return l.WithFields(logrus.Fields{errorKey: err})
}

func (l *log) SetOutput(output io.Writer) {
	l.entry.Logger.SetOutput(output)
}

func (l *log) AddHook(hook logrus.Hook) {
	l.entry.Logger.AddHook(hook)
}

// 兼容grpc log，暂无实际运用
func (l *log) Fatal(args ...interface{})                 {}
func (l *log) Fatalln(args ...interface{})               {}
func (l *log) Warning(args ...interface{})               {}
func (l *log) Fatalf(format string, args ...interface{}) {}
func (l *log) Warningln(args ...interface{}) {
}
func (l *log) Infoln(args ...interface{}) {
}

func (l *log) Errorln(args ...interface{}) {
}

func (l *log) Warningf(format string, args ...interface{}) {
	l.Logf(logrus.WarnLevel, format, args...)
}
func (l *log) V(n int) bool {
	return true
}

// caller的显示形式为 File:Line
func caller(depth int) string {
	_, f, n, ok := runtime.Caller(1 + depth)
	if !ok {
		return ""
	}
	if ok {
		idx := strings.LastIndex(f, "easyPay") // 不显示项目关键目录
		if idx >= 0 {
			f = f[idx+18:]
		}
	}
	return fmt.Sprintf("%s:%d", f, n)
}
