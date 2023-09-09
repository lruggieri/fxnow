package logic

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAPIKeyIDFromContext(t *testing.T) {
	apiKey := "api_key"

	ctx := context.WithValue(context.Background(), ContextKeyAPIKey, apiKey)

	//nolint:staticcheck
	assert.Equal(t, "", GetAPIKeyIDFromContext(nil))
	assert.Equal(t, "", GetAPIKeyIDFromContext(context.Background()))
	assert.Equal(t, apiKey, GetAPIKeyIDFromContext(ctx))
}
