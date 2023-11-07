package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Logger    *logrus.Logger
	GinLogger *logrus.Logger
)

// CustomFormatter 自定义日志格式器
type CustomFormatter struct {
	logrus.TextFormatter
}

// Format 格式化日志
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	level := strings.ToUpper(entry.Level.String())
	msg := entry.Message
	if entry.HasCaller() {
		file := fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		fileList := strings.Split(file, "/")
		file = fmt.Sprintf("%s/%s", fileList[len(fileList)-2], fileList[len(fileList)-1])
		msg = fmt.Sprintf("%s | %s", file, msg)
	}
	return []byte(fmt.Sprintf("%s | %8s | %s\n", timestamp, level, msg)), nil
}

func InitLogger(logLevel string, logPath string) {
	Logger = logrus.New()

	logFile := &lumberjack.Logger{
		Filename:   generateLogFilePath(logPath, "vmq"),
		MaxSize:    10, // MB
		MaxBackups: 10,
		MaxAge:     30, // days
	}

	Logger.SetFormatter(&CustomFormatter{
		logrus.TextFormatter{
			DisableColors: true,
		},
	})

	Logger.SetReportCaller(true)

	// 设置输出到终端和文件
	mw := io.MultiWriter(os.Stdout, logFile)
	Logger.SetOutput(mw)

	// 设置日志级别
	switch strings.ToLower(logLevel) {
	case "debug":
		Logger.SetLevel(logrus.DebugLevel)
	case "info":
		Logger.SetLevel(logrus.InfoLevel)
	case "warn":
		Logger.SetLevel(logrus.WarnLevel)
	case "error":
		Logger.SetLevel(logrus.ErrorLevel)
	default:
		Logger.SetLevel(logrus.InfoLevel)
	}

	GinLogger = logrus.New()
	logFile = &lumberjack.Logger{
		Filename:   generateLogFilePath(logPath, "gin"),
		MaxSize:    10, // MB
		MaxBackups: 10,
		MaxAge:     30, // days
	}
	mw = io.MultiWriter(os.Stdout, logFile)
	GinLogger.SetFormatter(&CustomFormatter{})

	GinLogger.SetLevel(logrus.InfoLevel)
	GinLogger.SetOutput(mw)
	switch strings.ToLower(logLevel) {
	case "debug":
		GinLogger.SetLevel(logrus.DebugLevel)
	case "info":
		GinLogger.SetLevel(logrus.InfoLevel)
	case "warn":
		GinLogger.SetLevel(logrus.WarnLevel)
	case "error":
		GinLogger.SetLevel(logrus.ErrorLevel)
	default:
		GinLogger.SetLevel(logrus.InfoLevel)
	}

	Logger.Debug("Logger init success!")
}

func generateLogFilePath(logsDir string, pefix ...string) string {
	if len(pefix) == 0 {
		pefix = append(pefix, "vmq")
	}
	err := os.MkdirAll(logsDir, os.ModePerm)
	if err != nil {
		panic("Create logs directory failed!")
	}
	currentTime := time.Now()
	logFileName := fmt.Sprintf("%s-%s.log", pefix, currentTime.Format("2006-01-02_15-04-05"))
	return path.Join(logsDir, logFileName)
}
