package main

import (
	"context"
	"crypto"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"net/http"
)

type OIDCConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	ProviderURL  string
	Scopes       []string
}

func initOIDCProvider(ctx context.Context, config OIDCConfig) (*oidc.Provider, *oauth2.Config, error) {
	provider, err := oidc.NewProvider(ctx, config.ProviderURL)
	if err != nil {
		return nil, nil, err
	}

	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       append(config.Scopes, oidc.ScopeOpenID),
	}

	return provider, oauth2Config, nil
}

func startOIDCFlow(oauth2Config *oauth2.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		state := crypto.Hash(0).HashFunc().String()
		url := oauth2Config.AuthCodeURL(state)
		c.Redirect(http.StatusFound, url)
	}
}

func oidcCallback(oauth2Config *oauth2.Config, verifier *oidc.IDTokenVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {

		code := c.Query("code")

		token, err := oauth2Config.Exchange(context.Background(), code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
			return
		}

		idToken, ok := token.Extra("id_token").(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No id_token field in oauth2 token."})
			return
		}

		_, err = verifier.Verify(context.Background(), idToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify ID Token"})
			return
		}

		// Successful authentication
		c.JSON(http.StatusOK, gin.H{"status": "Authenticated"})
	}
}

func oidcLogout(oauth2Config *oauth2.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := oauth2Config.AuthCodeURL("logout")
		c.Redirect(http.StatusFound, url)
	}
}
