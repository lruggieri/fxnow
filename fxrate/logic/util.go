package logic

import "context"

const (
	ContextKeyAPIKey = "api-key"
)

func GetAPIKeyIDFromContext(c context.Context) string {
	if c == nil {
		return ""
	}

	v, ok := c.Value(ContextKeyAPIKey).(string)
	if !ok {
		return ""
	}

	return v
}
