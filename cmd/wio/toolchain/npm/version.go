package npm

import (
    "fmt"
    "sort"
    "strconv"
    "strings"
    "wio/cmd/wio/errors"
)

type version struct {
    major int
    minor int
    patch int
}

type versionList []version

func (l versionList) Len() int {
    return len(l)
}

func (l versionList) Swap(i int, j int) {
    l[i], l[j] = l[j], l[i]
}

func (l versionList) Less(i int, j int) bool {
    return versionLess(l[i], l[j])
}

func versionLess(a version, b version) bool {
    if a.major != b.major {
        return a.major < b.major
    }
    if a.minor != b.minor {
        return a.minor < b.minor
    }
    return a.patch < b.patch
}

func iabs(a int64) int64 {
    if a < 0 {
        return -a
    }
    return a
}

func (l versionList) findAtLeast(min version) (version, error) {
    // assumes versionList sorted ascending and non-empty
    var ret version
    max := l[l.Len()-1]
    if versionLess(max, min) {
        return ret, errors.Stringf("failed to find version at least: %s", vertostr(min))
    }
    return max, nil
}

func (l versionList) findNearest(near version) version {
    // assumes versionList sorted ascending and non-empty
    dist := iabs(l[0].mag() - near.mag())
    index := 0
    for i, val := range l {
        curDist := iabs(val.mag() - near.mag())
        if curDist < dist {
            dist = curDist
            index = i
        }
        if curDist > dist {
            break
        }
    }
    return l[index]
}

func (v *version) mag() int64 {
    // assumes that no version value exceeds 1 << 20
    major := uint64(v.major)
    minor := uint64(v.minor)
    patch := uint64(v.patch)
    mag := (major << 40) | (minor << 20) | patch
    return int64(mag)
}

func sortedVersionList(versions []string) (versionList, error) {
    list := make(versionList, 0, len(versions))
    for _, version := range versions {
        val, err := strtover(version)
        if err != nil {
            return nil, err
        }
        list = append(list, val)
    }
    sort.Sort(list)
    return list, nil
}

func strtover(value string) (version, error) {
    var ret version
    versions := [3]int{0, 0, 0}
    values := strings.Split(value, ".")
    if len(values) > 3 {
        return ret, errors.Stringf("too many version delimiters: %s", value)
    }
    for i := 0; i < len(values); i++ {
        version, err := strconv.Atoi(values[i])
        if err != nil || version < 0 {
            return ret, errors.Stringf("invalid version string: %s", value)
        }
        versions[i] = version
    }
    return version{
        major: versions[0],
        minor: versions[1],
        patch: versions[2],
    }, nil
}

func vertostr(value version) string {
    return fmt.Sprintf("%d.%d.%d", value.major, value.minor, value.patch)
}
