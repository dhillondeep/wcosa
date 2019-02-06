package semver

import (
    "testing"

    "github.com/blang/semver"

    "github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
    assert.Equal(t, Parse("4.3.2"), &semver.Version{Major: 4, Minor: 3, Patch: 2})
    assert.Equal(t, Parse("2.2.2"), &semver.Version{Major: 2, Minor: 2, Patch: 2})
    assert.Equal(t, Parse("0.3.3"), &semver.Version{Major: 0, Minor: 3, Patch: 3})

    assert.Nil(t, Parse("0.3"))
    assert.Nil(t, Parse("3.3"))
}

func TestVersion_Str(t *testing.T) {
    assert.Equal(t, "4.2.12", (&semver.Version{Major: 4, Minor: 2, Patch: 12}).String())
    assert.Equal(t, "0.0.2", (&semver.Version{Patch: 2}).String())
    assert.Equal(t, "0.0.0", (&semver.Version{}).String())
}

func TestVersion_eq(t *testing.T) {
    var a *semver.Version
    var b *semver.Version

    a = &semver.Version{Major: 4, Minor: 5, Patch: 6}
    b = &semver.Version{Major: 4, Minor: 5, Patch: 6}
    assert.True(t, a.EQ(*b))
    assert.True(t, b.EQ(*a))

    a = &semver.Version{Major: 4, Minor: 4, Patch: 5}
    b = &semver.Version{Major: 4, Minor: 4, Patch: 6}
    assert.False(t, a.EQ(*b))
    assert.False(t, b.EQ(*a))
}

func TestVersion_less(t *testing.T) {
    var a *semver.Version
    var b *semver.Version

    a = &semver.Version{Major: 5}
    b = &semver.Version{Major: 4}
    assert.True(t, b.LT(*a))
    assert.False(t, a.LT(*b))

    b = &semver.Version{Major: 4, Minor: 99, Patch: 99}
    assert.True(t, b.LT(*a))
    assert.False(t, a.LT(*b))

    a = &semver.Version{Major: 4, Minor: 100, Patch: 99}
    assert.True(t, b.LT(*a))
    assert.False(t, a.LT(*b))

    a = &semver.Version{Major: 4, Minor: 89, Patch: 150}
    assert.False(t, b.LT(*a))
    assert.True(t, a.LT(*b))

    a = b
    assert.False(t, a.LT(*b))
    assert.False(t, b.LT(*a))

    assert.False(t, a.LT(*a))
    assert.False(t, b.LT(*b))
}
