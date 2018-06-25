package main

import (
	"encoding/pem"
	"flag"
	"github.com/hunkeelin/pki"
	"log"
	"os"
)

var (
	gconfdir  = flag.String("config", "/etc/superCA/CA.conf", "location of the genkins.conf")
	genCa     = flag.Bool("genCA", false, "Generate ca")
	server    = flag.Bool("server", false, "Run the CA server")
	caOutName = flag.String("caOutName", "", "Name of the ca without extensions")
)

func main() {
	flag.Parse()
	c := readconfig(*gconfdir)
	if *genCa {
		if *caOutName == "" {
			log.Fatal("Please specify caOutName.")
		}
		j := &klinpki.CAConfig{
			EmailAddress: "support@" + c.org + ".com",
			EcdsaCurve:   "",
			Certpath:     c.capath + *caOutName + ".crt",
			Keypath:      c.cakeypath + *caOutName + ".key",
			MaxDays:      7200,
			RsaBits:      4096,
			Organization: c.org,
		}
		klinpki.GenCA(j)
		return
	}
	if *server {
		if !Exist(c.keypath) {
			j := &klinpki.CSRConfig{
				RsaBits: 4096,
			}
			// generate key
			csr, key := klinpki.GenCSRv2(j)
			keyOut, err := os.OpenFile(c.keypath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				panic(err)
			}
			pem.Encode(keyOut, key)
			keyOut.Close()
			// generate cert
			f := &klinpki.SignConfig{
				Crtpath:  c.capath + "rootca.crt",
				Keypath:  c.cakeypath + "rootca.key",
				CsrBytes: csr.Bytes,
				Days:     365,
				IsCA:     false,
			}
			rawcert, err := klinpki.SignCSRv2(f)
			if err != nil {
				panic(err)
			}

			clientCRTFile, err := os.Create(c.certpath)
			if err != nil {
				panic(err)
			}
			pem.Encode(clientCRTFile, &pem.Block{Type: "CERTIFICATE", Bytes: rawcert})
			clientCRTFile.Close()

		}
		runServer(&c)
	}
}
