package secret

import "errors"

type SecretProviderProvider struct{}

// example secret provider
func (p *SecretProviderProvider) GetSecret(userID string) (string, error) {
	if userID == "jim" {
		return "password", nil
	}

	if userID == "john" {
		return "cheese", nil
	}

	return "", errors.New("unable to find user secret")
}
