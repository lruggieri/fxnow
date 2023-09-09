package util

import (
	"github.com/lithammer/shortuuid/v4"
)

func NewUUID() string {
	return shortuuid.New()
}
