package resolve

/*
import (
    "testing"
    "wio/cmd/wio/log"
    "wio/cmd/wio/toolchain/npm/semver"
)

func TestTest(t *testing.T) {
    root := &Node{name: "package", ver: "1.0.0", resolve: semver.Parse("1.0.0")}
    root.deps = append(root.deps, &Node{name: "wlib-json", ver: "1.0.4"})
    root.deps = append(root.deps, &Node{name: "wlib-memory", ver: "1.0.0"})
    root.deps = append(root.deps, &Node{name: "react", ver: "16"})
    root.deps = append(root.deps, &Node{name: "ember", ver: "1"})
    info := NewInfo("/home/jeff/Code/gopath/src/wio/tests/project-app/app-stdout/")
    for _, dep := range root.deps {
        if err := info.ResolveTree(dep); err != nil {
            t.Fatal(err.Error())
        }
    }
    log.Infoln()
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
*/
