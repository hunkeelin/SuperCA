package caserver

import (
	"errors"
	"fmt"
	"github.com/hunkeelin/SuperCA/utils"
	"github.com/hunkeelin/klinenv"
	"github.com/hunkeelin/klinutils"
	"strconv"
	"strings"
)

func recursePrint(h []string, p string) (klinenv.AppConfig, error) {
	if len(h) == 0 {
		if cautils.FileExist(p + "config") {
			return klinenv.NewAppConfig(p + "config"), nil
		} else {
			f := klinenv.AppConfig{}
			return f, errors.New("no such config file")
		}
	} else {
		var s string
		for i := range h {
			s += h[len(h)-1-i] + "/"
		}
		if cautils.FileExist(p + s + "config") {
			return klinenv.NewAppConfig(p + s + "config"), nil
		} else {
			return recursePrint(h[1:], p)
		}
	}
}
func iterativePrint(h []string, p string) (klinenv.AppConfig, error) {
	var s string
	for i := range h {
		s += h[len(h)-1-i] + "/"
	}
	for !cautils.FileExist(p + s + "config") {
		if len(h) == 0 {
			if cautils.FileExist(p + "config") {
				return klinenv.NewAppConfig(p + "config"), nil
			} else {
				f := klinenv.AppConfig{}
				return f, errors.New("no such config file")
			}
		}
		h = h[1:]
		s = ""
		for i := range h {
			s += h[len(h)-1-i] + "/"
		}
	}
	return klinenv.NewAppConfig(p + s + "config"), nil
}
func crtkeyDeterm(h, p, wca string) (string, string, float64, bool, error) {
	var cacrt, cakey string
	cfg, err := iterativePrint(strings.Split(h, "."), p)
	if err != nil {
		fmt.Println(err)
		return "", "", 0, false, errors.New("Server no defaults")
	} else {
		signca, err := cfg.Get("signca")
		if err != nil {
			fmt.Println(err)
			return "", "", 0, false, errors.New("Server no default")
		}
		listca := strings.Split(signca, ",")
		if listca[0] == "" {
			fmt.Println("no Signca specified")
			return "", "", 0, false, errors.New("Server no default")
		}
		if wca == "" {
			cacrt = listca[0] + ".crt"
			cakey = listca[0] + ".key"
		} else {
			if klinutils.StringInSlice(wca, listca) {
				cacrt = wca + ".crt"
				cakey = wca + ".key"
			} else {
				fmt.Println("Request SignCA not allowed")
				return "", "", 0, false, errors.New("Server no default")
			}
		}

		isca, err := cfg.Get("isca")
		if err != nil {
			fmt.Println(err)
			return "", "", 0, false, errors.New("Server no default")
		}

		allow, err := cfg.Get("allow")
		if err != nil {
			fmt.Println(err)
			return "", "", 0, false, errors.New("Server no default")
		}
		if strings.TrimSpace(strings.ToLower(allow)) != "true" {
			fmt.Println(err)
			return "", "", 0, false, errors.New("Server not allowed to get crt")
		}
		isCA := false
		if strings.TrimSpace(strings.ToLower(isca)) == "true" {
			isCA = true
		}

		validdays, err := cfg.Get("validdays")
		if err != nil {
			fmt.Println(err)
			return "", "", 0, false, errors.New("Server no default")
		}
		vdays, err := strconv.ParseFloat(validdays, 32)
		if err != nil {
			fmt.Println(err)
			return "", "", 0, false, errors.New("Server no default")
		}
		return cacrt, cakey, vdays, isCA, nil
	}
}
