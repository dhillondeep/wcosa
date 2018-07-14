package semver

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestList_Len(t *testing.T) {
    list := List{
        &Version{1, 2, 3},
        &Version{3, 2, 1},
        &Version{4, 5, 6},
    }
    assert.Equal(t, list.Len(), len(list))
    assert.Equal(t, list.Len(), 3)
}

func TestList_Swap(t *testing.T) {
    list := List{
        &Version{0, 0, 0},
        &Version{1, 1, 1},
    }
    assert.Equal(t, list[0].Str(), "0.0.0")
    assert.Equal(t, list[1].Str(), "1.1.1")

    list.Swap(0, 1)
    assert.Equal(t, list[0].Str(), "1.1.1")
    assert.Equal(t, list[1].Str(), "0.0.0")
}

func TestList_Less(t *testing.T) {
    list := List{
        &Version{0, 0, 0},
        &Version{1, 1, 1},
    }
    assert.True(t, list.Less(0, 1))
    assert.False(t, list.Less(1, 0))
}

func TestList_Sort(t *testing.T) {
    expected := List{
        &Version{0, 0, 0}, // 0
        &Version{0, 0, 5}, // 1
        &Version{0, 1, 0}, // 2
        &Version{0, 2, 2}, // 3
        &Version{0, 2, 3}, // 4
        &Version{1, 0, 0}, // 5
        &Version{2, 0, 0}, // 6
        &Version{2, 2, 0}, // 7
        &Version{2, 2, 6}, // 8
        &Version{5, 0, 0}, // 9
        &Version{5, 6, 0}, // 10
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
