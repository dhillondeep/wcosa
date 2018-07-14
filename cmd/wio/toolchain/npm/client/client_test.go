package client

import "testing"

func TestFindFirstSlash(t *testing.T) {
    val1 := findFirstSlash("")
    if val1 != 0 {
        t.Errorf("Expected val1 == 0")
    }
    val2 := findFirstSlash("/////a////")
    if val2 != 5 {
        t.Errorf("Expected val2 == 5")
    }
    val3 := findFirstSlash("////////")
    if val3 != 8 {
        t.Errorf("Expected val3 == 8")
    }
    val4 := findFirstSlash("hello")
    if val4 != 0 {
        t.Errorf("Expected val4 == 0")
    }
}

func TestFindLastSlash(t *testing.T) {
    val1 := findLastSlash("")
    if val1 != -1 {
        t.Errorf("Expected val1 == -1")
    }
    val2 := findLastSlash("////a////")
    if val2 != 4 {
        t.Errorf("Expected val2 == 4")
    }
    val3 := findLastSlash("/////////")
    if val3 != -1 {
        t.Errorf("Expected val3 == -1")
    }
    val4 := findLastSlash("hello")
    if val4 != 4 {
        t.Errorf("Expected val4 == 4")
    }
}

func TestUrlResolve(t *testing.T) {
    res1 := urlResolve("https://github.com/", "/wio/", "/toolchain.git")
    exp1 := "https://github.com/wio/toolchain.git"
    if res1 != exp1 {
        t.Errorf("TestResolveUrl() -- failed!")
    }
}
