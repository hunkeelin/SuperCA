package caserver

import (
	"sync"
)

type s_respBody struct {
	Cert         []byte `json:"cert"`
	ChainOfTrust []byte `json:"chainoftrust"`
}

type Conn struct {
	regex     string
	capath    string
	cakeypath string
	apikey    string
	concur    int
	jobdir    string
	workdir   string
	mu        sync.Mutex
	monorun   chan struct{}
}
type Config struct {
	apikey    string
	bindaddr  string
	org       string
	port      string
	workdir   string
	certpath  string
	keypath   string
	jobdir    string
	concur    int
	capath    string
	cakeypath string
}
