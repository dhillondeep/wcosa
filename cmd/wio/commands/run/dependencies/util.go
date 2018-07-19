package dependencies

import (
    "strconv"
    "wio/cmd/wio/toolchain/npm/semver"
)

// Converts semver version to a string "Major.Minor.Patch"
func SemverVersionToString(version *semver.Version) string {
    return strconv.Itoa(version.Major) + "." + strconv.Itoa(version.Minor) + "." + strconv.Itoa(version.Patch)
}
