package util

type ContextKey string

func (ck ContextKey) String() string {
	return string(ck)
}
