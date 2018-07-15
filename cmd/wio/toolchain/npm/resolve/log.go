package resolve

import (
    "wio/cmd/wio/log"
    "wio/cmd/wio/types"
)

var line = log.NewLine(log.INFO)

func logResolveStart(config types.IConfig) {
    log.Info(log.Cyan, "Resolving dependencies of: ")
    log.Infoln(log.Green, "%s@%s", config.Name(), config.Version())
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

func logResolveDone(root *Node) {
    line.Begin()
    line.End()
    printTree(root, "")
}

func printTree(node *Node, pre string) {
    log.Infoln(log.Green, "%s@%s", node.name, node.resolve.Str())
    for i := 0; i < len(node.deps)-1; i++ {
        log.Info("%s|_ ", pre)
        printTree(node.deps[i], pre+"|  ")
    }
    if len(node.deps) > 0 {
        log.Info("%s\\_ ", pre)
        printTree(node.deps[len(node.deps)-1], pre+"   ")
    }
}
