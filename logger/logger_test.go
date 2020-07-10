package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestCreateLogFile(t *testing.T) {
	CreateLogFile("/Users/hy/", "logger/app.log")
}

func TestInit(t *testing.T) {

	Init(&Config{
		Output:              "file",
		LoggerLevel:         "debug",
		LoggerFile:          "/tmp/app.log",
		LogFormatText:       false, // true=>text , false=>json
		DisableReportCaller: false,
		LogMaxSize:          1,
		LogMaxAge:           1,
		LogBackupCount:      2,
	})

	logInstance.Infof("test logger: %v", logFilePath)

	// 。。。

	timer := time.NewTicker(time.Second)
	timer2 := time.NewTicker(3600 * time.Second)

	for {
		select {
		case <-timer.C:
			logInstance.Infof("test logger: %v", logFilePath)
			logInstance.Errorf("test logger errorf: %v", "logFilePath")
			logInstance.Debug("test logger debug: %v", time.Now())
			logInstance.WithFields(logrus.Fields{"a": 1, "b": "111", "c": 1.1}).Warn("warning")
		case <-timer2.C:
			logInstance.Infof("Stop !!!")
			goto End
		}
	}

End:
	fmt.Println("task over")
}
