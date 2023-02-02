package log

import (
	"fmt"
	"os"
	"time"
)

const TimeFormat = "|2006-01-02T15:04:05.999Z07:00|"

type Timer interface {
	Time() time.Time
}

type LogContenter interface {
	fmt.Stringer
	Timer
}

type logger struct {
	file     *os.File
	filename string
	channel  chan LogContenter
}

func (l *logger) run() {
	for {
		contenter := <-l.channel
		if contenter == nil {
			continue
		}
		currentFilename := getLoggerFilename(l.filename)
		if l.file.Name() != currentFilename {
			f, err := os.OpenFile(currentFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
			if err != nil {
				println("create new log file", currentFilename, "error is", err.Error())
			} else {
				l.file = f
			}
		}
		l.file.WriteString(contenter.Time().Format(TimeFormat))
		l.file.Write([]byte{' '})
		l.file.WriteString(contenter.String())
		l.file.Write([]byte{'\n'})
	}
}

func NewLogger(filename string) *logger {
	f, err := os.OpenFile(getLoggerFilename(filename), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		panic(err)
	}
	l := &logger{f, filename, make(chan LogContenter)}
	go l.run()
	return l
}

func getLoggerFilename(filename string) string {
	return filename + time.Now().Format(".2006-01-02.log")
}

type LogString struct {
	time time.Time
	str  string
}

func (s *LogString) String() string {
	return s.str
}

func (s *LogString) Time() time.Time {
	return s.time
}

type LogFormat struct {
	time   time.Time
	format string
	args   []interface{}
}

func (s *LogFormat) String() string {
	return fmt.Sprintf(s.format, s.args...)
}

func (s *LogFormat) Time() time.Time {
	return s.time
}
