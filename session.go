package main

import (
	"github.com/robertkrimen/otto"
	"net/http"
	"sync"
	"time"
)

const SESSION_KEY = "JSSP-SESSION-ID"

type Sessions struct {
	data    *sync.Map
	expired time.Duration
}

func NewSessions() *Sessions {
	return &Sessions{new(sync.Map), time.Hour}
}

func (ss *Sessions) NewSession(id string) *Session {
	return &Session{id, time.Now().Add(ss.expired), new(sync.Mutex), make(map[string]*otto.Value)}
}

func (ss *Sessions) GetSession(r *http.Request, w http.ResponseWriter) *Session {
	c, err := r.Cookie(SESSION_KEY)
	if err != nil {
		c = &http.Cookie{Name: SESSION_KEY, Value: getUUID()}
		http.SetCookie(w, c)
	}
	s, ok := ss.data.Load(c.Value)
	if !ok || s.(*Session).isExpired() {
		s = ss.NewSession(c.Value)
		ss.data.Store(c.Value, s)
	}
	return s.(*Session)
}

type Session struct {
	id    string
	et    time.Time
	mutex *sync.Mutex
	data  map[string]*otto.Value
}

func (s *Session) isExpired() bool {
	return s.et.After(time.Now())
}
