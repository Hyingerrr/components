package logger

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

var logFilePath string

type Config struct {
	Output              string `yaml:"output"`
	LoggerLevel         string `yaml:"logger_level"`
	LoggerFile          string `yaml:"logger_file"`
	LogFormatText       bool   `yaml:"log_format_text"`
	DisableReportCaller bool   `yaml:"disable_report_caller"`
	LogRotateBool       bool   `yaml:"log_rotate_bool"`
	LogMaxSize          int    `yaml:"log_max_size"`
	LogMaxAge           int    `yaml:"log_max_age"`
	LogBackupCount      int    `yaml:"log_backup_count"`
	AccessLogOutput     string `yaml:"access_log_output"`
	AccessLogFile       string `yaml:"access_log_file"`
}

func Init(config *Config) {
	logs := NewLogClient(config)
	SetLogger(logs)

	logs.Info("init logger success")
}

func NewLogClient(opts *Config) Logger {
	var (
		output       = os.Stdout
		loggerLevel  logrus.Level
		reportCaller = true
		formatter    = JsonFormatter
	)

	if opts.LoggerLevel != "" {
		level, err := logrus.ParseLevel(opts.LoggerLevel)
		if err == nil {
			loggerLevel = level
		}
	}

	if opts.DisableReportCaller {
		reportCaller = false
	}

	if opts.Output == "file" {
		if opts.LoggerFile == "" {
			opts.LoggerFile = "~/logs/easy_pay.log"
		}

		if filepath.IsAbs(opts.LoggerFile) {
			CreateLogFile("", opts.LoggerFile)
			logFilePath = filepath.Join("", opts.LoggerFile)
		} else {
			CreateLogFile(os.Getenv("GOPATH"), opts.LoggerFile)
			logFilePath = filepath.Join(os.Getenv("GOPATH"), opts.LoggerFile)
		}

		file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			panic(err)
		}
		output = file
	}

	if opts.LogFormatText {
		formatter = TextFormatter
	}

	logs := NewWithOption(&Option{
		Output:       output,
		Level:        loggerLevel,
		Formatter:    formatter,
		ReportCaller: reportCaller,
	})

	// 滚动日志 切割
	if opts.LogRotateBool && opts.Output == "file" {
		logs.SetOutput(&lumberjack.Logger{
			Filename:   logFilePath,
			MaxSize:    opts.LogMaxSize,     // 最大文件 size, 默认 100MB
			MaxAge:     opts.LogMaxAge,      // 保留过期文件的最大时间间隔,单位是天(24h)
			MaxBackups: opts.LogBackupCount, // 最大过期日志保留的个数,默认都保留
			LocalTime:  false,               // 是否使用时间戳命名 backup 日志, 默认使用 UTC 格式
			Compress:   true,                // 是否压缩过期日志
		})
	}

	return logs
}

// 创建一个日志文件
func CreateLogFile(localPath, out string) {
	_, err := os.Stat(strings.Replace(filepath.Dir(filepath.Join(localPath, out)), "\\", "/", -1))
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(strings.Replace(filepath.Dir(filepath.Join(localPath, out)), "\\", "/", -1), os.ModePerm)
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	//f, err := os.OpenFile(strings.Replace(filepath.Join(localPath, out), "\\", "/", -1), os.O_CREATE, 0640)
	//if err != nil {
	//	panic(err)
	//}
	//defer f.Close()
}
