package resolve

import (
    "strings"
    "wio/cmd/wio/log"
    "wio/cmd/wio/types"
)

var line = log.NewLine(log.INFO)

const barLength = 30

func logResolveStart(config types.Config) {
    log.Info(log.Cyan, "Resolving dependencies of: ")
    log.Infoln(log.Green, "%s@%s", config.GetName(), config.GetVersion())
}

func logResolve(n *Node) {
    line.Begin()
    line.Write(" ")
    line.Write("[", log.Cyan)
    line.Write("resolve", log.Magenta)
    line.Write("]", log.Cyan)
    line.Write(" ")
    line.Write("%s@%s", n.Name, n.ConfigVersion, log.Green)
    line.End()
}

func logResolveDone(root *Node) {
    line.Begin()
    line.End()
    printTree(root, "")
}

func printTree(node *Node, pre string) {
    log.Infoln(log.Green, "%s@%s", node.Name, node.ResolvedVersion.Str())
    for i := 0; i < len(node.Dependencies)-1; i++ {
        log.Info("%s|_ ", pre)
        printTree(node.Dependencies[i], pre+"|  ")
    }
    if len(node.Dependencies) > 0 {
        log.Info("%s\\_ ", pre)
        printTree(node.Dependencies[len(node.Dependencies)-1], pre+"   ")
    }
}

func logInstallStart() {
    log.Infoln(log.Cyan, "Installing dependencies")
}

func logInstall(name string, ver string, curr uint64, total uint64) {
    line.Begin()
    line.Write(" ")
    line.Write("[", log.Cyan)
    line.Write("install", log.Magenta)
    line.Write("]", log.Cyan)
    line.Write(" ")
    logProgress(curr, total)
    line.Write(" %s@%s", name, ver, log.Green)
    line.End()
}

func logProgress(curr uint64, total uint64) {
    prog := float64(curr) / float64(total)
    fill := int(prog * float64(barLength))
    line.Write("[", log.Cyan)
    line.Write(strings.Repeat("=", fill), log.Blue)
    line.Write(">", log.Blue)
    line.Write(strings.Repeat(" ", barLength-fill))
    line.Write("]", log.Cyan)
}

func logInstallDone() {
    line.Begin()
    line.End()
    log.Infoln(log.Cyan, "|> Done!")
}

type callback func(curr uint64, total uint64)

type counter struct {
    total uint64
    curr  uint64
    cb    callback
}

func (w *counter) Write(p []byte) (int, error) {
    n := len(p)
    w.curr += uint64(n)
    w.cb(w.curr, w.total)
    return n, nil
}
