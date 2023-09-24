package logic

import "github.com/lruggieri/fxnow/common/model"

type ListAPIKeysRequest struct{}

type ListAPIKeysResponse struct {
	APIKeys []*model.APIKey
}

type CreateAPIKeyRequest struct{}

type CreateAPIKeyResponse struct {
	APIKeyID string
}

type DeleteAPIKeyRequest struct {
	APIKeyID string
}

type DeleteAPIKeyResponse struct{}

type CreateUserRequest struct {
	FirstName string
	LastName  string
	Email     string
}

type CreateUserResponse struct {
	UserID string
}
