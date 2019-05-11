package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

type RequestLog struct {
	Start      time.Time     // req start time
	Method     string        // req http method
	Path       string        // req path
	Duration   time.Duration // req duration
	RemoteAddr string        // client ip addr
}

func (rl *RequestLog) String() string {
	t := rl.Start.Format("2006-01-02 15:04:05")
	return fmt.Sprintf("%s Comleted %s %s in %v from %s\n", t, rl.Method, rl.Path, rl.Duration, rl.RemoteAddr)
}

type Logging struct {
	logFile           string
	requestLogChannel chan *RequestLog
}

// get *Logging
func NewLogging(logFile string) *Logging {
	l := &Logging{logFile, make(chan *RequestLog)}
	go l.run()
	return l
}

// get log output stream
func (l *Logging) getOutputStream() *os.File {
	var f *os.File
	var err error
	if !fileExists(l.logFile) {
		f, err = os.Create(l.logFile)
	} else {
		f, err = os.OpenFile(l.logFile, os.O_WRONLY|os.O_APPEND, 0644)
	}
	if err != nil {
		println(err.Error())
		f = os.Stderr
	}
	return f
}

func (l *Logging) run() {
	out := l.getOutputStream()
	for {
		rl := <-l.requestLogChannel
		if rl == nil {
			continue
		}
		out.WriteString(rl.String())
	}
}

// http request log
func (l *Logging) RequestLogHandlerFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		l.requestLogChannel <- &RequestLog{
			Start:      start,
			Method:     r.Method,
			Path:       r.URL.Path,
			Duration:   time.Since(start),
			RemoteAddr: r.RemoteAddr,
		}
	}
}
