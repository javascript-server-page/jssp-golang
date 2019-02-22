package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)

const JSSP = ".jssp"
const JSJS = ".jsjs"
const INDEX_JSSP = "index" + JSSP
const INDEX_JSJS = "index" + JSJS

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
	if index == nil {
		s.static.ServeHTTP(w, r)
	} else {
		data , err := readFile(index)
		if err != nil {
			println(err.Error())
		}
		if ext == JSSP {
			 data = jssp_jsjs(data)
		}
		GenerateJsspEnv(w, r).Run(data)
	}
}

func (s *JsspServer) getJsIndexAndExt(u *url.URL) (*http.File, string) {
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

// run Jssp server
func (s *JsspServer) Run(paras *Parameter) {
	err := http.ListenAndServe(":"+paras.Port, s)
	if err != nil {
		log.Fatal(err)
	}
}
