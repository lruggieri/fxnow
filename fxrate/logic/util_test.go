package logic

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAPIKeyIDFromContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), ContextKeyAPIKey, "api_key")

	assert.Equal(t, "", GetAPIKeyIDFromContext(nil))
	assert.Equal(t, "", GetAPIKeyIDFromContext(context.Background()))
	assert.Equal(t, "api_key", GetAPIKeyIDFromContext(ctx))
}
