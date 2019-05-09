package main

import (
	"fmt"
	"net/http"
	"time"
)

type RequestLog struct {
	Start      time.Time     // req start time
	Method     string        // req http method
	Path       string        // req path
	Duration   time.Duration // req duration
	RemoteAddr string        // client ip addr
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

func (l *Logging) run() {
	for {
		rl := <-l.requestLogChannel
		if rl == nil {
			continue
		}
		fmt.Println(*rl)
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
