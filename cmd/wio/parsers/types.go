package parsers

// Type for Dependency tree used for parsing libraries
type DependencyTree struct {
    Config DependencyTag
    Child  map[string]*DependencyTree
}

// Structure to handle individual dependency inside dependencies
type DependencyTag struct {
    Name          string
    Hash          string
    Path          string
    Source        string
    CompileFlags []string
    HeaderOnly bool
}
