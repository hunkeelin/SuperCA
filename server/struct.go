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
	Apikey    string
	Bindaddr  string
	Org       string
	Port      string
	Workdir   string
	Certpath  string
	Keypath   string
	Jobdir    string
	Concur    int
	Capath    string
	Cakeypath string
}
