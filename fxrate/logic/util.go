package logic

import (
	"context"

	"github.com/lruggieri/fxnow/common/util"
)

const (
	ContextKeyAPIKey util.ContextKey = "api-key"
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
