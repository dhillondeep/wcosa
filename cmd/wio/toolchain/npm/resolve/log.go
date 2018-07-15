package resolve

import (
    "wio/cmd/wio/constants"
    "wio/cmd/wio/log"
    "wio/cmd/wio/types"
)

var line = log.NewLine(log.INFO)

func logResolveStart(config types.IConfig) {
    log.Info(log.Cyan, "Resolving dependencies of: ")
    log.Info(log.Green, "%s", config.Name())
    if config.GetType() == constants.PKG {
        log.Info(log.Green, "@%s", config.Version())
    }
    log.Infoln()
}

func logResolve(n *Node) {
    line.Begin()
    line.Write(" ")
    line.Write("[", log.Cyan)
    line.Write("resolve", log.Magenta)
    line.Write("]", log.Cyan)
    line.Write(" ")
    line.Write("%s@%s", n.name, n.ver, log.Green)
    line.End()
}

func logResolveDone() {
    line.Begin()
    line.Write("--> Done!", log.Green)
    line.End()
    log.Infoln()
}
