package model

type DigestHeader struct {
	Realm     string `httpparam:"realm,omitempty"`
	Algorithm string `httpparam:"algorithm,unq,omitempty"`
	Qop       string `httpparam:"qop"`
	Opaque    string `httpparam:"opaque"`
	Nonce     string `httpparam:"nonce"`
	UserHash  bool   `httpparam:"userhash,omitempty"`
}
