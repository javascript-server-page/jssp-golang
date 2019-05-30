package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const Version = "0.1"
const Server = "jssp-golang-" + Version

const JSSP = ".jssp"
const JSJS = ".jsjs"
const INDEX_JSSP = "index" + JSSP
const INDEX_JSJS = "index" + JSJS

type JsspServer struct {
	http.ServeMux
	eng    *Engine
	log    *Logging
	set    *Setting
	ses    *Sessions
	static http.Handler
	root   http.FileSystem
}

// init JsspServer
func (s *JsspServer) Init() {
	s.set = new(Setting)
	s.set.Init()
	s.eng = NewEngine()
	s.ses = NewSessions()
	s.root = http.Dir(s.set.Dir)
	s.log = NewLogging(s.set.Log)
	s.static = http.FileServer(s.root)
	s.HandleFunc("/", s.log.RequestLogHandlerFunc(s.ServeAll))
}

// handler func
func (s *JsspServer) ServeAll(w http.ResponseWriter, r *http.Request) {
	index, ext := s.getJsIndexAndExt(r.URL)
	if index == nil {
		s.static.ServeHTTP(w, r)
	} else {
		s.header(w)
		data, err := readFile(index)
		if err != nil {
			s.error(w, err)
			return
		}
		if ext == JSSP {
			data = jssp_jsjs(data)
		}
		js := s.eng.GenJsspEnv(s, w, r)
		ast, err := js.Parse(data)
		if err != nil {
			s.error(w, err)
			return
		}
		_, err = js.Run(ast)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
	}
}

func (s *JsspServer) getJsIndexAndExt(u *url.URL) (http.File, string) {
	if u.Path[0] != '/' {
		u.Path = "/" + u.Path
	}
	if u.Path[len(u.Path)-1] == '/' {
		if f := getFile(s.root, u.Path+INDEX_JSSP); f != nil {
			return f, JSSP
		}
		if f := getFile(s.root, u.Path+INDEX_JSJS); f != nil {
			return f, JSJS
		}
	} else {
		if strings.HasSuffix(u.Path, JSSP) {
			if f := getFile(s.root, u.Path); f != nil {
				return f, JSSP
			}
		}
		if strings.HasSuffix(u.Path, JSJS) {
			if f := getFile(s.root, u.Path); f != nil {
				return f, JSJS
			}
		}
	}
	return nil, ""
}

func (s *JsspServer) error(w http.ResponseWriter, e error) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, e.Error())
}

func (s *JsspServer) header(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Server", Server)
}

// run Jssp server
func (s *JsspServer) Run() {
	err := http.ListenAndServe(":"+s.set.Port, s)
	if err != nil {
		println(err.Error())
	}
}

func main() {
	js := new(JsspServer)
	js.Init()
	js.Run()
}
