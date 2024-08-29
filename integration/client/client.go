package client

import (
	"fmt"
	"net/http"

	"github.com/NathMcBride/web-authentication/digest/authentication/digest"
	"github.com/NathMcBride/web-authentication/digest/authentication/hasher"
	"github.com/NathMcBride/web-authentication/digest/authentication/model"
	"github.com/NathMcBride/web-authentication/digest/headers/paramlist"
	"github.com/NathMcBride/web-authentication/digest/providers/credential"
)

type Client struct {
	Endpoint           string
	ShouldAuthenticate bool
	Username           string
	Password           string
	nonce              string
	nc                 int
	opaque             string
}

func (client *Client) addDigest(username string, password string, dh model.DigestHeader, request *http.Request) error {
	if client.nc == 0 {
		client.nc = 1
	}

	if client.nonce == "" {
		client.nonce = dh.Nonce
	}

	if client.opaque == "" {
		client.opaque = dh.Opaque
	}

	hasher := hasher.Hasher{}
	userhash, err := hasher.Hash(username + ":" + dh.Realm)
	if err != nil {
		return err
	}

	randomKey := digest.RandomKey{}
	authHeader := model.AuthHeader{
		UserID:    userhash,
		Realm:     dh.Realm,
		Algorithm: dh.Algorithm,
		Qop:       dh.Qop,
		Cnonce:    randomKey.Create(),
		Nc:        fmt.Sprintf("%d", client.nc),
		Opaque:    client.opaque,
		Uri:       request.RequestURI,
		Nonce:     client.nonce,
		UserHash:  true,
	}

	digest := digest.Digest{Sha256: &hasher}
	cr := credential.Credentials{Username: username, Password: password}

	result, err := digest.Calculate(cr, authHeader, request.Method)
	if err != nil {
		return err
	}
	authHeader.Response = result

	marshalled, err := paramlist.Marshal(authHeader)
	if err != nil {
		return err
	}

	request.Header.Add("Authorization", "Digest "+string(marshalled[:]))
	client.nc++

	return nil
}

func (client *Client) GetProtected() (*http.Response, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/protected", client.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	return client.doRequest(request)
}

func (client *Client) GetHealth() (*http.Response, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/health", client.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	return client.doRequest(request)
}

func (client *Client) GetRoot() (*http.Response, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/", client.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	return client.doRequest(request)
}

func (client *Client) doRequest(request *http.Request) (*http.Response, error) {
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized && client.ShouldAuthenticate {
		dh := model.DigestHeader{}
		paramlist.Unmarshal([]byte(resp.Header.Get("WWW-Authenticate")), &dh)
		client.addDigest(client.Username, client.Password, dh, request)

		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			return nil, err
		}
		return resp, err
	}

	return resp, err
}
