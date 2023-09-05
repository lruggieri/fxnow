package store

import "context"

type Store interface {
	// User
	GetUser(ctx context.Context, req GetUserRequest) (*GetUserResponse, error)
	CreateUser(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error)

	// API Key
	GetAPIKey(ctx context.Context, req GetAPIKeyRequest) (*GetAPIKeyResponse, error)
	ListAPIKeys(ctx context.Context, req ListAPIKeysRequest) (*ListAPIKeysResponse, error)
	CreateAPIKey(ctx context.Context, req CreateAPIKeyRequest) (*CreateAPIKeyResponse, error)
	DeleteAPIKey(ctx context.Context, req DeleteAPIKeyRequest) (*DeleteAPIKeyResponse, error)
}
