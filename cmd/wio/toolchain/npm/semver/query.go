package semver

import (
    "regexp"
    "strings"
    "fmt"
    "bytes"
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

type queryList []Query

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

func (ql queryList) Matches(ver *Version) bool {
    for _, q := range ql {
        if q.Matches(ver) {
            return true
        }
    }
    return false
}

func (ql queryList) FindBest(list List) *Version {
    res := make(List, 0, len(ql))
    for _, q := range ql {
        if ver := q.FindBest(list); ver != nil {
            res = append(res, ver)
        }
    }
    if len(res) == 0 {
        return nil
    }
    res.Sort()
    return res[len(res)-1]
}

func (ql queryList) Str() string {
    var buf bytes.Buffer
    for _, q := range ql {
        buf.WriteString(q.Str())
        buf.WriteString(" || ")
    }
    res := buf.String()
    return res[:len(res)-4]
}

func trimX(str string) string {
    loc := anyMatch.FindStringIndex(str)
    if loc != nil {
        str = str[:loc[0]]
    }
    if len(str) > 0 && str[len(str)-1] == '.' {
        str = str[:len(str)-1]
    }
    return str
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

var cmpMatch = regexp.MustCompile(`^(>=|<=|>|<)v?([*xX]|[0-9]+)(\.([*xX]|[0-9]+)){0,2}(-[a-zA-Z0-9]+(\.[0-9]+)?)?$`)
var misMatch = regexp.MustCompile(`^=?v?(([*xX]|[0-9]+)(\.([*xX]|[0-9]+)){0,2})?(-[a-zA-Z0-9]+(\.[0-9]+)?)?$`)
var tildeMatch = regexp.MustCompile(`^~=?v?(([*xX]|[0-9]+)(\.([*xX]|[0-9]+)){0,2}(-[0-9a-zA-Z]+(\.[0-9]+)?)?)?$`)
var caretMatch = regexp.MustCompile(`^\^=?v?(([*xX]|[0-9]+)(\.([*xX]|[0-9]+)){0,2}(-[0-9a-zA-Z]+(\.[0-9]+)?)?)?$`)
var rangeMatch = regexp.MustCompile(`^[0-9]+(\.([*xX]|[0-9]+)){0,2}\s+-\s+[0-9]+(\.([*xX]|[0-9]+)){0,2}$`)
var andMatch = regexp.MustCompile(`^(>=|>)[0-9]+(\.([*xX]|[0-9]+)){0,2}\s+(<=|<)[0-9]+(\.([*xX]|[0-9]+)){0,2}$`)
var opMatch = regexp.MustCompile(`^(>=|<=|>|<)`)
var anyMatch = regexp.MustCompile(`[*xX]`)
var betMatch = regexp.MustCompile(`\s+-\s+`)
var spaceMatch = regexp.MustCompile(`\s+`)
var orMatch = regexp.MustCompile(`\s+\|\|\s+`)

var queryMap = map[string]queryOp{
    "=":  queryEq,
    "<":  queryLt,
    ">":  queryGt,
    "<=": queryLe,
    ">=": queryGe,
}

var queryInv = [...]string{"=", "<", ">", "<=", ">="}

func parseIncompl(str string) *Version {
    str = trimX(str)
    if str == "" {
        str = "0"
    }
    str += strings.Repeat(".0", 2-strings.Count(str, "."))
    return Parse(str)
}

func parseCmpQuery(str string) *singleBound {
    loc := opMatch.FindStringIndex(str)
    opStr := str[loc[0]:loc[1]]
    return &singleBound{op: queryMap[opStr], ver: parseIncompl(str[loc[1]:])}
}

func parseMisQuery(str string) Query {
    lower := parseIncompl(str)
    str = trimX(str)
    ver := strings.Split(str, ".")
    switch {
    case len(str) == 0:
        return &singleBound{op: queryGe, ver: lower}

    case len(ver) == 1:
        upper := &Version{lower.Major + 1, 0, 0}
        return &dualBound{
            lower: &singleBound{op: queryGe, ver: lower},
            upper: &singleBound{op: queryLt, ver: upper},
        }

    case len(ver) == 2:
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
    str = trimX(str)
    ver := strings.Split(str, ".")
    switch {
    case len(str) == 0:
        return parseMisQuery(str)

    case len(ver) == 1:
        upper := &Version{lower.Major + 1, 0, 0}
        return &dualBound{
            lower: &singleBound{op: queryGe, ver: lower},
            upper: &singleBound{op: queryLt, ver: upper},
        }

    case len(ver) == 2:
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

func parseRangeQuery(str string) *dualBound {
    bounds := betMatch.Split(str, -1)
    lower := parseIncompl(bounds[0])
    upper := bounds[1]
    if IsValid(upper) {
        return &dualBound{
            lower: &singleBound{op: queryGe, ver: lower},
            upper: &singleBound{op: queryLe, ver: Parse(upper)},
        }
    }
    upper = trimX(upper)
    ver := parseIncompl(upper)
    switch strings.Count(upper, ".") {
    case 0:
        ver.Major++
    case 1:
        ver.Minor++
    }
    return &dualBound{
        lower: &singleBound{op: queryGe, ver: lower},
        upper: &singleBound{op: queryLt, ver: ver},
    }
}

func parseAndQuery(str string) *dualBound {
    bounds := spaceMatch.Split(str, -1)
    return &dualBound{
        lower: parseCmpQuery(bounds[0]),
        upper: parseCmpQuery(bounds[1]),
    }
}

func parseOrQuery(str string) queryList {
    queries := orMatch.Split(str, -1)
    ql := make(queryList, 0, len(queries))
    for _, query := range queries {
        if q := MakeQuery(query); q != nil {
            ql = append(ql, q)
        } else {
            return nil
        }
    }
    return ql
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
    case rangeMatch.MatchString(str):
        return parseRangeQuery(str)
    case andMatch.MatchString(str):
        return parseAndQuery(str)
    case orMatch.MatchString(str):
        return parseOrQuery(str)
    default:
        return nil
    }
}
