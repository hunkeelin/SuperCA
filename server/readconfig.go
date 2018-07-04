package caserver

import (
	"github.com/hunkeelin/klinenv"
	"log"
	"strconv"
)

func checkerr(err error) {
	if err == nil {
		panic(err)
	}
	return
}

func Readconfig(p string) Config {
	var c Config

	config := klinenv.NewAppConfig(p)
	rconcur, err := config.Get("concur")
	if err != nil {
		log.Fatal("unable to retrieve the value of concur check config file")
	}
	concur, err := strconv.Atoi(rconcur)
	if err != nil {
		log.Fatal("can't convert string to int for concur")
	}
	c.Concur = concur
	apikey, err := config.Get("apikey")
	checkerr(err)
	c.Apikey = apikey

	capath, err := config.Get("capath")
	checkerr(err)
	if len(capath) == 0 {
		log.Fatal("Please specify capath in config")
	} else {
		if string(capath[len(capath)-1]) != "/" {
			capath += "/"
		}
		c.Capath = capath
	}

	cakeypath, err := config.Get("cakeypath")
	checkerr(err)
	if len(cakeypath) == 0 {
		log.Fatal("Please specify cakeypath in config")
	} else {
		if string(cakeypath[len(cakeypath)-1]) != "/" {
			cakeypath += "/"
		}
		c.Cakeypath = cakeypath
	}

	workdir, err := config.Get("workdir")
	checkerr(err)
	if len(workdir) == 0 {
		log.Fatal("Please specify workdir in config")
	} else {
		if string(workdir[len(workdir)-1]) != "/" {
			workdir += "/"
		}
		c.Workdir = workdir
	}

	org, err := config.Get("org")
	checkerr(err)
	c.Org = org

	bindaddr, err := config.Get("bindaddr")
	checkerr(err)
	c.Bindaddr = bindaddr

	port, err := config.Get("port")
	checkerr(err)
	c.Port = port

	certpath, err := config.Get("certpath")
	checkerr(err)
	c.Certpath = certpath

	keypath, err := config.Get("keypath")
	checkerr(err)
	c.Keypath = keypath

	return c
}
