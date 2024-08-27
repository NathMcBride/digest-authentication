package model

type AuthHeader struct {
	Response  string `hparam:"response"`
	UserID    string `hparam:"username"`
	Realm     string `hparam:"realm,omitempty"`
	Algorithm string `hparam:"algorithm,unq"`
	Qop       string `hparam:"qop"`
	Cnonce    string `hparam:"cnonce"`
	Nc        string `hparam:"nc"`
	Opaque    string `hparam:"opaque"`
	Uri       string `hparam:"uri"`
	Nonce     string `hparam:"nonce,unq,omitempty"`
	UserHash  bool   `hparam:"userhash,omitempty"`
	// AThing    string `hparam:"-"`
}

//should the nonce be unquoted?
