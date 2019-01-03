package caserver

import (
	"net/http"
)

func (c *conn) get(w http.ResponseWriter, r *http.Request) error {
    if strings.HasPrefix(r.Header.Get("content-type"), "application/x-www-form-urlencoded") {
        err := c.websigncsr(w, r)
        if err != nil {
            return err
        }
    }
	return nil
}
