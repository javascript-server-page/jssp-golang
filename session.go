package main

import "net/http"

const SESSION_KEY = "JSSP-SESSION-ID"

var sessions = make(map[string]*Session)

type Session struct {
	data map[string]string
}

func GetSession(r *http.Request) *Session {
	c, err := r.Cookie(SESSION_KEY)
	if err != nil {
		c = &http.Cookie{Name: SESSION_KEY, Value: getUUID()}
		r.AddCookie(c)
	}
	s, ok := sessions[c.Value]
	if !ok {
		s = &Session{}
		sessions[c.Value] = s
	}
	return s
}
