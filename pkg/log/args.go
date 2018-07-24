package log

import (
    "io"

    "github.com/fatih/color"
)

type Args struct {
    queue   *Queue
    level   Type
    color   *color.Color
    message string
    newline bool
    args    []interface{}
    writer  io.Writer
}

func (a *Args) Append(v interface{}) {
    a.args = append(a.args, v)
}

func NewArgs(num int) *Args {
    return &Args{
        queue:   nil,
        level:   NONE,
        color:   Default,
        message: "",
        newline: false,
        args:    make([]interface{}, 0, num),
        writer:  nil,
    }
}

func GetArgs(args ...interface{}) *Args {
    ret := NewArgs(len(args))
    for _, arg := range args {
        switch val := arg.(type) {
        case *Queue:
            ret.queue = val
        case Type:
            ret.level = val
        case *color.Color:
            ret.color = val
        case string:
            if ret.message == "" {
                ret.message = val
            } else {
                ret.Append(val)
            }
        case bool:
            ret.newline = val
        case io.Writer:
            ret.writer = val
        case error:
            if ret.message == "" {
                ret.message = val.Error()
            } else {
                ret.Append(val.Error())
            }
        default:
            ret.Append(val)
        }
    }
    return ret
}
