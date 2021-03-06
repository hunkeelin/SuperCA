package caserver

import (
	"github.com/hunkeelin/SuperCA/utils"
	"github.com/hunkeelin/mtls/server"
	"log"
	"net/http"
)

func Server(c *Config, rootca string) (error){
	newcon := new(Conn)
	// define config params
	sema := make(chan struct{}, 1)
	newcon.monorun = sema
	newcon.apikey = c.Apikey
	newcon.concur = c.Concur
	newcon.workdir = c.Workdir
	newcon.capath = c.Capath
	newcon.cakeypath = c.Cakeypath

	if !cautils.Exist(c.Certpath) || !cautils.Exist(c.Keypath) {
		log.Fatal("key cert path for https does not exist!")
	}
	con := http.NewServeMux()
	con.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		newcon.handleWebHook(w, r)
	})
	con.Handle("/cacerts/", http.StripPrefix("/cacerts/", http.FileServer(http.Dir(c.Capath))))

	s := &klinserver.ServerConfig{
		BindPort: c.Port,
		Cert:     c.Certpath,
		Key:      c.Keypath,
		Trust:    "program/" + rootca + ".crt",
		Https:    true,
		//  Verify:   true,
		ServeMux: con,
	}
	return klinserver.Server(s)
}
