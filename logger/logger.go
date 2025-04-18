package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()
	
	customTimeFormat := "2006/01/02 15:04:05"
	consoleFormatter := &logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: customTimeFormat,
	}
	
	fileFormatter := &logrus.JSONFormatter{
		TimestampFormat: customTimeFormat,
	}

	logRotation := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10,    
		MaxBackups: 30,    
		MaxAge:     28,    
		Compress:   true,  
	}

	
	mw := io.MultiWriter(os.Stdout, logRotation)
	
	
	Log.SetOutput(mw)
	
	Log.SetFormatter(consoleFormatter)
	
	Log.AddHook(&fileHook{
		writer:    logRotation,
		formatter: fileFormatter,
	})

	Log.SetLevel(logrus.InfoLevel)
}


type fileHook struct {
	writer    io.Writer
	formatter logrus.Formatter
}

func (h *fileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *fileHook) Fire(entry *logrus.Entry) error {
	line, err := h.formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = h.writer.Write(line)
	return err
}

func GetLogger() *logrus.Logger {
	return Log
}