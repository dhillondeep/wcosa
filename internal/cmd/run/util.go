package run

import (
    "os"
    "strconv"
    "wio/internal/cmd/run/cmake"
    "wio/internal/constants"
    "wio/internal/types"
    "wio/pkg/util/sys"
)

func buildPath(info *runInfo) string {
    return cmake.BuildPath(info.directory)
}

func targetPath(info *runInfo, target types.Target) string {
    return sys.Path(buildPath(info), target.GetName())
}

func binaryPath(info *runInfo, target types.Target) string {
    return sys.Path(targetPath(info, target), constants.BinDir)
}

func nativeExtension() string {
    switch sys.GetOS() {
    case sys.WINDOWS:
        return ".exe"
    default:
        return ""
    }
}

func platformExtension(platform string) string {
    switch platform {
    case constants.Avr:
        return ".elf"
    case constants.Native:
        return nativeExtension()
    default:
        return ""
    }
}

func shouldCreateBuildFiles(projectDir string, targetName string) (bool, error) {
    wioTimeFile := sys.Path(cmake.BuildPath(projectDir), targetName, "wio.time")

    // check if wio.time file exists
    if sys.Exists(wioTimeFile) {
        data, err := sys.NormalIO.ReadFile(wioTimeFile)
        if err != nil {
            return true, err
        }
        lastModTime, err := strconv.Atoi(string(data))
        if err != nil {
            return true, err
        }

        info, err := os.Stat(sys.Path(projectDir, sys.Config))
        if err != nil {
            return true, err
        }

        currModTime := info.ModTime().Nanosecond()

        if err := sys.NormalIO.WriteFile(wioTimeFile, []byte(strconv.Itoa(currModTime))); err != nil {
            return true, err
        }

        if currModTime == lastModTime {
            return false, nil
        }
    } else {
        info, err := os.Stat(sys.Path(projectDir, "wio.yml"))
        if err != nil {
            return true, err
        }

        if err := os.MkdirAll(sys.Path(wioTimeFile, "../"), os.ModePerm); err != nil {
            return true, err
        }

        currModTime := info.ModTime().Nanosecond()

        if err := sys.NormalIO.WriteFile(wioTimeFile, []byte(strconv.Itoa(currModTime))); err != nil {
            return true, err
        }
    }

    return true, nil
}
