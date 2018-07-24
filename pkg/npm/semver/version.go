package semver

import (
    "fmt"
    "regexp"
    "strconv"
    "strings"
)

var match = regexp.MustCompile(`^=?v?[0-9]+\.[0-9]+\.[0-9]+(-[0-9a-zA-Z]+(\.[0-9]+)?)?$`)

// Version describes a string literal in the form
// [Query][Major].[Minor].[Patch]-[Tag].[Num]
type Version struct {
    Major int
    Minor int
    Patch int
}

func IsValid(str string) bool {
    return match.MatchString(str)
}

func Parse(str string) *Version {
    if !IsValid(str) {
        return nil
    }
    if str[0] == '=' {
        str = str[1:]
    }
    if str[0] == 'v' {
        str = str[1:]
    }

    var ret Version
    vers := strings.Split(str, ".")
    ret.Major, _ = strconv.Atoi(vers[0])
    ret.Minor, _ = strconv.Atoi(vers[1])
    ret.Patch, _ = strconv.Atoi(vers[2])
    return &ret
}

func (a *Version) Str() string {
    return fmt.Sprintf("%d.%d.%d", a.Major, a.Minor, a.Patch)
}

func (a *Version) eq(b *Version) bool {
    return a.Major == b.Major &&
        a.Minor == b.Minor &&
        a.Patch == b.Patch
}

func (a *Version) less(b *Version) bool {
    if a.Major != b.Major {
        return a.Major < b.Major
    }
    if a.Minor != b.Minor {
        return a.Minor < b.Minor
    }
    return a.Patch < b.Patch
}
