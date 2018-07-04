package main

import (
	"encoding/pem"
	"flag"
	"github.com/hunkeelin/SuperCA/server"
	"github.com/hunkeelin/SuperCA/utils"
	"github.com/hunkeelin/pki"
	"log"
	"os"
)

var (
	gconfdir  = flag.String("config", "/etc/superca/CA.conf", "location of the genkins.conf")
	genCa     = flag.Bool("genCA", false, "Generate ca")
	server    = flag.Bool("server", false, "Run the CA server")
	caOutName = flag.String("caOutName", "", "Name of the ca without extensions")
	caRsabits = flag.Int("caRsaBits", 4096, "rsabits when generating ca")
	rootCA    = flag.String("Rootca", "rootca", "name of the rootca crt")
)

func main() {
	flag.Parse()
	c := caserver.Readconfig(*gconfdir)
	if *genCa {
		if *caOutName == "" {
			log.Fatal("Please specify caOutName.")
		}
		j := &klinpki.CAConfig{
			EmailAddress: "support@" + c.Org + ".com",
			EcdsaCurve:   "",
			Certpath:     c.Capath + *caOutName + ".crt",
			Keypath:      c.Cakeypath + *caOutName + ".key",
			MaxDays:      7200,
			RsaBits:      *caRsabits,
			Organization: c.Org,
		}
		klinpki.GenCA(j)
		return
	}
	if *server {
		if !cautils.Exist(c.Capath+*rootCA+".crt") && !cautils.Exist(c.Cakeypath+*rootCA+".key") {
			rootcsr := &klinpki.CAConfig{
				EmailAddress: "support@" + c.Org + ".com",
				EcdsaCurve:   "",
				Certpath:     c.Capath + *rootCA + ".crt",
				Keypath:      c.Cakeypath + *rootCA + ".key",
				MaxDays:      7200,
				RsaBits:      *caRsabits,
				Organization: c.Org,
			}
			klinpki.GenCA(rootcsr)
		}
		if !cautils.Exist(c.Keypath) {
			j := &klinpki.CSRConfig{
				RsaBits: 4096,
			}
			// generate key
			csr, key := klinpki.GenCSRv2(j)
			keyOut, err := os.OpenFile(c.Keypath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				panic(err)
			}
			pem.Encode(keyOut, key)
			keyOut.Close()
			// generate cert
			f := &klinpki.SignConfig{
				Crtpath:  c.Capath + *rootCA + ".crt",
				Keypath:  c.Cakeypath + *rootCA + ".key",
				CsrBytes: csr.Bytes,
				Days:     365,
				IsCA:     false,
			}
			rawcert, err := klinpki.SignCSRv2(f)
			if err != nil {
				panic(err)
			}

			clientCRTFile, err := os.Create(c.Certpath)
			if err != nil {
				panic(err)
			}
			pem.Encode(clientCRTFile, &pem.Block{Type: "CERTIFICATE", Bytes: rawcert})
			clientCRTFile.Close()

		}
		caserver.Server(&c, *rootCA)
	}
}
