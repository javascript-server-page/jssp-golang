package main

import (
	"log"
	"net/http"
	"path/filepath"
)

type JsspServer struct {
	http.ServeMux
	static http.Handler
	root   http.FileSystem
}

// init JsspServer
func (s *JsspServer) Init(paras *Parameter) {
	s.root = http.Dir(paras.Dir)
	s.static = http.FileServer(s.root)
	s.HandleFunc("/", s.ServeAll)
}

// handler func
func (s *JsspServer) ServeAll(w http.ResponseWriter, r *http.Request) {
	switch filepath.Ext(r.URL.Path) {
	case "jssp":
		s.ServeJssp(w, r)
	case "jsjs":
		s.ServeJsjs(w, r)
	default:
		s.ServeStatic(w, r)
	}
}

func (s *JsspServer) ServeStatic(w http.ResponseWriter, r *http.Request) {
	s.static.ServeHTTP(w, r)
}

func (s *JsspServer) ServeJssp(w http.ResponseWriter, r *http.Request) {
}

func (s *JsspServer) ServeJsjs(w http.ResponseWriter, r *http.Request) {
}

// run Jssp server
func (s *JsspServer) Run(paras *Parameter) {
	err := http.ListenAndServe(":"+paras.Port, s)
	if err != nil {
		log.Fatal(err)
	}
}
