package caserver

import (
	"fmt"
	"github.com/json-iterator/go"
	"net/http"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (c *Conn) MainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("received request from ", r.RemoteAddr)
	switch r.Method {
	case "GET":
		err := c.get(w, r)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
		}
	default:
		http.Error(w, "invalid method: endpoint is GET only", http.StatusMethodNotAllowed)
	}
}
