package dao

import (
	"database/sql"
)

type APIKey struct {
	ID         uint64       `gorm:"column:id"`
	APIKeyID   string       `gorm:"column:api_key_id"`
	UserID     string       `gorm:"column:user_id"`
	Type       uint8        `gorm:"column:type"`
	Expiration sql.NullTime `gorm:"column:expiration"`

	User *User `gorm:"foreignKey:UserID;references:UserID"`
}

func (*APIKey) TableName() string {
	return "api_key"
}
