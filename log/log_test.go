package log

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	fakeConfig = &Config{
		OutFile:      false,
		Level:        "debug",
		Format:       "json",
		ReportCaller: false,
		Stacktrace:   false,
		File:         "/Users/hy/logs/zap_test.log",
		MaxSize:      2,
		MaxAge:       2,
		BackupCount:  3,
		Compress:     true,
	}
)

type Test struct {
	A string
	B int
	C map[string]string
}

func TestNewLogger(t *testing.T) {
	var (
		it = assert.New(t)
		l Logger
	)

	it.NotPanics(func() {
		l = NewLogger(WithLogConf(fakeConfig))
	})

	// with map[string]interface{} and msg
	l.WithFields(Field{"x":123, "y":345, "z": "hkhkh"}).Info("WithFields")

	// with msg
	tx := Test{
		A: "aaa",
		B: 123,
		C: map[string]string{"D": "888", "E": "999"},
	}
	l.Info("Info ...")

	// with msg and variable
	l.Infof("infof: %+v", tx)

	// with msg and []interface{}
	l.InfoW("infoW", []interface{}{"baz", false, "xxx", tx}...)

	//l.Errorfo("msgsss", zap.Int("a", 111), zap.String("b", "hhhhh"))
	l.Errorf("msg: %+v", 423545)

	//l.Debugf("print debug log ...")

	//l.Errorf("print error log ..., time[%v]", 333)
	//
	//l.Info("print logs ...", zap.Reflect("aaaa", map[string]interface{}{
	//			"aaa":1,
	//			"bbb":2,
	//			"nnn":"789",
	//		}))

	// test 切割
	goto End
	{
		timer := time.NewTicker(time.Millisecond*10)
		timer2 := time.NewTicker(3600 * time.Second)

		for {
			select {
			case <-timer.C:
				l.Infof("test logger: %v", 666)
				l.Errorf("test logger errorf: %v", "logFilePath")
				l.Debugf("test logger debug: %v", time.Now())
			case <-timer2.C:
				l.Infof("Stop !!!")
				goto End
			}
		}
	}


End:
	fmt.Println("task over")
}

func TestGetDefaultConfig(t *testing.T) {
	c := GetDefaultConfig()
	fmt.Printf("%+v\n", c)
}