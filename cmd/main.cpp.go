package main

import (
    "wio/cmd/wio/commands/run/dependencies"
    "fmt"
)


func main() {
    t1 := &dependencies.Target{
        Name: "hello",
        Version: "0.0.1",
        Flags: []string{"-DHello"},
        Definitions: nil,
    }

    t2 := &dependencies.Target{
        Name: "hello",
        Version: "0.0.2",
        Flags: []string{"-DHello"},
        Definitions: nil,
    }

    set := dependencies.NewTargetSet()
    set.Add(t1)
    set.Add(t2)

    for value := range set.Iterator() {
        fmt.Println(value)
    }
}





















