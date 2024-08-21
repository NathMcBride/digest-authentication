package credential

type SecretProvider interface {
	GetSecret(userID string) (string, error)
}

type UsernameProvider interface {
	GetUserName(usernameHash string) (string, error)
}

type Credentials struct {
	Username string
	Password string
}

type CredentialProvider struct {
	UsernameProvider UsernameProvider
	SecretProvider   SecretProvider
}

// example credential provider
func (c *CredentialProvider) GetCredentials(userID string, useHash bool) (*Credentials, error) {
	username := userID
	if useHash {
		plainUsername, err := c.UsernameProvider.GetUserName(userID)
		if err != nil {
			return nil, err
		}
		username = plainUsername
	}

	passwd, err := c.SecretProvider.GetSecret(username)
	if err != nil {
		return nil, err
	}

	cred := &Credentials{
		Username: username,
		Password: passwd,
	}
	return cred, nil
}
