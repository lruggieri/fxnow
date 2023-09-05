package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserInfoFromContext(t *testing.T) {
	uInfo := &UserInfo{
		Email:      "user@domain.com",
		GivenName:  "given_name",
		FamilyName: "family_name",
	}
	uInfoCtx := context.WithValue(context.Background(), ContextUserInfoKey, uInfo)
	wrongUInfoCtx := context.WithValue(context.Background(), ContextUserInfoKey, 42)

	assert.Equal(t, (*UserInfo)(nil), GetUserInfoFromContext(context.Background()))
	assert.Equal(t, (*UserInfo)(nil), GetUserInfoFromContext(wrongUInfoCtx))
	assert.Equal(t, uInfo, GetUserInfoFromContext(uInfoCtx))
}
