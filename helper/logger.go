package helper

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func NewLogger() *logrus.Logger {
	log := logrus.New()
	log.SetReportCaller(true)
	runPath, _ := os.Getwd()

	fmt.Println("os.Getwd()", runPath)
	logFilename := filepath.Join(runPath, "logs", "logrus.log")
	if NewFile(logFilename) {
		file, err := os.OpenFile(logFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("初始化日志文件失败：", err.Error())
			log.Info("Failed to log to file, using default stderr")
		}
		log.Out = file
	}

	log.Formatter = &logrus.JSONFormatter{
		DisableTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			_, filename := path.Split(f.File)
			return funcname, filename
		},
	}
	log.Info("example of custom format caller")
	return log
}
