package credential

type SecretProvider interface {
	GetSecret(userID string) (string, bool, error)
}

type UsernameProvider interface {
	GetUserName(usernameHash string) (string, bool, error)
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
func (c *CredentialProvider) GetCredentials(userID string, useHash bool) (*Credentials, bool, error) {
	username := userID
	if useHash {
		plainUsername, found, err := c.UsernameProvider.GetUserName(userID)
		if err != nil || !found {
			return nil, found, err
		}
		username = plainUsername
	}

	passwd, found, err := c.SecretProvider.GetSecret(username)
	if err != nil || !found {
		return nil, found, err
	}

	cred := &Credentials{
		Username: username,
		Password: passwd,
	}
	return cred, true, nil
}
