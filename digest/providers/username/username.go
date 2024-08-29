package username

import (
	"github.com/NathMcBride/web-authentication/digest/authentication/hasher"
)

type UsernameProvider struct {
	Realm string
}

// example UsernameProvider
func (u *UsernameProvider) GetUserName(usernameHash string) (string, bool, error) {
	hasher := hasher.Hasher{}
	hj, err := hasher.Hash("jim" + ":" + u.Realm)
	if err != nil {
		return "", false, err
	}

	if usernameHash == hj {
		return "jim", true, nil
	}

	hjn, err := hasher.Hash("john" + ":" + u.Realm)
	if err != nil {
		return "", false, err
	}

	if usernameHash == hjn {
		return "john", true, nil
	}

	return "", false, nil
}
