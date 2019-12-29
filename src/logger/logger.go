package logger

import (
	"io"
	"log"
	"os"
	"sync"
)

type logger struct {
	multiWriter io.Writer
	fpLog       *os.File
	err         error
}

var instance *logger
var once sync.Once

func Logger() *logger {
	once.Do(func() {
		instance = &logger{}
		instance.fpLog, instance.err = os.OpenFile("application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if instance.err != nil {
			panic(instance.err)
		}
		instance.multiWriter = io.MultiWriter(instance.fpLog, os.Stdout)
		log.SetOutput(instance.multiWriter)
		log.SetPrefix("Logger: ")
	})
	return instance
}

func (e *logger) Close() {
	e.fpLog.Close()
}

func Info(format string, v ...interface{}) {
	log.Printf(format, v)
}

func Error(format string, v ...interface{}) {
	log.Fatalf(format, v)
}
