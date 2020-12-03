package main

import (
	"fmt"
	"runtime"
	"time"
)

var Log *Logger

type LogData struct {
	logType uint
	file    string
	time    string
	log     string
}

type Logger struct {
	Database *DataBase
}

func setupLogger(db *DataBase) {
	Log = &Logger{
		Database: db,
	}
}

func (l *Logger) Info(log string) {
	_, file, line, _ := runtime.Caller(1)
	l.log(LogData{
		logType: 0,
		file:    fmt.Sprintf("%s:%d", file, line),
		time:    time.Now().Format("2006-01-02 15:04:05"),
		log:     log,
	})
}

func (l *Logger) Warning(log string) {
	_, file, line, _ := runtime.Caller(1)
	l.log(LogData{
		logType: 1,
		file:    fmt.Sprintf("%s:%d", file, line),
		time:    time.Now().Format("2006-01-02 15:04:05"),
		log:     log,
	})
}

func (l *Logger) Error(log string) {
	_, file, line, _ := runtime.Caller(1)
	l.log(LogData{
		logType: 2,
		file:    fmt.Sprintf("%s:%d", file, line),
		time:    time.Now().Format("2006-01-02 15:04:05"),
		log:     log,
	})
}

func (l *Logger) FatalError(log string) {
	_, file, line, _ := runtime.Caller(1)
	l.log(LogData{
		logType: 3,
		file:    fmt.Sprintf("%s:%d", file, line),
		time:    time.Now().Format("2006-01-02 15:04:05"),
		log:     log,
	})
}

func (l *Logger) log(log LogData) {
	err := l.Database.createLog(log)
	if err != nil {
		fmt.Errorf(log.time, log.file, log.log)
	}
}
