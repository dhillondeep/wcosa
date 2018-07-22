package config

import "wio/cmd/wio/constants"

type defaults struct {
    Version       string
    Port          string
    AVRBoard      string
    Baud          int
    AppKeywords   []string
    PkgKeywords   []string
    AppTargetName string
    PkgTargetName string
    AppTargetPath string
    PkgTargetPath string
}

var ProjectDefaults = defaults{
    Version:       "0.0.1",
    Port:          "none",
    AVRBoard:      "uno",
    Baud:          9600,
    AppKeywords:   []string{constants.WIO, constants.APP},
    PkgKeywords:   []string{constants.WIO, constants.PKG},
    AppTargetName: "main",
    PkgTargetName: "tests",
    AppTargetPath: "src",
    PkgTargetPath: "tests",
}
