package logger

import (
	"io"
	"log"
	"os"
	"sync"
)

type logger struct {
	InfoWriter  io.Writer
	ErrorWriter io.Writer
	InfoLog     *os.File
	ErrorLog    *os.File
}

var instance *logger
var once sync.Once

func Logger() *logger {
	once.Do(func() {
		instance = &logger{}
		var err error

		log.SetFlags(log.Ldate)
		log.SetFlags(log.Ltime)
		log.SetFlags(log.Lmicroseconds)

		instance.InfoLog, err = os.OpenFile("application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		instance.ErrorLog, err = os.OpenFile("error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		instance.InfoWriter = io.MultiWriter(instance.InfoLog, os.Stdout)
		instance.ErrorWriter = io.MultiWriter(instance.ErrorLog, os.Stdout)

	})
	return instance
}

func Close() {
	instance.InfoLog.Close()
	instance.ErrorLog.Close()
}

func Info(format string, v ...interface{}) {
	log.SetOutput(instance.InfoWriter)
	log.SetPrefix("Info ")
	log.Printf(format, v...)
}

func Error(format string, v ...interface{}) {
	log.SetOutput(instance.ErrorWriter)
	log.SetPrefix("Error ")
	log.Printf(format, v...)
}
