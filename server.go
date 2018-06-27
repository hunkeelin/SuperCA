package main

import (
	"github.com/hunkeelin/mtls/server"
	"log"
	"net/http"
)

func runServer(c *Config) {
	newcon := new(Conn)
	// define config params
	sema := make(chan struct{}, 1)
	newcon.monorun = sema
	newcon.apikey = c.apikey
	newcon.concur = c.concur
	newcon.workdir = c.workdir
	newcon.capath = c.capath
	newcon.cakeypath = c.cakeypath

	if !Exist(c.certpath) || !Exist(c.keypath) {
		log.Fatal("please generate csr and sign it and put it in the correct directory")
	}
	con := http.NewServeMux()
	con.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		newcon.handleWebHook(w, r)
	})
	con.Handle("/cacerts/", http.StripPrefix("/cacerts/", http.FileServer(http.Dir(c.capath))))

	s := &klinserver.ServerConfig{
		BindPort: c.port,
		Cert:     c.certpath,
		Key:      c.keypath,
		Trust:    "program/mtls.crt",
		Https:    true,
		//  Verify:   true,
		ServeMux: con,
	}
	klinserver.Server(s)
}
