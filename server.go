package main

import (
	"log"
	"net/http"
)

type JsspServer struct {
	http.ServeMux
}


// init JsspServer
func (s *JsspServer) Init(paras *Parameter) {
	s.Handle("/", s)
}

// handler func
func (s *JsspServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

// run Jssp server
func (s *JsspServer) Run(paras *Parameter) {
	err := http.ListenAndServe(":"+paras.Port, s)
	if err != nil {
		log.Fatal(err)
	}
}
