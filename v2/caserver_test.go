package caserver

import (
	"fmt"
	"github.com/hunkeelin/klinutils"
	"github.com/hunkeelin/mtls/klinserver"
	"net/http"
	"strings"
	"testing"
)

func TestServer(t *testing.T) {
	c := conn{
		workdir:   "/home/bgops/files/golesson/SuperCA/program/work/",
		cakeypath: "/home/bgops/files/golesson/SuperCA/program/CA/keys/",
		capath:    "/home/bgops/files/golesson/SuperCA/program/CA/certs/",
	}
	con := http.NewServeMux()
	con.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c.mainHandler(w, r)
	})
	con.Handle("/cacerts/", http.StripPrefix("/cacerts/", http.FileServer(http.Dir("/home/bgops/files/golesson/SuperCA/program/CA/certs"))))
	j := &klinserver.ServerConfig{
		BindPort: klinutils.Stringtoport("superca"),
		Cert:     "/home/bgops/files/golesson/SuperCA/program/certs/test1.klin-pro.com.crt",
		Key:      "/home/bgops/files/golesson/SuperCA/program/keys/test1.klin-pro.com.key",
		ServeMux: con,
		Https:    true,
	}
	err := klinserver.Server(j)
	if err != nil {
		panic(err)
	}
}

func TestItprint(t *testing.T) {
	h := "test1.klin-pro.com"
	path := "/home/bgops/files/golesson/SuperCA/program/work/"
	hostname := strings.Split(h, ".")
	cfg, err := itprint(hostname, path)
	if err != nil {
		panic(err)
	}
	cacrt, err := cfg.Get("signca")
	if err != nil {
		panic(err)
	}
	fmt.Println(cacrt)
}

func TestDeterm(t *testing.T) {
	fmt.Println("testing determ")
	c := conn{
		workdir: "/home/bgops/files/golesson/SuperCA/program/work/",
	}
	fmt.Println(c.crtkeyDeterm("test1.klin-pro.com", "intermca"))
}
