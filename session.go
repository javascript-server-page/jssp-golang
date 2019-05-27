package main

import (
	"github.com/robertkrimen/otto"
	"net/http"
	"sync"
)

const SESSION_KEY = "JSSP-SESSION-ID"

var sessions_mutex = new(sync.Mutex)
var sessions = make(map[string]*Session)

type Session struct {
	id    string
	mutex *sync.Mutex
	data  map[string]*otto.Value
}

func NewSession(id string) *Session {
	return &Session{id, new(sync.Mutex), make(map[string]*otto.Value)}
}

func GetSession(r *http.Request) *Session {
	sessions_mutex.Lock()
	defer sessions_mutex.Unlock()
	c, err := r.Cookie(SESSION_KEY)
	if err != nil {
		c = &http.Cookie{Name: SESSION_KEY, Value: getUUID()}
		r.AddCookie(c)
	}
	s, ok := sessions[c.Value]
	if !ok {
		s = NewSession(c.Value)
		sessions[c.Value] = s
	}
	return s
}
