package monitor

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"go.bug.st/serial.v1"
	"strconv"
	"time"
	"unicode/utf8"
)

var editBox EditBox
var screenBox = ScreenBox{Text: make(map[int][]string), PageNum: 0, StartX: 0, StartY: 3}

func tbPrint(x, y int, fg, bg termbox.Attribute, msg string) int {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}

	return x
}

// draws the top bar
func drawTopBar(port string, baud string, writeModeOn bool, pause bool) {
	const coldef = termbox.ColorDefault
	w, _ := termbox.Size()

	// print title
	x := tbPrint(0, 0, termbox.ColorCyan, coldef, "Wio Serial Monitor")
	x = tbPrint(x, 0, termbox.ColorYellow, coldef, "  @  ")
	x = tbPrint(x, 0, coldef, coldef, port)
	x = tbPrint(x, 0, termbox.ColorYellow, coldef, "  @  ")
	x = tbPrint(x, 0, coldef, coldef, baud)

	pageMode := "write"

	if !writeModeOn {
		pageMode = "use"
	}

	// print mode, pause, and page number
	modeAndPageStr := "    Mode: " + pageMode + " | Serial Pause: " + strconv.FormatBool(pause) + " | Page Num: " +
		strconv.Itoa(screenBox.PageNum+1)
	tbPrint(w-len(modeAndPageStr), 0, coldef, coldef, modeAndPageStr)

	// print instructions
	tbPrint(0, 1, coldef, coldef, "-- Quit: ESC | Mode: Ctrl + W | Write: Enter --")

	termbox.Flush()
}

// draws the send bar at the bottom
func drawSendBar() {
	const coldef = termbox.ColorDefault
	w, h := termbox.Size()

	x := tbPrint(0, h-1, termbox.ColorCyan, coldef, "Write: ")
	editBox.Draw(x, h-1, w-x, 1)
	termbox.SetCursor(x+editBox.CursorX(), h-1)

	termbox.Flush()
}

func drawModeInstructionsBar(writeModeOn bool) {
	const coldef = termbox.ColorDefault
	_, h := termbox.Size()

	if writeModeOn {
		termbox.SetCursor(7+editBox.CursorX(), h-1)

		tbPrint(0, 2, coldef, coldef, "-- Front: Home/Ctrl+A | End: End/Ctr+E | Delete Rest Line: Ctrl+K --")
	} else {
		tbPrint(0, 2, coldef, coldef, "-- Page Number: <- and -> | Pause: Ctrl+P --                        ")
		termbox.HideCursor()
	}

	termbox.Flush()
}

func drawSerialScreen() {
	const coldef = termbox.ColorDefault

	screenBox.Draw(coldef, coldef)
	termbox.Flush()
}

func clearSerialScreen() {
	terminalWidth, terminalHeight := termbox.Size()
	const coldef = termbox.ColorDefault

	for h := screenBox.StartY; h < terminalHeight-1; h++ {
		for w := screenBox.StartX; w < terminalWidth; w++ {
			tbPrint(w, h, coldef, coldef, " ")
		}
	}
}

func fill(x, y, w, h int, cell termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, cell.Ch, cell.Fg, cell.Bg)
		}
	}
}

func rune_advance_len(r rune, pos int) int {
	if r == '\t' {
		return tabstop_length - pos%tabstop_length
	}
	return runewidth.RuneWidth(r)
}

func voffset_coffset(text []byte, boffset int) (voffset, coffset int) {
	text = text[:boffset]
	for len(text) > 0 {
		r, size := utf8.DecodeRune(text)
		text = text[size:]
		coffset += 1
		voffset += rune_advance_len(r, voffset)
	}
	return
}

func byte_slice_grow(s []byte, desired_cap int) []byte {
	if cap(s) < desired_cap {
		ns := make([]byte, len(s), desired_cap)
		copy(ns, s)
		return ns
	}
	return s
}

func byte_slice_remove(text []byte, from, to int) []byte {
	size := to - from
	copy(text[from:], text[to:])
	text = text[:len(text)-size]
	return text
}

func byte_slice_insert(text []byte, offset int, what []byte) []byte {
	n := len(text) + len(what)
	text = byte_slice_grow(text, n)
	text = text[:n]
	copy(text[offset+len(what):], text[offset:])
	copy(text[offset:], what)
	return text
}

const preferred_horizontal_threshold = 5
const tabstop_length = 4

func main() {
	// Open the first serial port detected at 9600bps N81
	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	serialPort, err := serial.Open("/dev/cu.usbmodem14141", mode)
	if err != nil {
		panic(err)
	}

	// Read and print the response
	buff := make([]byte, 100)
	for {
		// Reads up to 100 bytes
		n, err := serialPort.Read(buff)
		if err != nil {
			panic(err)
			break
		}
		if n == 0 {
			fmt.Println("\nEOF")
			break
		}
		fmt.Printf("%v", string(buff[:n]))

	}

	return

	return

	writeModeOn := false
	pause := false
	port := "dev.usb"
	baud := "9600"

	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	drawTopBar(port, baud, writeModeOn, pause)
	drawModeInstructionsBar(writeModeOn)
	drawSendBar()

	quit := make(chan bool)

	// Read and print the response
	buff = make([]byte, 100)
	for {
		n, err := serialPort.Read(buff)
		if err != nil {
			panic(err)
			break
		}
		if n == 0 {
			screenBox.WriteString("\nEOF")
			break
		}

		tbPrint(screenBox.StartX, screenBox.StartY, termbox.ColorDefault, termbox.ColorDefault,
			fmt.Sprintf("%v", string(buff[:n])))
		time.Sleep(time.Millisecond * 1)

		/*
		   // Reads up to 100 bytes
		   n, err := serialPort.Read(buff)
		   if err != nil {
		       panic(err)
		       break
		   }
		   if n == 0 {
		       screenBox.WriteString("\nEOF")
		       break
		   }

		   if screenBox.WriteString(fmt.Sprintf("%v", string(buff[:n]))) {
		       drawTopBar(port, baud, writeModeOn, pause)
		       clearSerialScreen()
		       drawSerialScreen()
		   }
		*/
	}

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventResize:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			drawTopBar(port, baud, writeModeOn, pause)
			drawModeInstructionsBar(writeModeOn)
			drawSerialScreen()
			drawSendBar()
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				quit <- true
				break mainloop
			case termbox.KeyEnter:
				editBox.Reset()
			case termbox.KeyCtrlW:
				writeModeOn = !writeModeOn
				drawModeInstructionsBar(writeModeOn)
				drawTopBar(port, baud, writeModeOn, pause)
			case termbox.KeyCtrlP:
				pause = !pause
				drawTopBar(port, baud, writeModeOn, pause)
			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
				if writeModeOn {
					editBox.MoveCursorOneRuneBackward()
					drawSendBar()
				} else {
					screenBox.ChangePageNum(-1)
					drawTopBar(port, baud, writeModeOn, pause)
					drawSerialScreen()
					clearSerialScreen()
				}
			case termbox.KeyArrowRight, termbox.KeyCtrlF:
				if writeModeOn {
					editBox.MoveCursorOneRuneForward()
					drawSendBar()
				} else {
					screenBox.ChangePageNum(1)
					drawTopBar(port, baud, writeModeOn, pause)
					drawSerialScreen()
					clearSerialScreen()
				}
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				if writeModeOn {
					editBox.DeleteRuneBackward()
					drawSendBar()
				}
			case termbox.KeyDelete, termbox.KeyCtrlD:
				if writeModeOn {
					editBox.DeleteRuneForward()
					drawSendBar()
				}
			case termbox.KeyTab:
				if writeModeOn {
					editBox.InsertRune('\t')
					drawSendBar()
				}
			case termbox.KeySpace:
				if writeModeOn {
					editBox.InsertRune(' ')
					drawSendBar()
				}
			case termbox.KeyCtrlK:
				if writeModeOn {
					editBox.DeleteTheRestOfTheLine()
					drawSendBar()
				}
			case termbox.KeyHome, termbox.KeyCtrlA:
				if writeModeOn {
					editBox.MoveCursorToBeginningOfTheLine()
					drawSendBar()
				}
			case termbox.KeyEnd, termbox.KeyCtrlE:
				if writeModeOn {
					editBox.MoveCursorToEndOfTheLine()
					drawSendBar()
				}
			default:
				if writeModeOn && ev.Ch != 0 {
					editBox.InsertRune(ev.Ch)
					drawSendBar()
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
