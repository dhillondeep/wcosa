package npm

import (
	"testing"
)

func TestStrtover(t *testing.T) {
	str1 := "1.5.3"
	res1, err := strtover(str1)
	exp1 := version{1, 5, 3}
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if res1 != exp1 {
		t.Errorf("invalid res1: %s", vertostr(res1))
	}
	str2 := "4.g"
	res2, err := strtover(str2)
	exp2 := version{0, 0, 0}
	if err == nil {
		t.Errorf("missing expected error")
	}
	if res2 != exp2 {
		t.Errorf("invalid res2: %s", vertostr(res2))
	}
	str3 := "2.4"
	res3, err := strtover(str3)
	exp3 := version{2, 4, 0}
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if res3 != exp3 {
		t.Errorf("invalid res3: %s", vertostr(res3))
	}
}

func TestVertostr(t *testing.T) {
	val1 := version{5, 3, 4}
	exp1 := "5.3.4"
	if vertostr(val1) != exp1 {
		t.Errorf("result is not: %s", exp1)
	}
}

var versions = []string{
	"5.2.4",
	"3.2.6",
	"5.2.5",
	"5.3.2",
	"1.67.3",
	"9.23.1",
	"5.3.33",
}

func TestVersionListSort(t *testing.T) {
	sorted, err := sortedVersionList(versions)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	expected := []version{
		{1, 67, 3},
		{3, 2, 6},
		{5, 2, 4},
		{5, 2, 5},
		{5, 3, 2},
		{5, 3, 33},
		{9, 23, 1},
	}
	if len(sorted) != len(expected) {
		t.Errorf("sorted list unequal length")
	}
	for i := 0; i < len(expected); i++ {
		if sorted[i] != expected[i] {
			t.Errorf("element at %d expected to be %s but is %s", i,
				vertostr(expected[i]),
				vertostr(sorted[i]))
		}
	}
}

func TestFindAtLeast(t *testing.T) {
	sorted, _ := sortedVersionList(versions)
	query1 := version{5, 5, 5}
	res1, err := sorted.findAtLeast(query1)
	exp1 := version{9, 23, 1}
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if res1 != exp1 {
		t.Errorf("expected %s but got %s", vertostr(exp1), vertostr(res1))
	}
	query2 := version{10, 10, 10}
	res2, err := sorted.findAtLeast(query2)
	exp2 := version{0, 0, 0}
	if err == nil {
		t.Errorf("missing expected error")
	}
	if res2 != exp2 {
		t.Errorf("expected %s but got %s", vertostr(exp1), vertostr(res1))
	}
}

func TestFindNearest(t *testing.T) {
	sorted, _ := sortedVersionList(versions)
	query1 := version{6, 1, 1}
	res1 := sorted.findNearest(query1)
	exp1 := version{5, 3, 33}
	if exp1 != res1 {
		t.Errorf("expected %s but got %s", vertostr(exp1), vertostr(res1))
	}
	query2 := version{5, 3, 4}
	res2 := sorted.findNearest(query2)
	exp2 := version{5, 3, 2}
	if exp2 != res2 {
		t.Errorf("expected %s but got %s", vertostr(exp2), vertostr(res2))
	}
}
