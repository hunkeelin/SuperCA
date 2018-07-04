package caserver

import (
	"errors"
	"github.com/hunkeelin/SuperCA/utils"
	"github.com/hunkeelin/klinenv"
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
func crtkeyDeterm(h, p string) (string, string, float64, bool, error) {
	cfg, err := recursePrint(strings.Split(h, "."), p)
	if err != nil {
		return "", "", 0, false, errors.New("Server no defaults")
	} else {
		cacrt, err := cfg.Get("cacrt")
		if err != nil {
			return "", "", 0, false, errors.New("Server no default")
		}

		cakey, err := cfg.Get("cakey")
		if err != nil {
			return "", "", 0, false, errors.New("Server no default")
		}

		isca, err := cfg.Get("isca")
		if err != nil {
			return "", "", 0, false, errors.New("Server no default")
		}

		allow, err := cfg.Get("allow")
		if err != nil {
			return "", "", 0, false, errors.New("Server no default")
		}
		if strings.TrimSpace(strings.ToLower(allow)) != "true" {
			return "", "", 0, false, errors.New("Server not allowed to get crt")
		}
		isCA := false
		if strings.TrimSpace(strings.ToLower(isca)) == "true" {
			isCA = true
		}

		validdays, err := cfg.Get("validdays")
		if err != nil {
			return "", "", 0, false, errors.New("Server no default")
		}
		vdays, err := strconv.ParseFloat(validdays, 32)
		if err != nil {
			return "", "", 0, false, errors.New("Server no default")
		}
		return cacrt, cakey, vdays, isCA, nil
	}
}
