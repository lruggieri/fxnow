package util

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_All(t *testing.T) {
	ages := []int{2, 4, 6}

	assert.Equal(t, true, All(ages, func(item int) bool {
		return item%2 == 0
	}))

	assert.Equal(t, false, All(ages, func(item int) bool {
		return item%5 == 0
	}))
}

func Test_Some(t *testing.T) {
	ages := []int{2, 4, 6}

	assert.Equal(t, true, Some(ages, func(item int) bool {
		return item%3 == 0
	}))

	assert.Equal(t, false, Some(ages, func(item int) bool {
		return item%5 == 0
	}))
}

func Test_Map(t *testing.T) {
	nums := Map([]int{1, 3, 5}, func(item int) int {
		return item * 2
	})

	assert.Equal(t, []int{2, 6, 10}, nums)
}

func Test_Intersect(t *testing.T) {
	items := Intersect([]string{"1", "a", "c"}, []string{"4", "a", "b", "c"})

	assert.Equal(t, []string{"a", "c"}, items)
}

func Test_Reduce(t *testing.T) {
	sum := Reduce([]int{1, 4, 5}, func(prev int, item int, index int) int {
		return prev + item
	}, 0)
	assert.Equal(t, 10, sum)

	type Person struct {
		ID   int64
		Name string
	}

	persons := []Person{
		{
			ID:   1,
			Name: "John Doe",
		},
		{
			ID:   2,
			Name: "emiliy",
		},
	}

	indexed := Reduce(
		persons,
		func(prev map[int64]Person, item Person, index int) map[int64]Person {
			prev[item.ID] = item
			return prev
		},
		make(map[int64]Person),
	)
	assert.Equal(t,
		map[int64]Person{
			1: persons[0],
			2: persons[1],
		},
		indexed,
	)
}

func Test_SliceContains(t *testing.T) {
	assert.True(t,
		SliceContains([]int{1, 2, 5, 10}, []int{1, 5, 2}),
	)

	assert.False(t,
		SliceContains([]int{1, 2, 5, 10}, []int{1, 5, 2, 11}),
	)
}

func Test_SliceUnorderedEqual(t *testing.T) {
	assert.True(t,
		SliceUnorderedEqual([]int{1, 2, 5, 10}, []int{2, 1, 10, 5}),
	)

	assert.False(t,
		SliceUnorderedEqual([]int{1, 2, 5, 10}, []int{1, 5, 2}),
	)
}

func Test_Without(t *testing.T) {
	assert.Equal(t,
		[]int{1, 3},
		Without(
			[]int{1, 2, 3, 4, 5},
			[]int{2, 4, 5},
		),
	)

	assert.Equal(t,
		[]int{},
		Without(
			[]int{1, 2, 3, 4, 5},
			[]int{1, 2, 3, 4, 5},
		),
	)
}

func Test_ContainsFn(t *testing.T) {
	equal := func(a, b int) bool {
		return a == b
	}

	assert.True(t, ContainsFn([]int{1, 5, 3, 4}, 4, equal))
	assert.False(t, ContainsFn([]int{1, 5, 3, 4}, 2, equal))
}

func Test_DistinctFn(t *testing.T) {
	equal := func(a, b int) bool {
		return a == b
	}

	assert.Equal(t,
		[]int{1, 5, 3, 4},
		DistinctFn([]int{1, 5, 3, 5, 4}, equal),
	)

	assert.Equal(t,
		[]int{},
		DistinctFn([]int{}, equal),
	)

	assert.Equal(t,
		[]int{1, 2, 3, 4},
		DistinctFn([]int{1, 2, 3, 4}, equal),
	)
}

func Test_MapWithError(t *testing.T) {
	items, err := MapWithError(
		[]int{1, 5, 3},
		func(item int) (int, error) {
			return 0, errors.New("test")
		},
	)
	assert.Error(t, err)
	assert.Nil(t, items)

	items, err = MapWithError(
		[]int{1, 5, 3},
		func(item int) (int, error) {
			return item * 2, nil
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, []int{2, 10, 6}, items)
}

func Test_Filter(t *testing.T) {
	// only elements that are not an empty string once trimmed
	assert.Equal(t,
		[]string{"elem1", "elem2"},
		Filter(Map([]string{" elem1", "  ", "", "elem2 "}, func(item string) string {
			return strings.TrimSpace(item)
		}), func(item string) bool {
			return item != ""
		}),
	)

	// only even numbers
	assert.Equal(t,
		[]int{0, 2, 4, 6, 8, 10},
		Filter([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, func(item int) bool {
			return item%2 == 0
		}),
	)
}

func Test_PruneSlice(t *testing.T) {
	assert.Equal(t,
		[]int{1, 2, 4},
		PruneSlice([]int{1, 2, 0, 4, 0}),
	)

	assert.Equal(t,
		[]int{1, 2, 3, 4},
		PruneSlice([]int{1, 2, 3, 4}),
	)

	type A struct{}

	assert.Equal(t,
		[]*A{{}, {}},
		PruneSlice([]*A{nil, {}, nil, {}}),
	)
}

func Test_SliceLastElement(t *testing.T) {
	assert.Equal(t,
		2,
		SliceLastElement([]int{1, 2}),
	)

	assert.Panics(t, func() { SliceLastElement([]int{}) })
}

func TestSliceFirstElement(t *testing.T) {
	assert.Equal(t,
		1,
		SliceFirstElement([]int{1, 2}),
	)

	assert.Panics(t, func() { SliceFirstElement([]int{}) })
}

func TestSliceToMap(t *testing.T) {
	var nilSlice []string

	assert.Equal(t, (map[string]struct{})(nil), SliceToMap(nilSlice))

	assert.Equal(t, map[string]struct{}{
		"a": {},
		"b": {},
		"c": {},
	}, SliceToMap([]string{
		"a",
		"b",
		"c",
	}))
}
