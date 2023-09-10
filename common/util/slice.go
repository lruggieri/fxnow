package util

func All[T any](items []T, predicate func(item T) bool) bool {
	for i := range items {
		if !predicate(items[i]) {
			return false
		}
	}

	return true
}

func Some[T any](items []T, predicate func(item T) bool) bool {
	for i := range items {
		if predicate(items[i]) {
			return true
		}
	}

	return false
}

func Map[T any, K any](items []T, mapFunc func(item T) K) []K {
	result := make([]K, len(items))

	for i := range items {
		result[i] = mapFunc(items[i])
	}

	return result
}

func MapWithError[T any, K any](items []T, mapFunc func(item T) (K, error)) ([]K, error) {
	var (
		err    error
		result = make([]K, len(items))
	)

	for i := range items {
		result[i], err = mapFunc(items[i])
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func Intersect[T comparable](a, b []T) []T {
	intersect := make([]T, 0)

	for i := range a {
		for j := range b {
			if a[i] == b[j] {
				intersect = append(intersect, a[i])
				break
			}
		}
	}

	return intersect
}

// Without returns a new list without values in the first argument.
func Without[T comparable](items, without []T) []T {
	return Reduce(items, func(prev []T, item T, index int) []T {
		if Contains(without, item) {
			return prev
		}

		return append(prev, item)
	}, make([]T, 0))
}

func Reduce[T any, R any](items []T, reducer func(prev R, item T, index int) R, initialVal R) R {
	for i := range items {
		initialVal = reducer(initialVal, items[i], i)
	}

	return initialVal
}

func IndexOf[T comparable](items []T, item T) int {
	for i := range items {
		if items[i] == item {
			return i
		}
	}

	return -1
}

func IndexOfFn[T any](items []T, item T, equal func(a, b T) bool) int {
	for i := range items {
		if equal(items[i], item) {
			return i
		}
	}

	return -1
}

func Contains[T comparable](container []T, contain T) bool {
	return IndexOf(container, contain) > -1
}

func ContainsFn[T any](container []T, contain T, equal func(a, b T) bool) bool {
	return IndexOfFn(container, contain, equal) > -1
}

func SliceContains[T comparable](container []T, items []T) bool {
	return All(items, func(item T) bool {
		// container contains all element of items
		return Contains(container, item)
	})
}

// DistinctFn distince a slice with compare function
//
//nolint:revive
func DistinctFn[T any](items []T, equal func(a, b T) bool) []T {
	newItems := make([]T, 0)

	for _, a := range items {
		if !ContainsFn(newItems, a, equal) {
			newItems = append(newItems, a)
		}
	}

	return newItems
}

func PruneSlice[T comparable](src []T) []T {
	var (
		empty  T
		items2 = make([]T, 0)
	)

	for _, a := range src {
		if a != empty {
			items2 = append(items2, a)
		}
	}

	return items2
}

func Filter[T any](items []T, filterFn func(item T) bool) []T {
	newItems := make([]T, 0, len(items))

	for _, item := range items {
		if filterFn(item) {
			newItems = append(newItems, item)
		}
	}

	return newItems
}

func SliceLastElement[T any](items []T) T {
	return items[len(items)-1]
}

func SliceFirstElement[T any](items []T) T {
	return items[0]
}

// SliceUnorderedEqual : returns true if 2 slices have exactly the same elements, even if not in the same order
func SliceUnorderedEqual[T comparable](slice1, slice2 []T) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	return SliceContains(slice1, slice2)
}

func SliceToMap[T comparable](in []T) map[T]struct{} {
	if in == nil {
		return nil
	}

	res := make(map[T]struct{}, len(in))

	for _, e := range in {
		res[e] = struct{}{}
	}

	return res
}
