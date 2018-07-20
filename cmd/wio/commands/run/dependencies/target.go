package dependencies

import (
    "strings"
)

type linkNode struct {
    From     *Target
    To       *Target
    LinkInfo *TargetLinkInfo
}

type TargetLinkInfo struct {
    Visibility string
    Flags      []string
}

type Target struct {
    Name                  string
    Version               string
    Path                  string
    FromVendor            bool
    HeaderOnly            bool
    Flags                 []string
    FlagsVisibility       string
    Definitions           []string
    DefinitionsVisibility string
    hashValue             string
}

type TargetSet struct {
    tMap  map[string]*Target
    links []*linkNode
}

// creates a hash from target struct
func (target *Target) hash() string {
    structStr := target.Name + target.Version + strings.Join(target.Flags, "") +
        strings.Join(target.Definitions, "") + target.FlagsVisibility + target.DefinitionsVisibility

    return structStr
}

// Creates a hash set for dependency targets
func NewTargetSet() *TargetSet {
    return &TargetSet{
        tMap: make(map[string]*Target),
    }
}

// Add Target values to TargetSet
func (targetSet *TargetSet) Add(value *Target) {
    value.hashValue = value.hash()
    targetSet.tMap[value.hashValue] = value
}

// Links one target to another
func (targetSet *TargetSet) Link(fromTarget *Target, toTarget *Target, linkInfo *TargetLinkInfo) {
    linkNode := &linkNode{
        From:     fromTarget,
        To:       toTarget,
        LinkInfo: linkInfo,
    }

    targetSet.links = append(targetSet.links, linkNode)
}

// Function used to iterate over targets
func (targetSet *TargetSet) targetIterate(c chan<- *Target) {
    for _, b := range targetSet.tMap {
        c <- b
    }
    close(c)
}

// Function used to iterate over link nodes
func (targetSet *TargetSet) linkIterate(c chan<- *linkNode) {
    for _, b := range targetSet.links {
        c <- b
    }
    close(c)
}

// Public iterator uses channels to return targetValues
func (targetSet *TargetSet) TargetIterator() <-chan *Target {
    c := make(chan *Target)
    go targetSet.targetIterate(c)
    return c
}

// Public iterator uses channels to return Links between targets
func (targetSet *TargetSet) LinkIterator() <-chan *linkNode {
    c := make(chan *linkNode)
    go targetSet.linkIterate(c)
    return c
}
