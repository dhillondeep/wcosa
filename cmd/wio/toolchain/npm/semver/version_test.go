package semver

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestIsValid(t *testing.T) {
    assert.True(t, IsValid("1.0.0"))
    assert.True(t, IsValid("v4.23.12"))
    assert.True(t, IsValid("=55.632.123"))
    assert.True(t, IsValid("=v33.33.432"))

    assert.False(t, IsValid("3.3"))
    assert.False(t, IsValid(">22.22.22"))
    assert.False(t, IsValid("3"))
    assert.False(t, IsValid("66.66.3sx"))
}

func TestParse(t *testing.T) {
    assert.Equal(t, Parse("4.3.2"), &Version{4, 3, 2})
    assert.Equal(t, Parse("2.2.2"), &Version{2, 2, 2})
    assert.Equal(t, Parse("0.3.3"), &Version{0, 3, 3})

    assert.Nil(t, Parse("0.3"))
    assert.Nil(t, Parse("3.3"))
}

func TestVersion_Str(t *testing.T) {
    assert.Equal(t, "4.2.12", (&Version{4, 2, 12}).Str())
    assert.Equal(t, "0.0.2", (&Version{0, 0, 2}).Str())
    assert.Equal(t, "0.0.0", (&Version{}).Str())
}

func TestVersion_eq(t *testing.T) {
    var a *Version
    var b *Version

    a = &Version{4, 5, 6}
    b = &Version{4, 5, 6}
    assert.True(t, a.eq(b))
    assert.True(t, b.eq(a))

    a = &Version{4, 4, 5}
    b = &Version{4, 4, 6}
    assert.False(t, a.eq(b))
    assert.False(t, b.eq(a))
}

func TestVersion_less(t *testing.T) {
    var a *Version
    var b *Version

    a = &Version{5, 0, 0}
    b = &Version{4, 0, 0}
    assert.True(t, b.less(a))
    assert.False(t, a.less(b))

    b = &Version{4, 99, 99}
    assert.True(t, b.less(a))
    assert.False(t, a.less(b))

    a = &Version{4, 100, 99}
    assert.True(t, b.less(a))
    assert.False(t, a.less(b))

    a = &Version{4, 89, 150}
    assert.False(t, b.less(a))
    assert.True(t, a.less(b))

    a = b
    assert.False(t, a.less(b))
    assert.False(t, b.less(a))

    assert.False(t, a.less(a))
    assert.False(t, b.less(b))
}
