package rich

import (
	"image"
	"time"

	"github.com/gizak/termui"
)

const BlinkRate = time.Millisecond * 600

type lineInfo struct {
	pos, length int
}

type Widget struct {
	termui.Block
	textBuffer       TextBuffer
	WriteFg, WriteBg termui.Attribute

	MultiLine bool

	wrap bool

	cursorPos        int
	cursorEnabled    bool
	cursorBlinkState bool

	scrollX, scrollY int

	onChangeHandlers []func()
}

func New() *Widget {
	widget := Widget{
		Block:            *termui.NewBlock(),
		textBuffer:       *NewTextBuffer(),
		MultiLine:        true,
		onChangeHandlers: []func(){},
	}
	widget.CursorShow()
	//Log("Created widget")
	//Log("%v", widget.TextBuffer().Area)

	termui.Handle("/timer/1s", func(e termui.Event) {
		widget.cursorBlink()
	})

	return &widget
}

func (w *Widget) Resize(width int, height int) {
	w.Height = height
	w.Width = width
	Log("Min", w.InnerBounds().Min)
	Log("Max", w.InnerBounds().Max)
	w.textBuffer.SetArea(w.InnerBounds())
	Log("Min", w.textBuffer.Area.Min)
	Log("Max", w.textBuffer.Area.Max)
}

func (w *Widget) TextBuffer() TextBuffer {
	return w.textBuffer
}

func (w *Widget) onChange(h func()) {
	w.onChangeHandlers = append(w.onChangeHandlers, h)
}

func (w *Widget) dirty() {
	termui.Render(w)
	for _, h := range w.onChangeHandlers {
		h()
	}
}

func (w *Widget) shouldShowCursor(p image.Point) bool {
	// Log("===")
	// Log("y: %v, %v", w.CursorPosY(), p.Y)
	// Log("x: %v, %v", w.CursorPosX(), p.X)
	// Log("===")

	return w.cursorEnabled &&
		w.cursorIsVisible() &&
		w.cursorBlinkState &&
		w.CursorPosX() == p.X &&
		w.CursorPosY() == p.Y
}

func (w *Widget) cursorIsVisible() bool {
	return w.CursorPosY() < w.Height
}

func (w *Widget) isCellVisible(p image.Point) bool {
	return p.Y < w.Height
}

func (w *Widget) Buffer() termui.Buffer {
	textBuffer := w.TextBuffer()
	buffer := w.Block.Buffer()

	for p, c := range buffer.CellMap {
		if !w.isCellVisible(p) {
			continue
		}

		textBufferCell, ok := textBuffer.CellMap[p]
		if ok {
			c.Ch = textBufferCell.Ch
		}

		if w.shouldShowCursor(p) {
			c.Fg ^= termui.AttrReverse
			c.Bg ^= termui.AttrReverse
			if c.Ch == 0 {
				c.Ch = ' '
			}
		}

		buffer.Set(p.X, p.Y, c)
	}
	return buffer
}

func (w *Widget) IsOverflow() bool {
	buffer := w.TextBuffer()
	inner := buffer.Bounds()
	innerX := inner.Min.X
	innerY := inner.Min.Y
	innerWidth := w.InnerWidth()
	innerHeight := w.InnerHeight()
	col := innerX
	line := innerY
	overflow := false

	for _, c := range buffer.CellMap {
		if c.Ch == '\n' {
			c.Ch = 0
			line++
			if line >= innerY+innerHeight {
				overflow = true
				break
			}
			col = innerX
		}
		if c.Ch != 0 {
			col++
			if col >= innerX+innerWidth {
				col = innerX
				line++
				if line >= innerY+innerHeight {
					overflow = true
					break
				}
			}
		}
	}

	return overflow
}

func (w *Widget) WriteChar(c rune) {
	byteArray := []byte(string(c))
	w.Write(byteArray)
}

func (w *Widget) Write(p []byte) (n int, err error) {
	Log("Write %v", p)
	buffer := w.textBuffer
	l := 0
	for _, ch := range string(p) {
		if ch == '\n' && !w.MultiLine {
			continue
		}
		newCell := termui.NewCell(ch, w.WriteFg, w.WriteBg)
		//Log("char: %c, x: %v, y: %v, cell: %v", ch, w.CursorPosX(), w.CursorPosY(), newCell)
		//Log("buffer: %v", buffer)
		Log("Set(%v, %v, %v)", w.CursorPosX(), 2, newCell)
		time.Sleep(1)
		buffer.Set(w.CursorPosX(), w.CursorPosY(), newCell)
		// Log("buffer: %v", buffer)
		// Log(w.String())
		// Log("buffer.Area: %v", buffer.Area)
		l++
	}
	w.MoveCursor(l)
	w.dirty()
	return len(p), nil
}

func (w *Widget) CursorShow() {
	//w.Log("CursorShow")
	if w.cursorEnabled {
		return
	}
	w.cursorEnabled = true
	w.cursorBlinkState = false
	w.cursorBlink()
}

func (w *Widget) CursorHide() {
	if !w.cursorEnabled {
		return
	}
	w.cursorEnabled = false
}

func (w *Widget) SetCursorLoc(x, y int) bool {
	return false
}

func (w *Widget) SizeOf() int {
	return len(w.TextBuffer().CellMap)
}

func (w *Widget) SetCursorPos(pos int) bool {
	w.cursorPos = pos
	if w.cursorPos < 0 {
		w.cursorPos = 0
	}
	if l := w.SizeOf(); w.cursorPos > l {
		w.cursorPos = l
	}
	if w.cursorEnabled {
		w.CursorHide()
		w.CursorShow()
	}
	return true
}

func (w *Widget) Lines() [][]termui.Cell {
	buffer := w.TextBuffer()
	bounds := buffer.Bounds()

	lines := make([][]termui.Cell, bounds.Dx())

	for i := 0; i < bounds.Dy(); i++ {
		line := make([]termui.Cell, bounds.Dx())
		for p, c := range buffer.CellMap {
			if p.Y != i {
				continue
			}

			line[p.X] = c
		}

		lines[i] = line
	}

	return lines
}

func (w *Widget) CursorPosX() int {
	if w.TextBuffer().Area.Dx() == 0 {
		return 0
	}

	return w.CursorPos()%w.TextBuffer().Area.Dx() + 2
}

func (w *Widget) CursorPosY() int {
	if w.TextBuffer().Area.Dx() == 0 {
		return 0
	}

	return w.CursorPos()/w.TextBuffer().Area.Dx() + 2
}

func (w *Widget) CursorPos() int {
	return w.cursorPos
}

func (w *Widget) MoveCursor(n int) bool {
	return w.SetCursorPos(w.cursorPos + n)
}

func (w *Widget) cursorBlink() {
	//w.Log("cursorBlink")
	w.cursorBlinkState = !w.cursorBlinkState
	w.dirty()
}

func (w *Widget) WrapOn() {
	w.SetWrap(true)
}

func (w *Widget) WrapOff() {
	w.SetWrap(false)
}

func (w *Widget) SetWrap(wrap bool) {
	if wrap != w.wrap {
		defer w.dirty()
	}
	w.wrap = wrap
}

func (w *Widget) Delete(n int) {
	switch {
	case n > 0:
		if l := w.SizeOf(); w.cursorPos+n > l {
			n = l - w.cursorPos
		}
		//w.contents = append(w.contents[:w.cursorPos], w.contents[w.cursorPos+n:]...)
	case n < 0:
		if n < -w.cursorPos {
			n = -w.cursorPos
		}
		//w.contents = append(w.contents[:w.cursorPos+n], w.contents[w.cursorPos:]...)
		w.MoveCursor(n)
	}
	w.dirty()
}
