package npm

/*func printTree(node *depTreeNode, level log.Type, pre string) {
    log.Writeln(level, "%s@%s", node.name, node.version)
    for i := 0; i < len(node.children)-1; i++ {
        log.Write(level, "%s|_ ", pre)
        printTree(node.children[i], level, pre+"|  ")
    }
    if len(node.children) > 0 {
        log.Write(level, "%s\\_ ", pre)
        printTree(node.children[len(node.children)-1], level, pre+"   ")
    }
}*/

/*
pkg
|_ wlib-json@1.0.4
|  \_ wlib-wio@1.0.0
\_ wlib-memory@1.0.2
|  |_ wlib-tmp@1.0.0
|  |  \_ wlib-util@1.0.0
|  \_ wlib-malloc@1.0.2
|     \_ wlib-tlsf@1.0.1
\_ wlib-list@1.0.0
*/
