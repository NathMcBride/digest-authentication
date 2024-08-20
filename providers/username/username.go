package username

import (
	"errors"

	"github.com/NathMcBride/web-authentication/authentication/hasher"
)

type UsernameProvider struct {
	Realm string
}

// example UsernameProvider
func (u *UsernameProvider) GetUserName(usernameHash string) (string, error) {
	hj, _ := hasher.H("jim" + ":" + u.Realm)

	if usernameHash == hj {
		return "jim", nil
	}

	hjn, _ := hasher.H("john" + ":" + u.Realm)
	if usernameHash == hjn {
		return "john", nil
	}

	return "", errors.New("no matching user name")
}
