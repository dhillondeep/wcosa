package defaults

import "wio/cmd/wio/constants"

var (
    Version       = "0.0.1"
    Port          = "none"
    AVRBoard      = "uno"
    Baud          = 9600
    AppKeywords   = []string{constants.Wio, constants.App}
    PkgKeywords   = []string{constants.Wio, constants.Pkg}
    AppTargetName = "main"
    PkgTargetName = "tests"
    AppTargetPath = "src"
    PkgTargetPath = "tests"
)
