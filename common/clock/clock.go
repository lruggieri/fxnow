package clock

import (
	"time"
)

type Clock interface {
	Now() time.Time
}

type Default struct{}

func (Default) Now() time.Time {
	return time.Now()
}
