package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"
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

	tlsconfig := &tls.Config{
		PreferServerCipherSuites: true,
		// Only use curves which have assembly implementations
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		MinVersion: tls.VersionTLS12,
		//	ClientAuth: tls.RequireAndVerifyClientCert,
	}
	con := http.NewServeMux()
	con.HandleFunc("/", newcon.handleWebHook)
	s := &http.Server{
		Addr:         c.bindaddr + ":" + c.port,
		TLSConfig:    tlsconfig,
		Handler:      con,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	fmt.Println("listening to " + c.bindaddr + " " + c.port)
	err := s.ListenAndServeTLS(c.certpath, c.keypath)
	//err := s.ListenAndServe()
	if err != nil {
		log.Fatal("can't listen and serve check port and binding addr", err)
	}
}
