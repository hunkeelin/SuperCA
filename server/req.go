package caserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (f *Conn) handleWebHook(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.Header.Get("content-type"), "application/x-www-form-urlencoded") {
		msg, status, respb := websigncsr(w, r, f)
		if status != 200 {
			w.WriteHeader(status)
			w.Write(msg)
		} else {
			fmt.Println("Signing csr from ", r.RemoteAddr)
			w.WriteHeader(status)
			err := json.NewEncoder(w).Encode(respb)
			if err != nil {
				fmt.Println("unable to encode json to io writer")
			}
			fmt.Println("Sent cert back to", r.RemoteAddr)
		}
		return
	}
	w.WriteHeader(500)
	w.Write([]byte("Wrong type"))
	return
}
