package model

type AuthHeader struct {
	Response  string `httpparam:"response"`
	UserID    string `httpparam:"username"`
	Realm     string `httpparam:"realm,omitempty"`
	Algorithm string `httpparam:"algorithm,unq"`
	Qop       string `httpparam:"qop"`
	Cnonce    string `httpparam:"cnonce"`
	Nc        string `httpparam:"nc"`
	Opaque    string `httpparam:"opaque"`
	Uri       string `httpparam:"uri"`
	Nonce     string `httpparam:"nonce,unq,omitempty"`
	UserHash  bool   `httpparam:"userhash,omitempty"`
	// AThing    string `hparam:"-"`
}

//should the nonce be unquoted?
