package util

import (
	"github.com/lithammer/shortuuid/v4"
)

func NewUuid() string {
	return shortuuid.New()
}
