package log

import (
	"jssp/config"
	"time"
)

var (
	_info   *logger
	_error  *logger
	_access *logger
)

func init() {
	if config.Log.Info != "false" {
		_info = NewLogger(config.Log.Info)
	}
	if config.Log.Error != "false" {
		_error = NewLogger(config.Log.Error)
	}
	if config.Log.Access != "false" {
		_access = NewLogger(config.Log.Access)
	}
}

func Info(format string, a ...interface{}) {
	if _info != nil {
		_info.channel <- &LogFormat{time.Now(), format, a}
	}
}

func Error(format string, a ...interface{}) {
	if _error != nil {
		_error.channel <- &LogFormat{time.Now(), format, a}
	}
}

func InfoS(str string) {
	if _info != nil {
		_info.channel <- &LogString{time.Now(), str}
	}
}

func ErrorS(str string) {
	if _error != nil {
		_error.channel <- &LogString{time.Now(), str}
	}
}

func SInfo(logContenter LogContenter) {
	if _info != nil {
		_info.channel <- logContenter
	}
}

func SError(logContenter LogContenter) {
	if _error != nil {
		_error.channel <- logContenter
	}
}

func Access(Start time.Time, Method, Path, RemoteAddr string) {
	if _access != nil {
		_access.channel <- &LogFormat{Start, "[%s] <%s> in (%v) from {%s}", []interface{}{Method, Path, time.Since(Start), RemoteAddr}}
	}
}
