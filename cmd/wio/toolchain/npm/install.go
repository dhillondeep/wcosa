package npm

import "wio/cmd/wio/errors"

type versionQuery int

const (
	equal   versionQuery = 0
	atLeast versionQuery = 1
	near    versionQuery = 2
)

func getVersionQuery(versionStr string) (versionQuery, error) {
	if len(versionStr) < 1 {
		return 0, errors.Stringf("invalid version string: %s", versionStr)
	}
	leading := versionStr[0]
	switch leading {
	case '~':
		return near, nil
	case '^':
		return atLeast, nil
	default:
		return equal, nil
	}
}

func findPackage(name string, versionStr string) (*packageVersion, error) {
	pkgData, err := getPackageData(name)
	if err != nil {
		return nil, err
	}
	if len(pkgData.Versions) <= 0 {
		return nil, errors.Stringf("package %s found but no versions exist", name)
	}
	// check for dist tag version
	if distTag, exists := pkgData.DistTags[versionStr]; exists {
		versionStr = distTag
	}
	queryType, err := getVersionQuery(versionStr)
	if err != nil {
		return nil, err
	}
	if queryType == equal {
		if pkgVersion, exists := pkgData.Versions[versionStr]; exists {
			return &pkgVersion, nil
		}
		return nil, errors.Stringf("package %s@%s does not exist", name, versionStr)
	}
	versionStrList := make([]string, 0, len(pkgData.Versions))
	for versionKey := range pkgData.Versions {
		versionStrList = append(versionStrList, versionKey)
	}
	sortedVersions, err := sortedVersionList(versionStrList)
	queryVersion, err := strtover(versionStr[1:])
	if err != nil {
		return nil, err
	}
	var version version
	switch queryType{
	case atLeast:
		version, err = sortedVersions.findAtLeast(queryVersion)
		if err != nil {
			return nil, err
		}
	case near:
		version = sortedVersions.findNearest(queryVersion)
	}
	pkgVersion := pkgData.Versions[vertostr(version)]
	return &pkgVersion, nil
}
