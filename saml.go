package main

import (
	"crypto/rsa"
	"crypto/tls"
	"net/url"
	"os"

	"github.com/crewjam/saml/samlsp"
)

func initSAMLServiceProvider() (*samlsp.Middleware, error) {
	keyPair, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(os.Getenv("SERVER_URL"))
	if err != nil {
		return nil, err
	}

	sp, err := samlsp.New(samlsp.Options{
		URL:         *parsedURL,
		Key:         keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate: keyPair.Leaf,
	})
	if err != nil {
		return nil, err
	}
	return sp, nil
}
