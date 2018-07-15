package mylog

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
	"time"
)

var Logger *log.Logger

// 日志装饰器
func Log(fn func(w http.ResponseWriter, r *http.Request, param httprouter.Params)) func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		start := time.Now()
		Info(fmt.Sprintf("%s %s %s ", r.RemoteAddr, r.Method, r.URL.Path))
		fn(w, r, param)
		Info(fmt.Sprintf("%s Done in %v (%s %s)", r.RemoteAddr, time.Since(start), r.Method, r.URL.Path))
	}
}

func init() {
	logFile, err := os.Create("log.log")
	if err != nil {
		log.Fatalln("open log file error: ", err)
	}
	Logger = log.New(logFile, "[DEBUG]", log.LstdFlags)
}

func Debug(msg string) {
	Logger.SetPrefix("[DEBUG]")
	Logger.Println(msg)
}

func Info(msg string) {
	Logger.SetPrefix("[INFO]")
	Logger.Println(msg)
}

func Error(msg string) {
	Logger.SetPrefix("[ERROR]")
	Logger.Println(msg)
}
