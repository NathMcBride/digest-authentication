package secret

type SecretProviderProvider struct{}

// example secret provider
func (p *SecretProviderProvider) GetSecret(userID string) (string, bool, error) {
	if userID == "jim" {
		return "password", true, nil
	}

	if userID == "john" {
		return "cheese", true, nil
	}

	return "", false, nil
}
