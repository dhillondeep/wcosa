package defaults

import "wio/internal/constants"

type Defaults struct {
    Keywords []string
    Target   string
    Source   string
}

const (
    Version = "0.0.1"

    AVRBoard = "uno"
    Port     = "none"
    Baud     = 9600
)

var App = Defaults{
    Keywords: []string{constants.Wio, constants.App},
    Target:   "main",
    Source:   "src",
}

var Pkg = Defaults{
    Keywords: []string{constants.Wio, constants.Pkg},
    Target:   "tests",
    Source:   "tests",
}
