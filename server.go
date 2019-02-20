package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"
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
	index, ext := s.getJsIndexAndExt(r.URL)
	switch ext {
	case "jssp":
		s.ServeJssp(w, r, index)
	case "jsjs":
		s.ServeJsjs(w, r, index)
	default:
		s.static.ServeHTTP(w, r)
	}
}

func (s *JsspServer) ServeJssp(w http.ResponseWriter, r *http.Request, f *http.File) {
	w.Write([]byte("jssp"))
}

func (s *JsspServer) ServeJsjs(w http.ResponseWriter, r *http.Request, f *http.File) {
	w.Write([]byte("jsjs"))
}

func (s *JsspServer) getJsIndexAndExt(u *url.URL) (*http.File, string) {
	const JSSP = "index.jssp"
	const JSJS = "index.jsjs"
	if u.Path[len(u.Path)-1] == '/' {
		if !strings.HasPrefix(u.Path, "/") {
			u.Path = "/" + u.Path
		}
		f := getFile(s.root, u.Path+JSSP)
		if f != nil {
			return f, "jssp"
		}
		f = getFile(s.root, u.Path+JSJS)
		if f != nil {
			return f, "jsjs"
		}
	}
	return nil, ""
}

// run Jssp server
func (s *JsspServer) Run(paras *Parameter) {
	err := http.ListenAndServe(":"+paras.Port, s)
	if err != nil {
		log.Fatal(err)
	}
}
