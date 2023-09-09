package logic

type CachedAPIKey struct {
	APIKeyID string              `json:"api-key-id"`
	Type     uint8               `json:"type"`
	Usages   []CachedAPIKeyUsage `json:"usages"`
}

type CachedAPIKeyUsage struct {
	Timestamp int64 `json:"timestamp"` // unix (s)
}

type CachedRate struct {
	Rate      float64 `json:"rate"`
	Timestamp int64   `json:"timestamp"`
}
