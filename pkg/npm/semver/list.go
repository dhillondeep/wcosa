package semver

import (
    "sort"

    "github.com/blang/semver"
)

type List []*semver.Version

func (list List) Find(q Query) *semver.Version {
    // assumes list is sorted
    for i := len(list) - 1; i >= 0; i-- {
        if q.Matches(list[i]) {
            return list[i]
        }
    }
    return nil
}

func (list List) Len() int {
    return len(list)
}

func (list List) Swap(i int, j int) {
    list[i], list[j] = list[j], list[i]
}

func (list List) Less(i int, j int) bool {
    return list[i].LT(*list[j])
}

func (list List) Sort() {
    sort.Sort(list)
}

func (list List) Insert(v *semver.Version) List {
    for _, el := range list {
        if el.EQ(*v) {
            return list
        }
    }
    list = append(list, v)
    list.Sort()
    return list
}

func (list List) Last() *semver.Version {
    if len(list) == 0 {
        return nil
    }
    return list[len(list)-1]
}
