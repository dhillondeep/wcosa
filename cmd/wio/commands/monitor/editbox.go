package monitor

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"unicode/utf8"
)

type EditBox struct {
	text          []byte
	lineVOffset   int
	cursorBOffset int // cursor offset in bytes
	cursorVOffset int // visual cursor offset in termbox cells
	cursorCOffset int // cursor offset in unicode code points
}

// Draws the EditBox in the given location, 'h' is not used at the moment
func (eb *EditBox) Draw(x, y, w, h int) {
	eb.AdjustVOffset(w)

	const coldef = termbox.ColorDefault
	fill(x, y, w, h, termbox.Cell{Ch: ' '})

	t := eb.text
	lx := 0
	tabstop := 0
	for {
		rx := lx - eb.lineVOffset
		if len(t) == 0 {
			break
		}

		if lx == tabstop {
			tabstop += tabstop_length
		}

		if rx >= w {
			termbox.SetCell(x+w-1, y, '→',
				coldef, coldef)
			break
		}

		r, size := utf8.DecodeRune(t)
		if r == '\t' {
			for ; lx < tabstop; lx++ {
				rx = lx - eb.lineVOffset
				if rx >= w {
					goto next
				}

				if rx >= 0 {
					termbox.SetCell(x+rx, y, ' ', coldef, coldef)
				}
			}
		} else {
			if rx >= 0 {
				termbox.SetCell(x+rx, y, r, coldef, coldef)
			}
			lx += runewidth.RuneWidth(r)
		}
	next:
		t = t[size:]
	}

	if eb.lineVOffset != 0 {
		termbox.SetCell(x, y, '←', coldef, coldef)
	}
}

// Adjusts line visual offset to a proper value depending on width
func (eb *EditBox) AdjustVOffset(width int) {
	ht := preferred_horizontal_threshold
	max_h_threshold := (width - 1) / 2
	if ht > max_h_threshold {
		ht = max_h_threshold
	}

	threshold := width - 1
	if eb.lineVOffset != 0 {
		threshold = width - ht
	}
	if eb.cursorVOffset-eb.lineVOffset >= threshold {
		eb.lineVOffset = eb.cursorVOffset + (ht - width + 1)
	}

	if eb.lineVOffset != 0 && eb.cursorVOffset-eb.lineVOffset < ht {
		eb.lineVOffset = eb.cursorVOffset - ht
		if eb.lineVOffset < 0 {
			eb.lineVOffset = 0
		}
	}
}

func (eb *EditBox) MoveCursorTo(boffset int) {
	eb.cursorBOffset = boffset
	eb.cursorVOffset, eb.cursorCOffset = voffset_coffset(eb.text, boffset)
}

func (eb *EditBox) RuneUnderCursor() (rune, int) {
	return utf8.DecodeRune(eb.text[eb.cursorBOffset:])
}

func (eb *EditBox) RuneBeforeCursor() (rune, int) {
	return utf8.DecodeLastRune(eb.text[:eb.cursorBOffset])
}

func (eb *EditBox) MoveCursorOneRuneBackward() {
	if eb.cursorBOffset == 0 {
		return
	}
	_, size := eb.RuneBeforeCursor()
	eb.MoveCursorTo(eb.cursorBOffset - size)
}

func (eb *EditBox) MoveCursorOneRuneForward() {
	if eb.cursorBOffset == len(eb.text) {
		return
	}
	_, size := eb.RuneUnderCursor()
	eb.MoveCursorTo(eb.cursorBOffset + size)
}

func (eb *EditBox) MoveCursorToBeginningOfTheLine() {
	eb.MoveCursorTo(0)
}

func (eb *EditBox) MoveCursorToEndOfTheLine() {
	eb.MoveCursorTo(len(eb.text))
}

func (eb *EditBox) DeleteRuneBackward() {
	if eb.cursorBOffset == 0 {
		return
	}

	eb.MoveCursorOneRuneBackward()
	_, size := eb.RuneUnderCursor()
	eb.text = byte_slice_remove(eb.text, eb.cursorBOffset, eb.cursorBOffset+size)
}

func (eb *EditBox) DeleteRuneForward() {
	if eb.cursorBOffset == len(eb.text) {
		return
	}
	_, size := eb.RuneUnderCursor()
	eb.text = byte_slice_remove(eb.text, eb.cursorBOffset, eb.cursorBOffset+size)
}

func (eb *EditBox) DeleteTheRestOfTheLine() {
	eb.text = eb.text[:eb.cursorBOffset]
}

func (eb *EditBox) InsertRune(r rune) {
	var buf [utf8.UTFMax]byte
	n := utf8.EncodeRune(buf[:], r)
	eb.text = byte_slice_insert(eb.text, eb.cursorBOffset, buf[:n])
	eb.MoveCursorOneRuneForward()
}

func (eb *EditBox) Reset() {
	eb.MoveCursorToBeginningOfTheLine()
	eb.text = []byte{}
}

// Please, keep in mind that cursor depends on the value of line_voffset, which
// is being set on Draw() call, so.. call this method after Draw() one.
func (eb *EditBox) CursorX() int {
	return eb.cursorVOffset - eb.lineVOffset
}
