package semver

import (
    "regexp"
    "strings"
    "fmt"
)

type queryOp int

// An arithmetic operation requires a hard version query.
// Placeholders are expanded to zero in these cases
const (
    queryEq queryOp = 0
    queryLt queryOp = 1
    queryGt queryOp = 2
    queryLe queryOp = 3
    queryGe queryOp = 4
)

func (op queryOp) compare(a *Version, b *Version) bool {
    switch op {
    case queryEq:
        return a.eq(b)
    case queryLt:
        return a.less(b)
    case queryGt:
        return b.less(a)
    case queryLe:
        return !b.less(a)
    case queryGe:
        return !a.less(b)
    default:
        return false
    }
}

type Query interface {
    Matches(ver *Version) bool
    FindBest(list List) *Version
    Str() string
}

type singleBound struct {
    ver *Version
    op  queryOp
}

// lower.op always queryGe or queryGt
// upper.op always queryLe or queryLt
type dualBound struct {
    lower *singleBound
    upper *singleBound
}

func (q *singleBound) Matches(ver *Version) bool {
    return q.op.compare(ver, q.ver)
}

func (q *singleBound) FindBest(list List) *Version {
    // assumes list is sorted
    if len(list) <= 0 {
        return nil
    }

    switch q.op {
    case queryEq:
        for _, ver := range list {
            if q.ver.eq(ver) {
                return ver
            }
        }
        return nil

    case queryLt, queryLe:
        // find the highest matching version
        for i := len(list) - 1; i >= 0; i-- {
            if q.Matches(list[i]) {
                return list[i]
            }
        }
        return nil

    case queryGt, queryGe:
        // check if highest version matches
        last := list[len(list)-1]
        if q.Matches(last) {
            return last
        }
        return nil

    default:
        return nil
    }
}

func (q *singleBound) Str() string {
    return fmt.Sprintf("%s%s", queryInv[q.op], q.ver.Str())
}

func (q *dualBound) Matches(ver *Version) bool {
    return q.lower.Matches(ver) && q.upper.Matches(ver)
}

func (q *dualBound) FindBest(list List) *Version {
    // assumes list is sorted
    if len(list) <= 0 {
        return nil
    }
    // find highest allowable that meets lower bound
    top := q.upper.FindBest(list)
    if top != nil && q.lower.Matches(top) {
        return top
    }
    return nil
}

func (q *dualBound) Str() string {
    return fmt.Sprintf("%s %s", q.lower.Str(), q.upper.Str())
}

var cmpMatch = regexp.MustCompile(`^(>=|<=|>|<)v?([*xX]|[0-9]+)(\.([*xX]|[0-9]+)){0,2}(-[a-zA-Z0-9]+(\.[0-9]+)?)?$`)
var misMatch = regexp.MustCompile(`^=?v?(([*xX]|[0-9]+)(\.([*xX]|[0-9]+)){0,2})?(-[a-zA-Z0-9]+(\.[0-9]+)?)?$`)
var tildeMatch = regexp.MustCompile(`^~=?v?(([*xX]|[0-9]+)(\.([*xX]|[0-9]+)){0,2}(-[0-9a-zA-Z]+(\.[0-9]+)?)?)?$`)
var caretMatch = regexp.MustCompile(`^\^=?v?(([*xX]|[0-9]+)(\.([*xX]|[0-9]+)){0,2}(-[0-9a-zA-Z]+(\.[0-9]+)?)?)?$`)
var opMatch = regexp.MustCompile(`^(>=|<=|>|<)`)
var anyMatch = regexp.MustCompile(`[*xX]`)

var queryMap = map[string]queryOp{
    "=":  queryEq,
    "<":  queryLt,
    ">":  queryGt,
    "<=": queryLe,
    ">=": queryGe,
}

var queryInv = [...]string{"=", "<", ">", "<=", ">="}

func parseIncompl(str string) *Version {
    loc := anyMatch.FindStringIndex(str)
    if loc != nil {
        str = str[:loc[0]]
    }
    if str == "" {
        str = "0"
    } else if str[len(str)-1] == '.' {
        str = str[:len(str)-1]
    }
    str += strings.Repeat(".0", 2-strings.Count(str, "."))
    return Parse(str)
}

func parseCmpQuery(str string) Query {
    loc := opMatch.FindStringIndex(str)
    opStr := str[loc[0]:loc[1]]
    return &singleBound{op: queryMap[opStr], ver: parseIncompl(str[loc[1]:])}
}

func parseMisQuery(str string) Query {
    lower := parseIncompl(str)
    ver := strings.Split(str, ".")
    switch {
    case len(str) == 0 || anyMatch.MatchString(ver[0]):
        return &singleBound{op: queryGe, ver: lower}

    case len(ver) == 1 || anyMatch.MatchString(ver[1]):
        upper := &Version{lower.Major + 1, 0, 0}
        return &dualBound{
            lower: &singleBound{op: queryGe, ver: lower},
            upper: &singleBound{op: queryLt, ver: upper},
        }

    case len(ver) == 2 || anyMatch.MatchString(ver[2]):
        upper := &Version{lower.Major, lower.Minor + 1, 0}
        return &dualBound{
            lower: &singleBound{op: queryGe, ver: lower},
            upper: &singleBound{op: queryLt, ver: upper},
        }

    case len(ver) == 3:
        return &singleBound{op: queryEq, ver: lower}

    default:
        return nil
    }
}

func parseTildeQuery(str string) Query {
    switch {
    case IsValid(str):
        lower := Parse(str)
        upper := &Version{lower.Major, lower.Minor + 1, 0}
        return &dualBound{
            lower: &singleBound{op: queryGe, ver: lower},
            upper: &singleBound{op: queryLt, ver: upper},
        }

    default:
        return parseMisQuery(str)
    }
}

func parseCaretQuery(str string) Query {
    lower := parseIncompl(str)
    ver := strings.Split(str, ".")
    switch {
    case len(str) == 0 || anyMatch.MatchString(ver[0]):
        return parseMisQuery(str)

    case len(ver) == 1 || anyMatch.MatchString(ver[1]):
        upper := &Version{lower.Major + 1, 0, 0}
        return &dualBound{
            lower: &singleBound{op: queryGe, ver: lower},
            upper: &singleBound{op: queryLt, ver: upper},
        }

    case len(ver) == 2 || anyMatch.MatchString(ver[2]):
        upper := &Version{0, 0, 0}
        if lower.Major == 0 {
            upper.Minor = lower.Minor + 1
        } else {
            upper.Major = lower.Major + 1
        }
        return &dualBound{
            lower: &singleBound{op: queryGe, ver: lower},
            upper: &singleBound{op: queryLt, ver: upper},
        }

    case len(ver) == 3:
        upper := &Version{0, 0, 0}
        if lower.Major == 0 {
            if lower.Minor == 0 {
                upper.Patch = lower.Patch + 1
            } else {
                upper.Minor = lower.Minor + 1
            }
        } else {
            upper.Major = lower.Major + 1
        }
        return &dualBound{
            lower: &singleBound{op: queryGe, ver: lower},
            upper: &singleBound{op: queryLt, ver: upper},
        }

    default:
        return nil
    }
}

func stripEqV(str string) string {
    if len(str) > 0 && str[0] == '=' {
        str = str[1:]
    }
    if len(str) > 0 && str[0] == 'v' {
        str = str[1:]
    }
    return str
}

func MakeQuery(str string) Query {
    // should never be called if valid Version passed
    if IsValid(str) {
        return &singleBound{ver: Parse(str), op: queryEq}
    }
    switch {
    case cmpMatch.MatchString(str):
        return parseCmpQuery(str)

    case misMatch.MatchString(str):
        return parseMisQuery(stripEqV(str))

    case tildeMatch.MatchString(str):
        return parseTildeQuery(stripEqV(str[1:]))

    case caretMatch.MatchString(str):
        return parseCaretQuery(stripEqV(str[1:]))

    default:
        return nil
    }
}
