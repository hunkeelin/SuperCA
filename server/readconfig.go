package caserver

import (
	"github.com/hunkeelin/SuperCA/utils"
	"github.com/hunkeelin/klinenv"
	"log"
	"strconv"
)

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
	cautils.Checkerr(err)
	c.Apikey = apikey

	capath, err := config.Get("capath")
	cautils.Checkerr(err)
	if len(capath) == 0 {
		log.Fatal("Please specify capath in config")
	} else {
		if string(capath[len(capath)-1]) != "/" {
			capath += "/"
		}
		c.Capath = capath
	}

	cakeypath, err := config.Get("cakeypath")
	cautils.Checkerr(err)
	if len(cakeypath) == 0 {
		log.Fatal("Please specify cakeypath in config")
	} else {
		if string(cakeypath[len(cakeypath)-1]) != "/" {
			cakeypath += "/"
		}
		c.Cakeypath = cakeypath
	}

	workdir, err := config.Get("workdir")
	cautils.Checkerr(err)
	if len(workdir) == 0 {
		log.Fatal("Please specify workdir in config")
	} else {
		if string(workdir[len(workdir)-1]) != "/" {
			workdir += "/"
		}
		c.Workdir = workdir
	}

	org, err := config.Get("org")
	cautils.Checkerr(err)
	c.Org = org

	bindaddr, err := config.Get("bindaddr")
	cautils.Checkerr(err)
	c.Bindaddr = bindaddr

	port, err := config.Get("port")
	cautils.Checkerr(err)
	c.Port = port

	certpath, err := config.Get("certpath")
	cautils.Checkerr(err)
	c.Certpath = certpath

	keypath, err := config.Get("keypath")
	cautils.Checkerr(err)
	c.Keypath = keypath

	return c
}
