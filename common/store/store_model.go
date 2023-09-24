package store

import (
	"github.com/lruggieri/fxnow/common/model"
)

type GetUserRequest struct {
	UserID string
	Email  string
}

type GetUserResponse struct {
	User *model.User
}

type CreateUserRequest struct {
	FirstName string
	LastName  string
	Email     string
}

type CreateUserResponse struct {
	UserID string
}

type GetAPIKeyRequest struct {
	UserID   string
	APIKeyID string

	WithUsages bool
}

type GetAPIKeyResponse struct {
	APIKey *model.APIKey
}

type ListAPIKeysRequest struct {
	UserID string
}

type ListAPIKeysResponse struct {
	UserKeys []*model.APIKey
}

type CreateAPIKeyRequest struct {
	UserID     string
	Type       uint8
	Expiration int64 // Unix time (seconds)
}

type CreateAPIKeyResponse struct {
	APIKeyID string
}

type DeleteAPIKeyRequest struct {
	APIKeyID string
}

type DeleteAPIKeyResponse struct{}
