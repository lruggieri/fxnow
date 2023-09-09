package logic

import (
	"context"

	"github.com/pkg/errors"

	cError "github.com/lruggieri/fxnow/common/error"
	"github.com/lruggieri/fxnow/common/model"
	"github.com/lruggieri/fxnow/common/store"

	"github.com/lruggieri/fxnow/identity/auth"
)

type Logic interface {
	CreateAPIKey(context.Context, CreateAPIKeyRequest) (*CreateAPIKeyResponse, error)
	DeleteAPIKey(context.Context, DeleteAPIKeyRequest) (*DeleteAPIKeyResponse, error)
}

type Impl struct {
	Store store.Store
}

func (i *Impl) CreateAPIKey(ctx context.Context, _ CreateAPIKeyRequest) (*CreateAPIKeyResponse, error) {
	uInfo := auth.GetUserInfoFromContext(ctx)
	if uInfo == nil {
		return nil, cError.ErrNotAuthenticated
	}

	// create user if it doesn't exist
	uRes, err := i.createUser(ctx, CreateUserRequest{
		FirstName: uInfo.GivenName,
		LastName:  uInfo.FamilyName,
		Email:     uInfo.Email,
	})
	if err != nil {
		return nil, err
	}

	apiKey, err := i.Store.GetAPIKey(ctx, store.GetAPIKeyRequest{UserID: uRes.UserID})
	if err != nil && err != cError.ErrNotFound {
		return nil, err
	}

	// TODO Admin users can have multiple API keys

	if apiKey != nil {
		return nil, errors.Wrap(cError.ErrInvalidParameter, "users can only have 1 active API key")
	}

	// create API key
	akRes, err := i.Store.CreateAPIKey(ctx, store.CreateAPIKeyRequest{
		UserID: uRes.UserID,
		Type:   model.APIKeyTypeLimited.Uint8(),
	})
	if err != nil {
		return nil, err
	}

	return &CreateAPIKeyResponse{
		APIKeyID: akRes.APIKeyID,
	}, nil
}

func (i *Impl) DeleteAPIKey(ctx context.Context, req DeleteAPIKeyRequest) (*DeleteAPIKeyResponse, error) {
	uInfo := auth.GetUserInfoFromContext(ctx)
	if uInfo == nil {
		return nil, cError.ErrNotAuthenticated
	}

	dbUInfo, err := i.Store.GetUser(ctx, store.GetUserRequest{
		Email: uInfo.Email,
	})
	if err != nil {
		return nil, err
	}

	apiKey, err := i.Store.GetAPIKey(ctx, store.GetAPIKeyRequest{UserID: dbUInfo.User.UserID})
	if err != nil {
		return nil, err
	}

	// only the API Key owners can delete their own key
	if apiKey.APIKey.UserID != dbUInfo.User.UserID {
		return nil, errors.Wrap(cError.ErrNotAuthorized, "only API Key owners can delete their own key")
	}

	_, err = i.Store.DeleteAPIKey(ctx, store.DeleteAPIKeyRequest{
		APIKeyID: req.APIKeyID,
	})
	if err != nil {
		return nil, err
	}

	return &DeleteAPIKeyResponse{}, nil
}

// createUser : idempotent call, create user if it doesn't already exist
func (i *Impl) createUser(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error) {
	user, err := i.Store.GetUser(ctx, store.GetUserRequest{
		Email: req.Email,
	})

	var userID string

	if errors.Is(err, cError.ErrNotFound) {
		var res *store.CreateUserResponse

		res, err = i.Store.CreateUser(ctx, store.CreateUserRequest{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
		})

		if err != nil {
			return nil, err
		}

		userID = res.UserID
	} else if err != nil {
		return nil, err
	} else {
		userID = user.User.UserID
	}

	return &CreateUserResponse{UserID: userID}, nil
}
