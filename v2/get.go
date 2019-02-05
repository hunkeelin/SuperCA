package caserver

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func (c *Conn) get(w http.ResponseWriter, r *http.Request) error {
	if !strings.HasPrefix(r.Header.Get("content-type"), "application/x-www-form-urlencoded") {
		return fmt.Errorf("get classification error: %v", errors.New("Bad content type"))
	}
	err := c.websigncsr(w, r)
	if err != nil {
		return err
	}
	return nil
}
