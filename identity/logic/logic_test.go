package logic

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	cError "github.com/lruggieri/fxnow/common/error"
	mockstore "github.com/lruggieri/fxnow/common/mock/store"
	"github.com/lruggieri/fxnow/common/model"
	"github.com/lruggieri/fxnow/common/store"
	"github.com/lruggieri/fxnow/identity/auth"
)

func TestImpl_CreateAPIKey(t *testing.T) {
	testErr := errors.New("error")

	type deps struct {
		store *mockstore.Store
	}

	type args struct {
		ctx context.Context
		req CreateAPIKeyRequest
	}

	uInfo := auth.UserInfo{
		Email:      "user@domain.com",
		GivenName:  "name",
		FamilyName: "surname",
	}
	uInfoCtx := context.WithValue(context.Background(), auth.ContextUserInfoKey, &uInfo)

	tests := []struct {
		name      string
		deps      deps
		args      args
		mock      func(args args, d deps)
		assertion func(
			t *testing.T,
			res *CreateAPIKeyResponse,
			err error,
		)
	}{
		{
			name: "error-no-user-info",
			args: args{
				ctx: context.Background(),
				req: CreateAPIKeyRequest{},
			},
			mock: func(args args, d deps) {},
			assertion: func(t *testing.T, res *CreateAPIKeyResponse, err error) {
				assert.Nil(t, res)
				assert.ErrorIs(t, err, cError.NotAuthenticated)
			},
		},
		{
			name: "error-get-user",
			args: args{
				ctx: uInfoCtx,
				req: CreateAPIKeyRequest{},
			},
			mock: func(args args, d deps) {
				d.store.EXPECT().GetUser(args.ctx, store.GetUserRequest{
					Email: uInfo.Email,
				}).Return(nil, testErr).Once()
			},
			assertion: func(t *testing.T, res *CreateAPIKeyResponse, err error) {
				assert.Nil(t, res)
				assert.ErrorIs(t, err, testErr)
			},
		},
		{
			name: "error-create-user",
			args: args{
				ctx: uInfoCtx,
				req: CreateAPIKeyRequest{},
			},
			mock: func(args args, d deps) {
				d.store.EXPECT().GetUser(args.ctx, store.GetUserRequest{
					Email: uInfo.Email,
				}).Return(nil, cError.NotFound).Once()

				d.store.EXPECT().CreateUser(args.ctx, store.CreateUserRequest{
					FirstName: uInfo.GivenName,
					LastName:  uInfo.FamilyName,
					Email:     uInfo.Email,
				}).Return(nil, testErr).Once()
			},
			assertion: func(t *testing.T, res *CreateAPIKeyResponse, err error) {
				assert.Nil(t, res)
				assert.ErrorIs(t, err, testErr)
			},
		},
		{
			name: "error-get-api-key",
			args: args{
				ctx: uInfoCtx,
				req: CreateAPIKeyRequest{},
			},
			mock: func(args args, d deps) {
				d.store.EXPECT().GetUser(args.ctx, store.GetUserRequest{
					Email: uInfo.Email,
				}).Return(&store.GetUserResponse{User: &model.User{
					UserID: "user_id",
				}}, nil).Once()

				d.store.EXPECT().GetAPIKey(args.ctx, store.GetAPIKeyRequest{UserID: "user_id"}).
					Return(nil, testErr).Once()
			},
			assertion: func(t *testing.T, res *CreateAPIKeyResponse, err error) {
				assert.Nil(t, res)
				assert.ErrorIs(t, err, testErr)
			},
		},
		{
			name: "error-create-api-key",
			args: args{
				ctx: uInfoCtx,
				req: CreateAPIKeyRequest{},
			},
			mock: func(args args, d deps) {
				d.store.EXPECT().GetUser(args.ctx, store.GetUserRequest{
					Email: uInfo.Email,
				}).Return(&store.GetUserResponse{User: &model.User{
					UserID: "user_id",
				}}, nil).Once()

				d.store.EXPECT().GetAPIKey(args.ctx, store.GetAPIKeyRequest{UserID: "user_id"}).
					Return(nil, cError.NotFound).Once()

				d.store.EXPECT().CreateAPIKey(args.ctx, store.CreateAPIKeyRequest{
					UserID: "user_id",
					Type:   model.APIKeyTypeLimited.Uint8(),
				}).Return(nil, testErr).Once()
			},
			assertion: func(t *testing.T, res *CreateAPIKeyResponse, err error) {
				assert.Nil(t, res)
				assert.ErrorIs(t, err, testErr)
			},
		},
		{
			name: "error-multiple-key-normal-user",
			args: args{
				ctx: uInfoCtx,
				req: CreateAPIKeyRequest{},
			},
			mock: func(args args, d deps) {
				d.store.EXPECT().GetUser(args.ctx, store.GetUserRequest{
					Email: uInfo.Email,
				}).Return(&store.GetUserResponse{User: &model.User{
					UserID: "user_id",
				}}, nil).Once()

				d.store.EXPECT().GetAPIKey(args.ctx, store.GetAPIKeyRequest{UserID: "user_id"}).
					Return(&store.GetAPIKeyResponse{APIKey: &model.APIKey{APIKeyID: "api_key"}}, nil).Once()
			},
			assertion: func(t *testing.T, res *CreateAPIKeyResponse, err error) {
				assert.Nil(t, res)
				assert.ErrorIs(t, err, cError.InvalidParameter)
			},
		},
		{
			name: "happy-path-user-exists",
			args: args{
				ctx: uInfoCtx,
				req: CreateAPIKeyRequest{},
			},
			mock: func(args args, d deps) {
				d.store.EXPECT().GetUser(args.ctx, store.GetUserRequest{
					Email: uInfo.Email,
				}).Return(&store.GetUserResponse{User: &model.User{
					UserID: "user_id",
				}}, nil).Once()

				d.store.EXPECT().GetAPIKey(args.ctx, store.GetAPIKeyRequest{UserID: "user_id"}).
					Return(nil, cError.NotFound).Once()

				d.store.EXPECT().CreateAPIKey(args.ctx, store.CreateAPIKeyRequest{
					UserID: "user_id",
					Type:   model.APIKeyTypeLimited.Uint8(),
				}).Return(&store.CreateAPIKeyResponse{
					APIKeyID: "api_key",
				}, nil).Once()
			},
			assertion: func(t *testing.T, res *CreateAPIKeyResponse, err error) {
				assert.Nil(t, err)
				assert.Equal(t, &CreateAPIKeyResponse{APIKeyID: "api_key"}, res)
			},
		},
		{
			name: "happy-path-user-not-exist",
			args: args{
				ctx: uInfoCtx,
				req: CreateAPIKeyRequest{},
			},
			mock: func(args args, d deps) {
				d.store.EXPECT().GetUser(args.ctx, store.GetUserRequest{
					Email: uInfo.Email,
				}).Return(nil, cError.NotFound).Once()

				d.store.EXPECT().CreateUser(args.ctx, store.CreateUserRequest{
					FirstName: uInfo.GivenName,
					LastName:  uInfo.FamilyName,
					Email:     uInfo.Email,
				}).Return(&store.CreateUserResponse{
					UserID: "user_id",
				}, nil).Once()

				d.store.EXPECT().GetAPIKey(args.ctx, store.GetAPIKeyRequest{UserID: "user_id"}).
					Return(nil, cError.NotFound).Once()

				d.store.EXPECT().CreateAPIKey(args.ctx, store.CreateAPIKeyRequest{
					UserID: "user_id",
					Type:   model.APIKeyTypeLimited.Uint8(),
				}).Return(&store.CreateAPIKeyResponse{
					APIKeyID: "api_key",
				}, nil).Once()
			},
			assertion: func(t *testing.T, res *CreateAPIKeyResponse, err error) {
				assert.Nil(t, err)
				assert.Equal(t, &CreateAPIKeyResponse{APIKeyID: "api_key"}, res)
			},
		},
	}

	for _, tt := range tests {
		tc := tt // avoid loop closure issue
		t.Run(tc.name, func(t *testing.T) {
			d := deps{
				store: mockstore.NewStore(t),
			}

			l := Impl{
				Store: d.store,
			}

			tc.mock(tc.args, d)

			res, err := l.CreateAPIKey(tc.args.ctx, tt.args.req)

			tc.assertion(t, res, err)
		})
	}

}

func TestImpl_DeleteAPIKey(t *testing.T) {
	testErr := errors.New("error")

	type deps struct {
		store *mockstore.Store
	}

	type args struct {
		ctx context.Context
		req DeleteAPIKeyRequest
	}

	uInfo := auth.UserInfo{
		Email:      "user@domain.com",
		GivenName:  "name",
		FamilyName: "surname",
	}
	uInfoCtx := context.WithValue(context.Background(), auth.ContextUserInfoKey, &uInfo)

	tests := []struct {
		name      string
		deps      deps
		args      args
		mock      func(args args, d deps)
		assertion func(
			t *testing.T,
			res *DeleteAPIKeyResponse,
			err error,
		)
	}{
		{
			name: "error-no-user-info",
			args: args{
				ctx: context.Background(),
				req: DeleteAPIKeyRequest{APIKeyID: "api_key"},
			},
			mock: func(args args, d deps) {},
			assertion: func(t *testing.T, res *DeleteAPIKeyResponse, err error) {
				assert.Nil(t, res)
				assert.ErrorIs(t, err, cError.NotAuthenticated)
			},
		},
		{
			name: "error-get-user",
			args: args{
				ctx: uInfoCtx,
				req: DeleteAPIKeyRequest{APIKeyID: "api_key"},
			},
			mock: func(args args, d deps) {
				d.store.EXPECT().GetUser(args.ctx, store.GetUserRequest{
					Email: uInfo.Email,
				}).Return(nil, testErr).Once()
			},
			assertion: func(t *testing.T, res *DeleteAPIKeyResponse, err error) {
				assert.Nil(t, res)
				assert.ErrorIs(t, err, testErr)
			},
		},
		{
			name: "error-get-api-key",
			args: args{
				ctx: uInfoCtx,
				req: DeleteAPIKeyRequest{APIKeyID: "api_key"},
			},
			mock: func(args args, d deps) {
				d.store.EXPECT().GetUser(args.ctx, store.GetUserRequest{
					Email: uInfo.Email,
				}).Return(&store.GetUserResponse{User: &model.User{
					UserID: "user_id",
				}}, nil).Once()

				d.store.EXPECT().GetAPIKey(args.ctx, store.GetAPIKeyRequest{UserID: "user_id"}).
					Return(nil, testErr).Once()
			},
			assertion: func(t *testing.T, res *DeleteAPIKeyResponse, err error) {
				assert.Nil(t, res)
				assert.ErrorIs(t, err, testErr)
			},
		},
		{
			name: "error-delete-api-key",
			args: args{
				ctx: uInfoCtx,
				req: DeleteAPIKeyRequest{APIKeyID: "api_key"},
			},
			mock: func(args args, d deps) {
				d.store.EXPECT().GetUser(args.ctx, store.GetUserRequest{
					Email: uInfo.Email,
				}).Return(&store.GetUserResponse{User: &model.User{
					UserID: "user_id",
				}}, nil).Once()

				d.store.EXPECT().GetAPIKey(args.ctx, store.GetAPIKeyRequest{UserID: "user_id"}).
					Return(&store.GetAPIKeyResponse{
						APIKey: &model.APIKey{
							APIKeyID: "api_key",
							UserID:   "user_id",
						}}, nil).Once()

				d.store.EXPECT().DeleteAPIKey(args.ctx, store.DeleteAPIKeyRequest{APIKeyID: "api_key"}).
					Return(nil, testErr).Once()
			},
			assertion: func(t *testing.T, res *DeleteAPIKeyResponse, err error) {
				assert.Nil(t, res)
				assert.ErrorIs(t, err, testErr)
			},
		},
		{
			// only the API Key owners can delete their own key
			name: "error-wrong-user",
			args: args{
				ctx: uInfoCtx,
				req: DeleteAPIKeyRequest{APIKeyID: "api_key"},
			},
			mock: func(args args, d deps) {
				d.store.EXPECT().GetUser(args.ctx, store.GetUserRequest{
					Email: uInfo.Email,
				}).Return(&store.GetUserResponse{User: &model.User{
					UserID: "user_id",
				}}, nil).Once()

				d.store.EXPECT().GetAPIKey(args.ctx, store.GetAPIKeyRequest{UserID: "user_id"}).
					Return(&store.GetAPIKeyResponse{
						APIKey: &model.APIKey{
							APIKeyID: "api_key",
							UserID:   "user_id_2",
						}}, nil).Once()
			},
			assertion: func(t *testing.T, res *DeleteAPIKeyResponse, err error) {
				assert.Nil(t, res)
				assert.ErrorIs(t, err, cError.NotAuthorized)
			},
		},
		{
			name: "happy-path",
			args: args{
				ctx: uInfoCtx,
				req: DeleteAPIKeyRequest{APIKeyID: "api_key"},
			},
			mock: func(args args, d deps) {
				d.store.EXPECT().GetUser(args.ctx, store.GetUserRequest{
					Email: uInfo.Email,
				}).Return(&store.GetUserResponse{User: &model.User{
					UserID: "user_id",
				}}, nil).Once()

				d.store.EXPECT().GetAPIKey(args.ctx, store.GetAPIKeyRequest{UserID: "user_id"}).
					Return(&store.GetAPIKeyResponse{
						APIKey: &model.APIKey{
							APIKeyID: "api_key",
							UserID:   "user_id",
						}}, nil).Once()

				d.store.EXPECT().DeleteAPIKey(args.ctx, store.DeleteAPIKeyRequest{APIKeyID: "api_key"}).
					Return(&store.DeleteAPIKeyResponse{}, nil).Once()
			},
			assertion: func(t *testing.T, res *DeleteAPIKeyResponse, err error) {
				assert.Nil(t, err)
				assert.Equal(t, &DeleteAPIKeyResponse{}, res)
			},
		},
	}

	for _, tt := range tests {
		tc := tt // avoid loop closure issue
		t.Run(tc.name, func(t *testing.T) {
			d := deps{
				store: mockstore.NewStore(t),
			}

			l := Impl{
				Store: d.store,
			}

			tc.mock(tc.args, d)

			res, err := l.DeleteAPIKey(tc.args.ctx, tt.args.req)

			tc.assertion(t, res, err)
		})
	}

}
