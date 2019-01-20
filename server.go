package main

import (
	"log"
	"net/http"
)

type JsspServer struct {
	http.ServeMux
}


// register handlers
func (s *JsspServer) Init(paras *Parameter) {
}


// run Jssp server
func (s *JsspServer) Run(paras *Parameter) {
	err := http.ListenAndServe(":"+paras.Port, s)
	if err != nil {
		log.Fatal(err)
	}
}
