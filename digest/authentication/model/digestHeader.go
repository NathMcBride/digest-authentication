package model

type DigestHeader struct {
	Realm     string `hparam:"realm,omitempty"`
	Algorithm string `hparam:"algorithm,unq,omitempty"`
	Qop       string `hparam:"qop"`
	Opaque    string `hparam:"opaque"`
	Nonce     string `hparam:"nonce"`
	UserHash  bool   `hparam:"userhash,omitempty"`
}
