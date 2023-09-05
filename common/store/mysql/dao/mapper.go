package dao

import (
	"github.com/lruggieri/fxnow/common/model"
	"github.com/lruggieri/fxnow/common/store/mysql/util"
)

func APIKeyToModel(in *APIKey) *model.APIKey {
	if in == nil {
		return nil
	}

	return &model.APIKey{
		ID:         in.ID,
		APIKeyID:   in.APIKeyID,
		UserID:     in.UserID,
		Expiration: util.SQLTimeToUnix(in.Expiration),

		User: UserToModel(in.User),
	}
}

func UserToModel(in *User) *model.User {
	if in == nil {
		return nil
	}

	return &model.User{
		ID:        in.ID,
		UserID:    in.UserID,
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Email:     in.Email,
	}
}
