package main

import (
	"github.com/robertkrimen/otto"
	"net/http"
	"sync"
	"time"
)

const SESSION_KEY = "JSSP-SESSION-ID"

type Sessions struct {
	mutex *sync.Mutex
	data  map[string]*Session
}

func NewSessions() *Sessions {
	return &Sessions{new(sync.Mutex), make(map[string]*Session)}
}

func (ss *Sessions) NewSession(id string) *Session {
	return &Session{id, time.Now(), new(sync.Mutex), make(map[string]*otto.Value)}
}

func (ss *Sessions) GetSession(r *http.Request) *Session {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()
	c, err := r.Cookie(SESSION_KEY)
	if err != nil {
		c = &http.Cookie{Name: SESSION_KEY, Value: getUUID()}
		r.AddCookie(c)
	}
	s, ok := ss.data[c.Value]
	if !ok {
		s = ss.NewSession(c.Value)
		ss.data[c.Value] = s
	}
	return s
}

type Session struct {
	id    string
	ct    time.Time
	mutex *sync.Mutex
	data  map[string]*otto.Value
}

func (s *Session) isExpired() bool {
	return false
}
