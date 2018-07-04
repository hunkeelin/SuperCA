package caserver

import (
	"github.com/hunkeelin/SuperCA/utils"
	"github.com/hunkeelin/mtls/server"
	"log"
	"net/http"
)

func Server(c *Config, rootca string) {
	newcon := new(Conn)
	// define config params
	sema := make(chan struct{}, 1)
	newcon.monorun = sema
	newcon.apikey = c.apikey
	newcon.concur = c.concur
	newcon.workdir = c.workdir
	newcon.capath = c.capath
	newcon.cakeypath = c.cakeypath

	if !cautils.Exist(c.certpath) || !cautils.Exist(c.keypath) {
		log.Fatal("key cert path for https does not exist!")
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
		Trust:    "program/" + rootca + ".crt",
		Https:    true,
		//  Verify:   true,
		ServeMux: con,
	}
	klinserver.Server(s)
}
