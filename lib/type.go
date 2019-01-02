package supercalib

type ReturnPayload struct {
	Cert         []byte `json:"cert"`
	ChainOfTrust []byte `json:"chainoftrust"`
}
