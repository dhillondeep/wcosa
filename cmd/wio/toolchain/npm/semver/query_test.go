package semver

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestParseIncompl(t *testing.T) {
    values := map[string]*Version{
        "":        {0, 0, 0},
        "0":       {0, 0, 0},
        "1":       {1, 0, 0},
        "2":       {2, 0, 0},
        "2.0":     {2, 0, 0},
        "2.5":     {2, 5, 0},
        "2.8":     {2, 8, 0},
        "6.7":     {6, 7, 0},
        "1.5.6":   {1, 5, 6},
        "6.6.7":   {6, 6, 7},
        "x":       {0, 0, 0},
        "*":       {0, 0, 0},
        "X":       {0, 0, 0},
        "5.x":     {5, 0, 0},
        "7.*":     {7, 0, 0},
        "7.x.X":   {7, 0, 0},
        "7.*.4":   {7, 0, 0},
        "10.10.x": {10, 10, 0},
    }
    for arg, exp := range values {
        assert.Equal(t, parseIncompl(arg), exp)
    }
}

func TestMakeQuery(t *testing.T) {
    singleEq := func(b *singleBound, q Query) bool {
        switch val := q.(type) {
        case *singleBound:
            return b.op == val.op && b.ver.eq(val.ver)
        default:
            return false
        }
    }
    dualEq := func(b *dualBound, q Query) bool {
        switch val := q.(type) {
        case *dualBound:
            return singleEq(b.lower, val.lower) && singleEq(b.upper, val.upper)
        default:
            return false
        }
    }
    eq := func(a Query, b Query) bool {
        switch val := a.(type) {
        case *singleBound:
            return singleEq(val, b)
        case *dualBound:
            return dualEq(val, b)
        default:
            return false
        }
    }
    single := func(op queryOp, major int, minor int, patch int) *singleBound {
        return &singleBound{op: op, ver: &Version{major, minor, patch}}
    }
    dual := func(lmajor int, lminor int, lpatch int, umajor int, uminor int, upatch int) *dualBound {
        return &dualBound{
            lower: single(queryGe, lmajor, lminor, lpatch),
            upper: single(queryLt, umajor, uminor, upatch),
        }
    }
    doTest := func(q Query, str string) {
        res := MakeQuery(str)
        if !assert.True(t, eq(q, res)) {
            t.Errorf("Equality test [%s] == [%s] failed!", q.Str(), res.Str())
        }
    }

    // Ranges
    doTest(single(queryGe, 567, 765, 890), ">=567.765.890")
    doTest(single(queryLe, 12, 34, 56), "<=12.34.56")
    doTest(single(queryGt, 6, 7, 8), ">6.7.8")
    doTest(single(queryLt, 12, 13, 14), "<12.13.14")

    doTest(single(queryGe, 50, 0, 0), ">=50")
    doTest(single(queryLe, 50, 0, 0), "<=50")
    doTest(single(queryGt, 60, 0, 0), ">60")
    doTest(single(queryLt, 60, 0, 0), "<60")

    doTest(single(queryGe, 50, 0, 0), ">=50.x")
    doTest(single(queryLe, 50, 0, 0), "<=50.x.55")
    doTest(single(queryGt, 50, 0, 0), ">50.x.x")
    doTest(single(queryLt, 50, 0, 0), "<50.0.x")

    doTest(single(queryGe, 5, 5, 0), ">=5.5")
    doTest(single(queryLe, 5, 5, 0), "<=5.5.x")
    doTest(single(queryGt, 5, 5, 0), ">5.5")
    doTest(single(queryLt, 5, 5, 0), "<5.5.x")

    // Missing
    doTest(single(queryGe, 0, 0, 0), "*")
    doTest(dual(1, 0, 0, 2, 0, 0), "1.x")
    doTest(dual(1, 2, 0, 1, 3, 0), "1.2.x")
    doTest(single(queryGe, 0, 0, 0), "")
    doTest(dual(1, 0, 0, 2, 0, 0), "1")
    doTest(dual(1, 2, 0, 1, 3, 0), "1.2")

    doTest(single(queryGe, 0, 0, 0), "x")
    doTest(single(queryGe, 0, 0, 0), "x.4.4")
    doTest(single(queryGe, 0, 0, 0), "*.x.X")
    doTest(single(queryGe, 0, 0, 0), "x.2.2")

    doTest(dual(16, 0, 0, 17, 0, 0), "16.x.x")
    doTest(dual(15, 15, 0, 15, 16, 0), "15.15.x")
    doTest(dual(15, 0, 0, 16, 0, 0), "15")
    doTest(dual(15, 15, 0, 15, 16, 0), "15.15")

    // Tilde
    doTest(dual(1, 2, 3, 1, 3, 0), "~1.2.3")
    doTest(dual(1, 2, 0, 1, 3, 0), "~1.2")
    doTest(dual(1, 2, 0, 1, 3, 0), "~1.2.x")
    doTest(dual(1, 0, 0, 2, 0, 0), "~1")
    doTest(dual(1, 0, 0, 2, 0, 0), "~1.x.x")
    doTest(dual(1, 0, 0, 2, 0, 0), "~1.x")
    doTest(dual(1, 0, 0, 2, 0, 0), "~1.x.6")
    doTest(dual(0, 2, 3, 0, 3, 0), "~0.2.3")
    doTest(dual(0, 2, 0, 0, 3, 0), "~0.2")
    doTest(dual(0, 0, 0, 1, 0, 0), "~0")

    doTest(dual(55, 55, 0, 55, 56, 0), "~55.55.x")
    doTest(dual(0, 0, 0, 0, 1, 0), "~0.0.x")
    doTest(dual(0, 0, 0, 1, 0, 0), "~0.x")

    // Caret
    doTest(dual(1, 2, 3, 2, 0, 0), "^1.2.3")
    doTest(dual(0, 2, 3, 0, 3, 0), "^0.2.3")
    doTest(dual(0, 0, 3, 0, 0, 4), "^0.0.3")
    doTest(dual(1, 2, 0, 2, 0, 0), "^1.2.x")
    doTest(dual(0, 0, 0, 0, 1, 0), "^0.0.x")
    doTest(dual(0, 0, 0, 0, 1, 0), "^0.0")
    doTest(dual(1, 0, 0, 2, 0, 0), "^1.x")
    doTest(dual(0, 0, 0, 1, 0, 0), "^0.x")
}
