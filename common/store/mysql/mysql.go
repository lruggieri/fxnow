package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	cError "github.com/lruggieri/fxnow/common/error"
	"github.com/lruggieri/fxnow/common/model"
	"github.com/lruggieri/fxnow/common/store"
	"github.com/lruggieri/fxnow/common/store/mysql/dao"
	"github.com/lruggieri/fxnow/common/util"
)

type MySQL struct {
	db *gorm.DB
}

func (m *MySQL) GetUser(ctx context.Context, req store.GetUserRequest) (*store.GetUserResponse, error) {
	tx := m.db.Model(&dao.User{})
	if req.UserID != "" {
		tx = tx.Where("`user_id` = ?", req.UserID)
	}
	if req.Email != "" {
		tx = tx.Where("`email` = ?", req.Email)
	}

	var res dao.User
	if tx = tx.First(&res); tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, cError.NotFound
		}

		return nil, tx.Error
	}

	return &store.GetUserResponse{
		User: dao.UserToModel(&res),
	}, nil
}

func (m *MySQL) CreateUser(ctx context.Context, req store.CreateUserRequest) (*store.CreateUserResponse, error) {
	d := &dao.User{
		UserID:    util.NewUuid(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}
	if tx := m.db.Create(d); tx.Error != nil {
		return nil, tx.Error
	}

	return &store.CreateUserResponse{
		UserID: d.UserID,
	}, nil
}

func (m *MySQL) GetAPIKey(ctx context.Context, req store.GetAPIKeyRequest) (*store.GetAPIKeyResponse, error) {
	tx := m.db.Model(&dao.APIKey{}).Where("`disabled` = 0")

	if req.UserID != "" {
		tx = tx.Where("`user_id` = ?", req.UserID)
	}
	if req.APIKeyID != "" {
		tx = tx.Where("`api_key_id` = ?", req.APIKeyID)
	}

	if req.WithUsages {
		tx = tx.Preload("Usages")
	}

	var res dao.APIKey
	if tx = tx.First(&res); tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, errors.Wrap(cError.NotFound, "API Key not found")
		}

		return nil, tx.Error
	}

	return &store.GetAPIKeyResponse{
		APIKey: dao.APIKeyToModel(&res),
	}, nil
}

func (m *MySQL) ListAPIKeys(ctx context.Context, req store.ListAPIKeysRequest) (*store.ListAPIKeysResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *MySQL) CreateAPIKey(ctx context.Context, req store.CreateAPIKeyRequest) (*store.CreateAPIKeyResponse, error) {
	if req.Type == model.APIKeyTypeUndefined.Uint8() {
		req.Type = model.APIKeyTypeLimited.Uint8()
	}

	d := &dao.APIKey{
		APIKeyID: util.NewUuid(),
		UserID:   req.UserID,
		Type:     req.Type,
	}
	if tx := m.db.Create(d); tx.Error != nil {
		return nil, tx.Error
	}

	return &store.CreateAPIKeyResponse{
		APIKeyID: d.APIKeyID,
	}, nil
}

func (m *MySQL) DeleteAPIKey(ctx context.Context, req store.DeleteAPIKeyRequest) (*store.DeleteAPIKeyResponse, error) {
	tx := m.db.Model(&dao.APIKey{}).
		Where("`api_key_id` = ?", req.APIKeyID).
		Updates(map[string]interface{}{
			"disabled":    true,
			"disabled_at": sql.NullTime{Time: time.Now(), Valid: true},
		})

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, cError.NotFound
	}

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &store.DeleteAPIKeyResponse{}, nil
}

type Config struct {
	Username string
	Password string
	Host     string
	Port     int
	DBName   string
}

func New(c Config) (store.Store, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		c.Username, c.Password, c.Host, c.Port, c.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &MySQL{
		db: db,
	}, nil
}
