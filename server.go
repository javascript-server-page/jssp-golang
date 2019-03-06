package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
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
			msg, status := s.toHTTPError(err)
			s.error(w, msg, status)
		}
		if ext == JSSP {
			 data = jssp_jsjs(data)
		}
		GenerateJsspEnv(w, r).Run(data)
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

func (s *JsspServer) toHTTPError(err error) (msg string, httpStatus int) {
	if os.IsNotExist(err) {
		return "404 page not found", http.StatusNotFound
	}
	if os.IsPermission(err) {
		return "403 Forbidden", http.StatusForbidden
	}
	// Default:
	return "500 Internal Server Error", http.StatusInternalServerError
}

func (s *JsspServer) error(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
}

// run Jssp server
func (s *JsspServer) Run(paras *Parameter) {
	err := http.ListenAndServe(":"+paras.Port, s)
	if err != nil {
		log.Fatal(err)
	}
}
