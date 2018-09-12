package caserver

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/hunkeelin/SuperCA/utils"
	"github.com/hunkeelin/pki"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

func websigncsr(w http.ResponseWriter, r *http.Request, p *Conn) ([]byte, int, *s_respBody) {
	var respb s_respBody
	d, err := ioutil.ReadAll(r.Body)
	clientCSR, err := x509.ParseCertificateRequest(d)
	if err != nil {
		return []byte("unable to parse csr"), 500, &respb
	}
	hostname, err := net.LookupAddr(strings.Split(r.RemoteAddr, ":")[0])
	if err != nil {
		return []byte("Failed DNS validation"), 500, &respb
	}
	//	if clientCSR.DNSNames[0]+"." == hostname[0] {
	if cautils.StringInSlice(string(hostname[0][:len(hostname[0])-1]), clientCSR.DNSNames) {
		crtp, keyp, days, isCA, err := crtkeyDeterm(clientCSR.DNSNames[0], p.workdir, r.Header.Get("SignCA"))
		if err != nil {
			return []byte(err.Error()), 500, &respb
		}
		f := &klinpki.SignConfig{
			Crtpath:  p.capath + crtp,
			Keypath:  p.cakeypath + keyp,
			CsrBytes: d,
			Days:     days,
			IsCA:     isCA,
		}
		chainOfTrust, err := ioutil.ReadFile(p.capath + crtp)
		if err != nil {
			fmt.Println("Unable to read chain of trust")
			return []byte(""), 500, &respb
		}
		rawcert, err := klinpki.SignCSRv2(f)
		if err != nil {
			fmt.Println("Unable to Sign csr")
			return []byte(""), 500, &respb
		}
		respb := s_respBody{
			Cert:         rawcert,
			ChainOfTrust: chainOfTrust,
		}
		encodebody, err := json.Marshal(respb)
		if err != nil {
			fmt.Println("unable to marshall the json")
			return []byte(""), 500, &respb
		}
		return encodebody, 200, &respb
	} else {
		return []byte("Failed at dns validation"), 500, &respb
	}
}
