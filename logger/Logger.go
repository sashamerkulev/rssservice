package logger

import (
	"fmt"
	"time"
)

type Logger interface {
	Log(severity string, tag string, message string)
}

type ConsoleLogger struct {
}

func (cl ConsoleLogger) Log(severity string, tag string, message string) {
	go ConsoleLog(severity, "", "", tag, message)
}

func ConsoleLog(severity string, userId string, userIp string, tag string, message string) {
	fmt.Println(time.Now().Format(time.RFC3339), severity, tag, "USER: "+userId+"("+userIp+")", message)
}
