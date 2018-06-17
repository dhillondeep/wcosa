package monitor

import (
	"github.com/nsf/termbox-go"
	"strings"
)

const (
	maxBuffer = 10
)

var maxPageNum = 0

type ScreenBox struct {
	Text    map[int][]string
	PageNum int
	StartY  int
	StartX  int
}

func (screenBox *ScreenBox) Draw(fg, bg termbox.Attribute) {
	if screenBox.PageNum < 0 {
		screenBox.PageNum = 0
	} else if screenBox.PageNum%maxBuffer > len(screenBox.Text) {
		screenBox.PageNum = len(screenBox.Text)
	}

	tbPrint(screenBox.StartX, screenBox.StartY, fg, bg, strings.Join(screenBox.Text[screenBox.PageNum%maxBuffer], ""))
}

func (screenBox *ScreenBox) ChangePageNum(increment int) {
	pageNum := screenBox.PageNum + increment

	if pageNum < 0 {
		pageNum = 0
	} else if pageNum < (maxPageNum-maxBuffer)+1 {
		screenBox.PageNum = (maxPageNum - maxBuffer) + 1
	} else if pageNum > maxPageNum {
		screenBox.PageNum = maxPageNum
	} else {
		screenBox.PageNum = pageNum
	}
}

func (screenBox *ScreenBox) WriteString(str string) bool {
	pageNumChange := false

	_, h := termbox.Size()

	if len(screenBox.Text[screenBox.PageNum%maxBuffer]) >= h-(screenBox.StartY+1) {
		screenBox.PageNum++
		maxPageNum++
		pageNumChange = true
	}

	if _, exists := screenBox.Text[screenBox.PageNum%maxBuffer]; !exists {
		screenBox.Text[screenBox.PageNum%maxBuffer] = make([]string, 0)
	}

	// clear the buffer if we come around
	if pageNumChange && screenBox.PageNum/maxBuffer >= 1 {
		screenBox.Text[screenBox.PageNum%maxBuffer] = []string{}
	}

	screenBox.Text[screenBox.PageNum%maxBuffer] = append(screenBox.Text[screenBox.PageNum%maxBuffer], str)
	return pageNumChange
}
