package dependencies

import (
    "hash/fnv"
    "strconv"
    "strings"
    "wio/cmd/wio/toolchain/npm/semver"
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
    Version               *semver.Version
    Path                  string
    Vendor                bool
    HeaderOnly            bool
    Flags                 []string
    FlagsVisibility       string
    Definitions           []string
    DefinitionsVisibility string
    hashValue             uint64
}

type targetSet struct {
    tMap  map[uint64]*Target
    links []*linkNode
}

// It creates 64-bit FNV-1a hash from name, version, flags, definitions, and headerOnly
func (target *Target) hash() uint64 {
    structStr := target.Name + SemverVersionToString(target.Version) + strings.Join(target.Flags, "") +
        strings.Join(target.Definitions, "") + target.FlagsVisibility + target.DefinitionsVisibility +
        strconv.FormatBool(target.HeaderOnly)

    h := fnv.New64a()
    h.Write([]byte(structStr))

    return h.Sum64()
}

// Creates a hash set for dependency targets
func NewTargetSet() *targetSet {
    return &targetSet{
        tMap: make(map[uint64]*Target),
    }
}

// Add Target values to TargetSet
func (targetSet *targetSet) Add(value *Target) {
    value.hashValue = value.hash()
    targetSet.tMap[value.hashValue] = value
}

// Links one target to another
func (targetSet *targetSet) Link(fromTarget *Target, toTarget *Target, linkInfo *TargetLinkInfo) {
    linkNode := &linkNode{
        From:     fromTarget,
        To:       toTarget,
        LinkInfo: linkInfo,
    }

    targetSet.links = append(targetSet.links, linkNode)
}

// Function used to iterate over targets
func (targetSet *targetSet) targetIterate(c chan<- *Target) {
    for _, b := range targetSet.tMap {
        c <- b
    }
    close(c)
}

// Function used to iterate over link nodes
func (targetSet *targetSet) linkIterate(c chan<- *linkNode) {
    for _, b := range targetSet.links {
        c <- b
    }
    close(c)
}

// Public iterator uses channels to return targetValues
func (targetSet *targetSet) TargetIterator() <-chan *Target {
    c := make(chan *Target)
    go targetSet.targetIterate(c)
    return c
}

// Public iterator uses channels to return Links between targets
func (targetSet *targetSet) LinkIterator() <-chan *linkNode {
    c := make(chan *linkNode)
    go targetSet.linkIterate(c)
    return c
}
