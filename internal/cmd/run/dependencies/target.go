package dependencies

import (
	"strconv"
	"strings"
	"wio/internal/types"
	"wio/pkg/util"
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
	Name        string
	Version     string
	Path        string
	FromVendor  bool
	HeaderOnly  bool
	Flags       []string
	Definitions map[string][]string
	CXXStandard string
	CStandard   string
	hashValue   string
	ParentPath  string

	Library types.Library
}

type TargetSet struct {
	tMap    map[string]*Target
	nameMap map[string]*int
	links   []*linkNode
}

// creates a hash from target struct
func (target *Target) hash() string {
	structStr := target.Name + target.Version + strings.Join(target.Flags, "") +
		strings.Join(util.AppendIfMissing(target.Definitions[types.Private],
			target.Definitions[types.Public]), "")
	return structStr
}

// Creates a hash set for dependency targets
func NewTargetSet() *TargetSet {
	return &TargetSet{
		tMap:    make(map[string]*Target),
		nameMap: make(map[string]*int),
	}
}

// Add Target values to TargetSet
func (targetSet *TargetSet) Add(value *Target, isLibrary bool) {
	value.hashValue = value.hash()

	if dupTarget, exists := targetSet.tMap[value.hashValue]; !exists {
		namePostfix := 0

		// check if name exists
		if val, exists := targetSet.nameMap[value.Name]; exists {
			*val += 1
			namePostfix = *val
			value.Name += "__" + strconv.Itoa(namePostfix)
		} else {
			namePostfix = 0
			targetSet.nameMap[value.Name] = &namePostfix
		}

		if isLibrary {
			value.Name = "LIB__" + value.Name
		}

		targetSet.tMap[value.hashValue] = value
	} else {
		// use name of previous target with same config
		value.Name = dupTarget.Name
	}
}

func GetOriginalName(target *Target, isLibrary bool) string {
	if isLibrary {
		newStr := strings.Split(target.Name, "__")[1]
		return strings.Replace(newStr, "LIB__"+newStr, "", 1)
	} else {
		return strings.Split(target.Name, "__")[0]
	}
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
