package util

func MapMultipleItems[A any, B any](mapperFunc func(A) B, items []A) []B {
	if items == nil {
		return nil
	}

	result := make([]B, len(items))
	for i, item := range items {
		result[i] = mapperFunc(item)
	}

	return result
}
