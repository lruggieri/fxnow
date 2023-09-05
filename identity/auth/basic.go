package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type BasicAuthenticator struct {
	ctx context.Context

	// OIDC
	oauthConfig *oauth2.Config
	provider    *oidc.Provider
	verifier    *oidc.IDTokenVerifier
}

func (b *BasicAuthenticator) GetOIDCConsentURL(redirectURL string) string {
	return b.oauthConfig.AuthCodeURL(redirectURL)
}

func (b *BasicAuthenticator) AuthenticateOIDC(code string) (string, error) {
	// TODO make constant
	ctx, cancel := context.WithTimeout(b.ctx, 2*time.Second)
	defer cancel()

	oauth2Token, err := b.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return "", err
	}

	// Extract the ID token from the OAuth2 token
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return "", fmt.Errorf("id_token not found")
	}

	// Verify the ID token
	_, err = b.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return "", fmt.Errorf("failed to verify ID token")
	}

	return rawIDToken, nil
}

func (b *BasicAuthenticator) IsJWTValid(token string) bool {
	idToken, err := b.verifier.Verify(b.ctx, token)
	if err != nil || idToken.Expiry.Before(time.Now()) {
		return false
	}

	return true
}

func (b *BasicAuthenticator) GetUserInfo(token string) *UserInfo {
	idToken, err := b.verifier.Verify(b.ctx, token)
	if err != nil || idToken.Expiry.Before(time.Now()) {
		return nil
	}

	var userInfo UserInfo

	if err = idToken.Claims(&userInfo); err != nil {
		return nil
	}

	return &userInfo
}

func NewBasic(ctx context.Context, config Config) (Authenticator, error) {
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return nil, err
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: config.OIDC.ClientID})

	oauthConfig := &oauth2.Config{
		ClientID:     config.OIDC.ClientID,
		ClientSecret: config.OIDC.ClientSecret,
		RedirectURL:  config.OIDC.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return &BasicAuthenticator{
		ctx:         ctx,
		oauthConfig: oauthConfig,
		provider:    provider,
		verifier:    verifier,
	}, nil
}
