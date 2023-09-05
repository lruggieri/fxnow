package model

const (
	APIKeyTypeUndefined APIKeyType = iota
	APIKeyTypeUnlimited
	APIKeyTypeLimited
)

type APIKeyType uint8

func (akt APIKeyType) Uint8() uint8 {
	return uint8(akt)
}

func (akt APIKeyType) String() string {
	switch akt {
	case APIKeyTypeUnlimited:
		return "unlimited"
	case APIKeyTypeLimited:
		return "limited"
	default:
		return "undefined"
	}
}

type APIKey struct {
	ID         uint64     `json:"id"`
	APIKeyID   string     `json:"api_key"`
	UserID     string     `json:"user_id"`
	Type       APIKeyType `json:"type"`
	Expiration int64      `json:"expiration"`

	User *User
}
