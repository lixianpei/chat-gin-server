package helper

import (
	"GoChatServer/consts"
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

// NewLogger 实例化Logger
func NewLogger() *logrus.Logger {
	log := logrus.New()
	log.SetReportCaller(true)
	log.SetFormatter(&MyFormatter{})
	log.AddHook(&CustomLogFile{})
	return log
}

// MyFormatter ================= 自定义日志内容格式 =================
type MyFormatter struct{}

func (m *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Local().Format("2006-01-02 15:04:05.000")
	var newLog string

	//HasCaller()为true才会有调用信息
	if entry.HasCaller() {
		fName := filepath.Base(entry.Caller.File)
		newLog = fmt.Sprintf("[%s][%s][%s:%d %s][%s]\n",
			timestamp, entry.Level, fName, entry.Caller.Line, entry.Caller.Function, entry.Message)
	} else {
		newLog = fmt.Sprintf("[%s][%s][%s]\n", timestamp, entry.Level, entry.Message)
	}

	b.WriteString(newLog)
	return b.Bytes(), nil
}

// CustomLogFile ================= 自定义文件 =================
type CustomLogFile struct{}

func (hook *CustomLogFile) Fire(entry *logrus.Entry) error {
	entry.Logger.Out = logFileOut()
	return nil
}
func (hook *CustomLogFile) Levels() []logrus.Level {
	return logrus.AllLevels
}

func logFileOut() (file *os.File) {
	runPath, _ := os.Getwd()
	dateDay := time.Now().Local().Format(consts.DateYMD)
	logFilename := filepath.Join(runPath, "logs", dateDay+".log")
	if NewFile(logFilename) {
		var err error
		file, err = os.OpenFile(logFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("初始化日志文件失败：", err.Error())
		}
	}
	return
}
