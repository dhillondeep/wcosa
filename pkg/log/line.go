package log

import (
    "fmt"
    "strings"
)

type buffer []byte

func (w *buffer) Write(p []byte) (n int, err error) {
    *w = append(*w, p...)
    return len(p), nil
}

type Line struct {
    max    int
    count  int
    level  Type
    writer *buffer
}

func NewLine(level Type) *Line {
    return &Line{
        max:    0,
        count:  0,
        level:  level,
        writer: &buffer{},
    }
}

func (l *Line) Begin() *Line {
    l.count = 0
    l.writer = &buffer{}
    Write(l.writer, l.level, "\r")
    return l
}

func (l *Line) End() {
    if l.count > l.max {
        l.max = l.count
    }
    pad := l.max - l.count
    Write(l.writer, l.level, strings.Repeat(" ", pad))
    Write(l.writer, l.level, "\r")
    logOut.Write([]byte(*l.writer))
}

func (l *Line) Write(args ...interface{}) *Line {
    a := GetArgs(args...)
    str := fmt.Sprintf(a.message, a.args...)
    l.count += len(str)
    Write(l.writer, l.level, a.color, str)
    return l
}
