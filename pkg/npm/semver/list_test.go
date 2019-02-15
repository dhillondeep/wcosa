package semver

import (
	"testing"

	"github.com/blang/semver"
	"github.com/stretchr/testify/assert"
)

func TestList_Len(t *testing.T) {
	list := List{
		&semver.Version{Major: 1, Minor: 2, Patch: 3},
		&semver.Version{Major: 3, Minor: 2, Patch: 1},
		&semver.Version{Major: 4, Minor: 5, Patch: 6},
	}
	assert.Equal(t, list.Len(), len(list))
	assert.Equal(t, list.Len(), 3)
}

func TestList_Swap(t *testing.T) {
	list := List{
		&semver.Version{},
		&semver.Version{Major: 1, Minor: 1, Patch: 1},
	}
	assert.Equal(t, list[0].String(), "0.0.0")
	assert.Equal(t, list[1].String(), "1.1.1")

	list.Swap(0, 1)
	assert.Equal(t, list[0].String(), "1.1.1")
	assert.Equal(t, list[1].String(), "0.0.0")
}

func TestList_Less(t *testing.T) {
	list := List{
		&semver.Version{},
		&semver.Version{Major: 1, Minor: 1, Patch: 1},
	}
	assert.True(t, list.Less(0, 1))
	assert.False(t, list.Less(1, 0))
}

func TestList_Sort(t *testing.T) {
	expected := List{
		&semver.Version{},                             // 0
		&semver.Version{Patch: 5},                     // 1
		&semver.Version{Minor: 1},                     // 2
		&semver.Version{Minor: 2, Patch: 2},           // 3
		&semver.Version{Minor: 2, Patch: 3},           // 4
		&semver.Version{Major: 1},                     // 5
		&semver.Version{Major: 2},                     // 6
		&semver.Version{Major: 2, Minor: 2},           // 7
		&semver.Version{Major: 2, Minor: 2, Patch: 6}, // 8
		&semver.Version{Major: 5},                     // 9
		&semver.Version{Major: 5, Minor: 6},           // 10
	}
	scramble := []int{10, 5, 4, 7, 8, 1, 9, 0, 3, 6, 2}
	list := make(List, 0, len(scramble))
	for _, i := range scramble {
		list = append(list, expected[i])
	}
	list.Sort()
	for i, ver := range list {
		assert.Equal(t, ver, expected[i])
	}
}
