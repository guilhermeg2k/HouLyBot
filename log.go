package main

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

type Logger struct {
	Database *DataBase
}

func setupLogger(db *DataBase) {
	logger = &Logger{
		Database: db,
	}
}

func (l *Logger) Info(log string) {
	_, file, line, _ := runtime.Caller(1)
	l.log(Log{
		logType: 0,
		file:    fmt.Sprintf("%s:%d", file, line),
		time:    time.Now().Format("2006-01-02 15:04:05"),
		text:    log,
	})
}

func (l *Logger) Warning(log string) {
	_, file, line, _ := runtime.Caller(1)
	l.log(Log{
		logType: 1,
		file:    fmt.Sprintf("%s:%d", file, line),
		time:    time.Now().Format("2006-01-02 15:04:05"),
		text:    log,
	})
}

func (l *Logger) Error(log string) {
	_, file, line, _ := runtime.Caller(1)
	l.log(Log{
		logType: 2,
		file:    fmt.Sprintf("%s:%d", file, line),
		time:    time.Now().Format("2006-01-02 15:04:05"),
		text:    log,
	})
}

func (l *Logger) FatalError(log string) {
	_, file, line, _ := runtime.Caller(1)
	l.log(Log{
		logType: 3,
		file:    fmt.Sprintf("%s:%d", file, line),
		time:    time.Now().Format("2006-01-02 15:04:05"),
		text:    log,
	})
	panic(0)
}

func (l *Logger) log(logData Log) {
	err := l.Database.createLog(logData)
	if err != nil {
		log.Println("Failed to create a database log: " + err.Error())
		log.Println(logData.time, logData.file, logData.text)
	}
}
