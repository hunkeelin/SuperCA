package caserver

import (
	"fmt"
	"strings"
	"testing"
)

func TestPrint(t *testing.T) {
	h := "test1.klin-pro.com"
	path := "/home/bgops/files/golesson/SuperCA/program/work/"
	hostname := strings.Split(h, ".")
	cfg, err := recursePrint(hostname, path)
	if err != nil {
		panic(err)
	}
	cacrt, err := cfg.Get("signca")
	if err != nil {
		panic(err)
	}
	fmt.Println(cacrt)
}
func TestPrintv2(t *testing.T) {
	fmt.Println("testing version 2")
	h := "test1.klin-pro.com"
	path := "/home/bgops/files/golesson/SuperCA/program/work/"
	hostname := strings.Split(h, ".")
	cfg, err := iterativePrint(hostname, path)
	if err != nil {
		panic(err)
	}
	cacrt, err := cfg.Get("signca")
	if err != nil {
		panic(err)
	}
	fmt.Println(cacrt)
}
