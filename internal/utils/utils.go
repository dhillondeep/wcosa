package utils

import "wio/pkg/util/sys"

func BuildPath(projectPath string) string {
	return sys.Path(projectPath, sys.WioFolder, sys.TargetDir)
}
