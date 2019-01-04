package main

import (
	"bytes"
	"encoding/pem"
	"flag"
	"fmt"
	"github.com/hunkeelin/SuperCA/v2"
	"github.com/hunkeelin/klinutils"
	"github.com/hunkeelin/mtls/klinserver"
	"github.com/hunkeelin/pki"
	"net/http"
)

var (
	confdir   = flag.String("confdir", "/var/lib/superca/work/", "the directory for node classification")
	cakeydir  = flag.String("cakeydir", "/var/lib/superca/CA/keys/", "the directory for ca certs private keys")
	cacertdir = flag.String("cacertdir", "/var/lib/superca/CA/certs/", "the directory for ca certs")
	org       = flag.String("org", "st", "organization")
)

func main() {
	flag.Parse()
	if string((*cakeydir)[len(*cakeydir)-1]) != "/" {
		*cakeydir = *cakeydir + "/"
	}
	if string((*cacertdir)[len(*cacertdir)-1]) != "/" {
		*cacertdir = *cacertdir + "/"
	}
	if string((*confdir)[len(*confdir)-1]) != "/" {
		*confdir = *confdir + "/"
	}
	if !klinutils.Exist(*cacertdir+"rootca.crt") && !klinutils.Exist(*cakeydir+"rootca.key") {
		rootcsr := &klinpki.CAConfig{
			EmailAddress: "support@" + *org + ".com",
			EcdsaCurve:   "",
			Certpath:     *cacertdir + "rootca.crt",
			Keypath:      *cakeydir + "rootca.key",
			MaxDays:      7200,
			RsaBits:      4096,
			Organization: *org,
		}
		klinpki.GenCA(rootcsr)
	}

	// create connection
	c := caserver.Conn{
		Workdir:    *confdir,
		Cakeypath:  *cakeydir,
		Cacertpath: *cacertdir,
	}
	con := http.NewServeMux()
	con.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c.MainHandler(w, r)
	})
	con.Handle("/cacerts/", http.StripPrefix("/cacerts/", http.FileServer(http.Dir(*cacertdir))))

	//gen certs
	j := &klinpki.CSRConfig{
		RsaBits: 4096,
	}
	csr, key := klinpki.GenCSRv2(j)
	f := &klinpki.SignConfig{
		Crtpath:  *cacertdir + "rootca.crt",
		Keypath:  *cakeydir + "rootca.key",
		CsrBytes: csr.Bytes,
		Days:     365,
		IsCA:     false,
	}
	rawcert, err := klinpki.SignCSRv2(f)
	if err != nil {
		panic(err)
	}

	//create server
	var pemcrt, pemkey bytes.Buffer
	pem.Encode(&pemcrt, &pem.Block{Type: "CERTIFICATE", Bytes: rawcert})
	pem.Encode(&pemkey, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: key.Bytes})
	fmt.Println(string(pemkey.Bytes()))
	fmt.Println(string(pemcrt.Bytes()))
	sconfig := &klinserver.ServerConfig{
		BindPort:  klinutils.Stringtoport("superca"),
		KeyBytes:  pemkey.Bytes(),
		CertBytes: pemcrt.Bytes(),
		ServeMux:  con,
		Https:     true,
	}
	panic(klinserver.Server(sconfig))
}
