package main

import (
	"crypto/rsa"
	"crypto/tls"
	"net/url"
	"github.com/crewjam/saml/samlsp"
)

func initSAMLServiceProvider() (*samlsp.Middleware, error) {
	keyPair, err := tls.LoadX509KeyPair("path/to/cert.pem", "path/to/key.pem")
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse("https://example.com/saml/acs")
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
