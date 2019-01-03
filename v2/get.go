package caserver

import (
	"net/http"
)

func (c *conn) get(w http.ResponseWriter, r *http.Request) error {
	err := c.websigncsr(w, r)
	if err != nil {
		return err
	}
	return nil
}
