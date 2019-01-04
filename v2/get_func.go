package caserver

import (
	"crypto/x509"
	"errors"
	"github.com/hunkeelin/SuperCA/lib"
	"github.com/hunkeelin/klinenv"
	"github.com/hunkeelin/klinutils"
	"github.com/hunkeelin/pki"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
)

func (c *Conn) websigncsr(w http.ResponseWriter, r *http.Request) error {
	csrbytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	// check if hostname in dns
	hostname, err := net.LookupAddr(strings.Split(r.RemoteAddr, ":")[0])
	if err != nil {
		return err
	}
	clientCSR, err := x509.ParseCertificateRequest(csrbytes)
	if err != nil {
		return err
	}
	// check if hostname and dns name match
	if !klinutils.StringInSlice(string(hostname[0][:len(hostname[0])-1]), clientCSR.DNSNames) {
		return errors.New("DNS validate between hostname and client DNSNames do not match")
	}
	crtp, keyp, days, isCA, err := c.crtkeyDeterm(clientCSR.DNSNames[0], r.Header.Get("SignCA"))
	if err != nil {
		return err
	}

	f := &klinpki.SignConfig{
		Crtpath:  c.Cacertpath + crtp,
		Keypath:  c.Cakeypath + keyp,
		CsrBytes: csrbytes,
		Days:     days,
		IsCA:     isCA,
	}
	rawcert, err := klinpki.SignCSRv2(f)
	if err != nil {
		return err
	}
	chainOfTrust, err := ioutil.ReadFile(c.Cacertpath + crtp)
	if err != nil {
		return err
	}
	returncrts := supercalib.ReturnPayload{
		Cert:         rawcert,
		ChainOfTrust: chainOfTrust,
	}
	err = json.NewEncoder(w).Encode(returncrts)
	if err != nil {
		return err
	}
	return nil
}

func itprint(h []string, p string) (klinenv.AppConfig, error) {
	var s string
	for i := range h {
		s += h[len(h)-1-i] + "/"
	}
	for !klinutils.Exist(p + s + "config") {
		if len(h) == 0 {
			if klinutils.Exist(p + "config") {
				return klinenv.NewAppConfig(p + "config"), nil
			} else {
				f := klinenv.AppConfig{}
				return f, errors.New("no such config file for ")
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

func (c *Conn) crtkeyDeterm(hostname, signca string) (string, string, float64, bool, error) {
	var cacrt, cakey string
	var isCA bool
	cfg, err := itprint(strings.Split(hostname, "."), c.Workdir)
	if err != nil {
		return "", "", 0, false, err
	}

	// check the config
	tosign, err := cfg.Get("signca")
	if err != nil {
		return "", "", 0, false, err
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
	// end of checking config
	listca := strings.Split(tosign, ",")
	if listca[0] == "" {
		return "", "", 0, false, errors.New("no signca specified in the config for host " + hostname)
	}
	if signca == "" {
		cacrt = listca[0] + ".crt"
		cakey = listca[0] + ".key"
	} else {
		if klinutils.StringInSlice(signca, listca) {
			cacrt = signca + ".crt"
			cakey = signca + ".key"
		} else {
			return "", "", 0, false, errors.New("Requested SignCA " + signca + " not allowed for host" + hostname)
		}
	}
	return cacrt, cakey, vdays, isCA, nil
}
