package clock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Default_Now(t *testing.T) {
	now := time.Now()
	r := Default{}.Now()
	assert.False(t, r.IsZero())
	assert.True(t, r.Equal(now) || r.After(now))
}
