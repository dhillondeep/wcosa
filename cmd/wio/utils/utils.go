package utils

import (
    "os"
    "path/filepath"

    "strings"
)

// It takes in a slice and an element and then ut appends that element to the slice only
// if that element in not already in the slice
func AppendIfMissingElem(slice []string, i string) []string {
    for _, ele := range slice {
        if ele == i {
            return slice
        }
    }
    return append(slice, i)
}

// It takes two slices and appends the second one onto the first one. It does
// not allow duplicates
func AppendIfMissing(slice []string, slice2 []string) []string {
    newSlice := make([]string, 0)

    for _, ele1 := range slice {
        newSlice = AppendIfMissingElem(newSlice, ele1)
    }

    for _, ele2 := range slice2 {
        newSlice = AppendIfMissingElem(newSlice, ele2)
    }

    return newSlice
}

//  Eeturns elements in a that aren't in b
func Difference(a, b []string) []string {
    mb := map[string]bool{}
    for _, x := range b {
        mb[x] = true
    }
    var ab []string
    for _, x := range a {
        if _, ok := mb[x]; !ok {
            ab = append(ab, x)
        }
    }
    return ab
}

func Contains(slice []string, value string) bool {
    for _, element := range slice {
        if element == value {
            return true
        }
    }
    return false
}

func ContainsNoCase(slice []string, value string) bool {
    for _, element := range slice {
        if strings.ToLower(element) == strings.ToLower(value) {
            return true
        }
    }
    return false
}

// Deletes all the files from the directory
func RemoveContents(dir string) error {
    d, err := os.Open(dir)
    if err != nil {
        return err
    }
    defer d.Close()
    names, err := d.Readdirnames(-1)
    if err != nil {
        return err
    }
    for _, name := range names {
        err = os.RemoveAll(filepath.Join(dir, name))
        if err != nil {
            return err
        }
    }
    return nil
}
