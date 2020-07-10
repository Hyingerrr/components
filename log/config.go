package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

var (
	defaultConfig = &Config{
		OutFile:      true,
		Level:        "INFO",
		Format:       "json",
		ReportCaller: true,
		Stacktrace:   true,
		File:         "/",
		MaxSize:      1000,
		MaxAge:       15,
		BackupCount:  20,
		Compress:     true,
	}
)

type Config struct {
	OutFile      bool   `yaml:"log_out_file"` // 日志的位置  是否在文件，true 文件 false stdout
	Level        string `yaml:"log_level"`
	Format       string `yaml:"log_format"`
	ReportCaller bool   `yaml:"log_report_caller"`
	Stacktrace   bool   `yaml:"log_stack_trace"`
	File         string `yaml:"log_file"`
	MaxSize      int    `yaml:"log_max_size"`     // 单个文件最大size
	MaxAge       int    `yaml:"log_max_age"`      // 保留旧文件的最大天数
	BackupCount  int    `yaml:"log_backup_count"` // 保留旧文件的最大个数
	Compress     bool   `yaml:"log_compress"`     // 是否压缩/归档旧文件
}

func ParseLevel(lvl string) zapcore.Level {
	switch strings.ToLower(lvl) {
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	case "error":
		return zap.ErrorLevel
	case "warn", "warning":
		return zap.WarnLevel
	case "info":
		return zap.InfoLevel
	case "debug":
		return zap.DebugLevel
	default:
		return zap.InfoLevel
	}
}

func GetDefaultConfig() *Config {
	var (
		dir = getCurrentPath()
		options = config.ViperConfOptions{}
	)

	conf := config.NewViperConfig(options.WithConfigType("yaml"),
		options.WithConfFile([]string{dir + "/conf/conf.yaml"}))
	return &Config{
		OutFile:      conf.GetBool("log_out_file"),
		Level:        conf.GetString("log_level"),
		Format:       conf.GetString("log_format"),
		ReportCaller: conf.GetBool("log_report_caller"),
		Stacktrace:   conf.GetBool("log_stack_trace"),
		File:         conf.GetString("log_file"),
		MaxSize:      conf.GetInt("log_max_size"),
		MaxAge:       conf.GetInt("log_max_age"),
		BackupCount:  conf.GetInt("log_backup_count"),
		Compress:     conf.GetBool("log_compress"),
	}
}

// 项目的目录
func getCurrentPath() string {
	s, err := os.Getwd()
	if err == nil && strings.Contains(s, "repaychnl"){
		i := strings.SplitAfterN(s, "repaychnl", 2)
		return i[0]
	}
	return "."
}